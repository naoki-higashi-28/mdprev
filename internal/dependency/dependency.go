package dependency

import (
	"fmt"
	"io/fs"
	"net/http"
	"strings"

	"github.com/naoki-higashi-28/mdprev/internal/application/middleware"
	fileuc "github.com/naoki-higashi-28/mdprev/internal/application/usecase/file"
	searchuc "github.com/naoki-higashi-28/mdprev/internal/application/usecase/search"
	treeuc "github.com/naoki-higashi-28/mdprev/internal/application/usecase/tree"
	"github.com/naoki-higashi-28/mdprev/internal/infrastructure"
	"github.com/naoki-higashi-28/mdprev/internal/infrastructure/handler"
)

// NewServeMux creates a fully configured http.ServeMux with all routes and SPA fallback.
func NewServeMux(root string, staticFS fs.FS) (*http.ServeMux, *infrastructure.FileWatcher, error) {
	validator, err := middleware.NewPathValidator(root)
	if err != nil {
		return nil, nil, fmt.Errorf("initializing path validator: %w", err)
	}

	fsRepo := infrastructure.NewFileSystemRepository(validator.Root())

	treeUseCase := treeuc.NewUseCase(validator, fsRepo)
	fileUseCase := fileuc.NewUseCase(validator, fsRepo)
	searchUseCase := searchuc.NewUseCase(fsRepo)

	watcher, err := infrastructure.NewFileWatcher(validator.Root())
	if err != nil {
		return nil, nil, fmt.Errorf("initializing file watcher: %w", err)
	}

	mux := http.NewServeMux()
	handler.SetupRoutes(mux, treeUseCase, fileUseCase, searchUseCase, watcher, validator.Root())
	setupSPA(mux, staticFS)

	return mux, watcher, nil
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
			f.Close()
			fileServer.ServeHTTP(w, r)
			return
		}

		// SPA fallback: serve index.html for non-API, non-raw paths
		r.URL.Path = "/"
		fileServer.ServeHTTP(w, r)
	})
}
