package log

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger is a simple wrapper for a slog Logger that writes to a file on disk.
type Logger struct {
	*slog.Logger
	logFile *lumberjack.Logger
	path    string
}

// Creates a new logger
func NewLogger(logFilePath string, debugMode bool, enableSourceLogging bool) (*Logger, error) {
	// Make the file
	err := os.MkdirAll(filepath.Dir(logFilePath), logDirMode)
	if err != nil {
		return nil, fmt.Errorf("error creating API log directory for [%s]: %w", logFilePath, err)
	}
	logFile := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    MaxLogSize,
		MaxBackups: MaxLogBackups,
		MaxAge:     MaxLogAge,
	}

	// Create the logging options
	logOptions := &slog.HandlerOptions{
		ReplaceAttr: ReplaceTime,
	}
	if debugMode {
		logOptions.Level = slog.LevelDebug
	} else {
		logOptions.Level = slog.LevelInfo
	}
	if enableSourceLogging {
		logOptions.AddSource = true
	}

	// Make the logger
	return &Logger{
		Logger:  slog.New(slog.NewTextHandler(logFile, logOptions)),
		logFile: logFile,
		path:    logFilePath,
	}, nil
}

// Get the path of the file this logger is writing to
func (l *Logger) GetFilePath() string {
	return l.path
}

// Rotate the log file, migrating the current file to an old backup and starting a new one
func (l *Logger) Rotate() error {
	return l.logFile.Rotate()
}

// Closes the log file
func (l *Logger) Close() {
	if l.logFile != nil {
		l.Info("Shutting down.")
		l.logFile.Close()
		l.logFile = nil
	}
}

// Create a clone of the logger that prints each message with the "origin" attribute.
// The underlying file handle isn't copied, so calling Close() on the sublogger won't do anything.
func (l *Logger) CreateSubLogger(origin string) *Logger {
	return &Logger{
		Logger:  l.With(slog.String(OriginKey, origin)),
		logFile: nil,
	}
}

// Creates a copy of the parent context with the logger put into the ContextLogKey value
func (l *Logger) CreateContextWithLogger(parent context.Context) context.Context {
	return context.WithValue(parent, ContextLogKey, l)
}

// Retrieves the logger from the context
func FromContext(ctx context.Context) (*Logger, bool) {
	log, ok := ctx.Value(ContextLogKey).(*Logger)
	return log, ok
}
