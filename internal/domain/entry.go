package domain

// EntryType represents the type of a file system entry.
type EntryType string

const (
	EntryTypeDir  EntryType = "dir"
	EntryTypeFile EntryType = "file"
)

// Entry represents a file or directory entry.
type Entry struct {
	Type  EntryType `json:"type"`
	Name  string    `json:"name"`
	Path  string    `json:"path"`
	Ext   string    `json:"ext,omitempty"`
	Size  int64     `json:"size,omitempty"`
	Mtime int64     `json:"mtime,omitempty"`
}

// TreeResult represents the result of a directory listing.
type TreeResult struct {
	Path    string  `json:"path"`
	Entries []Entry `json:"entries"`
}

// SearchResult represents the result of a file search.
type SearchResult struct {
	Query   string  `json:"query"`
	Entries []Entry `json:"entries"`
}
