package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

// Config holds logger configuration
type Config struct {
	Level       string // debug, info, warn, error
	Pretty      bool   // Pretty print for development
	ServiceName string
}

// InitLogger initializes the global logger
func InitLogger(cfg Config) {
	var output io.Writer = os.Stdout

	// Pretty printing for development
	if cfg.Pretty {
		output = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
	}

	// Parse log level
	level := zerolog.InfoLevel
	switch cfg.Level {
	case "debug":
		level = zerolog.DebugLevel
	case "info":
		level = zerolog.InfoLevel
	case "warn":
		level = zerolog.WarnLevel
	case "error":
		level = zerolog.ErrorLevel
	}

	// Create logger
	Log = zerolog.New(output).
		Level(level).
		With().
		Timestamp().
		Str("service", cfg.ServiceName).
		Logger()
}

// WithRequestID creates a logger with request ID
func WithRequestID(requestID string) zerolog.Logger {
	return Log.With().Str("request_id", requestID).Logger()
}

// WithUserID creates a logger with user ID
func WithUserID(userID string) zerolog.Logger {
	return Log.With().Str("user_id", userID).Logger()
}
