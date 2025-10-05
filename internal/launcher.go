package internal

import (
	"context"
	"fmt"
	"tmux-session-launcher/internal/fuzzyfinder"
	"tmux-session-launcher/internal/mode"

	"github.com/urfave/cli/v3"
)

func HandlerLauncer(ctx context.Context, cmd *cli.Command) error {
	srv := NewServer(SockAddress)
	lcr := NewLauncher(srv)

	return lcr.Handler(ctx, cmd)
}

type Launcher struct {
	Server *Server
}

func NewLauncher(server *Server) *Launcher {
	return &Launcher{
		Server: server,
	}
}

func (l *Launcher) Handler(ctx context.Context, cmd *cli.Command) error {
	err := l.Server.Start(ctx)
	if err != nil {
		return err
	}

	l.registerHandlers()

	defer l.Server.Stop()

	if err := fuzzyfinder.Run(ctx); err != nil {
		return err
	}

	return nil
}

func (l *Launcher) registerHandlers() {
	l.Server.RegisterHandler(RouteNextMode, func(ctx context.Context, message string) (string, error) {
		m := mode.Next()
		return fmt.Sprintf("current mode: %s", m), nil
	})

	l.Server.RegisterHandler(RoutePrevMode, func(ctx context.Context, message string) (string, error) {
		m := mode.Prev()
		return fmt.Sprintf("current mode: %s", m), nil
	})
}
