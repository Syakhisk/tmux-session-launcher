package internal

import (
	"context"
	"net"
	"os"
	"tmux-session-launcher/pkg/logger"
	"tmux-session-launcher/pkg/util"

	"emperror.dev/errors"
)

type ServerHandlerFunc func(ctx context.Context, message string) (string, error)
type ServerHandlerRoute string

type Server struct {
	Address  string
	Handlers map[*ServerHandlerRoute]ServerHandlerFunc
}

func NewServer(address string) *Server {
	return &Server{
		Address: address,
	}
}

func (s *Server) Start(ctx context.Context) error {
	logger := logger.WithPrefix("launcher.StartSocketServer")

	logger.Infof("Starting socket listener at %s", s.Address)

	var lstCfg net.ListenConfig
	listener, err := lstCfg.Listen(ctx, "unix", s.Address)
	if err != nil {
		return errors.Wrap(err, "failed to start socket listener")
	}

	// regular cleanup when function exits
	defer func() {
		listener.Close()
		os.Remove(s.Address)
	}()

	// cleanup when context is cancelled
	go func() {
		<-ctx.Done()
		logger.Debugf("Cleanup due to context cancellation")
		listener.Close()
		os.Remove(s.Address)
	}()

	logger.Info("Socket listener started, waiting for connections...")

	// Main loop to accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			if util.IsContextDone(ctx) {
				logger.Debugf("Context done, exiting socket listener, error: %v", err)
				return nil
			}

			logger.Errorf("Failed to accept connection: %v", err)
			continue
		}

		// TODO: how to handle connection? router-like
		go s.handleConnection(ctx, conn)
	}
}

func (s *Server) RegisterHandler(route ServerHandlerRoute, handler ServerHandlerFunc) {
	if s.Handlers == nil {
		s.Handlers = make(map[*ServerHandlerRoute]ServerHandlerFunc)
	}

	s.Handlers[&route] = handler
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	logger := logger.WithPrefix("launcher.handleConnection")
	logger.Debugf("Handling connection from %s", conn.LocalAddr())

	msg := make([]byte, 1024)
	n, err := conn.Read(msg)
	if err != nil {
		logger.Errorf("Failed to read from connection: %v", err)
		return
	}


	// handlerFunc(ctx, string(msg[:n]))

	// fmt.Fprintf(conn, "Echo: %s", string(msg[:n]))
	// logger.Infof("Response sent to %s", conn.LocalAddr())
}
