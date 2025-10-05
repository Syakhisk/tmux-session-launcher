package fuzzyfinder

import (
	"context"
	"fmt"
	"strings"
	"tmux-session-launcher/internal/dirs"
	"tmux-session-launcher/internal/mode"
	"tmux-session-launcher/internal/tmux"
	"tmux-session-launcher/pkg/logger"
	"unicode/utf8"

	"github.com/acarl005/stripansi"
	"github.com/fatih/color"
	"github.com/rodaine/table"
)

func buildHeader() string {
	c := color.New(color.Faint, color.Bold, color.Italic)

	header := color.New(color.Faint).Sprintf("Press %s/%s to switch mode\n", keyModeNext, keyModePrev)

	currentMode := mode.Get()

	mSlc := make([]string, 0, len(mode.Modes))
	for _, m := range mode.Modes {
		if m == currentMode {
			mSlc = append(mSlc, colorCurrentSession(fmt.Sprintf("[%s]", m)))
			continue
		}

		mSlc = append(mSlc, c.Sprint(m))
	}

	header += strings.Join(mSlc, " ")

	return header
}

func buildContent(ctx context.Context) (string, error) {
	currentMode := mode.Get()

	output := make([][]string, 0)

	if currentMode == mode.ModeSession || currentMode == mode.ModeAll {
		sessions, err := tmux.GetSessions(ctx)
		if err != nil {
			logger.Warnf("Failed to get tmux sessions: %v", err)
		}

		formattedSessions := formatEntryTmuxSessionsAsRows(sessions, fzfSeparator)
		output = append(output, formattedSessions...)
	}

	if currentMode == mode.ModeDirectory || currentMode == mode.ModeAll {
		dirs := dirs.GetDirectories()

		formattedDirs := formatEntryDirectoryAsRows(dirs, fzfSeparator)
		output = append(output, formattedDirs...)
	}

	return formatTable(output), nil
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
		cols = append(cols, colorCategorySession(categorySession))
		cols = append(cols, sessionName)
		cols = append(cols, colorPath(s.Path))
		cols = append(cols, colorMute(fzfSep, s.Name))
		cols = append(cols, colorMute(fzfSep+categorySession+fzfSep+s.ID))

		rows = append(rows, cols)
	}

	return rows
}

func formatEntryDirectoryAsRows(dirs []dirs.Directory, fzfSep string) [][]string {
	rows := make([][]string, 0)

	for _, d := range dirs {
		cols := make([]string, 0)

		cols = append(cols, colorCategoryDir(categoryDirectory))
		cols = append(cols, d.Label)
		cols = append(cols, colorPath(d.TruncatedHomePath))
		cols = append(cols, colorMute(fzfSep, d.Label))
		cols = append(cols, colorMute(fzfSep+categoryDirectory+fzfSep+d.FullPath))

		rows = append(rows, cols)
	}

	return rows
}

func parseSelectedOutput(output string, fzfSep string) (string, string, error) {
	parts := strings.Split(strings.TrimSpace(output), fzfSep)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("unexpected output format: %s", output)
	}

	itemType := parts[0]
	itemID := parts[1]

	return itemType, itemID, nil
}
