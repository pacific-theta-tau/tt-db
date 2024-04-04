package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/pacific-theta-tau/tt-db/api/models"
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