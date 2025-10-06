package fuzzyfinder

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"tmux-session-launcher/internal/fzf"
	"tmux-session-launcher/internal/tmux"
	"tmux-session-launcher/pkg/logger"

	"github.com/fatih/color"

	"emperror.dev/errors"
)

const (
	fzfPort      = 6266
	fzfSeparator = "|"

	categorySession   = "session"
	categoryDirectory = "directory"

	keyModeNext = "ctrl-j"
	keyModePrev = "ctrl-k"
)

var (
	colorDefault         = fmt.Sprint
	colorCategorySession = color.New(color.FgHiCyan, color.Italic).Sprint
	colorCategoryDir     = color.New(color.FgHiBlue, color.Italic).Sprint
	colorCurrentSession  = color.New(color.FgHiGreen, color.Bold).Sprint
	colorPath            = color.New(color.FgHiBlack, color.Italic).Sprint
	colorMute            = color.RGB(0, 0, 0).Sprint
)

func Run(ctx context.Context) error {
	logger := logger.WithPrefix("fuzzyfinder.Exec")

	execPath, err := os.Executable()
	if err != nil {
		return errors.WrapIf(err, "failed to get executable path")
	}

	args := []string{
		"--ansi",
		"--no-sort",
		"--no-hscroll",
		"--listen", fmt.Sprint(fzfPort),
		"--header", buildHeader(),
		"--delimiter", fzfSeparator, // used as nth delimiter
		"--with-nth", "1,2", // what to show in the list
		"--nth", "2", // what to search in (based on with-nth)
		"--accept-nth", "3,4", // what to output on accept
		fmt.Sprintf("--bind=%s:execute-silent(%s action mode-next)", keyModeNext, execPath),
		fmt.Sprintf("--bind=%s:execute-silent(%s action mode-previous)", keyModePrev, execPath),
	}

	input, err := buildContent(ctx)
	if err != nil {
		return errors.WrapIf(err, "failed to build fzf input")
	}

	outputBuf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}

	if err := fzf.Select(
		ctx,
		args,
		bytes.NewReader([]byte(input)),
		outputBuf,
		errBuf,
	); err != nil {
		if errors.Is(err, fzf.ErrUserCancelled) {
			return nil
		}

		return errors.WrapIff(err, "fzf selection failed: %+v", errBuf.String())
	}

	output := outputBuf.String()
	logger.Debugf("fzf output: %s", output)

	category, id, err := parseSelectedOutput(output, fzfSeparator)
	if err != nil {
		return errors.WrapIf(err, "failed to parse fzf output")
	}

	var errTmux error
	switch category {
	case categorySession:
		logger.Infof("Attaching to tmux session with ID: %s", id)
		errTmux = tmux.SessionAttach(ctx, id)

	case categoryDirectory:
		logger.Infof("Opening directory with path: %s", id)
		_, errTmux = tmux.SessionCreateOrAttach(ctx, tmux.BuildSessionNameFromPath(id), id)
	}

	if errTmux != nil {
		return errors.WrapIf(errTmux, "failed to attach or create tmux session")
	}

	return nil
}

func UpdateContentAndHeader(ctx context.Context) error {
	header := buildHeader()

	if err := fzf.UpdateContentAndHeader(ctx, fzfPort, header); err != nil {
		return errors.WrapIf(err, "failed to update fzf content and header")
	}

	return nil
}

func GetContent(ctx context.Context) (string, error) {
	content, err := buildContent(ctx)
	if err != nil {
		return "", errors.WrapIf(err, "failed to build fzf input")
	}

	return content, nil
}
