// Package retry предоставляет утилиты для повторного выполнения операций при ошибках.
package retry

import (
	"errors"
	"time"
)

// Config содержит настройки для повторных попыток.
type Config struct {
	// MaxRetries — максимальное количество повторных попыток.
	MaxRetries int
	// StartRetryInterval — начальный интервал между попытками.
	StartRetryInterval time.Duration
	// AddRetryPeriod — на сколько увеличивается интервал с каждой попыткой.
	AddRetryPeriod time.Duration
}

// DefaultConfig возвращает конфигурацию с разумными дефолтными значениями.
func DefaultConfig() *Config {
	return &Config{
		MaxRetries:         3,
		StartRetryInterval: 1 * time.Second,
		AddRetryPeriod:     2 * time.Second,
	}
}

// RetryableError оборачивает ошибку и сигнализирует, что операцию можно повторить.
type RetryableError struct {
	Err error
}

// NewRetryableError создает новую RetryableError, оборачивая переданную ошибку.
func NewRetryableError(err error) error {
	return &RetryableError{Err: err}
}

// Error возвращает текст ошибки.
func (e *RetryableError) Error() string {
	return e.Err.Error()
}

// Unwrap возвращает оригинальную ошибку.
func (e *RetryableError) Unwrap() error {
	return e.Err
}

// RetryWithSettings выполняет функцию f с повторными попытками согласно конфигурации cnf.
// Если cnf равен nil, используется DefaultConfig.
// Повтор происходит только если функция вернула RetryableError.
func RetryWithSettings(f func() error, cnf *Config) error {
	if cnf == nil {
		cnf = DefaultConfig()
	}

	retryInterval := cnf.StartRetryInterval
	attempt := 0

	for {
		err := f()
		if err == nil {
			return nil
		}

		var re *RetryableError
		if !errors.As(err, &re) {
			return err
		}
		if attempt >= cnf.MaxRetries {
			return err
		}

		time.Sleep(retryInterval)
		attempt++
		retryInterval += cnf.AddRetryPeriod
	}
}

// Retry выполняет функцию f с повторными попытками, используя DefaultConfig.
func Retry(f func() error) error {
	return RetryWithSettings(f, DefaultConfig())
}
