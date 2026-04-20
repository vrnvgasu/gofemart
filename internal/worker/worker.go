package worker

import (
	"context"
	"errors"
	"time"

	"github.com/vrnvgasu/gofemart/internal/accrual"
	"github.com/vrnvgasu/gofemart/internal/logger"
	"github.com/vrnvgasu/gofemart/internal/model"
	"github.com/vrnvgasu/gofemart/internal/repository"
)

const defaultPollInterval = 2 * time.Second

type Worker struct {
	storage       repository.Storage
	accrualClient *accrual.Client
	pollInterval  time.Duration
}

func New(storage repository.Storage, accrualClient *accrual.Client) *Worker {
	return &Worker{
		storage:       storage,
		accrualClient: accrualClient,
		pollInterval:  defaultPollInterval,
	}
}

func (w *Worker) WithPollInterval(d time.Duration) *Worker {
	w.pollInterval = d
	return w
}

func (w *Worker) Run(ctx context.Context) {
	ticker := time.NewTicker(w.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.poll(ctx)
		}
	}
}

func (w *Worker) poll(ctx context.Context) {
	orders, err := w.storage.GetPendingOrders(ctx)
	if err != nil {
		logger.Log.Errorw("worker.poll GetPendingOrders", "error", err)
		return
	}

	for _, o := range orders {
		w.processOrder(ctx, o)
	}
}

func (w *Worker) processOrder(ctx context.Context, o model.Order) {
	info, err := w.accrualClient.GetOrder(ctx, o.Number)
	if err != nil {
		var tooMany *accrual.ErrTooManyRequests
		if errors.As(err, &tooMany) {
			logger.Log.Infow("worker.processOrder rate limited", "retry_after", tooMany.RetryAfter)
			time.Sleep(tooMany.RetryAfter)
			return
		}
		logger.Log.Errorw("worker.processOrder GetOrder", "order", o.Number, "error", err)
		return
	}

	if info == nil {
		// Order not yet registered in accrual system — skip.
		return
	}

	var newStatus model.OrderStatus
	switch info.Status {
	case accrual.StatusProcessed:
		newStatus = model.OrderStatusProcessed
	case accrual.StatusInvalid:
		newStatus = model.OrderStatusInvalid
	case accrual.StatusProcessing:
		newStatus = model.OrderStatusProcessing
	default:
		return
	}

	if err = w.storage.UpdateOrderStatus(ctx, o.Number, newStatus, info.Accrual); err != nil {
		logger.Log.Errorw("worker.processOrder UpdateOrderStatus", "order", o.Number, "error", err)
	}
}
