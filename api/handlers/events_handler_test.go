package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/pacific-theta-tau/tt-db/api/models"
)

// Test GET request for /api/brothers
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