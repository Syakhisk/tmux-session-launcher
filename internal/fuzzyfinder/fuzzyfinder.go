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
	"unicode/utf8"

	"github.com/acarl005/stripansi"
	"github.com/fatih/color"
	"github.com/rodaine/table"

	"emperror.dev/errors"
)

const (
	separatorFzf = "|"
)

var (
	colorDefault         = fmt.Sprint
	colorCategorySession = color.New(color.FgHiCyan, color.Italic).Sprint
	colorCategoryDir     = color.New(color.FgHiBlue, color.Italic).Sprint
	colorCurrentSession  = color.New(color.FgHiGreen, color.Bold).Sprint
	colorPath            = color.New(color.FgHiBlack, color.Italic).Sprint
	colorMute            = color.RGB(0, 0, 0).Sprint
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
		"--accept-nth", "3,4", // what to output on accept
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

	formattedSessions := formatEntryTmuxSessionsAsRows(sessions, separatorFzf)
	formattedDirs := formatEntryDirectoryAsRows(dirs, separatorFzf)

	output := formatTable(append(formattedSessions, formattedDirs...))
	return []byte(output), nil
}

func formatTable(rows [][]string) string {
	var output strings.Builder
	// fzf metadata: display|searchable|type|id
	tbl := table.
		New("category", "name", "path", "searchable", "type+id").
		WithWriter(&output).
		WithWidthFunc(func(s string) int {
			return utf8.RuneCountInString(stripansi.Strip(s))
		}).
		SetRows(rows).
		WithPrintHeaders(false)

	tbl.Print()

	return output.String()
}

func formatEntryTmuxSessionsAsRows(sessions []tmux.Session, fzfSep string) [][]string {
	rows := make([][]string, 0)

	for _, s := range sessions {
		sessionName := s.Name
		if s.Current {
			sessionName = fmt.Sprintf("[%s]", colorCurrentSession(sessionName))
		} else {
			sessionName = colorDefault(sessionName)
		}

		cols := make([]string, 0)
		cols = append(cols, colorCategorySession("session"))
		cols = append(cols, sessionName)
		cols = append(cols, colorPath(s.Path))
		cols = append(cols, colorMute(fzfSep, s.Name))
		cols = append(cols, colorMute(fzfSep+"session"+fzfSep+s.ID))

		rows = append(rows, cols)
	}

	return rows
}

func formatEntryDirectoryAsRows(dirs []string, fzfSep string) [][]string {
	rows := make([][]string, 0)

	for _, d := range dirs {
		base := filepath.Base(d)
		truncatedHome := strings.Replace(d, os.ExpandEnv("$HOME"), "~", 1)
		cols := make([]string, 0)

		cols = append(cols, colorCategoryDir("directory"))
		cols = append(cols, base)
		cols = append(cols, colorPath(truncatedHome))
		cols = append(cols, colorMute(fzfSep, truncatedHome))
		cols = append(cols, colorMute(fzfSep+"directory"+fzfSep+d))

		rows = append(rows, cols)
	}

	return rows
}
