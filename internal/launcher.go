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

type Launcher struct {
	Server      *Server
	currentMode Mode
}

func ActionLauncherHandler(ctx context.Context, cmd *cli.Command) error {
	srv := NewServer("/tmp/launcher.sock")
	lcr := NewLauncher(srv)

	return lcr.Handler(ctx, cmd)
}

func NewLauncher(server *Server) *Launcher {
	return &Launcher{
		Server: server,
	}
}

func (l *Launcher) Handler(ctx context.Context, cmd *cli.Command) error {
	// TODO: handle connection here?
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

// func (l *Launcher) SendRequest(ctx context.Context, message string) (string, error) {
// 	logger := logger.WithPrefix("launcher.SendRequest")
//
// 	logger.Infof("[CLIENT] Connecting to socket at %s", l.sockPath)
// 	conn, err := net.Dial("unix", l.sockPath)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer conn.Close()
//
// 	logger.Infof("[CLIENT] Sending message: %s", message)
// 	conn.Write([]byte(message))
// 	logger.Infof("[CLIENT] Message sent, waiting for response...")
//
// 	msg := make([]byte, 512)
// 	n, err := conn.Read(msg)
// 	if err != nil {
// 		logger.Errorf("[CLIENT] Failed to read response: %v", err)
// 		return "", err
// 	}
//
// 	logger.Infof("[CLIENT] Received response: %s", string(msg[:n]))
//
// 	return string(msg[:n]), nil
// }

//	func (l *Launcher) StartSocketServer(ctx context.Context) error {
//		logger := logger.WithPrefix("launcher.StartSocketServer")
//
//		logger.Infof("Starting socket listener at %s", l.sockPath)
//
//		var lstCfg net.ListenConfig
//		listener, err := lstCfg.Listen(ctx, "unix", l.sockPath)
//		if err != nil {
//			return errors.Wrap(err, "failed to start socket listener")
//		}
//
//		// regular cleanup when function exits
//		defer func() {
//			listener.Close()
//			os.Remove(l.sockPath)
//		}()
//
//		// cleanup when context is cancelled
//		go func() {
//			<-ctx.Done()
//			logger.Debugf("Cleanup due to context cancellation")
//			listener.Close()
//			os.Remove(l.sockPath)
//		}()
//
//		logger.Info("Socket listener started, waiting for connections...")
//
//		// Main loop to accept connections
//		for {
//			conn, err := listener.Accept()
//			if err != nil {
//				if util.IsContextDone(ctx) {
//					logger.Debugf("Context done, exiting socket listener, error: %v", err)
//					return nil
//				}
//
//				logger.Errorf("Failed to accept connection: %v", err)
//				continue
//			}
//
//			go l.handleConnection(ctx, conn)
//		}
//	}

// func (l *Launcher) handleConnection(_ context.Context, conn net.Conn) {
// 	logger := logger.WithPrefix("launcher.handleConnection")
// 	logger.Infof("Handling connection from %s", conn.LocalAddr())
//
// 	msg := make([]byte, 512)
// 	n, err := conn.Read(msg)
// 	if err != nil {
// 		logger.Errorf("Failed to read from connection: %v", err)
// 		return
// 	}
//
// 	time.Sleep(1 * time.Second) // simulate some processing time
//
// 	fmt.Fprintf(conn, "Echo: %s", string(msg[:n]))
//
// 	logger.Infof("Response sent to %s", conn.LocalAddr())
//
// 	defer conn.Close()
// }
