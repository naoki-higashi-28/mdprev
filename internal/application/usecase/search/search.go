package search

import (
	"fmt"

	"github.com/naoki-higashi-28/mdprev/internal/domain"
)

// UseCase handles file search operations.
type UseCase struct {
	repo domain.TreeRepository
}

// NewUseCase creates a new search UseCase.
func NewUseCase(repo domain.TreeRepository) *UseCase {
	return &UseCase{repo: repo}
}

// Search returns markdown files matching the query.
func (uc *UseCase) Search(query string) (domain.SearchResult, error) {
	if query == "" {
		return domain.SearchResult{
			Query:   "",
			Entries: []domain.Entry{},
		}, nil
	}

	result, err := uc.repo.SearchEntries(query)
	if err != nil {
		return domain.SearchResult{}, fmt.Errorf("searching entries: %w", err)
	}

	return result, nil
}
