package file

import (
	"fmt"

	"github.com/naoki-higashi-28/mdprev/internal/application/middleware"
	"github.com/naoki-higashi-28/mdprev/internal/domain/model"
	"github.com/naoki-higashi-28/mdprev/internal/domain/repository"
)

// GetRawFileUseCase handles raw file path resolution for serving.
type GetRawFileUseCase struct {
	validator *middleware.PathValidator
	repo      repository.FileRepository
}

// NewGetRawFileUseCase creates a new GetRawFileUseCase.
func NewGetRawFileUseCase(validator *middleware.PathValidator, repo repository.FileRepository) *GetRawFileUseCase {
	return &GetRawFileUseCase{validator: validator, repo: repo}
}

// Execute returns the absolute path of a raw file for serving.
func (uc *GetRawFileUseCase) Execute(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path is required: %w", model.ErrBadRequest)
	}

	absPath, err := uc.validator.Validate(path)
	if err != nil {
		return "", fmt.Errorf("validating path: %w", err)
	}

	filePath, err := uc.repo.ServeRaw(absPath)
	if err != nil {
		return "", fmt.Errorf("serving raw file: %w", err)
	}

	return filePath, nil
}
