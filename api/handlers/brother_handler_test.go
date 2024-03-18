package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/pacific-theta-tau/tt-db/api/models"
	"github.com/pacific-theta-tau/tt-db/db"
)

// var testdb *db.PostgresDB
var handler *Handler

func TestMain(m *testing.M) {
	// Create DB connection to dev db
	err := godotenv.Load("../../dev.env")
	if err != nil {
		log.Fatal("Error loading .env file. Make sure to have setup the appropriate .env file")
	}
	testdb := db.NewPostgresDB()
	testdb.Connect(os.Getenv("DATABASE_URL"))
	defer testdb.Conn.Close()
	handler = NewHandler(testdb.Conn)

	// Run tests
	exitCode := m.Run()

	// Exit with the appropriate exit code
	os.Exit(exitCode)
}

// Helper function to test expected status code from actual
func checkResponseCode(t *testing.T, expected int, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

// Test GET request for /api/brothers
func TestGetAllBrothers(t *testing.T) {
	// Init chi router and handler function
	router := chi.NewRouter()
	router.Get("/api/brothers", handler.GetAllBrothers)

	// Create new request
	req, err := http.NewRequest("GET", "/api/brothers", nil)
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
	var response []*models.Brother
	if err = json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("failed to parse response body: %v", err)
	}
	fmt.Println(response)
}
