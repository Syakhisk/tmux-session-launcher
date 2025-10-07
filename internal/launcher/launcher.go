package launcher

import (
	"context"
	"tmux-session-launcher/internal/rpc"
	"tmux-session-launcher/internal/fuzzyfinder"
	"tmux-session-launcher/internal/server"

	"github.com/urfave/cli/v3"
)

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

	return fuzzyfinder.Launcher(ctx)
}

func HandlerLauncer(ctx context.Context, cmd *cli.Command) error {
	srv := server.NewServer(rpc.SockAddress)
	lcr := NewLauncher(srv)

	return lcr.Handler(ctx, cmd)
}
