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
		http.Error(w, err.Error(), http.StatusInternalServerError)
        log.Print(err.Error())
		return
	}

	// Read rows from query to create Brother instances
    var attendance []*models.Attendance
	for rows.Next() {
		record, err := createAttendanceFromRow(rows)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
func (h *Handler) GetAttendanceFromEventID(w http.ResponseWriter, r *http.Request) {
    eventIDStr := chi.URLParam(r, "eventID")
    eventID, err := strconv.Atoi(eventIDStr)
    if err != nil {
        errMsg := "Invalid event ID"
        log.Printf(errMsg, err)
        http.Error(w, errMsg, http.StatusBadRequest)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
        log.Print(err.Error())
		return
	}

	// Read rows from query to create Brother instances
    var attendance []*models.Attendance
	for rows.Next() {
		record, err := createAttendanceFromRow(rows)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    // Validate data provided in request body
    validate := validator.New()
	if err := validate.Struct(input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

    fmt.Printf("\tRequest Body: %+v", input)

    // Check for missing or zero values
	if input.BrotherID == 0 || input.EventID == 0 {
		http.Error(w, "Missing brotherID or eventID", http.StatusBadRequest)
		return
	}

    query := `
    INSERT INTO attendance (brotherID, eventID, attendanceStatus)
    VALUES ($1, $2, $3)
    RETURNING *
    `
    
    result, err := h.db.ExecContext(
        ctx,
        query,
        &input.BrotherID,
        &input.EventID,
        &input.AttendanceStatus,
    )
    if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

    fmt.Println(result)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Inserted new record into `attendance` table successfully!"))
}

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
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }
    // Validate data provided in request body
    validate := validator.New()
	if err := validate.Struct(requestBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

    fmt.Printf("\tRequest Body: %+v", requestBody)

    // Check for missing or zero values
	if requestBody.BrotherID == 0 || requestBody.EventID == 0 {
		http.Error(w, "Invalid brotherID or eventID", http.StatusBadRequest)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
	w.Write([]byte("Deleted attendance record successfully"))
}

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
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    // Validate data provided in request body
    validate := validator.New()
	if err := validate.Struct(requestBody); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

    fmt.Printf("\tRequest Body: %+v", requestBody)

    // Check for missing or zero values
	if requestBody.BrotherID == 0 || requestBody.EventID == 0 {
		http.Error(w, "Invalid brotherID or eventID", http.StatusBadRequest)
		return
	}
    // validate attendance status
    _, ok := models.AttendanceStatus[requestBody.AttendanceStatus]; if !ok {
        // TODO: print valid statues dynamically instead of hardcoding
        errMsg := "Invalid attendance status. Must be one of: 'present', 'absent', or 'excused'"
        log.Println(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
	w.Write([]byte("Updated attendance record successfully"))
}

