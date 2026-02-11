package middleware

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/naoki-higashi-28/mdprev/internal/domain/model"
)

func TestNewPathValidatorAndRoot(t *testing.T) {
	root := t.TempDir()
	v, err := NewPathValidator(root)
	if err != nil {
		t.Fatalf("NewPathValidator returned error: %v", err)
	}

	absRoot, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatalf("EvalSymlinks failed: %v", err)
	}
	if v.Root() != absRoot {
		t.Fatalf("Root() = %q, want %q", v.Root(), absRoot)
	}
}

func TestPathValidatorValidate(t *testing.T) {
	root := t.TempDir()
	v, err := NewPathValidator(root)
	if err != nil {
		t.Fatalf("NewPathValidator returned error: %v", err)
	}

	t.Run("normal path", func(t *testing.T) {
		p, err := v.Validate("docs/a.md")
		if err != nil {
			t.Fatalf("Validate returned error: %v", err)
		}
		want := filepath.Join(v.Root(), "docs", "a.md")
		if p != want {
			t.Fatalf("path = %q, want %q", p, want)
		}
	})

	t.Run("escape with dotdot", func(t *testing.T) {
		_, err := v.Validate("../secret.txt")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, model.ErrForbidden) {
			t.Fatalf("error = %v, want ErrForbidden", err)
		}
	})

	t.Run("symlink escape", func(t *testing.T) {
		outsideDir := t.TempDir()
		outsideFile := filepath.Join(outsideDir, "outside.md")
		if err := os.WriteFile(outsideFile, []byte("x"), 0o644); err != nil {
			t.Fatalf("failed to create outside file: %v", err)
		}

		symlinkPath := filepath.Join(root, "link.md")
		if err := os.Symlink(outsideFile, symlinkPath); err != nil {
			t.Skipf("symlink unsupported on this environment: %v", err)
		}

		_, err := v.Validate("link.md")
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if !errors.Is(err, model.ErrForbidden) {
			t.Fatalf("error = %v, want ErrForbidden", err)
		}
	})
}
