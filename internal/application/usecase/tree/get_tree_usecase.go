package tree

import (
	"fmt"
	"path/filepath"

	"github.com/naoki-higashi-28/mdprev/internal/application/middleware"
	"github.com/naoki-higashi-28/mdprev/internal/domain/model"
	"github.com/naoki-higashi-28/mdprev/internal/domain/repository"
)

// TreeResult represents the result of a directory listing.
type TreeResult struct {
	Path    string        `json:"path"`
	Entries []model.Entry `json:"entries"`
}

// GetTreeUseCase handles tree listing operations.
type GetTreeUseCase struct {
	validator *middleware.PathValidator
	repo      repository.TreeRepository
}

// NewGetTreeUseCase creates a new GetTreeUseCase.
func NewGetTreeUseCase(validator *middleware.PathValidator, repo repository.TreeRepository) *GetTreeUseCase {
	return &GetTreeUseCase{validator: validator, repo: repo}
}

// Execute returns the directory listing for the given path.
func (uc *GetTreeUseCase) Execute(path string) (TreeResult, error) {
	if path == "" {
		path = "."
	}

	absPath, err := uc.validator.Validate(path)
	if err != nil {
		return TreeResult{}, fmt.Errorf("validating path: %w", err)
	}

	entries, err := uc.repo.ListEntries(absPath)
	if err != nil {
		return TreeResult{}, fmt.Errorf("listing entries: %w", err)
	}

	relPath, err := filepath.Rel(uc.validator.Root(), absPath)
	if err != nil {
		return TreeResult{}, fmt.Errorf("computing relative path: %w", err)
	}
	if relPath == "." {
		relPath = ""
	}

	return TreeResult{Path: relPath, Entries: entries}, nil
}
