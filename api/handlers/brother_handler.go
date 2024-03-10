package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pacific-theta-tau/tt-db/api/models"
)

// Get data from all brothers in the Brother's table
func (h *Handler) GetAllBrothers(w http.ResponseWriter, r *http.Request) {
	brothers, err := models.GetAllBrothers(h.db)
	if err != nil {
		// send error response
		fmt.Println(err)
		return
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
