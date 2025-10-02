package main

import (
	"context"
	"os"
	"tmux-session-launcher/pkg/logger"

	"github.com/urfave/cli/v3"
)

func main() {
	logger.SetupLogger(os.Stderr)

	cmd := &cli.Command{
		Name: "tmux-session-launcher",
		Flags: []cli.Flag{
			&cli.Uint8Flag{
				Name:     "vebosity",
				Aliases:  []string{"v"},
				Sources:  cli.EnvVars("VERBOSITY_LEVEL"),
				Value:    1,
				OnlyOnce: true,
				Action: func(_ context.Context, _ *cli.Command, value uint8) error {
					return logger.SetVerbosity(int(value))
				},
			},
		},

		Commands: []*cli.Command{
			{
				Name: "launch",
				Action: func(ctx context.Context, c *cli.Command) error {
					// 1. start the socket
					// 2. wait for input
					// 3. reply if requested
					return nil
				},
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
