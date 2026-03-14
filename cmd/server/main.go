package main

import (
	"log/slog"
	"net/http"
	"os"

	_ "effective_mobile_test/docs"
	"effective_mobile_test/internal/platform/config"
	"effective_mobile_test/internal/platform/db"
	platformhttp "effective_mobile_test/internal/platform/http"
	"effective_mobile_test/internal/platform/logging"
	"effective_mobile_test/internal/subscriptions/controller"
	"effective_mobile_test/internal/subscriptions/repository"
	"effective_mobile_test/internal/subscriptions/usecase"
)

// @title Subscription Aggregation Service
// @version 1.0
// @description REST service for subscription CRUD and aggregation.
// @basePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	logger := logging.New()

	cfg, err := config.Load()
	if err != nil {
		logger.Error("config load failed", slog.String("error", err.Error()))
		os.Exit(1)
	}

	conn, err := db.Connect(cfg.Database.DSN)
	if err != nil {
		logger.Error("db connect failed", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if err = db.Migrate(cfg.Database.DSN); err != nil {
		logger.Error("migration failed", slog.String("error", err.Error()))
		os.Exit(1)
	}

	repo := repository.NewPostgresRepository(conn)
	uc := usecase.NewSubscriptionUsecase(repo)
	subController := controller.NewSubscriptionController(uc)

	router := platformhttp.NewRouter(logger, subController, platformhttp.NewSwaggerHandler(), cfg.Auth.APIKey)

	server := &http.Server{
		Addr:    cfg.HTTP.Address,
		Handler: router,
	}

	logger.Info("server started", slog.String("address", cfg.HTTP.Address))
	if err = server.ListenAndServe(); err != nil {
		logger.Error("stopping server", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
