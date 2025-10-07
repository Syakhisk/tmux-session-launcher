package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tmux-session-launcher/internal/action"
	"tmux-session-launcher/internal/config"
	"tmux-session-launcher/internal/launcher"
	"tmux-session-launcher/pkg/logger"

	"github.com/urfave/cli/v3"
)

func main() {
	logger.SetupLogger(os.Stderr)
	logger.SetVerbosity(1)

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
				Action: WithSignalHandling(launcher.HandlerLauncer),
			},
			{
				Name: "action",
				Commands: []*cli.Command{
					{
						Name:    "mode-next",
						Aliases: []string{"next"},
						Action:  WithSignalHandling(action.HandlerNextMode),
					},
					{
						Name:    "mode-previous",
						Aliases: []string{"previous", "prev"},
						Action:  WithSignalHandling(action.HandlerPrevMode),
					},
					{
						Name:   "mode-get",
						Action: WithSignalHandling(action.HandlerGetMode),
					},
					{
						Name:   "content-get",
						Action: WithSignalHandling(action.HandlerGetContent),
					},
					{
						Name:   "open-in",
						Action: WithSignalHandling(action.HandlerOpenIn),
					},
				},
			},
			{
				Name:    "config",
				Aliases: []string{"cfg", "c"},
				Commands: []*cli.Command{
					{
						Name:   "path",
						Usage:  "Show the path to the configuration file",
						Action: config.HandlerShowConfigPath,
					},
					{
						Name:   "edit",
						Usage:  "Edit the configuration file",
						Action: config.HandlerEditConfig,
					},
					{
						Name:   "init",
						Usage:  "Initialize a default configuration file",
						Action: config.HandlerInitConfig,
					},
					{
						Name:   "validate",
						Usage:  "Validate the configuration file",
						Action: config.HandlerValidateConfig,
					},
					{
						Name:   "list",
						Usage:  "List all configured directories",
						Action: config.HandlerListDirectories,
					},
					{
						Name:    "add",
						Aliases: []string{"a"},
						Usage:   "Add a directory to the configuration (current dir if no path specified)",
						Action:  config.HandlerAddDirectory,
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name:    "depth",
								Aliases: []string{"d"},
								Value:   0,
								Usage:   "Depth of subdirectories to include (0 = no subdirectories)",
							},
						},
					},
					{
						Name:    "remove",
						Aliases: []string{"rm"},
						Usage:   "Remove a directory from the configuration",
						Action:  config.HandlerRemoveDirectory,
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
