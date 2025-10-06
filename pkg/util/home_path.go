package util

import (
	"os"
	"strings"
)

func TruncateHomePath(path string) string {
	return strings.Replace(path, os.ExpandEnv("$HOME"), "~", 1)
}
