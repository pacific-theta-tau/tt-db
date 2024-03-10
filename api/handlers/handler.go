package handlers

import "database/sql"

// Used to share database connection across all handler functions
type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}
