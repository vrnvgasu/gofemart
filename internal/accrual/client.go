// Package accrual предоставляет клиент для взаимодействия с внешней системой расчета начислений.
package accrual

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// OrderStatus представляет статус обработки заказа в системе начислений.
type OrderStatus string

const (
	// StatusRegistered — заказ зарегистрирован, но вознаграждение еще не рассчитано.
	StatusRegistered OrderStatus = "REGISTERED"
	// StatusInvalid — заказ не принят к расчету.
	StatusInvalid OrderStatus = "INVALID"
	// StatusProcessing — расчет начисления в процессе.
	StatusProcessing OrderStatus = "PROCESSING"
	// StatusProcessed — расчет начисления завершен.
	StatusProcessed OrderStatus = "PROCESSED"
)

// OrderInfo содержит информацию о заказе из системы начислений.
type OrderInfo struct {
	Order   string      `json:"order"`
	Status  OrderStatus `json:"status"`
	Accrual *float64    `json:"accrual,omitempty"`
}

// ErrTooManyRequests возвращается, когда превышен лимит запросов к системе начислений.
type ErrTooManyRequests struct {
	// RetryAfter — через сколько нужно повторить запрос.
	RetryAfter time.Duration
}

// Error возвращает текст ошибки с указанием времени ожидания.
func (e *ErrTooManyRequests) Error() string {
	return fmt.Sprintf("accrual: too many requests, retry after %s", e.RetryAfter)
}

// Client — HTTP-клиент для системы начислений.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient создает новый Client с указанным базовым URL.
func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// GetOrder запрашивает информацию о заказе из системы начислений.
// Возвращает nil, если заказ не найден (204).
// Возвращает ErrTooManyRequests при превышении лимита запросов (429).
func (c *Client) GetOrder(ctx context.Context, number string) (*OrderInfo, error) {
	url := fmt.Sprintf("%s/api/orders/%s", c.baseURL, number)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("accrual.GetOrder NewRequest: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("accrual.GetOrder Do: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var info OrderInfo
		if err = json.NewDecoder(resp.Body).Decode(&info); err != nil {
			return nil, fmt.Errorf("accrual.GetOrder Decode: %w", err)
		}
		return &info, nil
	case http.StatusNoContent:
		return nil, nil
	case http.StatusTooManyRequests:
		retryAfter := parseRetryAfter(resp.Header.Get("Retry-After"))
		return nil, &ErrTooManyRequests{RetryAfter: retryAfter}
	default:
		return nil, fmt.Errorf("accrual.GetOrder unexpected status: %d", resp.StatusCode)
	}
}

func parseRetryAfter(value string) time.Duration {
	var seconds int
	if _, err := fmt.Sscanf(value, "%d", &seconds); err == nil && seconds > 0 {
		return time.Duration(seconds) * time.Second
	}
	return 60 * time.Second
}
