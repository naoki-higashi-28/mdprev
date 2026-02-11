package repository

// FileRepository provides file access operations.
type FileRepository interface {
	ReadMarkdown(path string) ([]byte, error)
	ServeRaw(path string) (string, error)
}
