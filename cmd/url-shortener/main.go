package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/StiVit/URL-shortener-API/internal/config"
	"github.com/StiVit/URL-shortener-API/internal/http-server/handlers/delete"
	"github.com/StiVit/URL-shortener-API/internal/http-server/handlers/redirect"
	"github.com/StiVit/URL-shortener-API/internal/http-server/handlers/url/save"
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

	// init config: cleanenv
	config := config.MustLoad()

	// init logger: slog
	log := setupLogger(config.Env)
	log.Info("Starting URL Shortener API", "env", config.Env)
	log.Debug("Debug messages are enabled")


	// init storage: sqlite
	storage, err := sqlite.New(config.StoragePath)
	if err != nil {
		log.Error("Failed to initialize storage", sl.Err(err))
		os.Exit(1)
	}

	_ = storage

	// init router: chi, " chi render"
	router := chi.NewRouter()

	// middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer) // If the handler panics, this recovers and writes a 500
	router.Use(middleware.URLFormat) // If the URL is not v alid, this middleware will return a 400 Bad Request
 
	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			config.HTTPServer.User: config.HTTPServer.Password,
		}))

		r.Post("/", save.New(log, storage)) 
		r.Delete("/{alias}", delete.New(log, storage))
	})

	router.Get("/{alias}", redirect.New(log, storage))
	// init server

	log.Info("starting server", slog.String("address", config.Address))

	srv := &http.Server{
		Addr: config.Address,
		Handler: router, // Even if router uses handlers to connect and process requests, the router by itself is s handler
		ReadTimeout: config.HTTPServer.Timeout,
		WriteTimeout: config.HTTPServer.Timeout,
		IdleTimeout: config.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("Failed to start server")
	}

	log.Error("Server Stopped")

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