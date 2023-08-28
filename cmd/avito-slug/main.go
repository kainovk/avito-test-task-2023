package main

import (
	httpSwagger "github.com/swaggo/http-swagger"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "avito-test-task-2023/docs"
	"avito-test-task-2023/internal/config"
	"avito-test-task-2023/internal/http-server/handlers/segments"
	"avito-test-task-2023/internal/http-server/handlers/users"
	mwLogger "avito-test-task-2023/internal/http-server/middleware/logger"
	"avito-test-task-2023/internal/lib/logger/handlers/slogpretty"
	"avito-test-task-2023/internal/lib/logger/sl"
	"avito-test-task-2023/internal/storage/postgres"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// @title			Avito Test Task
// @version			1.0
// @description		User Segments Service
// @host			localhost:8080
func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info(
		"starting application",
		slog.String("env", cfg.Env),
	)

	storage, err := postgres.New(cfg.Storage)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	_ = storage

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(mwLogger.New(log))
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	r.Route("/users", func(r chi.Router) {
		r.Post("/", users.NewUserSaver(log, storage))
		r.Post("/{user_id}/configure-segments", users.NewUserSegmentConfigurer(log, storage))
		r.Get("/{user_id}/segments", users.NewUserSegmentsGetter(log, storage))
	})

	r.Route("/segments", func(r chi.Router) {
		r.Post("/", segments.NewSegmentSaver(log, storage))
		r.Get("/", segments.NewSegmentGetter(log, storage))
		r.Delete("/{slug}", segments.NewSegmentDeleter(log, storage))
	})

	r.Get("/swagger/*", httpSwagger.Handler())

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      r,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err = srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
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

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
