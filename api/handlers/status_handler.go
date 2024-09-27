package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/pacific-theta-tau/tt-db/api/models"
)

// GET /semesters?semester=[optional]
func (h *Handler) GetAllBrotherStatuses(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
    semester := r.URL.Query().Get("semester")

    // Query for all brother statuses
    query := `
    SELECT b.brotherID, b.rollCall, b.firstName, b.lastName, bs.status, s.semesterLabel
    FROM brotherStatus bs
    JOIN brothers b ON b.brotherID = bs.brotherID
    JOIN semester s ON s.semesterID = bs.semesterID
    `
    var rows *sql.Rows
    var err error
    if semester == "" {
        // Query without filter
        rows, err = h.db.QueryContext(ctx, query)
    } else {
        // query filtering by semester
        query += " WHERE semesterLabel = $1"
        rows, err = h.db.QueryContext(
            ctx,
            query,
            semester,
        )
    } 

    // Error handling query
    fmt.Printf("Querying for all semester statuses:\n%s\n", query)
    if err != nil {
        errMsg := fmt.Sprintf("Error while querying for all brother statuses: %s", err.Error())
        log.Println(errMsg)
        http.Error(w, errMsg, http.StatusInternalServerError)
        return
    }

    log.Println("Parsing brother objects\n")
    // Parse query rows
    var brotherStatuses []*models.BrotherStatus
    for rows.Next(){
        brotherStatus, err := models.CreateBrotherStatusFromRow(rows)
		if err != nil {
            errMsg := fmt.Sprintf("Error while parsing brotherStatus query: '%s'", err.Error())
            log.Println(errMsg)
			http.Error(w, errMsg, http.StatusInternalServerError)
			return
        }
        brotherStatuses = append(brotherStatuses, &brotherStatus)
    }

    log.Println("Building response\n")
    // Build response
    if err := json.NewEncoder(w).Encode(brotherStatuses); err != nil {
        errMsg := fmt.Sprintf("Erro while encoding response: '%s'", err.Error())
        log.Println(errMsg)
        http.Error(w, errMsg, http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// TODO: this is the same as POST `/api/brothers/{id}/statuses endpoint. Refactor or remove one
/* POST /statuses */
func (h *Handler) CreateStatusForBrother(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

    var requestBody struct {
        BrotherID   int `json:"brotherID"`
        SemesterID  int `json:"semesterID"`
        Status      string `json:"status"`
    }
    // Parse body
    err := json.NewDecoder(r.Body).Decode(&requestBody)
    if err != nil {
        errMsg := fmt.Sprintf("Error while parsing request body data: %s", err.Error())
        log.Println(errMsg)
        http.Error(w, errMsg, http.StatusInternalServerError)
        return
    }
    // Validate received data
    validate := validator.New()
    if err := validate.Struct(requestBody); err != nil {
        errMsg := fmt.Sprintf("Invalid Input: %s", err.Error())
        log.Println(errMsg)
        http.Error(w, errMsg, http.StatusBadRequest)
        return
    }

    // Validate received data
    query := `
    INSERT INTO brotherStatus (brotherID, semesterID, status)
    VALUES ($1, $2, $3)
    `
    _, err = h.db.QueryContext(
        ctx,
        query,
        requestBody.BrotherID,
        requestBody.SemesterID,
        requestBody.Status,
    )
    if err != nil {
        errMsg := fmt.Sprintf("Query error:\n\t'%s'\n", err.Error())
        log.Println(errMsg)
        log.Println("Sending HTTP error response\n")
		http.Error(w, errMsg, http.StatusInternalServerError)
        log.Println("After sending error response\n")
		return
	}

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
	w.Write([]byte("Created Brother Status successfully"))

}
