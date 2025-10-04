package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
	i "tmux-session-launcher/internal"
	"tmux-session-launcher/pkg/logger"

	"github.com/urfave/cli/v3"
)

func main() {
	logger.SetupLogger(os.Stderr)
	logger.SetVerbosity(3)

	cmd := &cli.Command{
		Name: "tmux-session-launcher",
		Flags: []cli.Flag{
			&cli.Uint8Flag{
				Name:     "vebosity",
				Aliases:  []string{"v"},
				Sources:  cli.EnvVars("VERBOSITY_LEVEL"),
				OnlyOnce: true,
				Action: func(_ context.Context, _ *cli.Command, value uint8) error {
					return logger.SetVerbosity(int(value))
				},
			},
		},

		Commands: []*cli.Command{
			{
				Name:   "launch",
				Action: WithSignalHandling(i.HandlerLauncer),
			},
			{
				Name: "action",
				Commands: []*cli.Command{
					{
						Name:    "next-mode",
						Aliases: []string{"next"},
						Action:  WithSignalHandling(i.HandlerActionNextMode),
					},
					{
						Name:    "previous-mode",
						Aliases: []string{"previous", "prev"},
						Action:  WithSignalHandling(i.HandlerActionPrevMode),
					},
				},
			},
		},

		DefaultCommand: "launch",
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		logger.Fatal(err)
	}
}

// WithSignalHandling wraps a CLI action with graceful shutdown signal handling.
func WithSignalHandling(next cli.ActionFunc) cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
		defer stop()
		errCh := make(chan error, 1)
		go func() { errCh <- next(ctx, cmd) }()

		select {
		case <-ctx.Done():
			logger.Info("received shutdown signal, cleaning up...")
			// Give the action some time to clean up
			select {
			case err := <-errCh:
				return err
			case <-time.After(5 * time.Second):
				logger.Warn("cleanup timeout exceeded")
				return ctx.Err()
			}
		case err := <-errCh:
			if err != nil {
				logger.Errorf("action error: %v", err)
			}
			return err
		}
	}
}
