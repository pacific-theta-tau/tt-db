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

//	@Summary		Get all Brothers data
//	@Description	Get data from all Brother records in `Brothers` table
//	@Tags			Brothers
//	@Success		200		object		models.APIResponse{data=models.Brother}
//	@Failure		400		{object}	models.APIResponse
//	@Router			/api/brothers [get]
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
            errMsg := fmt.Sprintf("Error creating Brother object from row: '%s'\n", err.Error())
			log.Println(errMsg)
			http.Error(w, errMsg, http.StatusInternalServerError)
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
//	@Summary		Get Brother record by ID
//	@Description	Get Brother record by ID
//	@Tags			Brothers
//	@Param			id		path		int											true	"Brother ID"
//	@Success		200		object		models.APIResponse
//	@Failure		400		{object}	models.APIResponse
//	@Router			/api/brothers/{id} [get]
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
//	@Summary		Create Brother record
//	@Description	Create a new Brother record row for `Brothers` table
//	@Tags			Brothers
//	@Param			body_params body	models.Brother true	"Values for new record"
//	@Success		200		object		models.APIResponse
//	@Failure		400		{object}	models.APIResponse
//	@Router			/api/brothers [post]
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

//	@Summary		Delete Brother by Roll Call
//	@Description	Delete Brother with by Roll Call
//	@Tags			Brothers
//	@Param			body_params body	string  true	"RollCall of Brother"
//  @Success		200		object		models.APIResponse
//	@Failure		400		{object}	models.APIResponse
//	@Router			/api/brothers/{id} [delete]
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


//	@Summary		Update Brother record
//	@Description	Update one or more fields for Brother record 
//	@Tags			Brothers
//	@Param			body_params body    models.Brother  true	"Values to update for Brother"
//	@Success		200		object		models.APIResponse
//	@Failure		400		{object}	models.APIResponse
//	@Router			/api/brothers/{id} [put]
/* PUT /api/brothers/{id} */
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


// POST /api/brothers/{id}/statuses
//	@Summary		Get status history of a Brother
//	@Description	Get all status recorded for Brother
//	@Tags			Brothers
//	@Param			id		path		int     true	"Brother ID"
//	@Success		200		object		models.APIResponse
//	@Failure		400		{object}	models.APIResponse
//	@Router			/api/brothers/{id}/statuses [get]
/* /api/brothers/{id}/statuses */
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
//	@Summary		Create status record for Brother
//	@Description	Create a new status record for a Brother
//	@Tags			Brothers
//	@Accept			json
//	@Produce		json
//	@Param			body_params body		handlers.CreateBrotherStatus.RequestBody	true	"Values for new record"
//	@Param			id		path		int											true	"Brother ID"
//	@Success		200		object		models.APIResponse
//	@Failure		400		{object}	models.APIResponse
//	@Router			/api/brothers/{id}/statuses [post]
func (h *Handler) CreateBrotherStatus(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

    // Expected request body data
    type RequestBody struct {
        BrotherID   int `json:"brotherID"`
        SemesterID  int `json:"semesterID"`
        Status      string `json:"status"`
    }
    var requestBody RequestBody
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

//	@Tags			Brothers
//	@Summary		Get total Brothers count
//	@description	Get major distribution counts across all members
//	@Success		200	{object}	models.APIResponse{data=int}	"Success"
//	@failure		400	{string}	string							"Error"
//	@Router			/api/brothers/count [get]
/* GET /api/brothers/count */
func (h *Handler) GetBrothersCount(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

    // query for total row counts in brothers table
    query := `
    SELECT COUNT(*) AS count
    FROM brothers
    `
    row := h.db.QueryRowContext(ctx, query)

    // parse query result
    var count int 
    err := row.Scan(
        &count,
    )  
    if err != nil {
        errMsg := fmt.Sprintf("Error parsing brothers count query result from row: '%s'\n", err.Error())
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
        return
    }

    data := map[string]int{"count": count}
    models.RespondWithSuccess(w, http.StatusOK, data)
}

//  GET /api/brothers/majors/count
//	@Tags			Brothers
//	@Summary		Get major counts
//	@description	Get major distribution counts across all members
//	@Success		200	{object}	models.APIResponse{data=handlers.GetBrotherStatusCount.SemesterCount}	"desc"
//	@failure		400	{string}	string																	"error"
//	@Router			/api/brothers/majors/count [get]
func (h *Handler) GetBrothersMajorsCount(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

    query := `
    SELECT major, COUNT(*) AS count
    FROM brothers
    GROUP BY major;
    `
    rows, err := h.db.QueryContext(ctx, query)
    if err != nil {
        errMsg := fmt.Sprintf("Error during query: '%s'\n", err.Error())
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
        return
	}

    type MajorCount struct {
        Major   string `json:"major"`
        Count   int `json:"count"`
    }
    var majorCounts []*MajorCount
    for rows.Next() {
        var curRow MajorCount
        err = rows.Scan(
            &curRow.Major,
            &curRow.Count,
        )   
        if err != nil {
            errMsg := fmt.Sprintf("Error parsing major count query result from row: '%s'\n", err.Error())
            models.RespondWithError(w, http.StatusInternalServerError, errMsg)
			return
		}
        majorCounts = append(majorCounts, &curRow)
    }

    models.RespondWithSuccess(w, http.StatusOK, majorCounts)
}


//	@Tags			Brothers
//	@Summary		Get status counts
//	@description	Get status counts for all semesters
//	@Param			status	query		string																	false	"Status filter"	
//	@Success		200		{object}	models.APIResponse{data=handlers.GetBrotherStatusCount.SemesterCount}   
//	@failure		400		{string}	models.APIResponse														"error"
//	@Router			/api/brothers/statuses/count [get]
func (h *Handler) GetBrotherStatusCount(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

    // Get query params
    queryStatus := ""
    status := r.URL.Query().Get("status")
    log.Printf("Received query param status: %s", status)
    if status != "" {
        queryStatus = fmt.Sprintf("WHERE bs.status = '%s'", status)
    }

    querySemester := ""
    semester := r.URL.Query().Get("semester")
    if semester != "" {
        querySemester = fmt.Sprintf("WHERE s.semesterLabel = %s", semester)
    }
    // Check for "all"

    query := fmt.Sprintf(`
    SELECT s.semesterLabel, COUNT(*) AS count
    FROM brotherStatus bs
    JOIN semester s ON bs.semesterID = s.semesterID
    %s
    %s
    GROUP BY s.semesterLabel;
    `, queryStatus, querySemester)
    log.Printf("Query:\n%s", query)
    //{
    //    data: [
    //        {'semester': 'Fall 2022', actives: 20, co-op: 20, etc...}
    //    ]
    //}
    rows, err := h.db.QueryContext(ctx, query)
    if err != nil {
        errMsg := fmt.Sprintf("Error during query: '%s'\n", err.Error())
        log.Println(errMsg)
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
        return
	}

    type SemesterCount struct {
        Semester    string `json:"semester"`
        Count       int `json:"count"`
    }
    var semesterCounts []*SemesterCount
    for rows.Next() {
        var curRow SemesterCount
        err = rows.Scan(
            &curRow.Semester,
            &curRow.Count,
        )   
        if err != nil {
            errMsg := fmt.Sprintf("Error parsing major count query result from row: '%s'\n", err.Error())
            log.Println(errMsg)
            models.RespondWithError(w, http.StatusInternalServerError, errMsg)
			return
		}
        semesterCounts = append(semesterCounts, &curRow)
    }
    log.Printf("Semester counts: %v:", semesterCounts)

    models.RespondWithSuccess(w, http.StatusOK, semesterCounts)
}
