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

// FileChangeEvent represents a file change notification.
type FileChangeEvent struct {
	Path string // relative path from root (slash-delimited)
}

// WatchRepository provides file change watching operations.
type WatchRepository interface {
	Subscribe() (events <-chan FileChangeEvent, cancel func())
	Close() error
}
