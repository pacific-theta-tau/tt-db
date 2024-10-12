package handlers

import (
	"context"
	//"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/pacific-theta-tau/tt-db/api/models"
)

// GET /api/statuses
func (h *Handler) GetAllStatusLabels(w http.ResponseWriter, r *http.Request) {
    // Hardcoded since we don't plan to modify brotherStatus table in databse
    statusLabels := [6]string{"Active", "Pre-Alumnus", "Alumnus", "Co-op", "Transferred", "Expelled"}
    models.RespondWithSuccess(w, http.StatusOK, statusLabels)
}


// TODO: this is the same as POST `/api/brothers/{id}/statuses endpoint. Refactor or remove one
//      make POST /api/brohters/statuses body receive a `brotherID`, while POST /api/brothers/{id}/statuses uses urlParams
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
