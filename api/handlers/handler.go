// Handler is initialized and used in server.go for routing
package handlers

import (
	"database/sql"
)

type Handler struct {
	db *sql.DB
}

// Create a new Handler instance with the app's database connection
func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}
