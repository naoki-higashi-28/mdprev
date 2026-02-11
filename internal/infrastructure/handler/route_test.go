package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/naoki-higashi-28/mdprev/internal/application/middleware"
	appport "github.com/naoki-higashi-28/mdprev/internal/application/port"
	fileuc "github.com/naoki-higashi-28/mdprev/internal/application/usecase/file"
	searchuc "github.com/naoki-higashi-28/mdprev/internal/application/usecase/search"
	treeuc "github.com/naoki-higashi-28/mdprev/internal/application/usecase/tree"
	watchuc "github.com/naoki-higashi-28/mdprev/internal/application/usecase/watch"
)

func TestSetupRoutesInfoAndEndpoints(t *testing.T) {
	root := t.TempDir()
	validator, err := middleware.NewPathValidator(root)
	if err != nil {
		t.Fatalf("NewPathValidator error: %v", err)
	}

	routes := Routes{
		Tree:     NewTreeHandler(treeuc.NewGetTreeUseCase(validator, &mockTreeRepo{})),
		File:     NewFileHandler(fileuc.NewGetFileUseCase(validator, &mockFileRepo{}), fileuc.NewGetRawFileUseCase(validator, &mockFileRepo{})),
		Search:   NewSearchHandler(searchuc.NewSearchFilesUseCase(&mockTreeRepo{})),
		Watch:    NewWatchHandler(watchuc.NewSubscribeFileChangesUseCase(&mockFileChangeSubscriber{events: make(chan appport.FileChangeEvent)}), nil, nil),
		RootName: "mdprev",
	}

	mux := http.NewServeMux()
	SetupRoutes(mux, routes)

	t.Run("api info", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/info", nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rr.Code)
		}
		if !strings.Contains(rr.Body.String(), "mdprev") {
			t.Fatalf("unexpected body: %s", rr.Body.String())
		}
	})

	t.Run("api tree is wired", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/tree", nil)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rr.Code)
		}
	})
}
