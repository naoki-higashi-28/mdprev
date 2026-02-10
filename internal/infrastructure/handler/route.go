package handler

import (
	"encoding/json"
	"net/http"
	"path/filepath"

	fileuc "github.com/naoki-higashi-28/mdprev/internal/application/usecase/file"
	searchuc "github.com/naoki-higashi-28/mdprev/internal/application/usecase/search"
	treeuc "github.com/naoki-higashi-28/mdprev/internal/application/usecase/tree"
	"github.com/naoki-higashi-28/mdprev/internal/domain"
)

// SetupRoutes registers all API routes on the given mux.
func SetupRoutes(mux *http.ServeMux, treeUC *treeuc.UseCase, fileUC *fileuc.UseCase, searchUC *searchuc.UseCase, watchRepo domain.WatchRepository, rootPath string, onConnect, onDisconnect func()) {
	treeHandler := NewTreeHandler(treeUC)
	fileHandler := NewFileHandler(fileUC)
	searchHandler := NewSearchHandler(searchUC)
	watchHandler := NewWatchHandler(watchRepo, onConnect, onDisconnect)

	rootName := filepath.Base(rootPath)

	mux.HandleFunc("GET /api/tree", treeHandler.HandleGetTree)
	mux.HandleFunc("GET /api/file", fileHandler.HandleGetFile)
	mux.HandleFunc("GET /api/search", searchHandler.HandleSearch)
	mux.HandleFunc("GET /raw/{path...}", fileHandler.HandleGetRaw)
	mux.HandleFunc("GET /api/watch", watchHandler.HandleWatch)
	mux.HandleFunc("GET /api/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"name": rootName})
	})
}
