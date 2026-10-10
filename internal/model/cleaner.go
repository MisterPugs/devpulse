package model

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// CleanItem represents a heavy developer folder or cache target.
type CleanItem struct {
	Path     string
	Category string
	SizeByte int64
	Selected bool
}

// FormatBytes formats byte sizes into human-readable strings.
func FormatBytes(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.2f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

// Known bloat folder names
var devTargets = map[string]string{
	"node_modules": "Node.js Dependencies",
	".next":        "Next.js Build Cache",
	".nuxt":        "Nuxt Build Cache",
	".turbo":       "Turborepo Cache",
	"target":       "Rust/Maven Build Artifacts",
	".gradle":      "Gradle Cache",
	"__pycache__":  "Python Bytecode Cache",
	".pytest_cache": "Pytest Cache",
	".venv":        "Python Virtualenv",
	"venv":         "Python Virtualenv",
	".cache":       "Local Tool Cache",
	"build":        "Generic Build Folder",
	"dist":         "Distribution Build",
}

// ScanDevBloat scans the specified root directory up to maxDepth for bloat directories.
func ScanDevBloat(rootDir string, maxDepth int) ([]CleanItem, error) {
	if rootDir == "" {
		var err error
		rootDir, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}

	var results []CleanItem
	var mu sync.Mutex

	rootDepth := strings.Count(filepath.Clean(rootDir), string(filepath.Separator))

	filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() {
			return nil
		}

		currentDepth := strings.Count(filepath.Clean(path), string(filepath.Separator)) - rootDepth
		if currentDepth > maxDepth {
			return filepath.SkipDir
		}

		base := d.Name()
		if cat, exists := devTargets[base]; exists {
			// Don't scan root itself if named like target
			if path == rootDir {
				return nil
			}

			// Calculate folder size in a lightweight manner
			size := calculateDirSize(path)
			mu.Lock()
			results = append(results, CleanItem{
				Path:     path,
				Category: cat,
				SizeByte: size,
				Selected: false,
			})
			mu.Unlock()

			// Don't descend inside node_modules to find another node_modules
			return filepath.SkipDir
		}

		// Skip hidden git folders
		if base == ".git" {
			return filepath.SkipDir
		}

		return nil
	})

	return results, nil
}

func calculateDirSize(path string) int64 {
	var total int64
	filepath.WalkDir(path, func(_ string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() {
			if info, err := d.Info(); err == nil {
				total += info.Size()
			}
		}
		return nil
	})
	return total
}

// DeleteItem removes a bloat directory from disk.
func DeleteItem(item CleanItem) error {
	return os.RemoveAll(item.Path)
}
