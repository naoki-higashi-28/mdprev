package handler

import (
	"fmt"
	"net/http"

	"github.com/naoki-higashi-28/mdprev/internal/domain"
)

// WatchHandler handles SSE connections for file change notifications.
type WatchHandler struct {
	repo domain.WatchRepository
}

// NewWatchHandler creates a new WatchHandler.
func NewWatchHandler(repo domain.WatchRepository) *WatchHandler {
	return &WatchHandler{repo: repo}
}

// HandleWatch streams file change events via SSE.
func (h *WatchHandler) HandleWatch(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	flusher.Flush()

	events, cancel := h.repo.Subscribe()
	defer cancel()

	for {
		select {
		case <-r.Context().Done():
			return
		case evt := <-events:
			fmt.Fprintf(w, "data: %s\n\n", evt.Path)
			flusher.Flush()
		}
	}
}
