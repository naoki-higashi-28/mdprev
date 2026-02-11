// Package file
package file

import (
	"fmt"

	"github.com/naoki-higashi-28/mdprev/internal/application/middleware"
	"github.com/naoki-higashi-28/mdprev/internal/domain/model"
	"github.com/naoki-higashi-28/mdprev/internal/domain/repository"
)

// GetFileUseCase handles markdown file reading operations.
type GetFileUseCase struct {
	validator *middleware.PathValidator
	repo      repository.FileRepository
}

// NewGetFileUseCase creates a new GetFileUseCase.
func NewGetFileUseCase(validator *middleware.PathValidator, repo repository.FileRepository) *GetFileUseCase {
	return &GetFileUseCase{validator: validator, repo: repo}
}

// Execute reads a markdown file and returns its content.
func (uc *GetFileUseCase) Execute(path string) ([]byte, error) {
	if path == "" {
		return nil, fmt.Errorf("path is required: %w", model.ErrBadRequest)
	}

	absPath, err := uc.validator.Validate(path)
	if err != nil {
		return nil, fmt.Errorf("validating path: %w", err)
	}

	data, err := uc.repo.ReadMarkdown(absPath)
	if err != nil {
		return nil, fmt.Errorf("reading markdown: %w", err)
	}

	return data, nil
}
