package internal

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
}

func NewServer(address string) *Server {
	return &Server{
		Address: address,
	}
}

func (s *Server) Start(ctx context.Context) error {
	logger := logger.WithPrefix("launcher.StartSocketServer")

	if len(s.handlers) > 0 {
		logger.Infof("%d handlers registered", len(s.handlers))
	} else {
		logger.Warn("no handlers registered, server will not respond to any requests")
	}

	logger.Infof("Starting socket listener at %s", s.Address)

	var lstCfg net.ListenConfig
	listener, err := lstCfg.Listen(ctx, "unix", s.Address)
	if err != nil {
		return errors.Wrap(err, "failed to start socket listener")
	}

	logger.Info("Socket listener started, waiting for connections...")

	// cleanup when context is cancelled
	go func() {
		<-ctx.Done()
		logger.Debugf("Cleanup: context cancelled")
		listener.Close()
		os.Remove(s.Address)
	}()

	// Main loop to accept connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			if util.IsContextDone(ctx) {
				logger.Debugf("Context closed, exiting main loop")
				break
			}

			logger.Errorf("Failed to accept connection: %v", err)
			continue
		}

		go s.handleConnection(ctx, conn)
	}

	return nil
}

func (s *Server) RegisterHandler(route ServerHandlerRoute, handler ServerHandlerFunc) {
	if s.handlers == nil {
		s.handlers = make(map[ServerHandlerRoute]ServerHandlerFunc)
	}

	s.handlers[route] = handler
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	logger := logger.WithPrefix("launcher.handleConnection")
	logger.Debugf("Handling connection from %s", conn.LocalAddr())

	scanner := bufio.NewScanner(conn)

	// Read route
	if !scanner.Scan() {
		logger.Error("failed to read route")
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
		logger.Errorf("route not found: %s", route)
		return
	}

	if util.IsContextDone(ctx) {
		logger.Error("failed to handle connection due to:", ctx.Err())
		return
	}

	// Run handler and send response
	response, err := handler(ctx, payload)
	if err != nil {
		logger.Errorf("handler error: %v", err)
		conn.Write([]byte("ERROR: " + err.Error() + "\n"))
		return
	}

	if response != "" {
		conn.Write([]byte(response + "\n"))
	}
}
