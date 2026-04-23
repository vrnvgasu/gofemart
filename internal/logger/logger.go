// Package logger предоставляет глобальный логгер на основе zap.
package logger

import (
	"go.uber.org/zap"
)

// Log — глобальный логгер. По умолчанию используется nop-логгер (ничего не пишет).
// Нужно вызвать Initialize перед использованием.
var Log *zap.SugaredLogger = zap.NewNop().Sugar()

// Initialize настраивает глобальный логгер с указанным уровнем логирования.
func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}

	cfg := zap.NewProductionConfig()
	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		return err
	}

	Log = zl.Sugar()
	defer Log.Sync()

	return nil
}
