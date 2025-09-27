//go:build exclude

package app

import (
	"log/slog"
	"os"
	"strings"
)

var (
	logLevel slog.LevelVar
	dLog     *slog.Logger
)

func init() {
	// a variable will hold the current log level (dynamic changes)
	logLevel.Set(getLogLevelFromEnv())
	// an instance of the level logger configured to the dynamic variable
	dLog = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: &logLevel,
	}))
	// set it as the global logger
	slog.SetDefault(dLog)
	//slog.SetLogLoggerLevel()
}

func getLogLevelFromEnv() slog.Level {
	levelStr := os.Getenv("LOG_LEVEL")
	switch strings.ToLower(levelStr) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
