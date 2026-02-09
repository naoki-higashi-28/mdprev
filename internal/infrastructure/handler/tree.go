package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	treeuc "github.com/naoki-higashi-28/mdprev/internal/application/usecase/tree"
	"github.com/naoki-higashi-28/mdprev/internal/domain"
)

// TreeHandler handles tree API requests.
type TreeHandler struct {
	uc *treeuc.UseCase
}

// NewTreeHandler creates a new TreeHandler.
func NewTreeHandler(uc *treeuc.UseCase) *TreeHandler {
	return &TreeHandler{uc: uc}
}

// HandleGetTree handles GET /api/tree requests.
func (h *TreeHandler) HandleGetTree(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")

	result, err := h.uc.GetTree(path)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "encoding response", http.StatusInternalServerError)
	}
}

func writeError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrForbidden):
		http.Error(w, "Forbidden", http.StatusForbidden)
	case errors.Is(err, domain.ErrNotFound):
		http.Error(w, "Not Found", http.StatusNotFound)
	case errors.Is(err, domain.ErrBadRequest), errors.Is(err, domain.ErrInvalidExtFile):
		http.Error(w, "Bad Request", http.StatusBadRequest)
	default:
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
