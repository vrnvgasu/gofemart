// Package helper содержит вспомогательные утилиты.
package helper

// NewRefFloat64 возвращает указатель на переданное значение float64.
func NewRefFloat64(v float64) *float64 {
	return &v
}

// NewRefInt64 возвращает указатель на переданное значение int64.
func NewRefInt64(v int64) *int64 {
	return &v
}
