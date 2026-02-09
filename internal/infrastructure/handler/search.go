package handler

import (
	"encoding/json"
	"net/http"

	searchuc "github.com/naoki-higashi-28/mdprev/internal/application/usecase/search"
)

// SearchHandler handles search API requests.
type SearchHandler struct {
	uc *searchuc.UseCase
}

// NewSearchHandler creates a new SearchHandler.
func NewSearchHandler(uc *searchuc.UseCase) *SearchHandler {
	return &SearchHandler{uc: uc}
}

// HandleSearch handles GET /api/search requests.
func (h *SearchHandler) HandleSearch(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	result, err := h.uc.Search(query)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "encoding response", http.StatusInternalServerError)
	}
}
