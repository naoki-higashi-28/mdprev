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
	watchedDirs map[string]struct{}
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
		watchedDirs: make(map[string]struct{}),
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
			fw.watchedDirs[path] = struct{}{}
			return fw.watcher.Add(path)
		}
		return nil
	})
}

func (fw *FileWatcher) loop() {
	debounce := make(map[string]domain.ChangeType)
	timer := time.NewTimer(0)
	if !timer.Stop() {
		<-timer.C
	}

	resetTimer := func() {
		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
		timer.Reset(100 * time.Millisecond)
	}

	for {
		select {
		case <-fw.done:
			timer.Stop()
			return

		case <-timer.C:
			fw.mu.Lock()
			subs := make([]chan domain.FileChangeEvent, 0, len(fw.subscribers))
			for ch := range fw.subscribers {
				subs = append(subs, ch)
			}
			fw.mu.Unlock()

			for p, ct := range debounce {
				evt := domain.FileChangeEvent{Path: p, Type: ct}
				for _, ch := range subs {
					select {
					case ch <- evt:
					default:
					}
				}
			}
			debounce = make(map[string]domain.ChangeType)

		case event, ok := <-fw.watcher.Events:
			if !ok {
				return
			}

			// Add newly created directories to the watch list
			if event.Has(fsnotify.Create) {
				if fw.tryAddDir(event.Name) {
					rel, err := filepath.Rel(fw.root, event.Name)
					if err == nil {
						relSlash := filepath.ToSlash(rel)
						if prev, exists := debounce[relSlash]; !exists || prev == domain.ChangeWrite {
							debounce[relSlash] = domain.ChangeCreate
						}
						resetTimer()
					}
				}
			}

			// Detect directory removal
			if event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
				if _, wasDir := fw.watchedDirs[event.Name]; wasDir {
					delete(fw.watchedDirs, event.Name)
					rel, err := filepath.Rel(fw.root, event.Name)
					if err == nil {
						relSlash := filepath.ToSlash(rel)
						if prev, exists := debounce[relSlash]; !exists || prev == domain.ChangeWrite {
							debounce[relSlash] = domain.ChangeRemove
						}
						resetTimer()
					}
				}
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

				var ct domain.ChangeType
				switch {
				case event.Has(fsnotify.Create):
					ct = domain.ChangeCreate
				case event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename):
					ct = domain.ChangeRemove
				default:
					ct = domain.ChangeWrite
				}

				// Create/Remove takes priority over Write for the same path
				if prev, exists := debounce[relSlash]; !exists || prev == domain.ChangeWrite {
					debounce[relSlash] = ct
				}

				resetTimer()
			}

		case _, ok := <-fw.watcher.Errors:
			if !ok {
				return
			}
		}
	}
}

func (fw *FileWatcher) tryAddDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return false
	}
	name := filepath.Base(path)
	if strings.HasPrefix(name, ".") {
		return false
	}
	// Recursively add subdirectories (handles cp -r scenarios)
	if err := fw.addDirs(path); err != nil {
		_ = fw.watcher.Add(path)
		fw.watchedDirs[path] = struct{}{}
	}
	return true
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
