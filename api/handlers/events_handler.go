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
		// log.Fatal(err)
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
        log.Fatalf("Failed to query event with eventID %d: %s", eventID, err)
    }

	if event.EventID == 0 {
        log.Printf("EventID %d not found", eventID)
		http.Error(w, "Event not found", http.StatusNotFound)
		return
	}

	// Build HTTP response
	out, err := json.MarshalIndent(event, "", "\t")
	if err != nil {
		log.Fatal(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(out)

	if err != nil {
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
		log.Fatal(err)
		return models.Event{}, err
	}
	defer row.Close()
	
	// Scan rows to create Brother instance
	var event models.Event
	for row.Next() {
		event, err = createEventFromRow(row)
		if err != nil {
			log.Fatal("Error!", err)
			return models.Event{}, err
		}
	}

    return event, nil
}

// POST endpoint for Event table
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
    log.Printf("Querying for categoryID of %s", event.CategoryName)
    var categoryID int
    query := "SELECT categoryID FROM eventsCategory WHERE categoryName = $1"
    err = h.db.QueryRow(query, event.CategoryName).Scan(&categoryID)
    if err != nil {
        http.Error(w, "Category not found", http.StatusNotFound)
        return
    }

    // 2. Insert new event in `events` table
    log.Println("Inserting new event")
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

    log.Println("Successfully inserted event to `events` table")
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Inserted new Event entry to `events` table"))
}

// PUT request to update fields of event row by event ID
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
    log.Printf("request body: %+v", requestBody)

	parsedEventID, ok := requestBody["eventID"].(float64)
	if !ok {
		http.Error(w, "Missing eventID in request body", http.StatusBadRequest)
        log.Printf("\tError: Missing eventID in request body")
		return
	}
    eventID := int(parsedEventID)
    log.Printf("Request Body: %+v", requestBody)

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
        log.Fatalf("Failed to query event with eventID %d: %s", eventID, err)
    }

    // API Response
    w.WriteHeader(http.StatusOK)
    //w.Write([]byte(fmt.Sprintf("Updated event with eventID %s successfully", eventID)))
    json.NewEncoder(w).Encode(event)
}
