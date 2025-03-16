package logger

import (
	"log/slog"
	"os"
	"sync"
)

var (
	logger *slog.Logger
	once   sync.Once
)

// Get return global logger with JSON format.
func Get() *slog.Logger {
	once.Do(func() {
		// If you want enable Debug for teh stdout replace nil by:
		// &slog.HandlerOptions{
		// 	Level: slog.LevelDebug,
		// }
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	})
	return logger
}

// Debug logs a debug message.
func Debug(msg string, args ...any) { Get().Debug(msg, args...) }

// Info logs an info message.
func Info(msg string, args ...any) { Get().Info(msg, args...) }

// Info logs a warn message.
func Warn(msg string, args ...any) { Get().Warn(msg, args...) }

// Error logs an error message.
func Error(msg string, args ...any) { Get().Error(msg, args...) }
