package repository

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/naoki-higashi-28/mdprev/internal/domain/model"
)

var allowedMarkdownExts = map[string]bool{
	".md":       true,
	".markdown": true,
}

// FileSystemRepository implements domain tree/file repository interfaces.
type FileSystemRepository struct {
	root string
}

// NewFileSystemRepository creates a new FileSystemRepository.
func NewFileSystemRepository(root string) *FileSystemRepository {
	return &FileSystemRepository{root: root}
}

// ListEntries returns directory entries for the given absolute path.
func (r *FileSystemRepository) ListEntries(absPath string) ([]model.Entry, error) {
	relPath, err := filepath.Rel(r.root, absPath)
	if err != nil {
		return nil, fmt.Errorf("computing relative path: %w", err)
	}
	if relPath == "." {
		relPath = ""
	}

	dirEntries, err := os.ReadDir(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("directory not found: %w", model.ErrNotFound)
		}
		return nil, fmt.Errorf("reading directory: %w", err)
	}

	entries := []model.Entry{}
	for _, de := range dirEntries {
		name := de.Name()

		if strings.HasPrefix(name, ".") {
			continue
		}

		info, err := de.Info()
		if err != nil {
			continue
		}

		entryRelPath := name
		if relPath != "" {
			entryRelPath = relPath + "/" + name
		}

		if de.IsDir() {
			entries = append(entries, model.Entry{Type: model.EntryTypeDir, Name: name, Path: entryRelPath})
		} else {
			ext := strings.TrimPrefix(filepath.Ext(name), ".")
			entries = append(entries, model.Entry{
				Type:  model.EntryTypeFile,
				Name:  name,
				Path:  entryRelPath,
				Ext:   ext,
				Size:  info.Size(),
				Mtime: info.ModTime().Unix(),
			})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Type != entries[j].Type {
			return entries[i].Type == model.EntryTypeDir
		}
		return strings.ToLower(entries[i].Name) < strings.ToLower(entries[j].Name)
	})

	return entries, nil
}

// SearchEntries searches for markdown files matching all given terms (AND search).
// Each term is matched against the full relative path (including parent directories).
func (r *FileSystemRepository) SearchEntries(terms []string) ([]model.Entry, error) {
	const maxResults = 100
	lowerTerms := make([]string, len(terms))
	for i, t := range terms {
		lowerTerms[i] = strings.ToLower(t)
	}
	var entries []model.Entry

	err := filepath.WalkDir(r.root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		name := d.Name()
		if strings.HasPrefix(name, ".") {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if d.IsDir() {
			return nil
		}

		ext := filepath.Ext(name)
		if !allowedMarkdownExts[ext] {
			return nil
		}

		relPath, err := filepath.Rel(r.root, path)
		if err != nil {
			return nil
		}
		relPath = filepath.ToSlash(relPath)
		relPathLower := strings.ToLower(relPath)

		for _, term := range lowerTerms {
			if !strings.Contains(relPathLower, term) {
				return nil
			}
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}

		entries = append(entries, model.Entry{
			Type:  model.EntryTypeFile,
			Name:  name,
			Path:  relPath,
			Ext:   strings.TrimPrefix(ext, "."),
			Size:  info.Size(),
			Mtime: info.ModTime().Unix(),
		})
		if len(entries) >= maxResults {
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walking directory: %w", err)
	}

	if entries == nil {
		entries = []model.Entry{}
	}

	return entries, nil
}
