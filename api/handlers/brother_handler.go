package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/pacific-theta-tau/tt-db/api/models"
)

// Get data from all brothers in the Brother's table
func (h *Handler) GetAllBrothers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	// TODO: explicitly type columns
	query := "SELECT * FROM brothers"
	rows, err := h.db.QueryContext(ctx, query)

	// TODO: better error handling. maybe move query logic to models so we can catch all errors in
	// 1 place in this function.
	if err != nil {
		fmt.Println(err)
		return
	}

	// Scan rows from query to create Brother instances
	var brothers []*models.Brother
	for rows.Next() {
		var brother models.Brother
		err = rows.Scan(
			&brother.ID,
			&brother.First,
			&brother.Last,
			&brother.Status,
			&brother.Class,
			&brother.Email,
		)

		if err != nil {
			fmt.Println(err)
			return
		}
		brothers = append(brothers, &brother)
	}

	// Build HTTP response
	out, err := json.MarshalIndent(brothers, "", "\t")
	if err != nil {
		fmt.Println(err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(out)

	if err != nil {
		fmt.Println(err)
		return
	}
}
