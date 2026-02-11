package search

import (
	"errors"
	"testing"

	"github.com/naoki-higashi-28/mdprev/internal/domain/model"
)

type mockTreeRepo struct {
	listFn   func(path string) ([]model.Entry, error)
	searchFn func(terms []string) ([]model.Entry, error)
}

func (m *mockTreeRepo) ListEntries(path string) ([]model.Entry, error) {
	if m.listFn != nil {
		return m.listFn(path)
	}
	return nil, nil
}

func (m *mockTreeRepo) SearchEntries(terms []string) ([]model.Entry, error) {
	if m.searchFn != nil {
		return m.searchFn(terms)
	}
	return nil, nil
}

func TestSearchFilesUseCaseExecute(t *testing.T) {
	t.Run("empty query", func(t *testing.T) {
		called := false
		uc := NewSearchFilesUseCase(&mockTreeRepo{
			searchFn: func(terms []string) ([]model.Entry, error) {
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

	t.Run("spaces only query", func(t *testing.T) {
		called := false
		uc := NewSearchFilesUseCase(&mockTreeRepo{
			searchFn: func(terms []string) ([]model.Entry, error) {
				called = true
				return nil, nil
			},
		})

		got, err := uc.Execute("   ")
		if err != nil {
			t.Fatalf("Execute error: %v", err)
		}
		if called {
			t.Fatalf("SearchEntries should not be called for spaces-only query")
		}
		if len(got.Entries) != 0 {
			t.Fatalf("unexpected result: %#v", got)
		}
	})

	t.Run("success", func(t *testing.T) {
		want := []model.Entry{{Type: model.EntryTypeFile, Name: "a.md", Path: "a.md", Ext: "md"}}
		uc := NewSearchFilesUseCase(&mockTreeRepo{
			searchFn: func(terms []string) ([]model.Entry, error) {
				if len(terms) != 1 || terms[0] != "a" {
					t.Fatalf("terms = %v, want [a]", terms)
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

	t.Run("AND search splits on space", func(t *testing.T) {
		uc := NewSearchFilesUseCase(&mockTreeRepo{
			searchFn: func(terms []string) ([]model.Entry, error) {
				if len(terms) != 2 || terms[0] != "foo" || terms[1] != "bar" {
					t.Fatalf("terms = %v, want [foo bar]", terms)
				}
				return []model.Entry{}, nil
			},
		})

		_, err := uc.Execute("foo bar")
		if err != nil {
			t.Fatalf("Execute error: %v", err)
		}
	})

	t.Run("consecutive spaces produce no empty terms", func(t *testing.T) {
		uc := NewSearchFilesUseCase(&mockTreeRepo{
			searchFn: func(terms []string) ([]model.Entry, error) {
				if len(terms) != 2 {
					t.Fatalf("terms = %v, want exactly 2 terms", terms)
				}
				return []model.Entry{}, nil
			},
		})

		_, err := uc.Execute("foo  bar")
		if err != nil {
			t.Fatalf("Execute error: %v", err)
		}
	})

	t.Run("repository error", func(t *testing.T) {
		repoErr := errors.New("search failed")
		uc := NewSearchFilesUseCase(&mockTreeRepo{
			searchFn: func(terms []string) ([]model.Entry, error) {
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
