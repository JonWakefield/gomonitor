package logging

import (
	"log"
	"log/slog"
	"os"
	"strings"
)

func SetupLogger() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "debug", "d":
		slog.SetLogLoggerLevel(slog.LevelDebug)
	case "info", "i":
		slog.SetLogLoggerLevel(slog.LevelInfo)
	default:
		slog.SetLogLoggerLevel(slog.LevelWarn)
	}
}
