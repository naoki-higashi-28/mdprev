package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/naoki-higashi-28/mdprev/internal/application/middleware"
	treeuc "github.com/naoki-higashi-28/mdprev/internal/application/usecase/tree"
	"github.com/naoki-higashi-28/mdprev/internal/domain/model"
)

func TestTreeHandlerHandleGetTree(t *testing.T) {
	root := t.TempDir()
	validator, err := middleware.NewPathValidator(root)
	if err != nil {
		t.Fatalf("NewPathValidator error: %v", err)
	}

	h := NewTreeHandler(treeuc.NewGetTreeUseCase(validator, &mockTreeRepo{
		listFn: func(path string) ([]model.Entry, error) {
			return []model.Entry{{Type: model.EntryTypeDir, Name: "docs", Path: "docs"}}, nil
		},
	}))

	req := httptest.NewRequest(http.MethodGet, "/api/tree?path=.", nil)
	rr := httptest.NewRecorder()
	h.HandleGetTree(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
}
