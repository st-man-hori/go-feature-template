// Command api starts the HTTP server.
//
// Laravel comparison: this is the equivalent of `php artisan serve` in dev,
// and of the nginx + php-fpm pair the Laravel template's compose.yaml runs
// in "prod-like" local mode. Go has no separate dev-server command — a
// compiled binary *is* the server — which is also why this template's
// compose.yaml needs only one app service instead of Laravel's app+nginx
// split.
//
//	@title			Go Feature Template API
//	@version		1.0
//	@description	Sample CRUD API demonstrating a package-by-feature Go project, ported from laravel-feature-template.
//	@BasePath		/api
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/st-man-hori/go-feature-template/internal/api"
	"github.com/st-man-hori/go-feature-template/internal/config"
	"github.com/st-man-hori/go-feature-template/internal/database"
)

func main() {
	cfg := config.Load()

	db, err := database.Open(cfg)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	router := api.NewRouter(db)

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		slog.Info("starting server", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	slog.Info("shutting down server")
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
	}
}
