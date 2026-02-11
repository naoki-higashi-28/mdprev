package handler

import (
	"encoding/json"
	"net/http"

	treeuc "github.com/naoki-higashi-28/mdprev/internal/application/usecase/tree"
)

// TreeHandler handles tree API requests.
type TreeHandler struct {
	uc *treeuc.GetTreeUseCase
}

// NewTreeHandler creates a new TreeHandler.
func NewTreeHandler(uc *treeuc.GetTreeUseCase) *TreeHandler {
	return &TreeHandler{uc: uc}
}

// HandleGetTree handles GET /api/tree requests.
func (h *TreeHandler) HandleGetTree(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")

	result, err := h.uc.Execute(path)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, "encoding response", http.StatusInternalServerError)
	}
}
