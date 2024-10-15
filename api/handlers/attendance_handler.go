package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
    "strconv"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/pacific-theta-tau/tt-db/api/models"
)


const ATTENDANCE_TABLE = "attendance"

// Create a Attendance struct from row data
func createAttendanceFromRow(row *sql.Rows) (models.Attendance, error) {
    var attendance models.Attendance
    err := row.Scan(
        &attendance.BrotherID,
        &attendance.EventID,
        &attendance.AttendanceStatus,
        &attendance.RollCall,
        &attendance.FirstName,
        &attendance.LastName,
        &attendance.EventName,
        &attendance.EventLocation,
        &attendance.EventDate,
        &attendance.EventCategory,
    )
    if err != nil {
        return models.Attendance{}, err
    }

    return attendance, err
}


// GET /api/attendance
//	@Summary		Get all attendance records
//	@Description	Get attendance data for all events
//	@Tags			Attendance
//	@Success		200		object		models.APIResponse{data=models.Attendance}
//	@Failure		400		{object}	models.APIResponse
//	@Router			/api/attendance [get]
func (h *Handler) GetAllAttendanceRecords(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
    SELECT a.brotherID, a.eventID, a.attendanceStatus, b.rollCall, b.FirstName, b.LastName, e.EventName, e.eventLocation, e.eventDate, ec.categoryName
    FROM attendance a
    JOIN brothers b ON b.brotherID = a.brotherID
    JOIN events e ON e.eventID = a.eventID
    JOIN eventsCategory ec ON ec.categoryID = e.categoryID
    `
	rows, err := h.db.QueryContext(ctx, query)
	if err != nil {
        errMsg := fmt.Sprintf("Error while querying for Attendance records: %s", err.Error())
        log.Print(errMsg)
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}

	// Read rows from query to create Brother instances
    var attendance []*models.Attendance
	for rows.Next() {
		record, err := createAttendanceFromRow(rows)
		if err != nil {
            errMsg := fmt.Sprintf("Error while parsing Attendance records: %s", err.Error())
            models.RespondWithError(w, http.StatusInternalServerError, errMsg)
			return
		}
        attendance = append(attendance, &record)
	}

    log.Println("\tQuery Successful")
    data := attendance

	// Build HTTP response
    w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

    if len(data) == 0 {
        json.NewEncoder(w).Encode(struct{}{})
        return
    }
    
    json.NewEncoder(w).Encode(data)
}


// GET /api/attendance/{eventID}
//	@Summary		Get all Brothers data
//	@Description	Get data from all Brother records in `Brothers` table
//	@Tags			Attendance
//  @Param          eventID      path        int true "EventID"
//	@Success		200		object		models.APIResponse{data=models.Brother}
//	@Failure		400		{object}	models.APIResponse
//	@Router			/api/attendance/{id} [get]
func (h *Handler) GetAttendanceFromEventID(w http.ResponseWriter, r *http.Request) {
    eventIDStr := chi.URLParam(r, "eventID")
    eventID, err := strconv.Atoi(eventIDStr)
    if err != nil {
        errMsg := "Invalid event ID"
        log.Printf(errMsg, err)
        models.RespondWithFail(w, http.StatusBadRequest, errMsg)
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
    defer cancel()

	query := `
    SELECT a.brotherID, a.eventID, a.attendanceStatus, b.rollCall, b.FirstName, b.LastName, e.EventName, e.eventLocation, e.eventDate, ec.categoryName
    FROM attendance a
    JOIN brothers b ON b.brotherID = a.brotherID
    JOIN events e ON e.eventID = a.eventID
    JOIN eventsCategory ec ON ec.categoryID = e.categoryID
    WHERE a.eventID = $1
    `
    log.Printf("\tEventID: %d", eventID)
    rows, err := h.db.QueryContext(ctx, query, eventID)
	if err != nil {
        errMsg := fmt.Sprintf("Error while querying for Attendance Record: %s", err.Error())
        log.Println(errMsg)
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}

	// Read rows from query to create Brother instances
    var attendance []*models.Attendance
	for rows.Next() {
		record, err := createAttendanceFromRow(rows)
		if err != nil {
            errMsg := fmt.Sprintf("Error while parsing Attendance Record: %s", err.Error())
            log.Println(errMsg)
            models.RespondWithError(w, http.StatusInternalServerError, errMsg)
			return
		}
        attendance = append(attendance, &record)
	}

    log.Println("\tQuery Successful")
    data := attendance

	// Build HTTP response
    // TODO: test
    models.RespondWithSuccess(w, http.StatusOK, data)
}


//	@Summary		Create attendance record
//	@Description	Create attendance record
//	@Tags		    Attendance	
//	@Success		200		object		models.APIResponse
//	@Failure		400		{object}	models.APIResponse
//	@Router			/api/attendance [post]
func (h *Handler) CreateAttendance(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
    defer cancel()
        
    // Expected input in request body
    var input struct {
        BrotherID           int    `json:"brotherID"`
        EventID             int     `json:"eventID"`
        AttendanceStatus    string `json:"attendanceStatus"`
    }
    err := json.NewDecoder(r.Body).Decode(&input)
    if err != nil {
        errMsg := fmt.Sprintf("Error while parsing request body params: %s", err.Error())
        log.Println(errMsg)
        models.RespondWithFail(w, http.StatusBadRequest, errMsg)
        return
    }

    // Validate data provided in request body
    validate := validator.New()
	if err := validate.Struct(input); err != nil {
        errMsg := fmt.Sprintf("Invalid body params: %s", err.Error())
        log.Println(errMsg)
        models.RespondWithFail(w, http.StatusBadRequest, errMsg)
		return
	}

    fmt.Printf("\tRequest Body: %+v", input)

    // Check for missing or zero values
	if input.BrotherID == 0 || input.EventID == 0 {
        errMsg := fmt.Sprintf("Missing brotherID or eventID", err.Error())
        log.Println(errMsg)
        models.RespondWithFail(w, http.StatusBadRequest, errMsg)
		return
	}

    query := `
    INSERT INTO attendance (brotherID, eventID, attendanceStatus)
    VALUES ($1, $2, $3)
    RETURNING *
    `
    
    _, err = h.db.ExecContext(
        ctx,
        query,
        &input.BrotherID,
        &input.EventID,
        &input.AttendanceStatus,
    )
    if err != nil {
        errMsg := fmt.Sprintf("Error while inserting attendance record to table: %s", err.Error())
        log.Println(errMsg)
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}

    models.RespondWithSuccess(w, http.StatusCreated, "")
}


//	@Summary		Delete attendance record
//	@Description	Delete attendance record
//	@Tags		    Attendance	
//	@Success		200		object		models.APIResponse
//	@Failure		400		{object}	models.APIResponse
//	@Router			/api/attendance [delete]
func (h *Handler) DeleteAttendanceRecord(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
    defer cancel()

    // Expected data in request body
    var requestBody struct {
        BrotherID   int `json:"brotherID"`
        EventID     int `json:"eventID"`
    }   
    err := json.NewDecoder(r.Body).Decode(&requestBody)
    if err != nil {
        errMsg := fmt.Sprintf("Error while parsing request body params: %s", err.Error())
        log.Println(errMsg)
        models.RespondWithFail(w, http.StatusBadRequest, errMsg)
        return
    }
    // Validate data provided in request body
    validate := validator.New()
	if err := validate.Struct(requestBody); err != nil {
        errMsg := fmt.Sprintf("Invalid or missing params in request body params: %s", err.Error())
        log.Println(errMsg)
        models.RespondWithFail(w, http.StatusBadRequest, errMsg)
		return
	}

    fmt.Printf("\tRequest Body: %+v", requestBody)

    // Check for missing or zero values
	if requestBody.BrotherID == 0 || requestBody.EventID == 0 {
        errMsg := "Missing BrotherID and/or EventID in request body params"
        log.Println(errMsg)
        models.RespondWithFail(w, http.StatusBadRequest, errMsg)
		return
	}

    query := `
    DELETE FROM attendance
    WHERE brotherID = $1 AND eventID = $2
    `
    _, err = h.db.ExecContext(
        ctx,
        query,
        requestBody.BrotherID,
        requestBody.EventID,
    )
    if err != nil {
        errMsg := fmt.Sprintf("Error while deleting attendance record: %s", err.Error())
        log.Println(errMsg)
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}

    models.RespondWithSuccess(w, http.StatusOK, "")
}


//	@Summary		Update attendance record
//	@Description	Update attendance record
//	@Tags		    Attendance	
//	@Success		200		object		models.APIResponse
//	@Failure		400		{object}	models.APIResponse
//	@Router			/api/attendance [put]
func (h *Handler) UpdateAttendanceRecord(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
    defer cancel()

    var requestBody struct {
        BrotherID           int `json:"brotherID"`
        EventID             int `json:"eventID"`
        AttendanceStatus    string `json:"attendanceStatus"`
    }

    // Unmarshal request body data
    err := json.NewDecoder(r.Body).Decode(&requestBody)
    if err != nil {
        errMsg := fmt.Sprintf("Failed to parse request body params: %s", err.Error())
        log.Println(errMsg)
        models.RespondWithError(w, http.StatusBadRequest, errMsg)
        return
    }

    // Validate data provided in request body
    validate := validator.New()
	if err := validate.Struct(requestBody); err != nil {
        errMsg := fmt.Sprintf("Invalid request body params: %s", err.Error())
        log.Println(errMsg)
        models.RespondWithFail(w, http.StatusBadRequest, errMsg)
		return
	}

    fmt.Printf("\tRequest Body: %+v", requestBody)

    // Check for missing or zero values
	if requestBody.BrotherID == 0 || requestBody.EventID == 0 {
        errMsg := fmt.Sprintf("Invalid brotherID or eventID", err.Error())
        log.Println(errMsg)
        models.RespondWithFail(w, http.StatusBadRequest, errMsg)
		return
	}
    // validate attendance status
    _, ok := models.AttendanceStatus[requestBody.AttendanceStatus]; if !ok {
        // TODO: print valid statues dynamically instead of hardcoding
        errMsg := "Invalid attendance status. Must be one of: 'present', 'absent', or 'excused'"
        log.Println(errMsg)
        models.RespondWithFail(w, http.StatusBadRequest, errMsg)
		return
    }

    query := `
    UPDATE attendance
    SET attendanceStatus = $1
    WHERE brotherID = $2 AND eventID = $3
    `
    _, err = h.db.QueryContext(
        ctx,
        query,
        requestBody.AttendanceStatus,
        requestBody.BrotherID,
        requestBody.EventID,
    )
    if err != nil {
        errMsg := fmt.Sprintf("Error while updating attendance record: %s", err.Error())
        log.Println(errMsg)
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}

    models.RespondWithSuccess(w, http.StatusOK, "")
}

