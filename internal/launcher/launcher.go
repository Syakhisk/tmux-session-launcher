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

	done := make(chan struct{})
	go func() {
		defer close(done)

		connCh := make(chan net.Conn)
		errCh := make(chan error)

		// Main loop to accept connections
		for {
			logger.Info("Socket listener started, waiting for connections...")

			// Accept connections in a separate goroutine to allow for context cancellation
			go func() {
				conn, err := listener.Accept()
				if err != nil {
					errCh <- err
					return
				}

				connCh <- conn
			}()

			// Wait for either a new connection, an error, or context cancellation
			select {
			case conn := <-connCh:
				logger.Infof("Accepted connection from %s", conn.RemoteAddr())
				return
			case err := <-errCh:
				logger.Errorf("Failed to accept connection: %v", err)
				return
			case <-ctx.Done():
				logger.Info("Context canceled, stopping listener")
				return
			}
		}
	}()

	select {
	case <-done:
		logger.Info("Shutting down socket listener")
	case <-ctx.Done():
		logger.Info("Context canceled, shutting down socket listener")
	}

	return nil
}
