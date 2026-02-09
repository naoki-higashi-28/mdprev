package infrastructure

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/naoki-higashi-28/mdprev/internal/domain"
)

var allowedMarkdownExts = map[string]bool{
	".md":       true,
	".markdown": true,
}

var allowedRawExts = map[string]bool{
	".png":  true,
	".jpg":  true,
	".jpeg": true,
	".gif":  true,
	".svg":  true,
	".webp": true,
}

// FileSystemRepository implements domain.TreeRepository and domain.FileRepository.
type FileSystemRepository struct {
	root string
}

// NewFileSystemRepository creates a new FileSystemRepository.
func NewFileSystemRepository(root string) *FileSystemRepository {
	return &FileSystemRepository{root: root}
}

// ListEntries returns directory entries for the given relative path.
func (r *FileSystemRepository) ListEntries(absPath string) (domain.TreeResult, error) {
	relPath, err := filepath.Rel(r.root, absPath)
	if err != nil {
		return domain.TreeResult{}, fmt.Errorf("computing relative path: %w", err)
	}
	if relPath == "." {
		relPath = ""
	}

	dirEntries, err := os.ReadDir(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return domain.TreeResult{}, fmt.Errorf("directory not found: %w", domain.ErrNotFound)
		}
		return domain.TreeResult{}, fmt.Errorf("reading directory: %w", err)
	}

	entries := []domain.Entry{}
	for _, de := range dirEntries {
		name := de.Name()

		// Skip hidden files/directories
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
			entries = append(entries, domain.Entry{
				Type: domain.EntryTypeDir,
				Name: name,
				Path: entryRelPath,
			})
		} else {
			ext := strings.TrimPrefix(filepath.Ext(name), ".")
			entries = append(entries, domain.Entry{
				Type:  domain.EntryTypeFile,
				Name:  name,
				Path:  entryRelPath,
				Ext:   ext,
				Size:  info.Size(),
				Mtime: info.ModTime().Unix(),
			})
		}
	}

	// Sort: directories first, then alphabetical (case-insensitive)
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Type != entries[j].Type {
			return entries[i].Type == domain.EntryTypeDir
		}
		return strings.ToLower(entries[i].Name) < strings.ToLower(entries[j].Name)
	})

	return domain.TreeResult{
		Path:    relPath,
		Entries: entries,
	}, nil
}

// ReadMarkdown reads a markdown file and returns its content.
func (r *FileSystemRepository) ReadMarkdown(absPath string) ([]byte, error) {
	ext := filepath.Ext(absPath)
	if !allowedMarkdownExts[ext] {
		return nil, fmt.Errorf("extension %q not allowed: %w", ext, domain.ErrInvalidExtFile)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file not found: %w", domain.ErrNotFound)
		}
		return nil, fmt.Errorf("reading file: %w", err)
	}

	return data, nil
}

// SearchEntries searches for markdown files matching the query.
func (r *FileSystemRepository) SearchEntries(query string) (domain.SearchResult, error) {
	const maxResults = 100
	queryLower := strings.ToLower(query)
	var entries []domain.Entry

	err := filepath.WalkDir(r.root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		name := d.Name()

		// Skip hidden files/directories
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

		if !strings.Contains(strings.ToLower(name), queryLower) {
			return nil
		}

		relPath, err := filepath.Rel(r.root, path)
		if err != nil {
			return nil
		}
		relPath = filepath.ToSlash(relPath)

		info, err := d.Info()
		if err != nil {
			return nil
		}

		entries = append(entries, domain.Entry{
			Type:  domain.EntryTypeFile,
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
		return domain.SearchResult{}, fmt.Errorf("walking directory: %w", err)
	}

	if entries == nil {
		entries = []domain.Entry{}
	}

	return domain.SearchResult{
		Query:   query,
		Entries: entries,
	}, nil
}

// ServeRaw returns the absolute path of a raw file for serving.
func (r *FileSystemRepository) ServeRaw(absPath string) (string, error) {
	ext := filepath.Ext(absPath)
	if !allowedRawExts[ext] {
		return "", fmt.Errorf("extension %q not allowed for raw serving: %w", ext, domain.ErrInvalidExtFile)
	}

	if _, err := os.Stat(absPath); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("file not found: %w", domain.ErrNotFound)
		}
		return "", fmt.Errorf("checking file: %w", err)
	}

	return absPath, nil
}
