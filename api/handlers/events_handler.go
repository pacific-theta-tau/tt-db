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

//	@Summary		Get all event records
//	@Description	Get data from all rows in events table
//	@Tags			Events
//	@Success		200		object		models.APIResponse{data=models.Event}
//	@Failure		400		{object}	models.APIResponse
//	@Router			/api/events [get]
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
        errMsg := fmt.Sprintf("Error while querying events for events table: %s", err.Error())
        log.Println(errMsg)
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}

	// Read rows from query to create Event instances
	var events []*models.Event
	for rows.Next() {
		event, err := createEventFromRow(rows)
		if err != nil {
            errMsg := fmt.Sprintf("Error while parsing rows for events table: %s", err.Error())
            log.Println(errMsg)
            models.RespondWithError(w, http.StatusInternalServerError, errMsg)
			return
		}
		events = append(events, &event)
	}

    models.RespondWithSuccess(w, http.StatusOK, events)
}

// Query Event by their eventID
//	@Summary		Get event data 
//	@Description	Get event information by eventID
//	@Tags			Events
//	@Param			eventid		path		int											true	"Event ID"
//	@Success		200		object		models.APIResponse{data=models.Event}
//	@Failure		400		{object}	models.APIResponse
//	@Router			/api/events/{eventid} [get]
func (h *Handler) GetEventByEventID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	requestEventID := chi.URLParam(r, "eventID")
	if requestEventID == "" {
		// If eventID is empty, return an error response
        errMsg := fmt.Sprintf("Missing eventID parameter in query params")
        log.Println(errMsg)
        models.RespondWithFail(w, http.StatusBadRequest, errMsg)
        return
	}

    eventID, err := strconv.Atoi(requestEventID)
    if err != nil {
        errMsg := fmt.Sprintf("Error while parsing eventID in query params: %s", err.Error())
        log.Println(errMsg)
        models.RespondWithFail(w, http.StatusInternalServerError, errMsg)
        return
    }

    event, err := queryEvent(h, ctx, eventID)
    if err != nil {
        errMsg := fmt.Sprintf("Failed to query event with eventID %d: %s", eventID, err.Error())
        log.Println(errMsg)
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
        return
    }

	if event.EventID == 0 {
        errMsg := fmt.Sprintf("EventID %d not found", eventID)
        log.Println(errMsg)
        models.RespondWithFail(w, http.StatusBadRequest, errMsg)
		return
	}

	// Build HTTP response
    models.RespondWithSuccess(w, http.StatusOK, event)
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
//	@Summary		Create new event record
//	@Description	Create new event record
//	@Tags			Events
//	@Param			body body models.Event true	"Values for new event record"
//	@Success		200		{object}		models.APIResponse{data=models.Event}
//	@Failure		400		{object}	models.APIResponse
//	@Router			/api/events [post]
func (h* Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
    defer cancel()

    var event models.Event
    err := json.NewDecoder(r.Body).Decode(&event)
    if err != nil {
        errMsg := fmt.Sprintf("Error decoding request body", err.Error())
        log.Println(errMsg)
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
        return
    }

    // Validate events struct
    validate := validator.New() 
    if err := validate.Struct(event); err != nil {
        errMsg := fmt.Sprintf("Error validating body params. Missing values: %s", err.Error())
        log.Println(errMsg)
        models.RespondWithFail(w, http.StatusBadRequest, errMsg)
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
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
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
    models.RespondWithSuccess(w, http.StatusOK, "")
}

// Update fields of event row by event ID
//	@Summary		Update event record
//	@Description	Update event record by eventID
//	@Tags			Events
//	@Param			eventid		body int											true	"Event ID"
//	@Success		200		object		models.APIResponse
//	@Failure		400		{object}	models.APIResponse
//	@Router			/api/events [put]
func (h *Handler) UpdateEventByID(w http.ResponseWriter, r *http.Request) {
    ctx, cancel :=  context.WithTimeout(context.Background(), dbTimeout)
    defer cancel()

    body, err := io.ReadAll(r.Body)
	if err != nil {
        errMsg := fmt.Sprintf("Error while reading request body: %s", err.Error())
        log.Printf(errMsg)
        models.RespondWithError(w, http.StatusInternalServerError, errMsg) 
        return
	}

	var requestBody map[string]interface{}
	if err = json.Unmarshal(body, &requestBody); err != nil {
        errMsg := "Error decoding JSON"
        log.Printf(errMsg)
        models.RespondWithFail(w, http.StatusInternalServerError, errMsg)
		return
	}

	parsedEventID, ok := requestBody["eventID"].(float64)
	if !ok {
        errMsg := "Missing eventID in request body"
        log.Printf(errMsg)
        models.RespondWithFail(w, http.StatusBadRequest, errMsg)
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
                errMsg := fmt.Sprintf("\tError while querying for Category - not found: %s", err.Error())
                log.Println(errMsg)
                models.RespondWithFail(w, http.StatusBadRequest, errMsg)
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
        errMsg := fmt.Sprintf("Error while querying update: %s", err.Error())
        log.Printf(errMsg)
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
		return
    }
    log.Printf("query result: %+v", queryResult)

    // query modified row
    // TODO: use RETURNING row clause instead of querying for the modified row
    event, err := queryEvent(h, ctx, eventID)
    if err != nil {
        errMsg := fmt.Sprintf("Failed to query event with eventID %d after update query: %s", eventID, err)
        log.Printf(errMsg)
        models.RespondWithFail(w, http.StatusInternalServerError, errMsg)
        return
    }

    models.RespondWithSuccess(w, http.StatusOK, event)
}


//	@Summary		Get event and attendance data
//	@Description	Get event and attendance data by eventID
//	@Tags			Events
//	@Param			eventid		path		int											true	"Event ID"
//	@Success		200		object		models.APIResponse{data=handlers.GetEventAttendance.EventDataAndAttendance}
//	@Failure		400		{object}	models.APIResponse
//	@Router			/api/events/{eventid}/attendance [get]
func (h *Handler) GetEventAttendance(w http.ResponseWriter, r *http.Request) {
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

    // Query event data
    eventData, err := queryEvent(h, ctx, eventID)
    if err != nil {
        errMsg := fmt.Sprintf("Error while fetching event data for eventID %d: %s", eventID, err.Error())
        log.Println(errMsg)
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
        return
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
        log.Println(errMsg)
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}

    var attendanceList []*models.EventAttendance
	for rows.Next() {
        record, err := createEventAttendanceFromRow(rows)
		if err != nil {
            errMsg := fmt.Sprintf("Error while parsing attendance query for eventID %d: '%s'", eventID, err.Error())
            log.Println(errMsg)
            models.RespondWithError(w, http.StatusInternalServerError, errMsg)
			return
		}
        attendanceList = append(attendanceList, &record)
	}

    // TODO: include category name
    type EventDataAndAttendance struct {
        EventID			int			`json:"eventID"`  //primary
        EventName		string 		`json:"eventName"`
        EventCategory   string      `json:"eventCategory"` 
        EventLocation	string 		`json:"eventLocation"`
        EventDate		time.Time	`json:"eventDate"`
        Attendance      []*models.EventAttendance `json:"attendance"`
    }
    // Build response
    response := EventDataAndAttendance{
        EventID: eventData.EventID,
        EventName: eventData.EventName,
        EventCategory: eventData.CategoryName,
        EventDate: eventData.EventDate,
        EventLocation: eventData.EventLocation,
        Attendance: attendanceList,
    }
    models.RespondWithSuccess(w, http.StatusOK, response)
}


// Add new attendance record for event
//	@Summary		Create new event record
//	@Description	Create new event record
//	@Tags			Events
//	@Param			eventid path int true	"eventID"
//	@Success		200		{object}		models.APIResponse{data=models.Event}
//	@Failure		400		{object}	models.APIResponse
//	@Router			/api/events/{eventid}/attendance [post]
func (h* Handler) CreateAttendanceRecordForEvent(w http.ResponseWriter, r *http.Request) {
    // TODO: fix swagger docs
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
    defer cancel()

    // Get event from URL Params
    var eventID int
    eventIDStr := chi.URLParam(r, "eventID")
    eventID, err := strconv.Atoi(eventIDStr)
    if err != nil {
        errMsg := "Invalid event ID"
        log.Printf(errMsg, err)
        models.RespondWithFail(w, http.StatusBadRequest, errMsg)
        return
    }

    // Parse request body params
    var requestBody struct {
        RollCall int `json:"rollCall"`
        AttendanceStatus string `json:"attendanceStatus"`
    }
    err = json.NewDecoder(r.Body).Decode(&requestBody)
    if err != nil {
        errMsg := fmt.Sprintf("Error decoding request body", err.Error())
        log.Println(errMsg)
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
        return
    }
    fmt.Printf("Event %d; Request Body: %v\n", eventID, requestBody)

    // Validate request body params 
    validate := validator.New() 
    if err := validate.Struct(requestBody); err != nil {
        errMsg := fmt.Sprintf("Error validating body params. Missing values: %s", err.Error())
        log.Println(errMsg)
        models.RespondWithFail(w, http.StatusBadRequest, errMsg)
        return
	}

    // Insert new attendance record for eventID
    log.Println("\tInserting new attendance record\n")
    query := `
    INSERT INTO attendance (eventID, brotherID, attendanceStatus)
    SELECT $1, b.brotherID, $2
    FROM brothers b
    WHERE b.rollCall = $3
    `
    log.Printf("query: %s\n", query)
    _, err = h.db.ExecContext(
        ctx,
        query,
        eventID,
        requestBody.AttendanceStatus,
        requestBody.RollCall,
    )
    if err != nil {
        errMsg := fmt.Sprintf("Error while querying: %s", err.Error())
        log.Printf(errMsg)
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
		return
    }

    msg := "Created new attendance record to `attendance` table successfully"
    log.Println(msg)
    models.RespondWithSuccess(w, http.StatusCreated, "")
}

