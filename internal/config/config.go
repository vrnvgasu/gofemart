// Package config содержит конфигурацию приложения.
package config

// Config хранит настройки сервиса, которые задаются через переменные окружения или флаги CLI.
type Config struct {
	// RunAddress — адрес и порт, на котором запускается сервис.
	RunAddress string `env:"RUN_ADDRESS"`
	// DatabaseURI — строка подключения к базе данных.
	DatabaseURI string `env:"DATABASE_URI"`
	// AccrualSystemAddress — адрес внешней системы расчета начислений.
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	// JWTSecret — секрет для подписи JWT-токенов.
	JWTSecret string `env:"JWT_SECRET"`
	// LogLevel — уровень логирования (info, debug, error и т.д.).
	LogLevel string `env:"LOG_LEVEL"`
}
