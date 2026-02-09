package file

import (
	"fmt"

	"github.com/naoki-higashi-28/mdprev/internal/application/middleware"
	"github.com/naoki-higashi-28/mdprev/internal/domain"
)

// UseCase handles file reading operations.
type UseCase struct {
	validator *middleware.PathValidator
	repo      domain.FileRepository
}

// NewUseCase creates a new file UseCase.
func NewUseCase(validator *middleware.PathValidator, repo domain.FileRepository) *UseCase {
	return &UseCase{
		validator: validator,
		repo:      repo,
	}
}

// GetFile reads a markdown file and returns its content.
func (uc *UseCase) GetFile(path string) ([]byte, error) {
	if path == "" {
		return nil, fmt.Errorf("path is required: %w", domain.ErrBadRequest)
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

// GetRaw returns the absolute path of a raw file for serving.
func (uc *UseCase) GetRaw(path string) (string, error) {
	if path == "" {
		return "", fmt.Errorf("path is required: %w", domain.ErrBadRequest)
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
