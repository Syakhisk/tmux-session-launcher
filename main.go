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
			{Name: "launch"},
			{Name: "cmd"},
		},

		DefaultCommand: "launch",
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		logger.Fatal(err)
	}
}
