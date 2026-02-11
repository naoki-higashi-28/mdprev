package handler

import (
	"net/http"

	fileuc "github.com/naoki-higashi-28/mdprev/internal/application/usecase/file"
)

// FileHandler handles file API requests.
type FileHandler struct {
	getFileUC    *fileuc.GetFileUseCase
	getRawFileUC *fileuc.GetRawFileUseCase
}

// NewFileHandler creates a new FileHandler.
func NewFileHandler(getFileUC *fileuc.GetFileUseCase, getRawFileUC *fileuc.GetRawFileUseCase) *FileHandler {
	return &FileHandler{getFileUC: getFileUC, getRawFileUC: getRawFileUC}
}

// HandleGetFile handles GET /api/file requests.
func (h *FileHandler) HandleGetFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Query().Get("path")
	if path == "" {
		http.Error(w, "Bad Request: path is required", http.StatusBadRequest)
		return
	}

	data, err := h.getFileUC.Execute(path)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write(data)
}

// HandleGetRaw handles GET /raw/{path...} requests.
func (h *FileHandler) HandleGetRaw(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("path")
	if path == "" {
		http.Error(w, "Bad Request: path is required", http.StatusBadRequest)
		return
	}

	filePath, err := h.getRawFileUC.Execute(path)
	if err != nil {
		writeError(w, err)
		return
	}

	http.ServeFile(w, r, filePath)
}
