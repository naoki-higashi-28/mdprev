package port

// FileChangeType represents the kind of file system change.
type FileChangeType string

const (
	FileChangeTypeWrite  FileChangeType = "write"
	FileChangeTypeCreate FileChangeType = "create"
	FileChangeTypeRemove FileChangeType = "remove"
)

// FileChangeEvent represents a file change notification.
type FileChangeEvent struct {
	Path string
	Type FileChangeType
}

// FileChangeSubscriber provides file change watching operations.
type FileChangeSubscriber interface {
	Subscribe() (events <-chan FileChangeEvent, cancel func())
	Close() error
}
