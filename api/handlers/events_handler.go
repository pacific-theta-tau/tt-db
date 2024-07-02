package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"github.com/pacific-theta-tau/tt-db/api/models"
	"github.com/go-chi/chi"
)

const events_table = "events"

// Helper function to scan SQL row and create new event instance
func createEventFromRow(row *sql.Rows) (models.Event, error) {
	var events models.Event
	err := row.Scan(
		&events.EventID,
		&events.EventName,
		&events.Category,
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
	fmt.Println(" - Called GetAllEvents")
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// TODO: explicitly type columns
	query := "SELECT * FROM " + events_table
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
	fmt.Println(" - Called GetEventByEventID")
	eventID := chi.URLParam(r, "eventID")
	if eventID == "" {
		fmt.Println("eventID = ", eventID)
		// If eventID is empty, return an error response
		http.Error(w, "eventID parameter is required", http.StatusBadRequest)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := "SELECT * FROM events WHERE eventID = $1"
	row, err := h.db.QueryContext(ctx, query, eventID)
	fmt.Println("row= ", row)
	if err != nil {
		fmt.Println(err)
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
		fmt.Println(err)
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
