package sl

import (
	"log/slog"
	"os"
)

// for unit testing
func SetupLogger() (log *slog.Logger) {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
}
