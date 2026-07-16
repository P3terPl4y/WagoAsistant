package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Logger wraps zerolog.Logger providing structured, leveled logging.
type Logger struct {
	zl zerolog.Logger
}

// New creates a new Logger. In development mode, output is human-readable.
// In production mode, output is JSON for log aggregation.
func New(environment string) Logger {
	var w io.Writer
	if environment == "production" {
		w = os.Stdout
	} else {
		w = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.DateTime,
		}
	}

	zl := zerolog.New(w).With().Timestamp().Logger()
	return Logger{zl: zl}
}

// With returns a new Logger with the given key-value context fields.
func (l Logger) With() zerolog.Context {
	return l.zl.With()
}

// WithComponent returns a child logger tagged with a component name.
func (l Logger) WithComponent(name string) Logger {
	return Logger{zl: l.zl.With().Str("component", name).Logger()}
}

// WithBotID returns a child logger tagged with a bot ID.
func (l Logger) WithBotID(botID int) Logger {
	return Logger{zl: l.zl.With().Int("bot_id", botID).Logger()}
}

// Info logs an info-level message.
func (l Logger) Info() *zerolog.Event {
	return l.zl.Info()
}

// Warn logs a warning-level message.
func (l Logger) Warn() *zerolog.Event {
	return l.zl.Warn()
}

// Error logs an error-level message.
func (l Logger) Error() *zerolog.Event {
	return l.zl.Error()
}

// Fatal logs a fatal-level message and exits.
func (l Logger) Fatal() *zerolog.Event {
	return l.zl.Fatal()
}

// Debug logs a debug-level message.
func (l Logger) Debug() *zerolog.Event {
	return l.zl.Debug()
}

// Zerolog returns the underlying zerolog.Logger for advanced use.
func (l Logger) Zerolog() zerolog.Logger {
	return l.zl
}
