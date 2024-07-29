package handlers

import (
    "bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
    "time"

	"github.com/go-chi/chi"
	"github.com/pacific-theta-tau/tt-db/api/models"
)

// Test GET request for /api/events
func TestGetAllEvents(t *testing.T) {
	// Init chi router and handler function
	router := chi.NewRouter()
	router.Get("/api/events", handler.GetAllEvents)

	// Create new request
	req, err := http.NewRequest("GET", "/api/events", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Record response in a ResponseReqcorder
	rr := httptest.NewRecorder()

	// Serve HTTP request
	router.ServeHTTP(rr, req)

	// Check status code
	checkResponseCode(t, 200, rr.Code)

	// Parse body
	var response []*models.Event
	if err = json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("failed to parse response body: %v", err)
	}
	fmt.Println(response)
}

func TestGetEventByID(t *testing.T) {
	// Init chi router and handler function
	router := chi.NewRouter()
	router.Get("/api/events/1", handler.GetAllEvents)

	// Create new request
	req, err := http.NewRequest("GET", "/api/events/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Record response in a ResponseReqcorder
	rr := httptest.NewRecorder()

	// Serve HTTP request
	router.ServeHTTP(rr, req)

	// Check status code
	checkResponseCode(t, 200, rr.Code)

	// Parse body
	var response []*models.Event
	if err = json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("failed to parse response body: %v", err)
	}
	fmt.Println(response)
}

func TestUpdateEventByID(t *testing.T) {
	// Init chi router and handler function
	router := chi.NewRouter()
    endpoint := "/api/events"
	router.Put(endpoint, handler.UpdateEventByID)

    dateLayout := "01-02-2006"
    eventDate, err := time.Parse(dateLayout, "07-27-2024")
    if err != nil {
        t.Fatalf("Error parsing date: %v\n", err)
    }

    event := models.Event{
        EventID: 1,
        EventName: "New Name",
        CategoryName: "Brotherhood",
        EventLocation: "New Location",
        EventDate: eventDate.UTC(), 
    }

    body, err := json.Marshal(event)
    if err != nil {
        t.Fatalf("Failed to marshal item: %v", err)
    }

	// Create new request
	req, err := http.NewRequest("PUT", endpoint, bytes.NewBuffer(body))
    if err != nil {
        t.Fatalf("Failed to create request: %v", err)
    }

	// Record response in a ResponseReqcorder
	rr := httptest.NewRecorder()

	// Serve HTTP request
	router.ServeHTTP(rr, req)

	// Check status code
	checkResponseCode(t, 200, rr.Code)

	// Parse body
	var response *models.Event
	if err = json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
        t.Logf("response body: %v", rr.Body.String())
		t.Errorf("failed to parse response body: %v", err)
	}

    // Check if event was updated
    if event != *response {
        t.Errorf("Failed to update event. \nExpected:\n%+v \n\nActual:\n%+v", event, *response)
    }
}
