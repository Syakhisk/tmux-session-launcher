package fuzzyfinder

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"tmux-session-launcher/internal/tmux"
	"tmux-session-launcher/pkg/logger"

	"github.com/fatih/color"

	"emperror.dev/errors"
)

const (
	separatorFzf      = "|"
	categorySession   = "session"
	categoryDirectory = "directory"
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

	logger.Debugf("fzf output: %s", string(output))

	category, id, err := parseSelectedOutput(string(output), separatorFzf)
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
