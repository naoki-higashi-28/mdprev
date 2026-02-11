package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	searchuc "github.com/naoki-higashi-28/mdprev/internal/application/usecase/search"
	"github.com/naoki-higashi-28/mdprev/internal/domain/model"
)

func TestSearchHandlerHandleSearch(t *testing.T) {
	h := NewSearchHandler(searchuc.NewSearchFilesUseCase(&mockTreeRepo{
		searchFn: func(terms []string) ([]model.Entry, error) {
			if len(terms) == 0 {
				return []model.Entry{}, nil
			}
			return []model.Entry{{Type: model.EntryTypeFile, Name: "guide.md", Path: "guide.md", Ext: "md"}}, nil
		},
	}))

	t.Run("empty query", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/search?q=", nil)
		rr := httptest.NewRecorder()
		h.HandleSearch(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rr.Code)
		}
		if !strings.Contains(rr.Body.String(), "\"entries\":[]") {
			t.Fatalf("unexpected body: %s", rr.Body.String())
		}
	})

	t.Run("with result", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/search?q=guide", nil)
		rr := httptest.NewRecorder()
		h.HandleSearch(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rr.Code)
		}
		if !strings.Contains(rr.Body.String(), "guide.md") {
			t.Fatalf("unexpected body: %s", rr.Body.String())
		}
	})
}
