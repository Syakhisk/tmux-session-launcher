package fzf

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
	"tmux-session-launcher/pkg/logger"

	"emperror.dev/errors"
	"golang.org/x/sync/errgroup"
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
		if exitErr := handleExitCodeErr(err); exitErr != nil {
			return exitErr
		}

		return errors.WrapIf(err, "fzf command failed")
	}

	return nil
}

func UpdateContentAndHeader(ctx context.Context, port int, header string, content string) error {
	logger := logger.WithPrefix("fzf.UpdateContentAndHeader")
	logger.Infof("Content length: %d characters", len(content))

	// HACK: Use a temporary file to avoid issues with large content, but this is not ideal since it needs a file
	//  - maybe create a subcommand to be called by fzf
	// Write content to temporary file
	tmpFile, err := os.CreateTemp("", "fzf-content-*.txt")
	if err != nil {
		return errors.WrapIf(err, "failed to create temp file")
	}
	defer tmpFile.Close()

	if _, err := tmpFile.WriteString(content); err != nil {
		return errors.WrapIf(err, "failed to write to temp file")
	}

	bodyHeader := fmt.Sprintf("change-header(%s)", header)
	bodyContent := fmt.Sprintf("reload-sync(cat %s)", tmpFile.Name())
	bodyMove := fmt.Sprint("first")

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error { return sendRequest(gCtx, port, bodyHeader) })
	g.Go(func() error { return sendRequest(gCtx, port, bodyContent) })
	g.Go(func() error { return sendRequest(gCtx, port, bodyMove) })

	err = g.Wait()

	// Clean up temp file after requests complete
	// TODO: There's still a race condition here - fzf might not have read the file yet
	//  Consider keeping temp files and cleaning them up periodically
	go func() {
		time.Sleep(100 * time.Millisecond) // Give fzf time to read
		os.Remove(tmpFile.Name())
	}()

	return err
}

func sendRequest(ctx context.Context, port int, body string) error {
	logger := logger.WithPrefix("fzf.sendRequest")
	logger.Debugf("Sending request with body: %.200s", body)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("http://localhost:%d", port),
		strings.NewReader(body),
	)

	if err != nil {
		logger.Errorf("Failed to create request: %v", err)
		return errors.WrapIf(err, "failed to create HTTP request")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Errorf("Failed to send request: %v", err)
		return errors.WrapIf(err, "failed to send HTTP request")
	}
	defer resp.Body.Close()

	// todo: remove
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.WrapIf(err, "failed to read response body")
	}

	logger.Infof("Got response:\n%s", responseBody)

	logger.Debugf("Response status: %s", resp.Status)
	return nil
}
