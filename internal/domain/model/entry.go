package model

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
