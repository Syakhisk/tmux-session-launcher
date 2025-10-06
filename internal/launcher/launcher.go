package launcher

import (
	"context"
	"tmux-session-launcher/internal/constants"
	"tmux-session-launcher/internal/fuzzyfinder"
	"tmux-session-launcher/internal/server"

	"github.com/urfave/cli/v3"
)

func HandlerLauncer(ctx context.Context, cmd *cli.Command) error {
	srv := server.NewServer(constants.SockAddress)
	lcr := NewLauncher(srv)

	return lcr.Handler(ctx, cmd)
}

type Launcher struct {
	Server *server.Server
}

func NewLauncher(server *server.Server) *Launcher {
	return &Launcher{
		Server: server,
	}
}

func (l *Launcher) Handler(ctx context.Context, cmd *cli.Command) error {
	err := l.Server.Start(ctx)
	if err != nil {
		return err
	}

	l.setupHandlers()

	defer l.Server.Stop()

	if err := fuzzyfinder.Run(ctx); err != nil {
		return err
	}

	return nil
}
