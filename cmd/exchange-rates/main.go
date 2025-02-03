package main

import (
	"exchange-rates/internal/config"
	"exchange-rates/internal/http-server/handlers/currency-rates/all"
	"exchange-rates/internal/http-server/handlers/currency-rates/ondate"
	"exchange-rates/internal/lib/logger/sl"
	"exchange-rates/internal/storage/mysql"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cnf := config.MustLoad()
	log := setupLogger(cnf.Env)

	log.Info("starting up", slog.String("env", cnf.Env))

	storage, err := mysql.New(cnf)
	if err != nil {
		log.Error("failed to initialize storage", sl.Error(err))
		os.Exit(1)
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Get("/exrates/all", all.New(log, storage))
	router.Get("/exrates/ondate/{date}", ondate.New(log, storage))

	log.Info("starting server", slog.String("address", cnf.HTTPServer.Addr))

	srv := &http.Server{
		Addr:         cnf.HTTPServer.Addr,
		Handler:      router,
		ReadTimeout:  cnf.HTTPServer.Timeout,
		WriteTimeout: cnf.HTTPServer.Timeout,
		IdleTimeout:  cnf.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
