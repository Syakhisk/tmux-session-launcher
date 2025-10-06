package server

import (
	"bufio"
	"context"
	"net"
	"os"
	"strings"
	"tmux-session-launcher/pkg/logger"
	"tmux-session-launcher/pkg/util"

	"emperror.dev/errors"
)

type ServerHandlerFunc func(ctx context.Context, payload string) (string, error)
type ServerHandlerRoute string

type Server struct {
	Address  string
	handlers map[ServerHandlerRoute]ServerHandlerFunc
	listener net.Listener
}

func NewServer(address string) *Server {
	return &Server{
		Address: address,
	}
}

func (s *Server) Start(ctx context.Context) error {
	logger := logger.WithPrefix("launcher.Start")

	logger.Infof("Starting socket listener at %s", s.Address)

	var lstCfg net.ListenConfig
	var err error
	s.listener, err = lstCfg.Listen(ctx, "unix", s.Address)
	if err != nil {
		return errors.Wrap(err, "failed to start socket listener")
	}

	logger.Info("Socket listener started, waiting for connections...")

	// Main loop to accept connections
	go func() {
		for {
			conn, err := s.listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					logger.Debugf("Stopping accept loop")
					break
				}

				logger.Errorf("Failed to accept connection: %v", err)
				continue
			}

			go s.handleConnection(ctx, conn)
		}
	}()

	return nil
}

func (s *Server) Stop() error {
	var combErr error
	if s.listener != nil {
		err := s.listener.Close()
		combErr = errors.Combine(err, errors.Wrap(err, "failed to close listener"))
	}

	if _, err := os.Stat(s.Address); err == nil {
		err := os.Remove(s.Address)
		combErr = errors.Combine(err, errors.Wrap(err, "failed to remove socket file"))
	}

	return combErr
}

func (s *Server) RegisterHandler(route ServerHandlerRoute, handler ServerHandlerFunc) {
	if s.handlers == nil {
		s.handlers = make(map[ServerHandlerRoute]ServerHandlerFunc)
	}

	s.handlers[route] = handler
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	log := logger.WithPrefix("launcher.handleConnection")
	log.Debugf("Handling connection from %s", conn.LocalAddr())

	scanner := bufio.NewScanner(conn)

	// Read route
	if !scanner.Scan() {
		log.Error("failed to read route")
		return
	}
	route := scanner.Text()

	// Read multi-line payload until empty line
	var payloadLines []string
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" { // empty line terminates payload
			break
		}
		payloadLines = append(payloadLines, line)
	}
	payload := strings.Join(payloadLines, "\n")

	handler, ok := s.handlers[ServerHandlerRoute(route)]
	if !ok {
		log.Errorf("route not found: %s", route)
		return
	}

	if util.IsContextDone(ctx) {
		log.Error("failed to handle connection due to:", ctx.Err())
		return
	}

	// Run handler and send response
	response, err := handler(ctx, payload)
	if err != nil {
		log.Errorf("handler error: %v", err)
		conn.Write([]byte("ERROR: " + err.Error() + "\n"))
		return
	}

	if response != "" {
		conn.Write([]byte(response + "\n"))
	}
}
