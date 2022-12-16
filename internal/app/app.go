// Package app configures and runs application.
package app

import (
	"fmt"
	"lifo-rest-api/internal/config"
	"lifo-rest-api/internal/controller/http"
	"lifo-rest-api/internal/repository"
	"lifo-rest-api/internal/service"
	"lifo-rest-api/pkg/database/postgres"
	"lifo-rest-api/pkg/database/postgres/migration"
	"lifo-rest-api/pkg/httpserver"
	"lifo-rest-api/pkg/logger"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

// Run - running all processes
func Run(cfg *config.Config) {
	logger.InitZapLogger(cfg.Log.LogLevel)

	migration.Migrate(cfg.Postgres.URI)

	db, err := postgres.New(cfg.Postgres.URI)
	if err != nil {
		zap.S().Fatalf("app - Run - postgres.New: %s", err.Error())
	}

	repos := repository.NewRepositories(db)

	serv := service.NewServices(&service.Dependencies{
		Repos: repos,
	})

	handler := http.NewHandler(serv)

	// HTTP Run
	srv := httpserver.New(handler.Init(), httpserver.Port(cfg.HTTP.Port))
	zap.S().Infof("HTTP server started on address [%s]", fmt.Sprintf(":%s", cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		zap.S().Info("app - Run - signal: " + s.String())
	case err := <-srv.Notify():
		zap.S().Errorf("app - Run - httpServer.Notify: %s", err.Error())
	}

	// Close db
	db.Close()

	// Shutdown
	err = srv.Shutdown()
	if err != nil {
		zap.S().Errorf("app - Run - httpServer.Shutdown: %s", err.Error())
	}
}
