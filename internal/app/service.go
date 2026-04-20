package app

import (
	"github.com/vrnvgasu/gofemart/internal/config"
	"github.com/vrnvgasu/gofemart/internal/repository"
)

// App реализует бизнес-логику приложения.
type App struct {
	storage repository.Storage
	cfg     *config.Config
}

// NewApp создает новый экземпляр App с переданным хранилищем и конфигурацией.
func NewApp(storage repository.Storage, cfg *config.Config) *App {
	return &App{
		storage: storage,
		cfg:     cfg,
	}
}
