package worker

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/vrnvgasu/gofemart/internal/accrual"
	"github.com/vrnvgasu/gofemart/internal/model"
	"github.com/vrnvgasu/gofemart/internal/repository"
	mockrepository "github.com/vrnvgasu/gofemart/internal/repository/mocks"
)

func TestWorker_UpdatesProcessedOrder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		statusFromClient accrual.OrderStatus
		storage          func(*mockrepository.MockStorage) repository.Storage
	}{
		{
			name:             "StatusProcessed",
			statusFromClient: accrual.StatusProcessed,
			storage: func(store *mockrepository.MockStorage) repository.Storage {
				store.EXPECT().GetPendingOrders(gomock.Any()).
					Return([]model.Order{{}}, nil)
				store.EXPECT().UpdateOrderStatus(gomock.Any(), gomock.Any(), model.OrderStatusProcessed, gomock.Any()).
					Return(nil)
				return store
			},
		},
		{
			name:             "StatusInvalid",
			statusFromClient: accrual.StatusInvalid,
			storage: func(store *mockrepository.MockStorage) repository.Storage {
				store.EXPECT().GetPendingOrders(gomock.Any()).
					Return([]model.Order{{}}, nil)
				store.EXPECT().UpdateOrderStatus(gomock.Any(), gomock.Any(), model.OrderStatusInvalid, gomock.Any()).
					Return(nil)
				return store
			},
		},
		{
			name:             "StatusProcessing",
			statusFromClient: accrual.StatusProcessing,
			storage: func(store *mockrepository.MockStorage) repository.Storage {
				store.EXPECT().GetPendingOrders(gomock.Any()).
					Return([]model.Order{{}}, nil)
				store.EXPECT().UpdateOrderStatus(gomock.Any(), gomock.Any(), model.OrderStatusProcessing, gomock.Any()).
					Return(nil)
				return store
			},
		},
		{
			name:             "wrong status",
			statusFromClient: "dummy",
			storage: func(store *mockrepository.MockStorage) repository.Storage {
				store.EXPECT().GetPendingOrders(gomock.Any()).
					Return([]model.Order{{}}, nil)
				return store
			},
		},
		{
			name:             "no orders",
			statusFromClient: accrual.StatusProcessing,
			storage: func(store *mockrepository.MockStorage) repository.Storage {
				store.EXPECT().GetPendingOrders(gomock.Any()).
					Return([]model.Order{}, nil)
				return store
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			accrualVal := 100.0
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(accrual.OrderInfo{
					Order:   "12345678903",
					Status:  tt.statusFromClient,
					Accrual: &accrualVal,
				})
			}))
			defer srv.Close()

			w := Worker{
				storage:       tt.storage(mockrepository.NewMockStorage(ctrl)),
				accrualClient: accrual.NewClient(srv.URL),
			}
			w.poll(t.Context())
		})
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := mockrepository.NewMockStorage(ctrl)
	client := accrual.NewClient("http://localhost")

	w := New(storage, client)

	assert.NotNil(t, w)
	assert.Equal(t, defaultPollInterval, w.pollInterval)
}

func TestWithPollInterval(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := mockrepository.NewMockStorage(ctrl)
	w := New(storage, accrual.NewClient("http://localhost"))

	result := w.WithPollInterval(5 * time.Second)

	assert.Equal(t, 5*time.Second, w.pollInterval)
	assert.Equal(t, w, result)
}

func TestRun_StopsOnContextCancel(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	storage := mockrepository.NewMockStorage(ctrl)
	w := New(storage, accrual.NewClient("http://localhost")).
		WithPollInterval(10 * time.Second)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	done := make(chan struct{})
	go func() {
		w.Run(ctx)
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Run did not stop after context cancellation")
	}
}
