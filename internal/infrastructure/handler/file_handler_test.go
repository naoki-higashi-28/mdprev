package handler

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/naoki-higashi-28/mdprev/internal/application/middleware"
	fileuc "github.com/naoki-higashi-28/mdprev/internal/application/usecase/file"
	"github.com/naoki-higashi-28/mdprev/internal/infrastructure/repository"
)

func TestFileHandlerHandleGetFile(t *testing.T) {
	root := t.TempDir()
	validator, err := middleware.NewPathValidator(root)
	if err != nil {
		t.Fatalf("NewPathValidator error: %v", err)
	}

	h := NewFileHandler(
		fileuc.NewGetFileUseCase(validator, &mockFileRepo{
			readFn: func(path string) ([]byte, error) {
				return []byte("# hello"), nil
			},
		}),
		fileuc.NewGetRawFileUseCase(validator, &mockFileRepo{}),
	)

	t.Run("missing path", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/file", nil)
		rr := httptest.NewRecorder()
		h.HandleGetFile(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rr.Code)
		}
	})

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/file?path=docs/a.md", nil)
		rr := httptest.NewRecorder()
		h.HandleGetFile(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rr.Code)
		}
		if got := rr.Header().Get("Content-Type"); !strings.Contains(got, "text/plain") {
			t.Fatalf("content-type = %q", got)
		}
		if rr.Body.String() != "# hello" {
			t.Fatalf("body = %q, want %q", rr.Body.String(), "# hello")
		}
	})
}

func TestFileHandlerHandleGetRaw(t *testing.T) {
	root := t.TempDir()
	validator, err := middleware.NewPathValidator(root)
	if err != nil {
		t.Fatalf("NewPathValidator error: %v", err)
	}
	filePath := filepath.Join(root, "image.png")
	if err := os.WriteFile(filePath, []byte("pngdata"), 0o644); err != nil {
		t.Fatalf("WriteFile error: %v", err)
	}

	h := NewFileHandler(
		fileuc.NewGetFileUseCase(validator, &mockFileRepo{}),
		fileuc.NewGetRawFileUseCase(validator, repository.NewFileSystemRepository(root)),
	)

	t.Run("missing path", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/raw/", nil)
		rr := httptest.NewRecorder()
		h.HandleGetRaw(rr, req)
		if rr.Code != http.StatusBadRequest {
			t.Fatalf("status = %d, want 400", rr.Code)
		}
	})

	t.Run("success", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/raw/image.png", nil)
		req.SetPathValue("path", "image.png")
		rr := httptest.NewRecorder()
		h.HandleGetRaw(rr, req)
		if rr.Code != http.StatusOK {
			t.Fatalf("status = %d, want 200", rr.Code)
		}
		if rr.Body.String() != "pngdata" {
			t.Fatalf("body = %q, want %q", rr.Body.String(), "pngdata")
		}
	})
}
