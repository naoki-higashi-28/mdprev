package middleware

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/naoki-higashi-28/mdprev/internal/domain/model"
)

// PathValidator validates that paths don't escape the root directory.
type PathValidator struct {
	root string
}

// NewPathValidator creates a new PathValidator for the given root directory.
func NewPathValidator(root string) (*PathValidator, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, fmt.Errorf("resolving root path: %w", err)
	}
	evalRoot, err := filepath.EvalSymlinks(absRoot)
	if err != nil {
		return nil, fmt.Errorf("evaluating symlinks for root: %w", err)
	}
	return &PathValidator{root: evalRoot}, nil
}

// Root returns the resolved root directory path.
func (v *PathValidator) Root() string {
	return v.root
}

// Validate checks that the given relative path stays within the root directory.
// It returns the absolute file path if valid, or an error if the path escapes root.
func (v *PathValidator) Validate(relPath string) (string, error) {
	cleaned := filepath.Clean(relPath)
	absPath := filepath.Join(v.root, cleaned)

	// Ensure joined path stays within root.
	relFromRoot, err := filepath.Rel(v.root, absPath)
	if err != nil {
		return "", fmt.Errorf("computing relative path: %w", err)
	}
	if relFromRoot == ".." || strings.HasPrefix(relFromRoot, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("path escapes root directory: %w", model.ErrForbidden)
	}

	// Resolve symlinks and check again.
	evalPath, err := filepath.EvalSymlinks(absPath)
	if err != nil {
		// If the path doesn't exist, return the cleaned path.
		// The caller will handle the not-found error.
		if os.IsNotExist(err) {
			return absPath, nil
		}
		return "", fmt.Errorf("evaluating symlinks: %w", err)
	}

	relFromRoot, err = filepath.Rel(v.root, evalPath)
	if err != nil {
		return "", fmt.Errorf("computing relative path after symlink resolution: %w", err)
	}
	if relFromRoot == ".." || strings.HasPrefix(relFromRoot, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("symlink escapes root directory: %w", model.ErrForbidden)
	}

	return evalPath, nil
}
