package tmux

import (
	"context"
	"os/exec"
	"path/filepath"
	"strings"
	"tmux-session-launcher/pkg/util"

	"emperror.dev/errors"
)

const (
	sessionFormat = "#{session_id}|#{session_name}|#{session_path}"
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
		sessionFormat,
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if err := handleTmuxError(string(output)); err != nil {
			return nil, err
		}

		return nil, errors.WrapIf(err, "failed to get current tmux session")
	}

	return sessionParse(strings.TrimSpace(string(output)))
}

func GetSessions(ctx context.Context) ([]Session, error) {
	cmd := exec.CommandContext(
		ctx,
		"tmux",
		// "-L", "test",
		"list-sessions",
		"-F", sessionFormat,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		if err := handleTmuxError(string(output)); err != nil {
			return []Session{}, err
		}
	}

	var sessions []Session

	currentSession, _ := GetCurrentSession(ctx)
	if currentSession != nil {
		sessions = append(sessions, *currentSession)
	}

	// PERF: can use cmd.StdoutPipe and a scanner to avoid loading all output in memory
	for line := range strings.SplitSeq(string(output), "\n") {
		session, err := sessionParse(line)
		if err != nil {
			continue
		}

		if currentSession != nil && session.ID == currentSession.ID {
			continue
		}

		sessions = append(sessions, *session)
	}

	return sessions, nil
}

func SessionCreate(ctx context.Context, name, path string) (*Session, error) {
	cmd := exec.CommandContext(
		ctx,
		"tmux",
		"new-session",
		"-d",       // detached
		"-s", name, // session name
		"-c", path, // start directory
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		if err := handleTmuxError(string(output)); err != nil {
			return nil, err
		}

		return nil, errors.WrapIf(err, "failed to create tmux session")
	}

	cmd = exec.CommandContext(
		ctx,
		"tmux",
		"list-sessions",
		"-f", "#{==:#{session_name},"+name+"}",
		"-F", sessionFormat,
	)

	output, err = cmd.CombinedOutput()
	if err != nil {
		if err := handleTmuxError(string(output)); err != nil {
			return nil, err
		}

		return nil, errors.WrapIf(err, "failed to get created tmux session")
	}

	return sessionParse(strings.TrimSpace(string(output)))
}

func SessionAttach(ctx context.Context, id string) error {
	cmd := exec.CommandContext(
		ctx,
		"tmux",
		"switch-client",
		"-t", id,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		if err := handleTmuxError(string(output)); err != nil {
			return err
		}

		return errors.WrapIf(err, "failed to attach tmux session")
	}

	return nil
}

func SessionCreateOrAttach(ctx context.Context, name, path string) (*Session, error) {
	session, err := SessionCreate(ctx, name, path)
	if err != nil {
		if !errors.Is(err, ErrSessionExists) {
			return nil, errors.WrapIf(err, "failed to create tmux session")
		}
	}

	err = SessionAttach(ctx, name)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to attach tmux session")
	}

	return session, nil
}

func BuildSessionNameFromPath(path string) string {
	base := filepath.Base(path)

	replacer := strings.NewReplacer(
		" ", "_",
		".", "_",
	)

	return replacer.Replace(base)
}

func sessionParse(line string) (*Session, error) {
	parts := strings.SplitN(line, "|", 3)
	if len(parts) < 3 {
		return nil, errors.New("unexpected output from tmux list-sessions")
	}

	return &Session{
		ID:   strings.TrimSpace(parts[0]),
		Name: strings.TrimSpace(parts[1]),
		Path: util.TruncateHomePath(strings.TrimSpace(parts[2])),
	}, nil
}
