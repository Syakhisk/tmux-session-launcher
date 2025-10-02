package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
	"time"
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
				Name: "launch",
				Action: WithSignalHandling(func(ctx context.Context, c *cli.Command) error {
					logger.Info("Press Ctrl+C to exit")
					time.Sleep(1 * time.Second)
					return errors.New("not implemented")

					// // 1. start the socket
					// // 2. wait for input
					// // 3. reply if requested
					// l := launcher.New("/tmp/tmux-session-launcher.sock")
					// err := l.StartSocket(ctx)
					// return err
				}),
			},
			{
				Name:    "action",
				Aliases: []string{"act", "do"},
				Commands: []*cli.Command{
					{
						Name:    "next-mode",
						Aliases: []string{"next"},
						Action: func(ctx context.Context, c *cli.Command) error {
							// 1. get current mode from socket
							// 2. get next mode from config
							// 3. set next mode to socket
							return nil
						},
					},
					{Name: "previous-mode", Aliases: []string{"previous", "prev"}},
				},
			},
		},

		DefaultCommand: "launch",
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		logger.Fatal(err)
	}
}

// TODO: move to somewhere else
func WithSignalHandling(next cli.ActionFunc) cli.ActionFunc {
	return func(ctx context.Context, cmd *cli.Command) error {
		ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
		defer stop()
		errCh := make(chan error, 1)
		go func() { errCh <- next(ctx, cmd) }()

		select {
		case <-ctx.Done():
			logger.Info("shutting down gracefully")
			return nil
		case err := <-errCh:
			logger.Errorf("action error: %v", err)
			return err
		}
	}
}
