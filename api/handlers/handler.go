package handlers

import (
	"database/sql"
	"time"
)

// Threshold for waiting database response
const dbTimeout = time.Second * 5

// Handler contains methods to handle all API requests
type Handler struct {
	db *sql.DB
}

// Create a new Handler instance with the app's database connection
func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}
