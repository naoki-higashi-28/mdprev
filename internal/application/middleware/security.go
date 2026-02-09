package middleware

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/naoki-higashi-28/mdprev/internal/domain"
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

	// Check the cleaned path doesn't escape root
	if !strings.HasPrefix(absPath, v.root) {
		return "", fmt.Errorf("path escapes root directory: %w", domain.ErrForbidden)
	}

	// Resolve symlinks and check again
	evalPath, err := filepath.EvalSymlinks(absPath)
	if err != nil {
		// If the file doesn't exist, return the cleaned path
		// The caller will handle the not-found error
		return absPath, nil
	}

	if !strings.HasPrefix(evalPath, v.root) {
		return "", fmt.Errorf("symlink escapes root directory: %w", domain.ErrForbidden)
	}

	return evalPath, nil
}
