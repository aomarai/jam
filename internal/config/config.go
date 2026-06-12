package config

import (
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	TwitchOAuthToken    string
	TwitchBotUsername   string
	TwitchTargetChannel string
	MarkovStateSize     int
	PostInterval        time.Duration
	SaveInterval        time.Duration
	ModelPath           string
	SeedPath            string
	Environment         string
	LogLevel            slog.Level
}

// getEnv returns the value of an environment variable or a fallback if it is not set.
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func parseLogLevel(s string) slog.Level {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo // default if empty or unrecognized
	}
}

// Load retrieves the configuration from environment variables and loads it into the application.
func Load() (*Config, error) {
	slog.Info("Loading configuration...")
	stateSize, err := strconv.Atoi(getEnv("MARKOV_STATE_SIZE", "3"))
	if err != nil {
		return nil, err
	}

	postInterval, err := time.ParseDuration(getEnv("POST_INTERVAL", "1h"))
	if err != nil {
		return nil, err
	}

	saveInterval, err := time.ParseDuration(getEnv("SAVE_INTERVAL", "24h"))
	if err != nil {
		return nil, err
	}

	return &Config{
		TwitchOAuthToken:    getEnv("TWITCH_OAUTH_TOKEN", ""),
		TwitchBotUsername:   getEnv("TWITCH_BOT_USERNAME", ""),
		TwitchTargetChannel: getEnv("TWITCH_TARGET_CHANNEL", ""),
		MarkovStateSize:     stateSize,
		PostInterval:        postInterval,
		SaveInterval:        saveInterval,
		ModelPath:           getEnv("MODEL_PATH", "model.json"),
		SeedPath:            getEnv("SEED_PATH", "seed.txt"),
		Environment:         getEnv("ENVIRONMENT", "production"),
		LogLevel:            parseLogLevel(os.Getenv("LOG_LEVEL")),
	}, nil
}
