package repository

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/naoki-higashi-28/mdprev/internal/domain/model"
)

func TestFileSystemRepositoryListEntries(t *testing.T) {
	root := t.TempDir()
	repo := NewFileSystemRepository(root)

	if err := os.Mkdir(filepath.Join(root, "bdir"), 0o755); err != nil {
		t.Fatalf("mkdir bdir failed: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, "Adir"), 0o755); err != nil {
		t.Fatalf("mkdir Adir failed: %v", err)
	}
	if err := os.Mkdir(filepath.Join(root, ".hidden"), 0o755); err != nil {
		t.Fatalf("mkdir .hidden failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "z.md"), []byte("z"), 0o644); err != nil {
		t.Fatalf("write z.md failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "A.txt"), []byte("a"), 0o644); err != nil {
		t.Fatalf("write A.txt failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, ".secret.md"), []byte("s"), 0o644); err != nil {
		t.Fatalf("write .secret.md failed: %v", err)
	}

	entries, err := repo.ListEntries(root)
	if err != nil {
		t.Fatalf("ListEntries error: %v", err)
	}

	if len(entries) != 4 {
		t.Fatalf("len(entries) = %d, want 4", len(entries))
	}

	if entries[0].Type != model.EntryTypeDir || entries[0].Name != "Adir" {
		t.Fatalf("entries[0] = %#v", entries[0])
	}
	if entries[1].Type != model.EntryTypeDir || entries[1].Name != "bdir" {
		t.Fatalf("entries[1] = %#v", entries[1])
	}
	if entries[2].Type != model.EntryTypeFile || entries[2].Name != "A.txt" {
		t.Fatalf("entries[2] = %#v", entries[2])
	}
	if entries[3].Type != model.EntryTypeFile || entries[3].Name != "z.md" {
		t.Fatalf("entries[3] = %#v", entries[3])
	}
}

func TestFileSystemRepositorySearchEntries(t *testing.T) {
	root := t.TempDir()
	repo := NewFileSystemRepository(root)

	if err := os.MkdirAll(filepath.Join(root, "docs"), 0o755); err != nil {
		t.Fatalf("mkdir docs failed: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(root, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir .git failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "Guide.MD"), []byte("x"), 0o644); err != nil {
		t.Fatalf("write Guide.MD failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "guide.markdown"), []byte("x"), 0o644); err != nil {
		t.Fatalf("write guide.markdown failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "docs", "note.txt"), []byte("x"), 0o644); err != nil {
		t.Fatalf("write note.txt failed: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, ".git", "guide.md"), []byte("x"), 0o644); err != nil {
		t.Fatalf("write hidden guide.md failed: %v", err)
	}

	entries, err := repo.SearchEntries([]string{"guide"})
	if err != nil {
		t.Fatalf("SearchEntries error: %v", err)
	}

	if len(entries) != 1 {
		t.Fatalf("len(entries) = %d, want 1", len(entries))
	}
	if entries[0].Name != "guide.markdown" {
		t.Fatalf("entry name = %q, want guide.markdown", entries[0].Name)
	}
	if entries[0].Path != "docs/guide.markdown" {
		t.Fatalf("entry path = %q, want docs/guide.markdown", entries[0].Path)
	}

	t.Run("parent directory match", func(t *testing.T) {
		got, err := repo.SearchEntries([]string{"docs"})
		if err != nil {
			t.Fatalf("SearchEntries error: %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("len(got) = %d, want 1", len(got))
		}
		if got[0].Path != "docs/guide.markdown" {
			t.Fatalf("path = %q, want docs/guide.markdown", got[0].Path)
		}
	})

	t.Run("AND search matches both terms", func(t *testing.T) {
		got, err := repo.SearchEntries([]string{"docs", "guide"})
		if err != nil {
			t.Fatalf("SearchEntries error: %v", err)
		}
		if len(got) != 1 {
			t.Fatalf("len(got) = %d, want 1", len(got))
		}
	})

	t.Run("AND search returns empty when one term misses", func(t *testing.T) {
		got, err := repo.SearchEntries([]string{"docs", "zzz"})
		if err != nil {
			t.Fatalf("SearchEntries error: %v", err)
		}
		if len(got) != 0 {
			t.Fatalf("len(got) = %d, want 0", len(got))
		}
	})
}
