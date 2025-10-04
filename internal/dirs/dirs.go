package dirs

import (
	"os"
	"path/filepath"
	"strings"
)

type Directory struct {
	Path  string
	Depth int
}

// TODO: extract to config + different domain
var dirsConfig = []Directory{
	{Path: "$HOME/.config", Depth: 1},
	{Path: "$HOME/.local/share/nvim/lazy/LazyVim"},
	{Path: "$HOME/Codes", Depth: 1},
	{Path: "$HOME/Dotfiles", Depth: 1},
	{Path: "$HOME/GoVault"},
	{Path: "$HOME/GoVault/Scratch"},
	{Path: "$HOME/Work", Depth: 1},
	{Path: "$HOME/Work/_projects"},
}

func GetDirectories() []string {
	var result []string
	for _, dir := range dirsConfig {
		expandedPath := os.ExpandEnv(dir.Path)

		result = append(result, expandedPath)

		if dir.Depth > 0 {
			result = append(result, getSubDirectories(expandedPath, dir.Depth)...)
		}
	}

	return result
}

func getSubDirectories(basePath string, depth int) []string {
	var result []string

	entries, err := os.ReadDir(basePath)
	if err != nil {
		return result
	}

	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			fullPath := filepath.Join(basePath, entry.Name())
			result = append(result, fullPath)

			if depth > 1 {
				result = append(result, getSubDirectories(fullPath, depth-1)...)
			}
		}
	}

	return result
}
