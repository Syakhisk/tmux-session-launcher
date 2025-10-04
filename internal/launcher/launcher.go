package launcher

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"
	"tmux-session-launcher/pkg/logger"

	"github.com/urfave/cli/v3"
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
func (l *Launcher) Handler(ctx context.Context, cmd *cli.Command) error {
	// 1. start the socket
	// 2. wait for input
	// 3. reply if requested
	return l.StartSocketServer(ctx)
}

func (l *Launcher) StartSocketServer(ctx context.Context) error {
	logger.Infof("Starting socket listener at %s", l.sockPath)
	var lst net.ListenConfig
	listener, err := lst.Listen(ctx, "unix", l.sockPath)
	if err != nil {
		return err
	}

	defer func() {
		logger.Info("Cleaning up socket")

		listener.Close()
		if err := os.Remove(l.sockPath); err != nil {
			return
		}

		logger.Info("Socket file removed")
	}()

	logger.Info("Socket listener started, waiting for connections...")

	// Handle context cancellation in a separate goroutine
	go func() {
		<-ctx.Done()
		listener.Close()
	}()

	// Main loop to accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			// Check if this is due to context cancellation
			select {
			case <-ctx.Done():
				logger.Info("Listener closed due to context cancellation")
				return ctx.Err()
			default:
				logger.Errorf("Failed to accept connection: %v", err)
				continue
			}
		}

		logger.Infof("Accepted connection from %s", conn.LocalAddr())

		// Handle each connection in its own goroutine
		go l.handleConnection(ctx, conn)
	}
}

func (l *Launcher) SendRequest(ctx context.Context, message string) (string, error) {
	logger.Infof("[CLIENT] Connecting to socket at %s", l.sockPath)
	conn, err := net.Dial("unix", l.sockPath)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	logger.Infof("[CLIENT] Sending message: %s", message)
	conn.Write([]byte(message))
	logger.Infof("[CLIENT] Message sent, waiting for response...")

	msg := make([]byte, 512)
	n, err := conn.Read(msg)
	if err != nil {
		logger.Errorf("[CLIENT] Failed to read response: %v", err)
		return "", err
	}

	logger.Infof("[CLIENT] Received response: %s", string(msg[:n]))

	return string(msg[:n]), nil
}

func (l *Launcher) handleConnection(_ context.Context, conn net.Conn) {
	logger.Infof("Handling connection from %s", conn.LocalAddr())

	msg := make([]byte, 512)
	n, err := conn.Read(msg)
	if err != nil {
		logger.Errorf("Failed to read from connection: %v", err)
		return
	}

	time.Sleep(1 * time.Second) // simulate some processing time

	fmt.Fprintf(conn, "Echo: %s", string(msg[:n]))

	logger.Infof("Response sent to %s", conn.LocalAddr())

	defer conn.Close()
}
