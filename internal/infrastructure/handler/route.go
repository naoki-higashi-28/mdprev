package handler

import (
	"encoding/json"
	"net/http"
)

// Routes groups handlers used for route registration.
type Routes struct {
	Tree     *TreeHandler
	File     *FileHandler
	Search   *SearchHandler
	Watch    *WatchHandler
	RootName string
}

// SetupRoutes registers all API routes on the given mux.
func SetupRoutes(mux *http.ServeMux, routes Routes) {
	mux.HandleFunc("GET /api/tree", routes.Tree.HandleGetTree)
	mux.HandleFunc("GET /api/file", routes.File.HandleGetFile)
	mux.HandleFunc("GET /api/search", routes.Search.HandleSearch)
	mux.HandleFunc("GET /raw/{path...}", routes.File.HandleGetRaw)
	mux.HandleFunc("GET /api/watch", routes.Watch.HandleWatch)
	mux.HandleFunc("GET /api/info", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"name": routes.RootName})
	})
}
