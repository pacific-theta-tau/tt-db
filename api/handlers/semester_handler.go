package handlers

import (
	"context"
	//"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/pacific-theta-tau/tt-db/api/models"
)

// Get all semester labels. E.g.: "Spring 2024"
/* GET /semesters?semester=[optional] */
//	@Summary		Get semester labels
//	@Description	Get all semester labels (e.g. "Spring 2024")
//	@Tags		    Semesters
//	@Success		200		object		models.APIResponse{data=[]string}
//	@Failure		400		{object}	models.APIResponse
//	@Router			/api/semesters [get]
func (h *Handler) GetAllSemesterLabels(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

    query := `SELECT semesterLabel FROM semester`
    rows, err := h.db.QueryContext(ctx, query)
    log.Printf("Querying for semester labels:\n%s", query)
    if err != nil {
        errMsg := fmt.Sprintf("Error while querying for semester data: `%s`\n", err.Error())
        log.Println(errMsg)
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
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
            models.RespondWithError(w, http.StatusInternalServerError, errMsg)
			return
		}

        semesterLabels = append(semesterLabels, &label)
    }

    // Build response
    models.RespondWithSuccess(w, http.StatusOK, semesterLabels)
}


// Helper function to get semesterID related to a semesterLabel
func (h *Handler) GetSemesterIdBySemesterLabel(semesterLabel string) (int, error) {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
    defer cancel()

    query := `
    SELECT semesterID
    FROM semester
    WHERE semesterLabel = $1
    `
    var semesterID int
    err := h.db.QueryRowContext(ctx, query, semesterLabel).Scan(&semesterID)
    log.Printf("Querying for semester labels:\n%s", query)
    if err != nil {
        errMsg := fmt.Sprintf("Error while querying for semesterID: `%s`\n", err.Error())
        log.Println(errMsg)
        return -1, err
    }
    return semesterID, nil
}


// Create new semester label. E.g.: "Fall 2023", "Spring 2024"
/* endpoint: POST /api/semesters */
//	@Summary		Create semester label
//	@Description	Create semester label (e.g. Spring 2024)
//	@Tags		    Semesters 
//	@Param			semester body	string  true	"Semester Label (e.g. `Fall 2023`)"
//	@Success		200		object		models.APIResponse
//	@Failure		400		{object}	models.APIResponse
//	@Router			/api/semesters [post]
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
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
    }
    log.Printf("request body data: `%v%`", requestBody)

    query := `INSERT INTO semester (semesterLabel) VALUES ($1)`
    log.Printf("Query insert:\n%s\n", query)
    _, err = h.db.QueryContext(ctx, query, requestBody.Semester)
    if err != nil {
        errMsg := fmt.Sprintf("Query error:\n\t'%s'\n", err.Error())
        log.Println(errMsg)
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}

    models.RespondWithSuccess(w, http.StatusCreated, "")
} 


// TODO: move status-related endpoint to status-handler.go
//	@Summary		Get Brother statuses for a semester
//	@Description	Get all brother statuses for a semester
//	@Tags		    Semesters
//	@Success		200		object		models.APIResponse{data=models.Attendance}
//	@Failure		400		{object}	models.APIResponse
//	@Router			/api/semesters/{semesterLabel}/statuses [get]
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
    SELECT b.brotherID, b.rollCall, b.firstName, b.lastName, b.major, b.className, bs.status, s.semesterID, s.semesterLabel
    FROM brotherStatus bs
    JOIN brothers b ON b.brotherID = bs.brotherID
    JOIN semester s ON s.semesterID = bs.semesterID
    WHERE semesterLabel = $1
    `
    log.Printf("Query:\n%s\n", query)
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

//	@Summary		Create Brother statuses for a semester
//	@Description	Create all brother statuses for a semester
//	@Tags		    Semesters
//	@Param			semesterLabel path string true	"semesterLabel"
//	@Param			brotherID body int											true	"BrotherID"
//	@Param			status body string true	"Status"
//	@Success		200		object		models.APIResponse{data=models.Attendance}
//	@Failure		400		{object}	models.APIResponse
//	@Router			/api/semesters/{semesterLabel}/statuses [post]
func (h *Handler) CreateBrotherStatusForSemester(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
    defer cancel()

    // urlParams: if none provided, get for all semesters
    semesterLabel := chi.URLParam(r, "semester")
    log.Printf("Semester url param: %s", semesterLabel)
    if semesterLabel == "" {
        errMsg := "Missing semester in query params"
        log.Println(errMsg)
        models.RespondWithFail(w, http.StatusBadRequest, errMsg)
        return
    }

    // Get SemesterID related to semesterLabel
    semesterID, err := h.GetSemesterIdBySemesterLabel(semesterLabel)
    if err != nil {
        errMsg := fmt.Sprintf("Error while querying for semesterID", err.Error())
        log.Println(errMsg)
        models.RespondWithFail(w, http.StatusBadRequest, errMsg)
        return
    }

    // Parse request body params
    type RequestBody struct {
        BrotherID   int `json:"brotherID"`
        Status      string `json:"status"`
    }
    var bodyParams RequestBody
    err = json.NewDecoder(r.Body).Decode(&bodyParams)
    if err != nil {
        errMsg := fmt.Sprintf("Error while decoding body params: %s", err)
        log.Println(errMsg)
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
    }

    // Validate data provided in request body
    validate := validator.New()
	if err := validate.Struct(bodyParams); err != nil {
        errMsg := fmt.Sprintf("Missing or Invalid request body params: %s", err)
        log.Println(errMsg)
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}

    // Query INSERT
    query := `
    INSERT INTO brotherStatus (brotherID, semesterID, status)
    VALUES
        ($1, $2, $3)
    RETURNING *
    `
    _, err = h.db.QueryContext(
        ctx,
        query,
        bodyParams.BrotherID,
        semesterID,
        bodyParams.Status,
    )
    if err != nil {
        errMsg := fmt.Sprintf("Error while querying brother statuses for semester %s: %s", semesterLabel, err.Error())
        log.Println(errMsg)
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
        return
    }
    
    models.RespondWithSuccess(w, http.StatusCreated, "")
}

