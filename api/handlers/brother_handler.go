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
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	// TODO: explicitly type columns
	query := "SELECT * FROM " + brothers_table
    log.Printf("Query:\n%s", query)
	rows, err := h.db.QueryContext(ctx, query)
	if err != nil {
        errMsg := fmt.Sprintf("Error while querying rows in Brother's table: %s", err.Error())
        log.Println(errMsg)
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}

	// Read rows from query to create Brother instances
	var brothers []*models.Brother
	for rows.Next() {
		brother, err := createBrotherFromRow(rows)
        if err != nil {
            errMsg := fmt.Sprintf("Error creating Brother object from row: '%s'\n", err.Error())
			log.Println(errMsg)
            models.RespondWithError(w, http.StatusInternalServerError, errMsg)
			return
		}
		brothers = append(brothers, &brother)
	}

    data := brothers
    models.RespondWithSuccess(w, http.StatusOK, data)
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
        errMsg := fmt.Sprintf("Invalid brother ID: %v", err.Error())
        log.Printf(errMsg)
        models.RespondWithError(w, http.StatusBadRequest, errMsg)
        return
    }

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

    query := "SELECT * FROM brothers WHERE brotherID = $1"
    log.Printf("Query:\n%s", query)
    row, err := h.db.QueryContext(ctx, query, brotherID)
	if err != nil {
        errMsg := fmt.Sprintf("Erro while querying for Brother with ID %d: %s", brotherID, err.Error())
		log.Println(errMsg)
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}

    // Scan rows to create Brother instance
	var brother models.Brother
	for row.Next() {
		brother, err = createBrotherFromRow(row)
		if err != nil {
            errMsg := fmt.Sprintf("Erro while parsing rows: %s", err.Error())
			log.Println(errMsg)
            models.RespondWithError(w, http.StatusInternalServerError, errMsg)
			return
		}
	}

    // postgres returns 0 if row not found
    if brother.BrotherID == 0 {
        errMsg := fmt.Sprintf("Brother ID %d not found", brotherID)
        log.Printf(errMsg)
        models.RespondWithError(w, http.StatusBadRequest, errMsg)
        return
    }

	// Build HTTP response
    models.RespondWithSuccess(w, http.StatusOK, brother)
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
        errMsg := fmt.Sprint("Error decoding request body", err.Error())
        log.Println(errMsg)
        models.RespondWithFail(w, http.StatusInternalServerError, errMsg)
        return
	}

	// Validate brothers struct
	validate := validator.New()
	if err := validate.Struct(brother); err != nil {
        errMsg := fmt.Sprint("Error decoding request body: %s", err.Error())
        log.Println(errMsg)
        models.RespondWithFail(w, http.StatusBadRequest, errMsg)
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
        errMsg := fmt.Sprint("Error during query: %s", err.Error())
        log.Println(errMsg)
        models.RespondWithFail(w, http.StatusInternalServerError, errMsg)
		return
	}

    models.RespondWithSuccess(w, http.StatusCreated, "")
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
        errMsg := fmt.Sprint("Error parsing request body: %s", err.Error())
        log.Println(errMsg)
        models.RespondWithFail(w, http.StatusInternalServerError, errMsg)

		return
	}

	var requestBody map[string]interface{}
	if err = json.Unmarshal(body, &requestBody); err != nil {
        errMsg := fmt.Sprint("Error validating body params: %s", err.Error())
        log.Println(errMsg)
        models.RespondWithFail(w, http.StatusBadRequest, errMsg)
		return
	}

	rollCall, ok := requestBody["rollCall"]
	if !ok {
        errMsg := "Roll Call missing in body params: %s"
        models.RespondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE rollCall = $1", brothers_table)
	_, err = h.db.ExecContext(ctx, query, rollCall)
	if err != nil {
        errMsg := fmt.Sprint("Error while querying for brother with Roll Call %s: %s", rollCall, err.Error())
        log.Println(errMsg)
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}

    // TODO: return deleted row
    models.RespondWithSuccess(w, http.StatusOK, "")
}


//	@Summary		Update Brother record
//	@Description	Update one or more fields for Brother record 
//	@Tags			Brothers
//	@Param			body_params body    models.Brother  true	"Values to update for Brother"
//	@Success		200		object		models.APIResponse
//	@Failure		400		{object}	models.APIResponse
//	@Router			/api/brothers/{id} [patch]
/* PATCH /api/brothers/{id} */
func (h *Handler) UpdateBrother(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

    // TODO: parse body params using JSON NewDecoder()
	body, err := io.ReadAll(r.Body)
	if err != nil {
        errMsg := fmt.Sprint("Error reading request body %s", err.Error())
        log.Println(errMsg)
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}

	var requestBody map[string]interface{}
	if err = json.Unmarshal(body, &requestBody); err != nil {
        errMsg := fmt.Sprintf("Error decoding JSON: %s", err.Error())
        log.Printf(errMsg)
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}

	// rollCall, ok := requestBody["rollCall"]
	brotherID := chi.URLParam(r, "id")
	if brotherID == "0" {
        errMsg := fmt.Sprintf("Missing urlParam brotherID: %s", brotherID)
        log.Printf(errMsg)
        models.RespondWithFail(w, http.StatusBadRequest, errMsg)
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

        switch v := newColumnValue.(type) {
        case float64: // Numbers in JSON decode as float64 by default
            query += fmt.Sprintf(" %s = '%d',", column, int(v))
        case int:
            query += fmt.Sprintf(" %s = '%d',", column, v)
        case string:
            query += fmt.Sprintf(" %s = '%s',", column, v)
        default:
            log.Printf("Unsupported type for column %s: %T", column, v)
        }

        //query += fmt.Sprintf(" %s = '%s',", column, newColumnValue)
	}

	// remove trailling comma
	query = query[:len(query)-1] + " WHERE brotherID = $1"

	_, err = h.db.ExecContext(ctx, query, brotherID)
	if err != nil {
        errMsg := fmt.Sprintf("Error while querying `%s`: %s", query, err.Error())
        log.Printf(errMsg)
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}

	// w.Write([]byte(fmt.Sprintf("Updated brother with brotherID %s successfully", brotherID)))
    models.RespondWithSuccess(w, http.StatusOK, "")
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
        errMsg := fmt.Sprintf("Invalid brother ID: %v", err.Error())
        log.Println(errMsg)
        models.RespondWithFail(w, http.StatusBadRequest, errMsg)
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
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}
    defer row.Close()

	var brother models.Brother
	for row.Next() {
		brother, err = createBrotherFromRow(row)
		if err != nil {
            errMsg := fmt.Sprintf("Error creating Brother object from row: '%s'\n", err.Error())
			log.Println(errMsg)
            models.RespondWithError(w, http.StatusInternalServerError, errMsg)
			return
		}
	}
    if brother.BrotherID == 0 {
        errMsg := fmt.Sprintf("Brother ID %d not found", brother.BrotherID)
        log.Printf(errMsg)
        models.RespondWithFail(w, http.StatusBadRequest, errMsg)
        return
    }

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
        errMsg := fmt.Sprintf("Error while querying for status and semester: %s", err.Error())
		fmt.Println(errMsg)
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}

	var brotherStatuses []*models.Status
	for row.Next() {
        status, err := models.CreateStatusFromRow(row)
		if err != nil {
            errMsg := fmt.Sprintf("Error creating Status object from row: %s", err.Error())
			log.Println(errMsg)
            models.RespondWithError(w, http.StatusInternalServerError, errMsg)
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

    models.RespondWithSuccess(w, http.StatusOK, response)
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
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
        return
    }
    // Validate received data
    validate := validator.New()
    if err := validate.Struct(requestBody); err != nil {
        errMsg := fmt.Sprintf("Invalid Input: %s", err.Error())
        log.Println(errMsg)
        models.RespondWithFail(w, http.StatusBadRequest, errMsg)
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
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
		return
	}

    models.RespondWithSuccess(w, http.StatusCreated, "")
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
//	@Summary		Get all status records per brother
//	@description	Get all status records per brother
//	@Param			semester query		string																	false	"Semester filter"	
//	@Success		200		{object}	models.APIResponse{data=models.BrotherStatus}   
//	@failure		400		{string}	models.APIResponse														"error"
//	@Router			/api/brothers/statuses [get]
func (h *Handler) GetAllBrotherStatuses(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
    defer cancel()
    semester := r.URL.Query().Get("semester")

    // Query for all brother statuses
    query := `
    SELECT b.brotherID, b.rollCall, b.firstName, b.lastName, b.major, bs.status, s.semesterLabel
    FROM brotherStatus bs
    JOIN brothers b ON b.brotherID = bs.brotherID
    JOIN semester s ON s.semesterID = bs.semesterID
    `
    var rows *sql.Rows
    var err error
    if semester == "" {
        // Query without filter
        rows, err = h.db.QueryContext(ctx, query)
    } else {
        // query filtering by semester
        query += " WHERE semesterLabel = $1"
        rows, err = h.db.QueryContext(
            ctx,
            query,
            semester,
        )
    } 

    // Error handling query
    fmt.Printf("Querying for all semester statuses:\n%s\n", query)
    if err != nil {
        errMsg := fmt.Sprintf("Error while querying for all brother statuses: %s", err.Error())
        log.Println(errMsg)
        models.RespondWithError(w, http.StatusInternalServerError, errMsg)
        return
    }

    log.Println("Parsing brother objects\n")
    // Parse query rows
    var brotherStatuses []*models.BrotherStatus
    for rows.Next(){
        brotherStatus, err := models.CreateBrotherStatusFromRow(rows)
        if err != nil {
            errMsg := fmt.Sprintf("Error while parsing brotherStatus query: '%s'", err.Error())
            log.Println(errMsg)
            models.RespondWithError(w, http.StatusInternalServerError, errMsg)
            return
        }
        brotherStatuses = append(brotherStatuses, &brotherStatus)
    }

    // Build response
    models.RespondWithSuccess(w, http.StatusOK, brotherStatuses)
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

