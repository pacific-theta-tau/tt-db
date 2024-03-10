package handlers

import (
	"net/http"
)

func (h *Handler) GetAllBrothers(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get All Brothers Test"))
}
