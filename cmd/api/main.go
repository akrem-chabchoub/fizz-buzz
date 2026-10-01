package main

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/akrem-chabchoub/fizz-buzz/internal/app"
	"github.com/akrem-chabchoub/fizz-buzz/internal/config"
	"github.com/akrem-chabchoub/fizz-buzz/internal/httpapi"
	"github.com/akrem-chabchoub/fizz-buzz/internal/logger"
	"github.com/akrem-chabchoub/fizz-buzz/internal/stats/redistore"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("config", "err", err)
		os.Exit(1)
	}
	log, err := logger.New(os.Stdout, cfg.LogLevel)
	if err != nil {
		slog.Error("logger", "err", err)
		os.Exit(1)
	}
	slog.SetDefault(log)

	store := redistore.New(cfg.RedisAddr, cfg.RedisPassword)
	defer func() { _ = store.Close() }()
	pingCtx, pingCancel := context.WithTimeout(context.Background(), 2*time.Second)
	err = store.Ping(pingCtx)
	pingCancel()
	if err != nil {
		log.Error("redis", "err", err)
		os.Exit(1)
	}

	ln, err := net.Listen("tcp", cfg.HTTPAddr)
	if err != nil {
		log.Error("listen", "err", err)
		os.Exit(1)
	}
	api := httpapi.New(app.New(store, cfg.StatsTimeout), cfg.MaxLimit, cfg.MaxStringLen)
	srv := &http.Server{
		Handler:           api.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	log.Info("server starting", "addr", ln.Addr().String())
	if err := serve(ctx, srv, ln, cfg.ShutdownTimeout); err != nil {
		log.Error("serve", "err", err)
		os.Exit(1)
	}
	log.Info("server stopped")
}

func serve(ctx context.Context, srv *http.Server, ln net.Listener, timeout time.Duration) error {
	errCh := make(chan error, 1)
	go func() {
		err := srv.Serve(ln)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return err
	}
	return <-errCh
}
