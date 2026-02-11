package repository

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/naoki-higashi-28/mdprev/internal/domain/model"
)

func TestFileSystemRepositoryReadMarkdown(t *testing.T) {
	repo := NewFileSystemRepository(t.TempDir())

	t.Run("invalid extension", func(t *testing.T) {
		_, err := repo.ReadMarkdown("/tmp/a.txt")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, model.ErrInvalidExtFile) {
			t.Fatalf("error = %v, want ErrInvalidExtFile", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.ReadMarkdown("/tmp/not-found.md")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, model.ErrNotFound) {
			t.Fatalf("error = %v, want ErrNotFound", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		root := t.TempDir()
		path := filepath.Join(root, "a.md")
		if err := os.WriteFile(path, []byte("# hi"), 0o644); err != nil {
			t.Fatalf("WriteFile failed: %v", err)
		}
		data, err := repo.ReadMarkdown(path)
		if err != nil {
			t.Fatalf("ReadMarkdown error: %v", err)
		}
		if string(data) != "# hi" {
			t.Fatalf("data = %q, want %q", string(data), "# hi")
		}
	})
}

func TestFileSystemRepositoryServeRaw(t *testing.T) {
	repo := NewFileSystemRepository(t.TempDir())

	t.Run("invalid extension", func(t *testing.T) {
		_, err := repo.ServeRaw("/tmp/a.md")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, model.ErrInvalidExtFile) {
			t.Fatalf("error = %v, want ErrInvalidExtFile", err)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.ServeRaw("/tmp/not-found.png")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, model.ErrNotFound) {
			t.Fatalf("error = %v, want ErrNotFound", err)
		}
	})

	t.Run("success", func(t *testing.T) {
		root := t.TempDir()
		path := filepath.Join(root, "a.png")
		if err := os.WriteFile(path, []byte("png"), 0o644); err != nil {
			t.Fatalf("WriteFile failed: %v", err)
		}
		got, err := repo.ServeRaw(path)
		if err != nil {
			t.Fatalf("ServeRaw error: %v", err)
		}
		if got != path {
			t.Fatalf("got = %q, want %q", got, path)
		}
	})
}
