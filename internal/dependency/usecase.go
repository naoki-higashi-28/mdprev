package dependency

import (
	"github.com/naoki-higashi-28/mdprev/internal/application/middleware"
	appport "github.com/naoki-higashi-28/mdprev/internal/application/port"
	fileuc "github.com/naoki-higashi-28/mdprev/internal/application/usecase/file"
	searchuc "github.com/naoki-higashi-28/mdprev/internal/application/usecase/search"
	treeuc "github.com/naoki-higashi-28/mdprev/internal/application/usecase/tree"
	watchuc "github.com/naoki-higashi-28/mdprev/internal/application/usecase/watch"
	"github.com/naoki-higashi-28/mdprev/internal/domain/repository"
)

// UseCases groups all use cases used by the HTTP layer.
type UseCases struct {
	GetTree              *treeuc.GetTreeUseCase
	GetFile              *fileuc.GetFileUseCase
	GetRawFile           *fileuc.GetRawFileUseCase
	SearchFiles          *searchuc.SearchFilesUseCase
	SubscribeFileChanges *watchuc.SubscribeFileChangesUseCase
}

// NewUseCases builds all use cases.
func NewUseCases(validator *middleware.PathValidator, fileRepo repository.FileRepository, treeRepo repository.TreeRepository, subscriber appport.FileChangeSubscriber) UseCases {
	return UseCases{
		GetTree:              treeuc.NewGetTreeUseCase(validator, treeRepo),
		GetFile:              fileuc.NewGetFileUseCase(validator, fileRepo),
		GetRawFile:           fileuc.NewGetRawFileUseCase(validator, fileRepo),
		SearchFiles:          searchuc.NewSearchFilesUseCase(treeRepo),
		SubscribeFileChanges: watchuc.NewSubscribeFileChangesUseCase(subscriber),
	}
}
