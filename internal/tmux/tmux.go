package tmux

import (
	"context"
	"os/exec"
	"strings"

	"emperror.dev/errors"
)

const (
	ErrTmuxNotRunning = errors.Sentinel("tmux server is not running")
)

type Session struct {
	ID      string
	Name    string
	Path    string
	Current bool
}

func IsRunning(ctx context.Context) bool {
	cmd := exec.CommandContext(ctx, "tmux", "info")
	err := cmd.Run()
	return err == nil
}

func GetCurrentSession(ctx context.Context) (*Session, error) {
	cmd := exec.CommandContext(
		ctx,
		"tmux",
		"display-message",
		"-p",
		"#{session_id}|#{session_name}|#{session_path}",
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if isTmuxNotRunningErr(string(output)) {
			return nil, ErrTmuxNotRunning
		}

		return nil, errors.WrapIf(err, "failed to get current tmux session")
	}
	parts := strings.SplitN(string(output), "|", 3)
	if len(parts) < 3 {
		return nil, errors.New("unexpected output from tmux display-message")
	}

	return &Session{
		ID:      strings.TrimSpace(parts[0]),
		Name:    strings.TrimSpace(parts[1]),
		Path:    strings.TrimSpace(parts[2]),
		Current: true,
	}, nil
}

func GetSessions(ctx context.Context) ([]Session, error) {
	cmd := exec.CommandContext(
		ctx,
		"tmux",
		// "-L", "test",
		"list-sessions",
		"-F",
		"#{session_id}|#{session_name}|#{session_path}",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		if isTmuxNotRunningErr(string(output)) {
			return []Session{}, ErrTmuxNotRunning
		}
	}

	var sessions []Session

	currentSession, _ := GetCurrentSession(ctx)
	if currentSession != nil {
		sessions = append(sessions, *currentSession)
	}

	// PERF: can use cmd.StdoutPipe and a scanner to avoid loading all output in memory
	for line := range strings.SplitSeq(string(output), "\n") {
		parts := strings.SplitN(line, "|", 3)
		if len(parts) < 3 {
			continue
		}

		if currentSession != nil && parts[0] == currentSession.ID {
			continue
		}

		sessions = append(sessions, Session{
			ID:   strings.TrimSpace(parts[0]),
			Name: strings.TrimSpace(parts[1]),
			Path: strings.TrimSpace(parts[2]),
		})
	}

	return sessions, nil
}

func isTmuxNotRunningErr(output string) bool {
	return strings.HasPrefix(output, "no server running")
}
