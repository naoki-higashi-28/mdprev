package domain

// TreeRepository provides directory listing operations.
type TreeRepository interface {
	ListEntries(path string) (TreeResult, error)
	SearchEntries(query string) (SearchResult, error)
}

// FileRepository provides file access operations.
type FileRepository interface {
	ReadMarkdown(path string) ([]byte, error)
	ServeRaw(path string) (string, error)
}

// ChangeType represents the kind of file system change.
type ChangeType string

const (
	ChangeWrite  ChangeType = "write"
	ChangeCreate ChangeType = "create"
	ChangeRemove ChangeType = "remove"
)

// FileChangeEvent represents a file change notification.
type FileChangeEvent struct {
	Path string     // relative path from root (slash-delimited)
	Type ChangeType // kind of change
}

// WatchRepository provides file change watching operations.
type WatchRepository interface {
	Subscribe() (events <-chan FileChangeEvent, cancel func())
	Close() error
}
