package port

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	appport "github.com/naoki-higashi-28/mdprev/internal/application/port"
)

// FileSystemFileChangeSubscriber watches for markdown file changes and broadcasts events.
type FileSystemFileChangeSubscriber struct {
	root        string
	watcher     *fsnotify.Watcher
	done        chan struct{}
	mu          sync.Mutex
	subscribers map[chan appport.FileChangeEvent]struct{}
	watchedDirs map[string]struct{}
}

var newFSNotifyWatcher = fsnotify.NewWatcher

// NewFileSystemFileChangeSubscriber creates a new file change subscriber.
func NewFileSystemFileChangeSubscriber(root string) (*FileSystemFileChangeSubscriber, error) {
	w, err := newFSNotifyWatcher()
	if err != nil {
		return nil, err
	}

	s := &FileSystemFileChangeSubscriber{
		root:        root,
		watcher:     w,
		done:        make(chan struct{}),
		subscribers: make(map[chan appport.FileChangeEvent]struct{}),
		watchedDirs: make(map[string]struct{}),
	}

	if err := s.addDirs(root); err != nil {
		if closeErr := w.Close(); closeErr != nil {
			return nil, errors.Join(err, closeErr)
		}
		return nil, err
	}

	go s.loop()
	return s, nil
}

func (s *FileSystemFileChangeSubscriber) addDirs(root string) error {
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if path != root && strings.HasPrefix(d.Name(), ".") {
				return filepath.SkipDir
			}
			s.watchedDirs[path] = struct{}{}
			return s.watcher.Add(path)
		}
		return nil
	})
}

func (s *FileSystemFileChangeSubscriber) loop() {
	debounce := make(map[string]appport.FileChangeType)
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
		case <-s.done:
			timer.Stop()
			return
		case <-timer.C:
			s.mu.Lock()
			subs := make([]chan appport.FileChangeEvent, 0, len(s.subscribers))
			for ch := range s.subscribers {
				subs = append(subs, ch)
			}
			s.mu.Unlock()

			for p, ct := range debounce {
				evt := appport.FileChangeEvent{Path: p, Type: ct}
				for _, ch := range subs {
					select {
					case ch <- evt:
					default:
					}
				}
			}
			debounce = make(map[string]appport.FileChangeType)

		case event, ok := <-s.watcher.Events:
			if !ok {
				return
			}

			if event.Has(fsnotify.Create) {
				if s.tryAddDir(event.Name) {
					rel, err := filepath.Rel(s.root, event.Name)
					if err == nil {
						relSlash := filepath.ToSlash(rel)
						if prev, exists := debounce[relSlash]; !exists || prev == appport.FileChangeTypeWrite {
							debounce[relSlash] = appport.FileChangeTypeCreate
						}
						resetTimer()
					}
				}
			}

			if event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
				if _, wasDir := s.watchedDirs[event.Name]; wasDir {
					delete(s.watchedDirs, event.Name)
					rel, err := filepath.Rel(s.root, event.Name)
					if err == nil {
						relSlash := filepath.ToSlash(rel)
						if prev, exists := debounce[relSlash]; !exists || prev == appport.FileChangeTypeWrite {
							debounce[relSlash] = appport.FileChangeTypeRemove
						}
						resetTimer()
					}
				}
			}

			if !isMarkdown(event.Name) {
				continue
			}

			if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) || event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename) {
				rel, err := filepath.Rel(s.root, event.Name)
				if err != nil {
					continue
				}
				relSlash := filepath.ToSlash(rel)

				var ct appport.FileChangeType
				switch {
				case event.Has(fsnotify.Create):
					ct = appport.FileChangeTypeCreate
				case event.Has(fsnotify.Remove) || event.Has(fsnotify.Rename):
					ct = appport.FileChangeTypeRemove
				default:
					ct = appport.FileChangeTypeWrite
				}

				if prev, exists := debounce[relSlash]; !exists || prev == appport.FileChangeTypeWrite {
					debounce[relSlash] = ct
				}

				resetTimer()
			}

		case _, ok := <-s.watcher.Errors:
			if !ok {
				return
			}
		}
	}
}

func (s *FileSystemFileChangeSubscriber) tryAddDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return false
	}
	name := filepath.Base(path)
	if strings.HasPrefix(name, ".") {
		return false
	}
	if err := s.addDirs(path); err != nil {
		_ = s.watcher.Add(path)
		s.watchedDirs[path] = struct{}{}
	}
	return true
}

func isMarkdown(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".md" || ext == ".markdown"
}

// Subscribe returns a channel that receives file change events and a cancel function.
func (s *FileSystemFileChangeSubscriber) Subscribe() (<-chan appport.FileChangeEvent, func()) {
	ch := make(chan appport.FileChangeEvent, 16)
	s.mu.Lock()
	s.subscribers[ch] = struct{}{}
	s.mu.Unlock()

	cancel := func() {
		s.mu.Lock()
		delete(s.subscribers, ch)
		s.mu.Unlock()
	}

	return ch, cancel
}

// Close stops the file watcher.
func (s *FileSystemFileChangeSubscriber) Close() error {
	close(s.done)
	return s.watcher.Close()
}
