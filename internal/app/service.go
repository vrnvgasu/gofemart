package app

import (
	"github.com/vrnvgasu/gofemart/internal/config"
	"github.com/vrnvgasu/gofemart/internal/repository"
)

type App struct {
	storage repository.Storage
	cfg     *config.Config
}

func NewApp(storage repository.Storage, cfg *config.Config) *App {
	return &App{
		storage: storage,
		cfg:     cfg,
	}
}
