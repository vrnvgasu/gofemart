package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/spf13/pflag"

	"github.com/vrnvgasu/gofemart/internal/accrual"
	"github.com/vrnvgasu/gofemart/internal/app"
	"github.com/vrnvgasu/gofemart/internal/config"
	"github.com/vrnvgasu/gofemart/internal/handler"
	"github.com/vrnvgasu/gofemart/internal/logger"
	"github.com/vrnvgasu/gofemart/internal/repository/postgres"
	"github.com/vrnvgasu/gofemart/internal/worker"
)

func main() {
	cfg := parseConfig()

	if err := logger.Initialize(cfg.LogLevel); err != nil {
		log.Fatalf("logger.Initialize: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	storage := postgres.NewStorage()
	if err := storage.Start(ctx, cfg.DatabaseURI); err != nil {
		logger.Log.Fatalf("storage.Start: %v", err)
	}
	defer storage.Stop()

	accrualClient := accrual.NewClient(cfg.AccrualSystemAddress)

	w := worker.New(storage, accrualClient)
	go w.Run(ctx)

	appService := app.NewApp(storage, cfg)
	h := handler.NewHandler(appService)
	router := handler.NewRouter(h, cfg)

	srv := &http.Server{
		Addr:    cfg.RunAddress,
		Handler: router,
	}

	go func() {
		logger.Log.Infow("starting server", "address", cfg.RunAddress)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Fatalf("srv.ListenAndServe: %v", err)
		}
	}()

	<-ctx.Done()
	logger.Log.Info("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Log.Errorw("srv.Shutdown", "error", err)
	}
}

func parseConfig() *config.Config {
	cfg := &config.Config{}

	pflag.StringVarP(&cfg.RunAddress, "address", "a", ":8000", "address and port to listen on")
	pflag.StringVarP(&cfg.DatabaseURI, "database", "d", "postgres://gophermart:gophermart@localhost:5432/gophermart?sslmode=disable", "PostgreSQL connection URI")
	pflag.StringVarP(&cfg.AccrualSystemAddress, "accrual", "r", "", "accrual system base URL")
	pflag.StringVarP(&cfg.JWTSecret, "secret", "s", "secret", "JWT signing secret")
	pflag.StringVarP(&cfg.LogLevel, "loglevel", "l", "info", "log level")
	pflag.Parse()

	if err := env.Parse(cfg); err != nil {
		log.Fatalf("env.Parse: %v", err)
	}

	return cfg
}
