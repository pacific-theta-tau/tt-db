package handlers

import (
    "fmt"
	"context"
	"database/sql"
	"encoding/json"
    "io"
	"log"
	"net/http"
    "strconv"
    "strings"
	"time"

	"github.com/pacific-theta-tau/tt-db/api/models"
	"github.com/go-chi/chi"
    "github.com/go-playground/validator/v10"
)

const events_table = "events"

// Helper function to scan SQL row and create new event instance
func createEventFromRow(row *sql.Rows) (models.Event, error) {
	var events models.Event
	err := row.Scan(
		&events.EventID,
		&events.EventName,
		&events.CategoryName,
		&events.EventLocation,
		&events.EventDate,
	)
	if err != nil {
		return models.Event{}, err
	}

	return events, err
}

// Helper function to scan SQL row and create new event attendance instance
func createEventAttendanceFromRow(row *sql.Rows) (models.EventAttendance, error) {
	var eventAttendance models.EventAttendance
	err := row.Scan(
        &eventAttendance.BrotherID,
        &eventAttendance.FirstName,
        &eventAttendance.LastName,
        &eventAttendance.RollCall,
        &eventAttendance.AttendanceStatus,
	)
	if err != nil {
		return models.EventAttendance{}, err
	}

	return eventAttendance, err
}

// Get data events table
func (h *Handler) GetAllEvents(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
        SELECT e.eventid, e.eventName, ec.categoryName, e.eventLocation, e.eventDate
        FROM events e
        JOIN eventsCategory ec ON e.categoryID = ec.categoryID
    `
	rows, err := h.db.QueryContext(ctx, query)
	if err != nil {
		// return error status code
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Read rows from query to create Event instances
	var events []*models.Event
	for rows.Next() {
		event, err := createEventFromRow(rows)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		events = append(events, &event)
	}

	// Build HTTP response
	out, err := json.MarshalIndent(events, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(out)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Query Event by their eventID
func (h *Handler) GetEventByEventID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	requestEventID := chi.URLParam(r, "eventID")
	if requestEventID == "" {
		// If eventID is empty, return an error response
		http.Error(w, "eventID parameter is required", http.StatusBadRequest)
		return
	}
    eventID, err := strconv.Atoi(requestEventID)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    event, err := queryEvent(h, ctx, eventID)
    if err != nil {
        errMsg := fmt.Sprintf("Failed to query event with eventID %d: %s", eventID, err.Error())
        log.Println(errMsg)
        http.Error(w, errMsg, http.StatusNotFound)
        return
    }

	if event.EventID == 0 {
        errMsg := fmt.Sprintf("EventID %d not found", eventID)
        log.Println(errMsg)
		http.Error(w, errMsg, http.StatusNotFound)
		return
	}

	// Build HTTP response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
    err = json.NewEncoder(w).Encode(event)
	if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Helper function to query event by eventID
func queryEvent(h* Handler, ctx context.Context, eventID int) (models.Event, error) {
    // TODO: refactor so that we don't need to inject handler and context
    query := `
        SELECT e.eventid, e.eventName, ec.categoryName, e.eventLocation, e.eventDate
        FROM events e
        JOIN eventsCategory ec ON e.categoryID = ec.categoryID
        WHERE eventID = $1
    `

	row, err := h.db.QueryContext(ctx, query, eventID)
	if err != nil {
        errMsg := fmt.Sprintf("\t[queryEvent()] Error while querying event data for eventID %d: %s", eventID, err.Error())
		log.Println(errMsg)
		return models.Event{}, err
	}
	defer row.Close()
	
	// Scan rows to create Brother instance
	var event models.Event
	for row.Next() {
		event, err = createEventFromRow(row)
		if err != nil {
            log.Printf("\t[queryEvent()] Error while parsing rows from database: %s", err.Error())
			return models.Event{}, err
		}
	}

    return event, nil
}

// Add new event to events table
// POST /api/events
func (h* Handler) InsertEvent(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
    defer cancel()

    var event models.Event
    err := json.NewDecoder(r.Body).Decode(&event)
    if err != nil {
        http.Error(w, fmt.Sprint("Error decoding request body", err.Error()), http.StatusBadRequest)
    }

    // Validate events struct
    validate := validator.New() 
    if err := validate.Struct(event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

    // 1. fetch categoryID from categoryName
    log.Printf("\tQuerying for categoryID of %s", event.CategoryName)
    var categoryID int
    query := "SELECT categoryID FROM eventsCategory WHERE categoryName = $1"
    err = h.db.QueryRow(query, event.CategoryName).Scan(&categoryID)
    if err != nil {
        errMsg := fmt.Sprintf("Category not found. %v", err.Error())
        log.Println(errMsg)
        http.Error(w, errMsg, http.StatusNotFound)
        return
    }

    // 2. Insert new event in `events` table
    log.Println("\tInserting new event")
    query = `
    INSERT INTO events (eventName, categoryID, eventLocation, eventDate)
    VALUES ($1, $2, $3, $4) returning *
    `
    log.Println(query)
    _, err = h.db.ExecContext(
        ctx,
        query,
        event.EventName,
        categoryID,
        event.EventLocation,
        event.EventDate,
    )

    msg := "Created new event to `events` table successfully"
    log.Println(msg)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(msg))
}

// Update fields of event row by event ID
// PUT /api/events
func (h *Handler) UpdateEventByID(w http.ResponseWriter, r *http.Request) {
    ctx, cancel :=  context.WithTimeout(context.Background(), dbTimeout)
    defer cancel()

    body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var requestBody map[string]interface{}
	if err = json.Unmarshal(body, &requestBody); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	parsedEventID, ok := requestBody["eventID"].(float64)
	if !ok {
        errMsg := "Missing eventID in request body"
        log.Printf(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
    eventID := int(parsedEventID)

    updateQuery := fmt.Sprintf("UPDATE %s SET", events_table)
    columns := []string{
        "eventName",
        "categoryName",
        "eventLocation",
        "eventDate",
    }
    // Build query statement for each param in requestBody
    for _, column := range columns {
        newColumnValue, ok := requestBody[column]
        if !ok {
            continue
        }
        if column == "categoryName" {
            // 1. fetch categoryID from categoryName
            // TODO: refactor this into own function to avoid duplicate
            var categoryID int
            categoryIdQuery := "SELECT categoryID FROM eventsCategory WHERE categoryName = $1"
            err = h.db.QueryRow(categoryIdQuery, newColumnValue).Scan(&categoryID)
            if err != nil {
                http.Error(w, "Category not found", http.StatusNotFound)
                log.Printf("\tError while querying category: Category not found")
                return
            }
            updateQuery += fmt.Sprintf(" %s = %d,", "categoryID", categoryID)
            continue
        }

        updateQuery += fmt.Sprintf(" %s = '%s',", column, newColumnValue)
    }

    // Remove trailling comma
    updateQuery = strings.TrimRight(updateQuery, ",") 
    // UPDATE events SET (values) FROM eventsCategory ...
    updateQuery += `
        FROM eventsCategory ec
        WHERE eventID = $1 AND events.categoryID = ec.categoryID
        RETURNING events.eventID, events.eventName, events.eventDate, events.eventLocation, ec.categoryName
    `
    log.Printf("UPDATE query: %s", updateQuery)

    var queryResult models.Event
    err = h.db.QueryRowContext(ctx, updateQuery, eventID).Scan(
        &queryResult.EventID,
        &queryResult.EventName,
        &queryResult.EventDate,
        &queryResult.EventLocation,
        &queryResult.CategoryName,
    )
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
		return
    }
    log.Printf("query result: %+v", queryResult)

    // query modified row
    // TODO: refactor this
    event, err := queryEvent(h, ctx, eventID)
    if err != nil {
        log.Printf("Failed to query event with eventID %d: %s", eventID, err)
    }

    // API Response
    w.WriteHeader(http.StatusOK)
    //w.Write([]byte(fmt.Sprintf("Updated event with eventID %s successfully", eventID)))
    json.NewEncoder(w).Encode(event)
}



/* Get event data and attendance
endpoint: GET /api/events/{eventID}/attendance/
Response format:
    {
        eventID: str,
        eventName: str,
        eventCategory: str,
        eventDate: time.Time,
        eventLocation: str,
        attendance: [
            {
                brotherID: int,
                firstName: str,
                lastName: str,
                rollCall: int,
                attendanceStatus: str
            },
        ]
    }
*/
func (h *Handler) GetEventAttendance(w http.ResponseWriter, r *http.Request) {
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

    // Query event data
    eventData, err := queryEvent(h, ctx, eventID)
    if err != nil {
        errMsg := fmt.Sprintf("Error while fetching event data for eventID %d: %s", eventID, err.Error())
        log.Println(errMsg)
        http.Error(w, errMsg, http.StatusBadRequest)
    }

    // Query attendance data for eventID and parse into a list
	query := `
    SELECT b.brotherID, b.firstName, b.lastName, b.rollCall, a.attendanceStatus
    FROM attendance a
    JOIN brothers b ON b.brotherID = a.brotherID
    WHERE eventID = $1
    `
    log.Printf("\tEventID: %d", eventID)
    rows, err := h.db.QueryContext(ctx, query, eventID)
	if err != nil {
        errMsg := fmt.Sprintf("Error while querying for attendance for eventID %d: %s", eventID, err.Error())
		http.Error(w, errMsg, http.StatusInternalServerError)
        log.Println(errMsg)
		return
	}

    var attendanceList []*models.EventAttendance
	for rows.Next() {
        record, err := createEventAttendanceFromRow(rows)
		if err != nil {
            errMsg := fmt.Sprintf("Error while parsing attendance query for eventID %d: '%s'", eventID, err.Error())
            log.Println(errMsg)
			http.Error(w, errMsg, http.StatusInternalServerError)
			return
		}
        attendanceList = append(attendanceList, &record)
	}

    // Build response
    response := map[string]interface{}{
        "eventID": eventData.EventID,
        "eventName": eventData.EventName,
        "eventCategory": eventData.CategoryName,
        "eventDate": eventData.EventDate,
        "eventLocation": eventData.EventLocation,
        "attendance": attendanceList,
    }
    if err := json.NewEncoder(w).Encode(response); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

	// Build HTTP response
    w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

