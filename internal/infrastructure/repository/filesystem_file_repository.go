package repository

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/naoki-higashi-28/mdprev/internal/domain/model"
)

var allowedRawExts = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".gif":  true,
	".svg":  true,
	".webp": true,
}

// ReadMarkdown reads a markdown file from an absolute path and returns its content.
func (r *FileSystemRepository) ReadMarkdown(absPath string) ([]byte, error) {
	ext := filepath.Ext(absPath)
	if !allowedMarkdownExts[ext] {
		return nil, fmt.Errorf("extension %q not allowed: %w", ext, model.ErrInvalidExtFile)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %w", model.ErrNotFound)
		}
		return nil, fmt.Errorf("reading file: %w", err)
	}

	return data, nil
}

// ServeRaw validates a raw file at an absolute path and returns it for serving.
func (r *FileSystemRepository) ServeRaw(absPath string) (string, error) {
	ext := filepath.Ext(absPath)
	if !allowedRawExts[ext] {
		return "", fmt.Errorf("extension %q not allowed for raw serving: %w", ext, model.ErrInvalidExtFile)
	}

	if _, err := os.Stat(absPath); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("file not found: %w", model.ErrNotFound)
		}
		return "", fmt.Errorf("checking file: %w", err)
	}

	return absPath, nil
}
