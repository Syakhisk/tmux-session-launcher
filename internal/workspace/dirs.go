package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"tmux-session-launcher/pkg/util"
)

type DirectoryConfig struct {
	Path  string
	Depth int
}

type Directory struct {
	FullPath          string
	TruncatedHomePath string
	Parent            string
	Label             string
}

// TODO: Make this configurable via a config file
var dirsConfig = []DirectoryConfig{
	{Path: "$HOME/.config", Depth: 1},
	{Path: "$HOME/.local/share/nvim/lazy/LazyVim"},
	{Path: "$HOME/Codes", Depth: 2},
	{Path: "$HOME/Dotfiles", Depth: 1},
	{Path: "$HOME/GoVault"},
	{Path: "$HOME/GoVault/Scratch"},
	{Path: "$HOME/Work", Depth: 1},
	{Path: "$HOME/Work/_projects"},
}

// TODO: make the walking depth based on common project root files
func GetDirectories() []Directory {
	var result []Directory
	for _, dir := range dirsConfig {
		expandedPath := os.ExpandEnv(dir.Path)
		base := filepath.Base(expandedPath)
		truncatedHome := util.TruncateHomePath(expandedPath)

		result = append(result, Directory{
			FullPath:          expandedPath,
			Parent:            base,
			Label:             base,
			TruncatedHomePath: truncatedHome,
		})

		if dir.Depth > 0 {
			result = append(result, getSubDirectories(expandedPath, base, dir.Depth)...)
		}
	}

	return result
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
