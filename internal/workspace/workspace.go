package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"tmux-session-launcher/internal/config"
	"tmux-session-launcher/pkg/logger"
	"tmux-session-launcher/pkg/util"
)

type Directory struct {
	FullPath          string
	TruncatedHomePath string
	Parent            string
	Label             string
}

// GetDirectories loads directory configuration from YAML file and returns all directories
func GetDirectories() []Directory {
	var allDirs []Directory

	// Load configuration from YAML file
	cfg, err := config.Load()
	if err != nil {
		logger.Errorf("Failed to load configuration: %v", err)
		return allDirs
	}

	// Collect all directories first
	for _, dir := range cfg.Directories {
		expandedPath := os.ExpandEnv(dir.Path)
		base := filepath.Base(expandedPath)
		truncatedHome := util.TruncateHomePath(expandedPath)

		allDirs = append(allDirs, Directory{
			FullPath:          expandedPath,
			Parent:            base,
			Label:             base,
			TruncatedHomePath: truncatedHome,
		})

		if dir.Depth > 0 {
			allDirs = append(allDirs, getSubDirectories(expandedPath, base, dir.Depth)...)
		}
	}

	// Deduplicate and return
	return deduplicateDirectories(allDirs)
}

func getSubDirectories(basePath string, baseLabel string, depth int) []Directory {
	var result []Directory

	entries, err := os.ReadDir(basePath)
	if err != nil {
		return result
	}

	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			fullPath := filepath.Join(basePath, entry.Name())
			baseLabel := filepath.Join(baseLabel, entry.Name())
			truncatedHome := util.TruncateHomePath(fullPath)

			result = append(result, Directory{
				FullPath:          fullPath,
				Parent:            baseLabel,
				Label:             entry.Name(),
				TruncatedHomePath: truncatedHome,
			})

			if depth > 1 {
				result = append(result, getSubDirectories(fullPath, baseLabel, depth-1)...)
			}
		}
	}

	return result
}

// deduplicateDirectories removes duplicate directories based on FullPath
func deduplicateDirectories(dirs []Directory) []Directory {
	seen := make(map[string]struct{})
	var result []Directory

	for _, dir := range dirs {
		if _, exists := seen[dir.FullPath]; !exists {
			seen[dir.FullPath] = struct{}{}
			result = append(result, dir)
		}
	}

	return result
}
