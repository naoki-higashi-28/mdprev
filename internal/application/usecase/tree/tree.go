package tree

import (
	"fmt"

	"github.com/naoki-higashi-28/mdprev/internal/application/middleware"
	"github.com/naoki-higashi-28/mdprev/internal/domain"
)

// UseCase handles tree listing operations.
type UseCase struct {
	validator *middleware.PathValidator
	repo      domain.TreeRepository
}

// NewUseCase creates a new tree UseCase.
func NewUseCase(validator *middleware.PathValidator, repo domain.TreeRepository) *UseCase {
	return &UseCase{
		validator: validator,
		repo:      repo,
	}
}

// GetTree returns the directory listing for the given path.
func (uc *UseCase) GetTree(path string) (domain.TreeResult, error) {
	if path == "" {
		path = "."
	}

	absPath, err := uc.validator.Validate(path)
	if err != nil {
		return domain.TreeResult{}, fmt.Errorf("validating path: %w", err)
	}

	result, err := uc.repo.ListEntries(absPath)
	if err != nil {
		return domain.TreeResult{}, fmt.Errorf("listing entries: %w", err)
	}

	return result, nil
}
