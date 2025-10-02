package launcher

import (
	"context"
	"net"
	"os"
	"tmux-session-launcher/pkg/logger"
)

type Mode string

const (
	ModeAll      Mode = "all"
	ModeSessions Mode = "sessions"
	ModeDir      Mode = "dir"
)

// TODO: create urfave/cli actions to be called
type Launcher struct {
	sockPath    string
	currentMode Mode
}

func New(sockPath string) *Launcher {
	return &Launcher{
		sockPath: sockPath,
	}
}

func (l *Launcher) StartSocket(ctx context.Context) error {
	logger.Infof("Starting socket listener at %s", l.sockPath)

	// create the socket listener
	_ = os.Remove(l.sockPath)
	var lst net.ListenConfig
	listener, err := lst.Listen(ctx, "unix", l.sockPath)
	if err != nil {
		return err
	}

	defer func() {
		listener.Close()
		os.Remove(l.sockPath)
	}()

	logger.Info("Socket listener started, waiting for connections...")

	// Main loop to accept connections
	for {
		select {
		case <-ctx.Done():
			logger.Info("Context canceled, stopping main loop")
			return ctx.Err()
		default:
		}

		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				logger.Info("Context canceled, stopping accept loop")
				return ctx.Err()
			default:
				logger.Errorf("Failed to accept connection: %v", err)
				continue
			}
		}

		logger.Infof("Accepted connection from %s", conn.RemoteAddr())

		// // Handle each connection in its own goroutine
		// go l.handleConnection(ctx, conn)
	}
}
