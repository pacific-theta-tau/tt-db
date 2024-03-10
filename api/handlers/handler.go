package handlers

import (
	"database/sql"
	"time"
)

const dbTimeout = time.Second * 3

type Handler struct {
	db *sql.DB
}

// Create a new Handler instance with the app's database connection
func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}
