package server

import (
	"context"
	"net"
	"os"
	"tmux-session-launcher/pkg/logger"

	"emperror.dev/errors"
	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
	"github.com/creachadair/jrpc2/handler"
)

type Server struct {
	Address  string
	assigner handler.Map
	listener net.Listener
}

func NewServer(address string) *Server {
	return &Server{
		Address:  address,
		assigner: handler.Map{},
	}
}

func (s *Server) Start(ctx context.Context) error {
	logger := logger.WithPrefix("server.Start")

	logger.Infof("Starting JSON-RPC server at %s", s.Address)

	// Remove existing socket file if it exists
	if _, err := os.Stat(s.Address); err == nil {
		logger.Warnf("Socket file %s already exists, removing it", s.Address)
		if err := os.Remove(s.Address); err != nil {
			return errors.Wrap(err, "failed to remove existing socket file")
		}
	}

	listener, err := net.Listen("unix", s.Address)
	if err != nil {
		return errors.Wrap(err, "failed to start socket listener")
	}
	s.listener = listener

	logger.Info("JSON-RPC server started, waiting for connections...")

	// Accept connections in a goroutine
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					logger.Debugf("Stopping accept loop")
					break
				}
				logger.Errorf("Failed to accept connection: %v", err)
				continue
			}

			// Handle each connection in a separate goroutine
			go func(conn net.Conn) {
				defer conn.Close()

				// Create channel for the connection
				ch := channel.Line(conn, conn)

				// Create a new server instance for this connection
				srv := jrpc2.NewServer(s.assigner, &jrpc2.ServerOptions{
					Logger: func(text string) {
						logger.Debugf("jrpc2: %s", text)
					},
				})

				// Start serving and wait for completion
				srv.Start(ch).Wait()
			}(conn)
		}
	}()

	return nil
}

func (s *Server) Stop() error {
	var combErr error

	if s.listener != nil {
		err := s.listener.Close()
		combErr = errors.Combine(combErr, errors.Wrap(err, "failed to close listener"))
	}

	if _, err := os.Stat(s.Address); err == nil {
		err := os.Remove(s.Address)
		combErr = errors.Combine(combErr, errors.Wrap(err, "failed to remove socket file"))
	}

	return combErr
}

func (s *Server) RegisterHandler(method string, handler jrpc2.Handler) {
	s.assigner[method] = handler
}
