package handlers

import (
	"context"
	//"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pacific-theta-tau/tt-db/api/models"
)

// Get all semester labels. E.g.: "Spring 2024"
/* GET /semesters?semester=[optional] */
func (h *Handler) GetAllSemesterLabels(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

    query := `SELECT semesterID, semesterLabel FROM semester`
    rows, err := h.db.QueryContext(ctx, query)
    log.Printf("Querying for semester labels:\n%s", query)
    if err != nil {
        errMsg := fmt.Sprintf("Error while querying for semester data: `%s`\n", err.Error())
        log.Println(errMsg)
        http.Error(w, errMsg, http.StatusInternalServerError)
        return
    }

    var semesterLabels []*models.Semester
    for rows.Next() {
        semester, err := models.CreateSemesterFromRow(rows)
        if err != nil {
            errMsg := fmt.Sprintf("Error creating Brother object from row: '%s'\n", err.Error())
			log.Println(errMsg)
			http.Error(w, errMsg, http.StatusInternalServerError)
			return
		}

        semesterLabels = append(semesterLabels, &semester)
    }
    log.Println("Parsed semester data successfully")

    log.Println("Building response\n")
    // Build response
    if err := json.NewEncoder(w).Encode(semesterLabels); err != nil {
        errMsg := fmt.Sprintf("Erro while encoding response: '%s'", err.Error())
        log.Println(errMsg)
        http.Error(w, errMsg, http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
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
