// This file contains all functions that handle requests for the /api/brothers endpoint
package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
    "strconv"

	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
	"github.com/pacific-theta-tau/tt-db/api/models"
)

const brothers_table = "brothers"

// Helper function to scan SQL row and create new Brother instance
// TODO: move this to models/brother.go
func createBrotherFromRow(row *sql.Rows) (models.Brother, error) {
	var brother models.Brother
	err := row.Scan(
        &brother.BrotherID,
		&brother.RollCall,
		&brother.FirstName,
		&brother.LastName,
		&brother.Major,
		&brother.Status,
		&brother.Class,
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
	log.Println("-- Called GetAllBrothers() --")
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

    log.Println("Query Successful")
    // data := make(map[string]interface{})
    // data["data"] = brothers
    // data["status"] = "success"
    data := brothers

	// Build HTTP response
	out, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(out)

	if err != nil {
        log.Printf("Error while creating response: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Query brothers by ID
// GET /api/brothers/{id}
func (h *Handler) GetBrotherByID(w http.ResponseWriter, r *http.Request) {
    fmt.Println("\nGetBrrotherByID called")
    brotherIDStr := chi.URLParam(r, "id")
    brotherID, err := strconv.Atoi(brotherIDStr)
    if err != nil {
        log.Printf("Invalid brother ID: %v", err.Error())
        http.Error(w, "Invalid event ID", http.StatusBadRequest)
        return
    }

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

    query := "SELECT * FROM brothers WHERE brotherID = $1"
    row, err := h.db.QueryContext(ctx, query, brotherID)
	fmt.Printf("row", row)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

    // Scan rows to create Brother instance
	var brother models.Brother
	for row.Next() {
		brother, err = createBrotherFromRow(row)
		if err != nil {
			log.Println("Error!", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
    if brother.BrotherID == 0 {
        http.Error(w, "Brother ID %d not found", brotherID)
        log.Printf("Brother ID %d not found", brotherID)
        return
    }

	// Build HTTP response
	out, err := json.MarshalIndent(brother, "", "\t")
	if err != nil {
		fmt.Println(err.Error())
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
	INSERT INTO brothers (rollCall, firstName, lastName, major, status, className, email, phoneNumber, badStanding)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning *
	`
	_, err = h.db.ExecContext(
		ctx,
		query,
		brother.RollCall,
		brother.FirstName,
		brother.LastName,
		brother.Major,
		brother.Status,
		brother.Class,
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

func (h *Handler) UpdateBrother(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	body, err := io.ReadAll(r.Body)
	if err != nil {
        log.Printf("Error reading request body %s", err.Error())
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var requestBody map[string]interface{}
	if err = json.Unmarshal(body, &requestBody); err != nil {
        log.Printf("Error decoding JSON: %s", err.Error())
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	// rollCall, ok := requestBody["rollCall"]
	brotherID, ok := requestBody["brotherID"]
	if !ok {
        log.Printf("Key not found in request body: %s", err.Error())
		http.Error(w, "Key not found in request body", http.StatusBadRequest)
		return
	}

	// Format query with each param in request body
	// TODO: add validator checks for Body params
	query := fmt.Sprintf("UPDATE %s SET", brothers_table)
	columns := []string{
		"firstName",
		"lastName",
        "major",
		"status",
		"class",
		"email",
		"phoneNumber",
		"badStanding",
        "rollCall",
	}
	for _, column := range columns {
		newColumnValue, ok := requestBody[column]
		if !ok {
			continue
		}
		query += fmt.Sprintf(" %s = '%s',", column, newColumnValue)
	}

	// remove trailling comma
	query = query[:len(query)-1] + " WHERE brotherID = $1"

	_, err = h.db.ExecContext(ctx, query, brotherID)
	if err != nil {
        log.Printf("Error while querying `%s`: %s", query, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Updated brother with brotherID %s successfully", brotherID)))
}


// /api/brothers/{id}/statuses
func (h *Handler) GetBrotherStatusHistory(w http.ResponseWriter, r *http.Request) {
    log.Println("\n\nCalled get brother status history")
    brotherIDStr := chi.URLParam(r, "id")
    brotherID, err := strconv.Atoi(brotherIDStr)
    if err != nil {
        log.Printf("Invalid brother ID: %v", err.Error())
        http.Error(w, "Invalid event ID", http.StatusBadRequest)
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

    query := `
    SELECT
        brotherID, rollCall, firstName, lastName, major, status, className, email, phoneNumber, badStanding
    FROM brothers b
    WHERE b.brotherID = $1
    `
    log.Printf("Querying for brother:\n%s", query)
    row, err := h.db.QueryContext(ctx, query, brotherID)
	if err != nil {
        errMsg := fmt.Sprintf("Error while querying for brother data: `%s`\n", err.Error())
        log.Println(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
    defer row.Close()

	var brother models.Brother
	for row.Next() {
		brother, err = createBrotherFromRow(row)
		if err != nil {
            errMsg := fmt.Sprintf("Error creating Brother object from row: '%s'\n", err.Error())
			log.Println(errMsg)
			http.Error(w, errMsg, http.StatusInternalServerError)
			return
		}
	}
    if brother.BrotherID == 0 {
        http.Error(w, "Brother ID %d not found", http.StatusBadRequest)
        log.Printf("Brother ID %d not found", brotherID)
        return
    }
    log.Println("Parsed brother data successfully\n")

    // Query for status
    query = `
    SELECT s.semesterLabel, bs.status
    FROM brotherStatus bs
    JOIN semester s ON s.semesterID = bs.semesterID
    WHERE brotherID = $1
    `
    log.Printf("Querying for status and semester:\n%s\n", query)
    row, err = h.db.QueryContext(ctx, query, brotherID)
    if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var brotherStatuses []*models.Status
	for row.Next() {
        status, err := models.CreateStatusFromRow(row)
		if err != nil {
            errMsg := fmt.Sprintf("Error creating Status object from row: %s", err.Error())
			log.Println(errMsg)
			http.Error(w, errMsg, http.StatusInternalServerError)
			return
		}
        brotherStatuses= append(brotherStatuses, &status)
	}
    log.Println("Parsed semesterLabel and status successfully\n")
    
    // Write response
    response := map[string]interface{}{
        "brotherID": brother.BrotherID,
        "firstName": brother.FirstName,
        "lastName": brother.LastName,
        "rollCall": brother.RollCall,
        "class": brother.Class,
        "statuses": brotherStatuses,
    }
    log.Printf("Response: %v\n", response)

    w.Header().Set("Content-Type", "application/json")
    if err = json.NewEncoder(w).Encode(response); err != nil {
        errMsg := fmt.Sprintf("Error while encoding response: %s", err.Error())
        log.Println(errMsg)
        http.Error(w, errMsg, http.StatusInternalServerError)
        return
    }
	w.WriteHeader(http.StatusOK)
}

// POST /api/brothers/{id}/statuses
func (h *Handler) CreateBrotherStatus(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

    // Expected request body data
    var requestBody struct {
        BrotherID   int `json:"brotherID"`
        SemesterID  int `json:"semesterID"`
        Status      string `json:"status"`
    }
    // Parse body
    err := json.NewDecoder(r.Body).Decode(&requestBody)
    if err != nil {
        errMsg := fmt.Sprintf("Error while parsing request body data: %s", err.Error())
        log.Println(errMsg)
        http.Error(w, errMsg, http.StatusInternalServerError)
        return
    }
    // Validate received data
    validate := validator.New()
    if err := validate.Struct(requestBody); err != nil {
        errMsg := fmt.Sprintf("Invalid Input: %s", err.Error())
        log.Println(errMsg)
        http.Error(w, errMsg, http.StatusBadRequest)
        return
    }

    // Create new row for brotherStatus
    query := `
    INSERT INTO brotherStatus (brotherID, semesterID, status)
    VALUES ($1, $2, $3)
    `
    log.Printf("Insert query:\n%s\n", query)
    _, err = h.db.QueryContext(
        ctx,
        query,
        requestBody.BrotherID,
        requestBody.SemesterID,
        requestBody.Status,
    )

    if err != nil {
        errMsg := fmt.Sprintf("Query error:\n\t'%s'\n", err.Error())
        log.Println(errMsg)
        log.Println("Sending HTTP error response\n")
		http.Error(w, errMsg, http.StatusInternalServerError)
        log.Println("After sending error response\n")
		return
	}

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
	w.Write([]byte("Created Brother Status successfully"))
}
