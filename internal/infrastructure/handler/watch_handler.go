package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"

	appport "github.com/naoki-higashi-28/mdprev/internal/application/port"
	watchuc "github.com/naoki-higashi-28/mdprev/internal/application/usecase/watch"
)

type sseMessage struct {
	Type string `json:"type"`
	Path string `json:"path"`
}

// WatchHandler handles SSE connections for file change notifications.
type WatchHandler struct {
	uc           *watchuc.SubscribeFileChangesUseCase
	onConnect    func()
	onDisconnect func()
}

// NewWatchHandler creates a new WatchHandler.
func NewWatchHandler(uc *watchuc.SubscribeFileChangesUseCase, onConnect, onDisconnect func()) *WatchHandler {
	return &WatchHandler{uc: uc, onConnect: onConnect, onDisconnect: onDisconnect}
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

	if h.onConnect != nil {
		h.onConnect()
	}
	if h.onDisconnect != nil {
		defer h.onDisconnect()
	}

	events, cancel := h.uc.Execute()
	defer cancel()

	for {
		select {
		case <-r.Context().Done():
			return
		case evt := <-events:
			if evt.Type == appport.FileChangeTypeCreate || evt.Type == appport.FileChangeTypeRemove {
				dir := path.Dir(evt.Path)
				if dir == "." {
					dir = ""
				}
				treeMsg, _ := json.Marshal(sseMessage{Type: "tree_change", Path: dir})
				_, _ = fmt.Fprintf(w, "data: %s\n\n", treeMsg)
				flusher.Flush()
			}
			fileMsg, _ := json.Marshal(sseMessage{Type: "file_change", Path: evt.Path})
			_, _ = fmt.Fprintf(w, "data: %s\n\n", fileMsg)
			flusher.Flush()
		}
	}
}
