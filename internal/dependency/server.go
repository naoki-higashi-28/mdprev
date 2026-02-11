package dependency

import (
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/naoki-higashi-28/mdprev/internal/application/middleware"
	appport "github.com/naoki-higashi-28/mdprev/internal/application/port"
	infraAdapter "github.com/naoki-higashi-28/mdprev/internal/infrastructure/adapter"
	"github.com/naoki-higashi-28/mdprev/internal/infrastructure/handler"
	infraRepo "github.com/naoki-higashi-28/mdprev/internal/infrastructure/repository"
)

// NewServerMux creates a fully configured http.ServeMux with all routes and SPA fallback.
func NewServerMux(root string, staticFS fs.FS, onConnect, onDisconnect func()) (*http.ServeMux, appport.FileChangeSubscriber, error) {
	validator, err := middleware.NewPathValidator(root)
	if err != nil {
		return nil, nil, fmt.Errorf("initializing path validator: %w", err)
	}

	fsRepo := infraRepo.NewFileSystemRepository(validator.Root())
	subscriber, err := infraAdapter.NewFileSystemFileChangeSubscriber(validator.Root())
	if err != nil {
		return nil, nil, fmt.Errorf("initializing file change subscriber: %w", err)
	}

	useCases := NewUseCases(validator, fsRepo, fsRepo, subscriber)
	routes := NewRoutes(useCases, validator.Root(), onConnect, onDisconnect)

	mux := http.NewServeMux()
	handler.SetupRoutes(mux, routes)
	setupSPA(mux, staticFS)

	return mux, subscriber, nil
}

func setupSPA(mux *http.ServeMux, staticFS fs.FS) {
	fileServer := http.FileServer(http.FS(staticFS))

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" {
			fileServer.ServeHTTP(w, r)
			return
		}

		cleanPath := strings.TrimPrefix(path, "/")
		if f, err := staticFS.Open(cleanPath); err == nil {
			_ = f.Close()
			fileServer.ServeHTTP(w, r)
			return
		}

		r.URL.Path = "/"
		fileServer.ServeHTTP(w, r)
	})
}
