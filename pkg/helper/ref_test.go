package helper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRefFloat64(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		value         float64
		expectedValue float64
	}{
		{
			name:          "0",
			value:         0,
			expectedValue: 0,
		},
		{
			name:          "1.1",
			value:         1.1,
			expectedValue: 1.1,
		},
		{
			name:          "-1.1",
			value:         -1.1,
			expectedValue: -1.1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.expectedValue, *NewRefFloat64(tt.value))
		})
	}
}

func TestNewRefInt64(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		value         int64
		expectedValue int64
	}{
		{
			name:          "0",
			value:         0,
			expectedValue: 0,
		},
		{
			name:          "1",
			value:         1,
			expectedValue: 1,
		},
		{
			name:          "-1",
			value:         -1,
			expectedValue: -1,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			require.Equal(t, tt.expectedValue, *NewRefInt64(tt.value))
		})
	}
}
