package handler

import (
	"net/http"

	fileuc "github.com/naoki-higashi-28/mdprev/internal/application/usecase/file"
)

// FileHandler handles file API requests.
type FileHandler struct {
	uc *fileuc.UseCase
}

// NewFileHandler creates a new FileHandler.
func NewFileHandler(uc *fileuc.UseCase) *FileHandler {
	return &FileHandler{uc: uc}
}

// HandleGetFile handles GET /api/file requests.
func (h *FileHandler) HandleGetFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "Bad Request: path is required", http.StatusBadRequest)
		return
	}

	data, err := h.uc.GetFile(path)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write(data)
}

// HandleGetRaw handles GET /raw/{path...} requests.
func (h *FileHandler) HandleGetRaw(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("path")
	if path == "" {
		http.Error(w, "Bad Request: path is required", http.StatusBadRequest)
		return
	}

	filePath, err := h.uc.GetRaw(path)
	if err != nil {
		writeError(w, err)
		return
	}

	http.ServeFile(w, r, filePath)
}
