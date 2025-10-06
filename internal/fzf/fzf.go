package fzf

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
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

func UpdateContentAndHeader(ctx context.Context, port int, header string) error {
	executable, err := os.Executable()
	if err != nil {
		return errors.WrapIf(err, "failed to get executable path")
	}

	bodyHeader := fmt.Sprintf("change-header(%s)", header)
	bodyContent := fmt.Sprintf("reload-sync(%s action content-get)", executable)
	bodyMove := fmt.Sprint("first")

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error { return sendRequest(gCtx, port, bodyHeader) })
	g.Go(func() error { return sendRequest(gCtx, port, bodyContent) })
	g.Go(func() error { return sendRequest(gCtx, port, bodyMove) })

	return g.Wait()
}

func sendRequest(ctx context.Context, port int, body string) error {
	log := logger.WithPrefix("fzf.sendRequest")
	log.Debugf("Sending request with body: %.200s", body)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("http://localhost:%d", port),
		strings.NewReader(body),
	)

	if err != nil {
		log.Errorf("Failed to create request: %v", err)
		return errors.WrapIf(err, "failed to create HTTP request")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Errorf("Failed to send request: %v", err)
		return errors.WrapIf(err, "failed to send HTTP request")
	}
	defer resp.Body.Close()

	// todo: remove
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.WrapIf(err, "failed to read response body")
	}

	log.Infof("Got response:\n%s", responseBody)

	log.Debugf("Response status: %s", resp.Status)
	return nil
}
