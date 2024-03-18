// This file contains all functions that handle requests for the /api/brothers endpoint
package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/pacific-theta-tau/tt-db/api/models"
)

const brothers_table = "brothers"

// Helper function to scan SQL row and create new Brother instance
func createBrotherFromRow(row *sql.Rows) (models.Brother, error) {
	var brother models.Brother
	err := row.Scan(
		&brother.PacificId,
		&brother.FirstName,
		&brother.LastName,
		&brother.Status,
		&brother.Class,
		&brother.RollCall,
		&brother.Email,
		&brother.PhoneNumber,
		&brother.BadStanding,
	)
	if err != nil {
		return models.Brother{}, err
	}

	return brother, err
}

// Get data from all brothers in the Brother's table
func (h *Handler) GetAllBrothers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// TODO: explicitly type columns
	query := "SELECT * FROM " + brothers_table
	rows, err := h.db.QueryContext(ctx, query)
	if err != nil {
		// return error status code
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Read rows from query to create Brother instances
	var brothers []*models.Brother
	for rows.Next() {
		brother, err := createBrotherFromRow(rows)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		brothers = append(brothers, &brother)
	}

	// Build HTTP response
	out, err := json.MarshalIndent(brothers, "", "\t")
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

// Query Brothers by their PacificID
func (h *Handler) GetBrotherByPacificID(w http.ResponseWriter, r *http.Request) {
	// TODO: handle case when param is not provided
	pacificID := r.URL.Query().Get("pacificID")
	// brother, err := models.GetBrotherByPacificID(h.db, pacificID)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	query := "SELECT * FROM brothers WHERE pacificId = $1"
	row, err := h.db.QueryContext(ctx, query, pacificID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Scan rows to create Brother instance
	brother, err := createBrotherFromRow(row)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Build HTTP response
	out, err := json.MarshalIndent(brother, "", "\t")
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

// Add new brother entry to database
func (h *Handler) AddBrother(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var brother models.Brother
	err := json.NewDecoder(r.Body).Decode(&brother)
	if err != nil {
		http.Error(w, fmt.Sprint("Error decoding request body", err.Error()), http.StatusBadRequest)
	}

	// Validate brothers struct
	validate := validator.New()
	if err := validate.Struct(brother); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `
	INSERT INTO brothers (pacificId, firstName, lastName, status, className, rollCall, email, phoneNumber, badStanding)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning *
	`
	_, err = h.db.ExecContext(
		ctx,
		query,
		brother.PacificId,
		brother.FirstName,
		brother.LastName,
		brother.Status,
		brother.Class,
		brother.RollCall,
		brother.Email,
		brother.PhoneNumber,
		brother.BadStanding,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println(brother)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Inserted new Brother entry to 'brothers' table"))
}

// Delete brother from brothers table by Roll Call number
func (h *Handler) RemoveBrother(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var requestBody map[string]interface{}
	if err = json.Unmarshal(body, &requestBody); err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	rollCall, ok := requestBody["rollCall"]
	if !ok {
		http.Error(w, "Key not found in request body", http.StatusBadRequest)
		return
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE rollCall = $1", brothers_table)
	_, err = h.db.ExecContext(ctx, query, rollCall)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Removed brother from 'brothers' table successfully"))
}
