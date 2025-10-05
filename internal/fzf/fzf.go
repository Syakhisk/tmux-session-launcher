package fzf

import (
	"context"
	"io"
	"os/exec"

	"emperror.dev/errors"
)

var (
	ErrUserCancelled = errors.New("user cancelled the operation")
)

func Select(ctx context.Context, args []string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	cmd := exec.CommandContext(ctx, "fzf", args...)
	cmd.Stdin = stdin
	cmd.Stdout = stdout
	cmd.Stderr = stderr

	if err := cmd.Start(); err != nil {
		return errors.WrapIf(err, "failed to start fzf command")
	}

	if err := cmd.Wait(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode := exitErr.ExitCode()
			if exitCode == 130 {
				return ErrUserCancelled
			}

			return errors.Wrapf(err, "fzf failed with exit code: %d", exitCode)
		}

		return errors.WrapIf(err, "fzf command failed")
	}

	return nil
}
