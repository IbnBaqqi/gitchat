// Package main is the entry point for the book-me server.
package main

import (
	"context"
	"fmt"

	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IbnBaqqi/gitchat/internal/app"
	"github.com/IbnBaqqi/gitchat/internal/config"
	"github.com/IbnBaqqi/gitchat/internal/database"
	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}

func run() error {

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	log, _ := config.NewLogger(cfg.Server.LogLevel, cfg.App.Env)

	log.Info("starting gitchat server",
		zap.String("port", cfg.Server.Port),
		zap.String("log_level", cfg.Server.LogLevel),
	)

	ctx := context.Background()
	db, err := database.Connect(ctx, &cfg.Database, log)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			slog.Error("failed to close database connection", "error", err)
		}
	}()

	handler := app.NewRouter(cfg, db)
	// if err != nil {
	// 	return fmt.Errorf("failed to initialize api services: %w", err)
	// }

	// handler := app.NewRouter(apiCfg)

	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      handler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	errCh := make(chan error, 1)
	go func() {
		slog.Info("Server listening", "address", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case <-quit:
		slog.Info("Server is shutting down...")
	case err := <-errCh:
		return fmt.Errorf("server failed to start: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	slog.Info("Server exited gracefully")
	return nil
}