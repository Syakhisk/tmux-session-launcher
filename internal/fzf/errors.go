package fzf

import (
	"os/exec"

	"emperror.dev/errors"
)

var (
	ErrUserCancelled = errors.New("user cancelled the operation")
)

func handleExitCodeErr(err error) error {
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		exitCode := exitErr.ExitCode()
		if exitCode == 130 {
			return ErrUserCancelled
		}

		return errors.Wrapf(err, "fzf failed with exit code: %d", exitCode)
	}

	return nil
}
