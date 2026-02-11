package file

import (
	"errors"
	"testing"

	"github.com/naoki-higashi-28/mdprev/internal/application/middleware"
	"github.com/naoki-higashi-28/mdprev/internal/domain/model"
)

func TestGetRawFileUseCaseExecute(t *testing.T) {
	root := t.TempDir()
	validator, err := middleware.NewPathValidator(root)
	if err != nil {
		t.Fatalf("NewPathValidator error: %v", err)
	}

	t.Run("empty path", func(t *testing.T) {
		uc := NewGetRawFileUseCase(validator, &mockFileRepo{})
		_, err := uc.Execute("")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, model.ErrBadRequest) {
			t.Fatalf("error = %v, want ErrBadRequest", err)
		}
	})

	t.Run("repository success", func(t *testing.T) {
		uc := NewGetRawFileUseCase(validator, &mockFileRepo{
			rawFn: func(path string) (string, error) {
				return path, nil
			},
		})
		got, err := uc.Execute("images/a.png")
		if err != nil {
			t.Fatalf("Execute error: %v", err)
		}
		if got == "" {
			t.Fatalf("got empty path")
		}
	})

	t.Run("validation error", func(t *testing.T) {
		uc := NewGetRawFileUseCase(validator, &mockFileRepo{})
		_, err := uc.Execute("../a.png")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, model.ErrForbidden) {
			t.Fatalf("error = %v, want ErrForbidden", err)
		}
	})

	t.Run("repository error", func(t *testing.T) {
		repoErr := errors.New("raw failed")
		uc := NewGetRawFileUseCase(validator, &mockFileRepo{
			rawFn: func(path string) (string, error) {
				return "", repoErr
			},
		})
		_, err := uc.Execute("images/a.png")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, repoErr) {
			t.Fatalf("error = %v, want wrapped repo error", err)
		}
	})
}
