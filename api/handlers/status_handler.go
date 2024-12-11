// status_handler.go: Handle requests for getting active status from BrotherStatus table
package handlers

import (
	"context"
	"strconv"
	//"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/pacific-theta-tau/tt-db/api/models"
)

// GET /api/statuses
//	@Summary		Get status labels
//	@Description	Get all valid status labels (e.g.: "Active")
//	@Tags			Statuses
//	@Success		200		object		models.APIResponse{data=[]string}
//	@Failure		400		{object}	models.APIResponse
//	@Router			/api/statuses [get]
func (h *Handler) GetAllStatusLabels(w http.ResponseWriter, r *http.Request) {
    // Hardcoded since we don't plan to modify brotherStatus table in databse
    statusLabels := [8]string{"Active", "Pre-Alumnus", "Alumnus", "Co-op", "Transferred", "Expelled", "Inactive", "Out of Contact"}
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
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
        return
    }

    // Validate received data
    validate := validator.New()
    if err := validate.Struct(requestBody); err != nil {
        errMsg := fmt.Sprintf("Invalid Input: %s", err.Error())
        log.Println(errMsg)
        models.RespondWithFail(w, http.StatusBadRequest, errMsg)
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
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}

    models.RespondWithSuccess(w, http.StatusCreated, "")
}


// Delete /api/brothers/{brotherID}/statuses/{semesterID}
//	@Summary		Deletes the status of brother for specified semester
//	@Description	Deletes the status of the specified brother for the specified semester.
//	@Tags			Statuses
//  @Param  brotherID path   string 
//  @Param  semesterID path   string 
//	@Success		200		object		models.APIResponse{data=[]string}
//	@Failure		400		{object}	models.APIResponse
//	@Router			/api/brothers/{brotherID}/statuses/{semesterID} [delete]
func (h* Handler) DeleteStatusByMemberAndSemesterHandler(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
    defer cancel()

    brotherID := chi.URLParam(r, "brotherID")
    semesterID := chi.URLParam(r, "semesterID")

    if brotherID == "" || semesterID == "" {
        errMsg := "brotherID and semesterID are required"
        log.Println(errMsg)
        models.RespondWithError(w, http.StatusBadRequest, errMsg)
        return
    }

    query := "DELETE FROM brotherStatus WHERE brotherID = $1 AND semesterID = $2"
    _, err := h.db.ExecContext(ctx, query, brotherID, semesterID)
    if err != nil {
        errMsg := fmt.Sprintf("Error while querying:\n\t'%s'\n", err.Error())
        log.Println(errMsg)
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}

    models.RespondWithSuccess(w, http.StatusOK, "Deleted row successfully")
}


// Patch /api/brothers/{brotherID}/statuses
//	@Summary		Deletes the status of brother for specified semester
//	@Description	Deletes the status of the specified brother for the specified semester.
//	@Tags			Statuses
//  @Param  brotherID path   string 
//  @Param  semesterID body int
//  @Param  status body string 
//	@Success		200		object		models.APIResponse{data=[]string}
//	@Failure		400		{object}	models.APIResponse
//	@Router			/api/brothers/{brotherID}/statuses [patch]
func (h* Handler) UpdateBrotherStatusByBrotherID(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
    defer cancel()

    // parse url params
    brotherID := chi.URLParam(r, "id")
    if brotherID == "" {
        errMsg := "Missing brotherID in url params"
        log.Println(errMsg)
        models.RespondWithFail(w, http.StatusBadRequest, errMsg)
        return
    }
    brotherIDInt, err  := strconv.Atoi(brotherID)
    if err != nil {
        errMsg := fmt.Sprintf("Error converting brotherID to int. Error: %s", err.Error())
        log.Println(errMsg)
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
        return
    }

    // parse request body
    var requestBody struct {
        Status string `json:"status"`
        SemesterID int `json:"semesterID"`
        //SemesterLabel string `json:"semesterLabel"`
    }
    err = json.NewDecoder(r.Body).Decode(&requestBody)
    if err != nil {
        errMsg := fmt.Sprintf("Error while parsing request body. Error: %s", err.Error())
        log.Println(errMsg)
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
        return
    }
    
    // validate request body
    validate := validator.New()
    if err := validate.Struct(requestBody); err != nil {
        errMsg := fmt.Sprintf("Invalid Input: %s", err.Error())
        log.Println(errMsg)
        models.RespondWithFail(w, http.StatusBadRequest, errMsg)
        return
    }

    //// Get SemesterID
    //semesterID, err := h.GetSemesterIdBySemesterLabel(requestBody.SemesterLabel) 
    //if err != nil {
    //    errMsg := fmt.Sprintf("Failed to fetch semesterID for semester label `$s`", requestBody.SemesterLabel)
    //    log.Println(errMsg)
    //    models.RespondWithFail(w, http.StatusBadRequest, errMsg)
    //    return
    //}

    // Query Database
    query := "UPDATE brotherStatus SET status = $1 WHERE brotherID = $2 AND semesterID = $3"
    _, err = h.db.ExecContext(ctx, query, requestBody.Status, brotherIDInt, requestBody.SemesterID)
    if err != nil {
        errMsg := fmt.Sprintf("Error while querying:\n\t'%s'\n", err.Error())
        log.Println(errMsg)
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}

    models.RespondWithSuccess(w, http.StatusOK, "")
}
