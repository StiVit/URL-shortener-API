package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/StiVit/URL-shortener-API/internal/config"
	"github.com/StiVit/URL-shortener-API/internal/http-server/middleware/logger"
	"github.com/StiVit/URL-shortener-API/internal/lib/logger/sl"
	"github.com/StiVit/URL-shortener-API/internal/storage/sqlite"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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
	router := chi.NewRouter()

	// middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer) // If the handler panics, this recovers and writes a 500
	router.Use(middleware.URLFormat) // If the URL is not valid, this middleware will return a 400 Bad Request

	// TODO: init server

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