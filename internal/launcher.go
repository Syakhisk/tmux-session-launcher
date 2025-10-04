package internal

import (
	"context"

	"github.com/urfave/cli/v3"
)

type Mode string

const (
	ModeAll      Mode = "all"
	ModeSessions Mode = "sessions"
	ModeDir      Mode = "dir"
)

func LauncherHandler(ctx context.Context, cmd *cli.Command) error {
	srv := NewServer(SockAddress)
	lcr := NewLauncher(srv)

	return lcr.Handler(ctx, cmd)
}

type Launcher struct {
	Server      *Server
	currentMode Mode
}

func NewLauncher(server *Server) *Launcher {
	return &Launcher{
		Server: server,
	}
}

func (l *Launcher) Handler(ctx context.Context, cmd *cli.Command) error {
	l.Server.RegisterHandler("next", func(ctx context.Context, message string) (string, error) {
		return "next mode", nil
	})

	l.Server.RegisterHandler("prev", func(ctx context.Context, message string) (string, error) {
		return "prev mode", nil
	})

	err := l.Server.Start(ctx)
	if err != nil {
		return err
	}

	// server.RegisterHandler("/next-mode", l.NextModeHandler)

	// 1. start the socket
	// 2. wait for input
	// 3. reply if requested

	return nil
}
