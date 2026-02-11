package search

import (
	"errors"
	"testing"

	"github.com/naoki-higashi-28/mdprev/internal/domain/model"
)

type mockTreeRepo struct {
	listFn   func(path string) ([]model.Entry, error)
	searchFn func(query string) ([]model.Entry, error)
}

func (m *mockTreeRepo) ListEntries(path string) ([]model.Entry, error) {
	if m.listFn != nil {
		return m.listFn(path)
	}
	return nil, nil
}

func (m *mockTreeRepo) SearchEntries(query string) ([]model.Entry, error) {
	if m.searchFn != nil {
		return m.searchFn(query)
	}
	return nil, nil
}

func TestSearchFilesUseCaseExecute(t *testing.T) {
	t.Run("empty query", func(t *testing.T) {
		called := false
		uc := NewSearchFilesUseCase(&mockTreeRepo{
			searchFn: func(query string) ([]model.Entry, error) {
				called = true
				return nil, nil
			},
		})

		got, err := uc.Execute("")
		if err != nil {
			t.Fatalf("Execute error: %v", err)
		}
		if called {
			t.Fatalf("SearchEntries should not be called for empty query")
		}
		if got.Query != "" || len(got.Entries) != 0 {
			t.Fatalf("unexpected result: %#v", got)
		}
	})

	t.Run("success", func(t *testing.T) {
		want := []model.Entry{{Type: model.EntryTypeFile, Name: "a.md", Path: "a.md", Ext: "md"}}
		uc := NewSearchFilesUseCase(&mockTreeRepo{
			searchFn: func(query string) ([]model.Entry, error) {
				if query != "a" {
					t.Fatalf("query = %q, want a", query)
				}
				return want, nil
			},
		})

		got, err := uc.Execute("a")
		if err != nil {
			t.Fatalf("Execute error: %v", err)
		}
		if got.Query != "a" || len(got.Entries) != 1 {
			t.Fatalf("unexpected result: %#v", got)
		}
	})

	t.Run("repository error", func(t *testing.T) {
		repoErr := errors.New("search failed")
		uc := NewSearchFilesUseCase(&mockTreeRepo{
			searchFn: func(query string) ([]model.Entry, error) {
				return nil, repoErr
			},
		})

		_, err := uc.Execute("x")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, repoErr) {
			t.Fatalf("error = %v, want wrapped repo error", err)
		}
	})
}
