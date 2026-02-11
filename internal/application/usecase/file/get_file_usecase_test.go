package file

import (
	"errors"
	"testing"

	"github.com/naoki-higashi-28/mdprev/internal/application/middleware"
	"github.com/naoki-higashi-28/mdprev/internal/domain/model"
)

type mockFileRepo struct {
	readFn func(path string) ([]byte, error)
	rawFn  func(path string) (string, error)
}

func (m *mockFileRepo) ReadMarkdown(path string) ([]byte, error) {
	if m.readFn != nil {
		return m.readFn(path)
	}
	return nil, nil
}

func (m *mockFileRepo) ServeRaw(path string) (string, error) {
	if m.rawFn != nil {
		return m.rawFn(path)
	}
	return "", nil
}

func TestGetFileUseCaseExecute(t *testing.T) {
	root := t.TempDir()
	validator, err := middleware.NewPathValidator(root)
	if err != nil {
		t.Fatalf("NewPathValidator error: %v", err)
	}

	t.Run("empty path", func(t *testing.T) {
		uc := NewGetFileUseCase(validator, &mockFileRepo{})
		_, err := uc.Execute("")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, model.ErrBadRequest) {
			t.Fatalf("error = %v, want ErrBadRequest", err)
		}
	})

	t.Run("repository success", func(t *testing.T) {
		want := []byte("# hello")
		uc := NewGetFileUseCase(validator, &mockFileRepo{
			readFn: func(path string) ([]byte, error) {
				return want, nil
			},
		})
		got, err := uc.Execute("docs/a.md")
		if err != nil {
			t.Fatalf("Execute error: %v", err)
		}
		if string(got) != string(want) {
			t.Fatalf("got %q, want %q", string(got), string(want))
		}
	})

	t.Run("path validation error", func(t *testing.T) {
		uc := NewGetFileUseCase(validator, &mockFileRepo{})
		_, err := uc.Execute("../a.md")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, model.ErrForbidden) {
			t.Fatalf("error = %v, want ErrForbidden", err)
		}
	})

	t.Run("repository error", func(t *testing.T) {
		repoErr := errors.New("read failed")
		uc := NewGetFileUseCase(validator, &mockFileRepo{
			readFn: func(path string) ([]byte, error) {
				return nil, repoErr
			},
		})
		_, err := uc.Execute("docs/a.md")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, repoErr) {
			t.Fatalf("error = %v, want wrapped repo error", err)
		}
	})
}
