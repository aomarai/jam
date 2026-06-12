package bot

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aomarai/jam/internal/chain"
	"github.com/aomarai/jam/internal/config"
	"github.com/aomarai/jam/internal/twitch"
)

func main() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(handler)
	slog.SetDefault(logger)
	slog.Info("JAM Bot started...")
	defer slog.Info("JAM Bot stopped.")

	// Load environment variables
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load config", "error", err)
		return
	}
	slog.Info("Successfully loaded configuration.")
	slog.Debug("Loaded configuration", "config", cfg)

	c := chain.New(cfg.MMarkovStateSize)
	if err := c.Load(cfg.ModelPath); err != nil {
		slog.Warn("No existing model found. Creating new model...")
		if cfg.SeedPath != "" {
			data, err := os.ReadFile(cfg.SeedPath)
			if err != nil {
				slog.Error("Failed to read seed file", "error", err)
			}
			c.Seed(string(data))
			slog.Info("Seeded chain", "seedPath", cfg.SeedPath, "size", c.Size())
		} else {
			slog.Info("Loaded pre-existing model", "modelPath", cfg.ModelPath, "size", c.Size())
		}
	}

	client := twitch.NewClient(cfg, c)

	// Graceful shutdown for SIGINT/SIGTERM
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := client.Connect(); err != nil {
			slog.Error("Twitch client error", "error", err)
		}
	}()

	// Track time between posts
	postTicker := time.NewTicker(cfg.PostInterval)
	defer postTicker.Stop()

	// Save timer for saving the current model
	saveTicker := time.NewTicker(cfg.SaveInterval)
	defer saveTicker.Stop()

	for {
		select {
		case <-postTicker.C:
			client.Post()
		case <-saveTicker.C:
			if err := c.Save(cfg.ModelPath); err != nil {
				slog.Error("Failed to save model", "error", err, "modelPath", cfg.ModelPath, "size", c.Size())
			} else {
				slog.Info("Saved model", "modelPath", cfg.ModelPath, "size", c.Size())
			}
		case <-stop:
			slog.Info("Shutting down...")
			// Final save before shutdown
			if err := c.Save(cfg.ModelPath); err != nil {
				slog.Error("Failed to save model during shutdown", "error", err, "modelPath", cfg.ModelPath, "size", c.Size())
			}
			return
		}
	}
}
