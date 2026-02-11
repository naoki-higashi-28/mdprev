package search

import (
	"fmt"

	"github.com/naoki-higashi-28/mdprev/internal/domain/model"
	"github.com/naoki-higashi-28/mdprev/internal/domain/repository"
)

// SearchResult represents the result of a file search.
type SearchResult struct {
	Query   string        `json:"query"`
	Entries []model.Entry `json:"entries"`
}

// SearchFilesUseCase handles file search operations.
type SearchFilesUseCase struct {
	repo repository.TreeRepository
}

// NewSearchFilesUseCase creates a new SearchFilesUseCase.
func NewSearchFilesUseCase(repo repository.TreeRepository) *SearchFilesUseCase {
	return &SearchFilesUseCase{repo: repo}
}

// Execute returns markdown files matching the query.
func (uc *SearchFilesUseCase) Execute(query string) (SearchResult, error) {
	if query == "" {
		return SearchResult{Query: "", Entries: []model.Entry{}}, nil
	}

	entries, err := uc.repo.SearchEntries(query)
	if err != nil {
		return SearchResult{}, fmt.Errorf("searching entries: %w", err)
	}

	return SearchResult{Query: query, Entries: entries}, nil
}
