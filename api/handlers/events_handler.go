package handlers

import (
    "fmt"
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
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
	log.Println(" - Called GetAllEvents")
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	query := `
        SELECT e.eventid, e.eventName, ec.categoryName, e.eventLocation, e.eventDate
        FROM events e
        JOIN eventsCategory ec ON e.categoryID = ec.categoryID;
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
	log.Println(" - Called GetEventByEventID")
	eventID := chi.URLParam(r, "eventID")
	if eventID == "" {
		log.Println("eventID = ", eventID)
		// If eventID is empty, return an error response
		http.Error(w, "eventID parameter is required", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

    query := `
        SELECT e.eventid, e.eventName, ec.categoryName, e.eventLocation, e.eventDate
        FROM events e
        JOIN eventsCategory ec ON e.categoryID = ec.categoryID
        WHERE eventID = $1
    `

	row, err := h.db.QueryContext(ctx, query, eventID)
	log.Println("row= ", row)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer row.Close()
	
	// Scan rows to create Brother instance
	var event models.Event
	for row.Next() {
		event, err = createEventFromRow(row)
		if err != nil {
			log.Fatal("Error!", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if event.EventID == 0 {
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

func (h* Handler) InsertEvent(w http.ResponseWriter, r *http.Request) {
    log.Println("Called InsertEvent()")
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
    log.Println(query)
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
    log.Println(event)
    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Inserted new Event entry to `events` table"))
}
