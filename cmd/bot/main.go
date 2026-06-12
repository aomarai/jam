package bot

import (
	"log/slog"
	"os"
)

func main() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(handler)
	slog.SetDefault(logger)
	slog.Info("JAM Bot started...")
	defer slog.Info("JAM Bot stopped.")
}
