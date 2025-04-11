package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/StiVit/URL-shortener-API/internal/config"
	"github.com/StiVit/URL-shortener-API/internal/lib/logger/sl"
	"github.com/StiVit/URL-shortener-API/internal/storage/sqlite"
	"github.com/joho/godotenv"
)

const (
	envLocal = "local"
	envDev   = "development"
	envProd  = "production"
)

func main() {
    if err := godotenv.Load("local.env"); err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

	// TODO: init config: cleanenv
	config := config.MustLoad()

	// TODO: init logger: slog
	log := setupLogger(config.Env)
	log.Info("Starting URL Shortener API", "env", config.Env)
	log.Debug("Debug messages are enabled")


	// TODO: init storage: sqlite
	storage, err := sqlite.New(config.StoragePath)
	if err != nil {
		log.Error("Failed to initialize storage", sl.Err(err))
		os.Exit(1)
	}

	_ = storage

	// TODO: init router: chi, " chi render"

	// TODO: init server: 

}

func setupLogger(env string) *slog.Logger{
	var log *slog.Logger
	switch env { 
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}