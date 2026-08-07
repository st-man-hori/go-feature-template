// Command api starts the HTTP server.
//
// Laravel比較: 開発時の `php artisan serve` や、Laravel版の compose.yaml が
// 本番相当のローカルモードで動かす nginx + php-fpm の組み合わせに相当する。
// Go には別建ての開発サーバーコマンドが存在しない — コンパイル済みバイナリ
// そのものがサーバーになるため、この template の compose.yaml が Laravel の
// app+nginx 構成と違い app サービス1つで済むのもそれが理由。
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
