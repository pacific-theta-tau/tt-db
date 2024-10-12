package handlers

import (
	"context"
	//"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/pacific-theta-tau/tt-db/api/models"
)

// Get all semester labels. E.g.: "Spring 2024"
/* GET /semesters?semester=[optional] */
func (h *Handler) GetAllSemesterLabels(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

    query := `SELECT semesterLabel FROM semester`
    rows, err := h.db.QueryContext(ctx, query)
    log.Printf("Querying for semester labels:\n%s", query)
    if err != nil {
        errMsg := fmt.Sprintf("Error while querying for semester data: `%s`\n", err.Error())
        log.Println(errMsg)
        http.Error(w, errMsg, http.StatusInternalServerError)
        return
    }

    //var semesterLabels []*models.Semester
    var semesterLabels []*string
    for rows.Next() {
        var label string
        err := rows.Scan(&label)
        if err != nil {
            errMsg := fmt.Sprintf("Error creating Semester label slice from row: '%s'\n", err.Error())
			log.Println(errMsg)
			http.Error(w, errMsg, http.StatusInternalServerError)
			return
		}

        semesterLabels = append(semesterLabels, &label)
    }
    log.Println("Parsed semester data successfully")
    log.Println(semesterLabels)

    log.Println("Building response\n")
    // Build response
    models.RespondWithSuccess(w, http.StatusOK, semesterLabels)
}


// Create new semester label. E.g.: "Fall 2023", "Spring 2024"
/* endpoint: POST /api/semesters */
func (h *Handler) CreateSemesterLabel(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
    defer cancel()

    var requestBody struct {
        Semester    string `json:"semester"`
    }
    err := json.NewDecoder(r.Body).Decode(&requestBody)
    if err != nil {
        errMsg := fmt.Sprintf("Error while parsing request body:\n%s\n", err.Error())
        log.Println(err)
        http.Error(w, errMsg, http.StatusInternalServerError)
    }
    log.Printf("request body data: `%v%`", requestBody)

    query := `INSERT INTO semester (semesterLabel) VALUES ($1)`
    log.Printf("Query insert:\n%s\n", query)
    _, err = h.db.QueryContext(ctx, query, requestBody.Semester)
    if err != nil {
        errMsg := fmt.Sprintf("Query error:\n\t'%s'\n", err.Error())
        log.Println(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
	w.Write([]byte("Created new semester label successfully!"))
} 

func (h *Handler) GetAllBrotherStatusesForSemester(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
    defer cancel()

    // urlParams: if none provided, get for all semesters
    semester := chi.URLParam(r, "semester")
    log.Printf("Semester url param: %s", semester)
    if semester == "" {
        errMsg := "Missing semester in query params"
        log.Println(errMsg)
        models.RespondWithFail(w, http.StatusBadRequest, errMsg)
        return
    }

    query := `
    SELECT b.rollCall, b.firstName, b.lastName, b.major, b.className, bs.status
    FROM brotherStatus bs
    JOIN brothers b ON b.brotherID = bs.brotherID
    JOIN semester s ON s.semesterID = bs.semesterID
    WHERE semesterLabel = $1
    `
    rows, err := h.db.QueryContext(ctx, query, semester)
    if err != nil {
        errMsg := fmt.Sprintf("Error while querying brother statuses for semester %s: %s", semester, err.Error())
        log.Println(errMsg)
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
        return
    }

    var brotherStatuses []*models.BrotherStatusFromSemester
    for rows.Next() {
        brotherStatus, err := models.CreateBrotherStatusFromSemesterFromRow(rows)
        if err != nil {
            errMsg := fmt.Sprintf("Error while parsing query: %s", err.Error())
            log.Println(errMsg)
            models.RespondWithError(w, http.StatusInternalServerError, errMsg)
            return
        }
        brotherStatuses = append(brotherStatuses, &brotherStatus)
    }

    models.RespondWithSuccess(w, http.StatusOK, brotherStatuses)
}
