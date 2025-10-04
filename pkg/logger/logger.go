package logger

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/rs/zerolog"
)

var (
	logger    zerolog.Logger
	verbosity int
)

type Logger struct {
	logger zerolog.Logger
	prefix string
}

const (
	LevelError int = iota
	LevelWarn
	LevelInfo
	LevelDebug
)

func init() {
	// Setup pretty console logging by default
	output := zerolog.ConsoleWriter{Out: os.Stderr}
	logger = zerolog.New(output).With().Timestamp().Logger()
	verbosity = LevelInfo // Default to info level
}

func SetupLogger(w io.Writer) {
	output := zerolog.ConsoleWriter{Out: w}
	logger = zerolog.New(output).With().Timestamp().Logger()
}

func SetVerbosity(v int) error {
	verbosity = v

	if verbosity < LevelError {
		return errors.New("invalid verbosity level: must be >= 0")
	}

	if verbosity > LevelDebug {
		return errors.New("invalid verbosity level: must be <= 3")
	}

	// Set zerolog level based on verbosity
	switch verbosity {
	case LevelError:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case LevelWarn:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case LevelInfo:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case LevelDebug:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	return nil
}

// Helper function to format arguments into a string
func formatArgs(args ...any) string {
	if len(args) == 0 {
		return ""
	}

	var result string
	for i, arg := range args {
		if i > 0 {
			result += " "
		}
		switch v := arg.(type) {
		case string:
			result += v
		default:
			result += fmt.Sprintf("%v", v)
		}
	}
	return result
}

// Helper method to prepend prefix to message
func (l *Logger) prefixMsg(msg string) string {
	if l.prefix != "" {
		return fmt.Sprintf("[%s] %s", l.prefix, msg)
	}
	return msg
}

func WithPrefix(p string) *Logger {
	return &Logger{
		logger: logger,
		prefix: p,
	}
}

func (l *Logger) Errorf(format string, args ...any) {
	Errorf(l.prefixMsg(format), args...)
}

func (l *Logger) Error(args ...any) {
	Error(l.prefixMsg(formatArgs(args...)))
}

func (l *Logger) Warnf(format string, args ...any) {
	Warnf(l.prefixMsg(format), args...)
}

func (l *Logger) Warn(args ...any) {
	Warn(l.prefixMsg(formatArgs(args...)))
}

func (l *Logger) Infof(format string, args ...any) {
	Infof(l.prefixMsg(format), args...)
}

func (l *Logger) Info(args ...any) {
	Info(l.prefixMsg(formatArgs(args...)))
}

func (l *Logger) Debugf(format string, args ...any) {
	Debugf(l.prefixMsg(format), args...)
}

func (l *Logger) Debug(args ...any) {
	Debug(l.prefixMsg(formatArgs(args...)))
}

func (l *Logger) Fatalf(format string, args ...any) {
	Fatalf(l.prefixMsg(format), args...)
}

func (l *Logger) Fatal(args ...any) {
	Fatal(l.prefixMsg(formatArgs(args...)))
}

func Errorf(format string, args ...any) {
	if verbosity < LevelError {
		return
	}
	logger.Error().Msgf(format, args...)
}

func Error(args ...any) {
	if verbosity < LevelError {
		return
	}
	logger.Error().Msg(formatArgs(args...))
}

func Warnf(format string, args ...any) {
	if verbosity < LevelWarn {
		return
	}
	logger.Warn().Msgf(format, args...)
}

func Warn(args ...any) {
	if verbosity < LevelWarn {
		return
	}
	logger.Warn().Msg(formatArgs(args...))
}

func Infof(format string, args ...any) {
	if verbosity < LevelInfo {
		return
	}
	logger.Info().Msgf(format, args...)
}

func Info(args ...any) {
	if verbosity < LevelInfo {
		return
	}
	logger.Info().Msg(formatArgs(args...))
}

func Debugf(format string, args ...any) {
	if verbosity < LevelDebug {
		return
	}
	logger.Debug().Msgf(format, args...)
}

func Debug(args ...any) {
	if verbosity < LevelDebug {
		return
	}
	logger.Debug().Msg(formatArgs(args...))
}

func Fatalf(format string, args ...any) {
	logger.Fatal().Msgf(format, args...)
}

func Fatal(args ...any) {
	logger.Fatal().Msg(formatArgs(args...))
}
