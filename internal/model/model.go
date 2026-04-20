// Package model содержит основные модели данных приложения.
package model

import "time"

// OrderStatus представляет статус обработки заказа.
type OrderStatus string

const (
	// OrderStatusNew — заказ загружен, но еще не обработан.
	OrderStatusNew OrderStatus = "NEW"
	// OrderStatusProcessing — заказ находится в обработке.
	OrderStatusProcessing OrderStatus = "PROCESSING"
	// OrderStatusInvalid — заказ не прошел проверку.
	OrderStatusInvalid OrderStatus = "INVALID"
	// OrderStatusProcessed — заказ успешно обработан.
	OrderStatusProcessed OrderStatus = "PROCESSED"
)

// User представляет пользователя системы.
type User struct {
	ID           int64
	Login        string
	PasswordHash string
}

// Order представляет заказ пользователя.
type Order struct {
	ID         int64
	UserID     int64
	Number     string
	Status     OrderStatus
	Accrual    *float64
	UploadedAt time.Time
}

// Withdrawal представляет операцию списания баллов.
type Withdrawal struct {
	ID          int64
	UserID      int64
	OrderNumber string
	Sum         float64
	ProcessedAt time.Time
}

// Balance содержит информацию о балансе пользователя.
type Balance struct {
	// Current — текущий остаток баллов.
	Current float64
	// Withdrawn — сумма всех списаний за все время.
	Withdrawn float64
}
