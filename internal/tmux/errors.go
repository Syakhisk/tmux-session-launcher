package tmux

import (
	"strings"

	"emperror.dev/errors"
)

const (
	ErrTmuxNotRunning  = errors.Sentinel("tmux server is not running")
	ErrSessionExists   = errors.Sentinel("tmux session already exists")
	ErrSessionNotFound = errors.Sentinel("tmux session not found")
)

func isTmuxNotRunningErr(output string) bool {
	return strings.HasPrefix(output, "no server running")
}

func isSessionNotFoundErr(output string) bool {
	return strings.HasPrefix(output, "can't find session")
}

func isSessionExistsErr(output string) bool {
	return strings.HasPrefix(output, "duplicate session")
}

func handleTmuxError(output string) error {
	if isTmuxNotRunningErr(output) {
		return ErrTmuxNotRunning
	}

	if isSessionNotFoundErr(output) {
		return ErrSessionNotFound
	}

	if isSessionExistsErr(output) {
		return ErrSessionExists
	}

	return nil
}
