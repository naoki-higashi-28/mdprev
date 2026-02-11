package repository

import "github.com/naoki-higashi-28/mdprev/internal/domain/model"

// TreeRepository provides directory listing operations.
type TreeRepository interface {
	ListEntries(path string) ([]model.Entry, error)
	SearchEntries(query string) ([]model.Entry, error)
}
