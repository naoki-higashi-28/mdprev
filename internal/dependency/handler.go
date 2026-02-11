package dependency

import (
	"path/filepath"

	"github.com/naoki-higashi-28/mdprev/internal/infrastructure/handler"
)

// NewRoutes builds all HTTP handlers and route metadata.
func NewRoutes(useCases UseCases, rootPath string, onConnect, onDisconnect func()) handler.Routes {
	return handler.Routes{
		Tree:     handler.NewTreeHandler(useCases.GetTree),
		File:     handler.NewFileHandler(useCases.GetFile, useCases.GetRawFile),
		Search:   handler.NewSearchHandler(useCases.SearchFiles),
		Watch:    handler.NewWatchHandler(useCases.SubscribeFileChanges, onConnect, onDisconnect),
		RootName: filepath.Base(rootPath),
	}
}
