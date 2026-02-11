package port

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fsnotify/fsnotify"
	appport "github.com/naoki-higashi-28/mdprev/internal/application/port"
)

func TestIsMarkdown(t *testing.T) {
	tests := []struct {
		name string
		path string
		want bool
	}{
		{name: "md", path: "a.md", want: true},
		{name: "markdown", path: "a.markdown", want: true},
		{name: "upper case", path: "A.MD", want: true},
		{name: "non markdown", path: "a.txt", want: false},
		{name: "no extension", path: "README", want: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if got := isMarkdown(tt.path); got != tt.want {
				t.Fatalf("isMarkdown(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestAddDirsSkipsHiddenDirectories(t *testing.T) {
	root := t.TempDir()
	visibleDir := filepath.Join(root, "docs")
	hiddenDir := filepath.Join(root, ".git")

	if err := os.MkdirAll(visibleDir, 0o755); err != nil {
		t.Fatalf("failed to create visible dir: %v", err)
	}
	if err := os.MkdirAll(hiddenDir, 0o755); err != nil {
		t.Fatalf("failed to create hidden dir: %v", err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		t.Fatalf("failed to create watcher: %v", err)
	}
	t.Cleanup(func() {
		_ = watcher.Close()
	})

	s := &FileSystemFileChangeSubscriber{
		root:        root,
		watcher:     watcher,
		done:        make(chan struct{}),
		subscribers: make(map[chan appport.FileChangeEvent]struct{}),
		watchedDirs: map[string]struct{}{},
	}

	if err := s.addDirs(root); err != nil {
		t.Fatalf("addDirs returned error: %v", err)
	}

	if _, ok := s.watchedDirs[root]; !ok {
		t.Fatalf("root directory is not watched")
	}
	if _, ok := s.watchedDirs[visibleDir]; !ok {
		t.Fatalf("visible directory is not watched")
	}
	if _, ok := s.watchedDirs[hiddenDir]; ok {
		t.Fatalf("hidden directory should not be watched")
	}
}

func TestTryAddDir(t *testing.T) {
	root := t.TempDir()
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		t.Fatalf("failed to create watcher: %v", err)
	}
	t.Cleanup(func() {
		_ = watcher.Close()
	})

	s := &FileSystemFileChangeSubscriber{
		root:        root,
		watcher:     watcher,
		done:        make(chan struct{}),
		subscribers: make(map[chan appport.FileChangeEvent]struct{}),
		watchedDirs: make(map[string]struct{}),
	}

	visibleDir := filepath.Join(root, "notes")
	hiddenDir := filepath.Join(root, ".cache")
	if err := os.MkdirAll(visibleDir, 0o755); err != nil {
		t.Fatalf("failed to create visible dir: %v", err)
	}
	if err := os.MkdirAll(hiddenDir, 0o755); err != nil {
		t.Fatalf("failed to create hidden dir: %v", err)
	}

	if got := s.tryAddDir(hiddenDir); got {
		t.Fatalf("tryAddDir(%q) = true, want false", hiddenDir)
	}
	if got := s.tryAddDir(visibleDir); !got {
		t.Fatalf("tryAddDir(%q) = false, want true", visibleDir)
	}
	if _, ok := s.watchedDirs[visibleDir]; !ok {
		t.Fatalf("visible dir was not added to watchedDirs")
	}
}

func TestSubscribeAndCancel(t *testing.T) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		t.Fatalf("failed to create watcher: %v", err)
	}
	t.Cleanup(func() {
		_ = watcher.Close()
	})

	s := &FileSystemFileChangeSubscriber{
		watcher:     watcher,
		done:        make(chan struct{}),
		subscribers: make(map[chan appport.FileChangeEvent]struct{}),
		watchedDirs: make(map[string]struct{}),
	}

	ch, cancel := s.Subscribe()
	if ch == nil {
		t.Fatalf("Subscribe returned nil channel")
	}

	s.mu.Lock()
	if got := len(s.subscribers); got != 1 {
		s.mu.Unlock()
		t.Fatalf("subscriber count = %d, want 1", got)
	}
	s.mu.Unlock()

	cancel()

	s.mu.Lock()
	if got := len(s.subscribers); got != 0 {
		s.mu.Unlock()
		t.Fatalf("subscriber count after cancel = %d, want 0", got)
	}
	s.mu.Unlock()
}

func TestNewFileSystemFileChangeSubscriber_ClosesWatcherOnAddDirsError(t *testing.T) {
	root := t.TempDir()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		t.Fatalf("failed to create watcher: %v", err)
	}
	if err := watcher.Close(); err != nil {
		t.Fatalf("failed to close watcher for setup: %v", err)
	}

	original := newFSNotifyWatcher
	newFSNotifyWatcher = func() (*fsnotify.Watcher, error) {
		return watcher, nil
	}
	t.Cleanup(func() {
		newFSNotifyWatcher = original
	})

	sub, err := NewFileSystemFileChangeSubscriber(root)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if sub != nil {
		t.Fatalf("subscriber = %#v, want nil", sub)
	}

	if !strings.Contains(err.Error(), "closed") {
		t.Fatalf("error %q does not contain expected text", err.Error())
	}
}
