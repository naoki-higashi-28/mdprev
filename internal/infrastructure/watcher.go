package infrastructure

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/naoki-higashi-28/mdprev/internal/domain"
)

// FileWatcher watches for markdown file changes and broadcasts events to subscribers.
type FileWatcher struct {
	root        string
	watcher     *fsnotify.Watcher
	done        chan struct{}
	mu          sync.Mutex
	subscribers map[chan domain.FileChangeEvent]struct{}
}

// NewFileWatcher creates a new FileWatcher that monitors the given root directory.
func NewFileWatcher(root string) (*FileWatcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	fw := &FileWatcher{
		root:        root,
		watcher:     w,
		done:        make(chan struct{}),
		subscribers: make(map[chan domain.FileChangeEvent]struct{}),
	}

	if err := fw.addDirs(root); err != nil {
		w.Close()
		return nil, err
	}

	go fw.loop()

	return fw, nil
}

func (fw *FileWatcher) addDirs(root string) error {
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if path != root && strings.HasPrefix(d.Name(), ".") {
				return filepath.SkipDir
			}
			return fw.watcher.Add(path)
		}
		return nil
	})
}

func (fw *FileWatcher) loop() {
	debounce := make(map[string]struct{})
	var timer *time.Timer

	flush := func() {
		fw.mu.Lock()
		subs := make([]chan domain.FileChangeEvent, 0, len(fw.subscribers))
		for ch := range fw.subscribers {
			subs = append(subs, ch)
		}
		fw.mu.Unlock()

		for path := range debounce {
			evt := domain.FileChangeEvent{Path: path}
			for _, ch := range subs {
				select {
				case ch <- evt:
				default:
				}
			}
		}
		debounce = make(map[string]struct{})
	}

	for {
		select {
		case <-fw.done:
			if timer != nil {
				timer.Stop()
			}
			return
		case event, ok := <-fw.watcher.Events:
			if !ok {
				return
			}

			// Add newly created directories to the watch list
			if event.Has(fsnotify.Create) {
				fw.tryAddDir(event.Name)
			}

			if !isMarkdown(event.Name) {
				continue
			}

			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) || event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
				rel, err := filepath.Rel(fw.root, event.Name)
				if err != nil {
					continue
				}
				relSlash := filepath.ToSlash(rel)
				debounce[relSlash] = struct{}{}

				if timer != nil {
					timer.Stop()
				}
				timer = time.AfterFunc(100*time.Millisecond, flush)
			}
		case _, ok := <-fw.watcher.Errors:
			if !ok {
				return
			}
		}
	}
}

func (fw *FileWatcher) tryAddDir(path string) {
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return
	}
	name := filepath.Base(path)
	if strings.HasPrefix(name, ".") {
		return
	}
	_ = fw.watcher.Add(path)
}

func isMarkdown(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".md" || ext == ".markdown"
}

// Subscribe returns a channel that receives file change events and a cancel function.
func (fw *FileWatcher) Subscribe() (<-chan domain.FileChangeEvent, func()) {
	ch := make(chan domain.FileChangeEvent, 16)
	fw.mu.Lock()
	fw.subscribers[ch] = struct{}{}
	fw.mu.Unlock()

	cancel := func() {
		fw.mu.Lock()
		delete(fw.subscribers, ch)
		fw.mu.Unlock()
	}

	return ch, cancel
}

// Close stops the file watcher.
func (fw *FileWatcher) Close() error {
	close(fw.done)
	return fw.watcher.Close()
}
