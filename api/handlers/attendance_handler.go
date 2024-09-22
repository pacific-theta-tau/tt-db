package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

//	"github.com/go-chi/chi"
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
		// return error status code
		http.Error(w, err.Error(), http.StatusInternalServerError)
        fmt.Print()
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
    log.Printf("\tData: \n%v")

	// Build HTTP response
    // use json.MarshalIndent for pretty printing.
	// out, err := json.MarshalIndent(data, "", "\t")
    //if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
    //_, err = w.Write(out)
    //if err != nil {
    //    log.Printf("Error while creating response: %v", err)
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//}

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
		http.Error(w, "Invalid brotherID or eventID", http.StatusBadRequest)
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
