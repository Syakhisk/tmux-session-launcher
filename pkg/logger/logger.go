package logger

import (
	"errors"
	"io"
	"log"
)

var (
	logger    *log.Logger
	verbosity int
)

const ()

const (
	LevelError int = iota
	LevelWarn
	LevelInfo
	LevelDebug
)

func init() {
	// ensure logger is initialized
	logger = log.New(nil, "", 0)
}

func SetupLogger(w io.Writer) {
	logger = log.New(w, "", 0)
}

func SetVerbosity(v int) error {
	verbosity = v

	if verbosity < LevelError {
		return errors.New("invalid verbosity level: must be >= 0")
	}

	if verbosity > LevelDebug {
		return errors.New("invalid verbosity level: must be <= 3")
	}

	return nil
}

func Errorf(format string, args ...any) {
	if verbosity < LevelError {
		return
	}
	logger.Printf("error: "+format, args...)
}

func Error(args ...any) {
	if verbosity < LevelError {
		return
	}
	logger.Print(append([]any{"error: "}, args...)...)
}

func Warnf(format string, args ...any) {
	if verbosity < LevelWarn {
		return
	}
	logger.Printf("warn: "+format, args...)
}

func Warn(args ...any) {
	if verbosity < LevelWarn {
		return
	}
	logger.Print(append([]any{"warn: "}, args...)...)
}

func Infof(format string, args ...any) {
	if verbosity < LevelInfo {
		return
	}
	logger.Printf("info: "+format, args...)
}

func Info(args ...any) {
	if verbosity < LevelInfo {
		return
	}
	logger.Print(append([]any{"info: "}, args...)...)
}

func Debugf(format string, args ...any) {
	if verbosity < LevelDebug {
		return
	}
	logger.Printf("debug: "+format, args...)
}

func Debug(args ...any) {
	if verbosity < LevelDebug {
		return
	}
	logger.Print(append([]any{"debug: "}, args...)...)
}

func Fatalf(format string, args ...any) {
	logger.Fatalf("fatal: "+format, args...)
}

func Fatal(args ...any) {
	logger.Fatal(append([]any{"fatal: "}, args...)...)
}
