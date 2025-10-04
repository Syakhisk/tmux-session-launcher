package fuzzyfinder

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"tmux-session-launcher/internal/dirs"
	"tmux-session-launcher/internal/tmux"
	"tmux-session-launcher/pkg/logger"

	"github.com/fatih/color"

	"emperror.dev/errors"
)

const (
	separatorColumns = "::"
	separatorFzf     = "|"
)

var (
	colorCategory = color.New(color.FgHiBlue, color.Italic).Sprint
	// colorCurrentSession = color.New(color.FgHiGreen, color.Bold).Sprint
	colorPath = color.New(color.FgHiBlack, color.Italic).Sprint
	colorMute = color.RGB(0, 0, 0).Sprint
)

func Exec(ctx context.Context) error {
	logger := logger.WithPrefix("fuzzyfinder.Exec")

	args := []string{
		"--ansi",
		"--no-sort",
		"--no-hscroll",
		"--delimiter", separatorFzf, // used as nth delimiter
		"--with-nth", "1,2", // what to show in the list
		"--nth", "2", // what to search in (based on with-nth)
		"--accept-nth", "{3,4}", // what to output on accept
	}

	cmd := exec.CommandContext(ctx, "fzf", args...)

	input, err := buildInput(ctx)
	if err != nil {
		return errors.WrapIf(err, "failed to build fzf input")
	}

	reader := bytes.NewReader(input)
	cmd.Stdin = reader

	output, err := cmd.CombinedOutput()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode := exitErr.ExitCode()
			if exitCode == 130 {
				logger.Info("fzf was cancelled by user (exit code 130)")
				return nil
			}
		}

		return errors.Wrapf(err, "fzf failed with stderr: %s", string(output))
	}

	logger.Infof("fzf result: %s", string(output))

	return nil
}

func buildInput(ctx context.Context) ([]byte, error) {
	sessions, err := tmux.GetSessions(ctx)
	if err != nil {
		logger.Warnf("Failed to get tmux sessions: %v", err)
	}

	dirs := dirs.GetDirectories()

	formattedSessions := formatEntryTmuxSessions(sessions, separatorColumns, separatorFzf)
	formattedDirs := formatEntryDirectory(dirs, separatorColumns, separatorFzf)

	prettified, _ := prettifyColumns(formattedSessions+formattedDirs, separatorColumns)

	return []byte(prettified), nil
}

func prettifyColumns(input string, separator string) (string, error) {
	cmd := exec.Command(
		"column",
		"-t",
		"-s", separator,
	)
	cmd.Stdin = strings.NewReader(input)

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Errorf("Failed to prettify columns (exit code %v): %s", err, string(output))
		return input, nil // Return original input on error
	}

	return strings.TrimSpace(string(output)), nil
}

func formatEntryTmuxSessions(sessions []tmux.Session, columnSep, fzfSep string) string {
	var builder strings.Builder

	for _, s := range sessions {
		sessionName := s.Name
		if s.Current {
			// sessionName = colorCurrentSession(sessionName)
			sessionName = fmt.Sprintf("[%s]", sessionName)
		}

		builder.WriteString(colorCategory("session"))
		builder.WriteString(columnSep)
		builder.WriteString(sessionName)
		builder.WriteString(columnSep)
		builder.WriteString(colorPath(s.Path))

		// fzf metadata: display|searchable|type|id
		builder.WriteString(colorMute(
			columnSep + fzfSep + s.Name + columnSep + fzfSep,
		))

		builder.WriteString("session")
		builder.WriteString(fzfSep)
		builder.WriteString(s.ID)

		builder.WriteString("\n")
	}

	return builder.String()
}

func formatEntryDirectory(dirs []string, columnSep, fzfSep string) string {
	var builder strings.Builder

	for _, d := range dirs {
		base := filepath.Base(d)

		truncatedHome := strings.Replace(d, os.ExpandEnv("$HOME"), "~", 1)
		builder.WriteString(colorCategory("directory"))
		builder.WriteString(columnSep)
		builder.WriteString(base)
		builder.WriteString(columnSep)
		builder.WriteString(colorPath(truncatedHome))

		// fzf metadata: display|searchable|type|path
		builder.WriteString(colorMute(
			columnSep + fzfSep + truncatedHome + columnSep + fzfSep,
		))

		builder.WriteString("directory")
		builder.WriteString(fzfSep)
		builder.WriteString(d)

		builder.WriteString("\n")
	}

	return builder.String()
}
