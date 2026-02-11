package tree

import (
	"errors"
	"testing"

	"github.com/naoki-higashi-28/mdprev/internal/application/middleware"
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

func TestGetTreeUseCaseExecute(t *testing.T) {
	root := t.TempDir()
	validator, err := middleware.NewPathValidator(root)
	if err != nil {
		t.Fatalf("NewPathValidator error: %v", err)
	}

	t.Run("default path", func(t *testing.T) {
		uc := NewGetTreeUseCase(validator, &mockTreeRepo{
			listFn: func(path string) ([]model.Entry, error) {
				if path != validator.Root() {
					t.Fatalf("path = %q, want root %q", path, validator.Root())
				}
				return []model.Entry{{Type: model.EntryTypeDir, Name: "docs", Path: "docs"}}, nil
			},
		})

		got, err := uc.Execute("")
		if err != nil {
			t.Fatalf("Execute error: %v", err)
		}
		if got.Path != "" || len(got.Entries) != 1 {
			t.Fatalf("unexpected result: %#v", got)
		}
	})

	t.Run("nested path", func(t *testing.T) {
		uc := NewGetTreeUseCase(validator, &mockTreeRepo{
			listFn: func(path string) ([]model.Entry, error) {
				return []model.Entry{}, nil
			},
		})
		got, err := uc.Execute("docs")
		if err != nil {
			t.Fatalf("Execute error: %v", err)
		}
		if got.Path != "docs" {
			t.Fatalf("got path = %q, want docs", got.Path)
		}
	})

	t.Run("validation error", func(t *testing.T) {
		uc := NewGetTreeUseCase(validator, &mockTreeRepo{})
		_, err := uc.Execute("../etc")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, model.ErrForbidden) {
			t.Fatalf("error = %v, want ErrForbidden", err)
		}
	})

	t.Run("repository error", func(t *testing.T) {
		repoErr := errors.New("list failed")
		uc := NewGetTreeUseCase(validator, &mockTreeRepo{
			listFn: func(path string) ([]model.Entry, error) {
				return nil, repoErr
			},
		})
		_, err := uc.Execute(".")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, repoErr) {
			t.Fatalf("error = %v, want wrapped repo error", err)
		}
	})
}
