package retry

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRetry(t *testing.T) {
	t.Parallel()

	calls := 0

	tests := []struct {
		name     string
		fn       func() func() error
		expected require.ErrorAssertionFunc
	}{
		{
			name: "success",
			fn: func() func() error {
				return func() error {
					return nil
				}
			},
			expected: require.NoError,
		},
		{
			name: "not retryable",
			fn: func() func() error {
				return func() error {
					return errors.New("error")
				}
			},
			expected: require.Error,
		},
		{
			name: "retryable",
			fn: func() func() error {
				return func() error {
					calls++
					if calls > 2 {
						return nil
					} else {
						return NewRetryableError(errors.New("error"))
					}
				}
			},
			expected: require.NoError,
		},
		{
			name: "retryable max retries",
			fn: func() func() error {
				return func() error {
					return NewRetryableError(errors.New("error"))
				}
			},
			expected: require.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			tt.expected(t, RetryWithSettings(tt.fn(), &Config{
				MaxRetries:         3,
				StartRetryInterval: 10 * time.Millisecond,
				AddRetryPeriod:     20 * time.Millisecond,
			}))
			calls = 0
		})
	}
}
