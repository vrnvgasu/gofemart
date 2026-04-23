package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitialize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		level       string
		expectedErr require.ErrorAssertionFunc
	}{
		{
			level:       "info",
			expectedErr: require.NoError,
		},
		{
			level:       "debug",
			expectedErr: require.NoError,
		},
		{
			level:       "warn",
			expectedErr: require.NoError,
		},
		{
			level:       "error",
			expectedErr: require.NoError,
		},
		{
			level:       "invalid",
			expectedErr: require.Error,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.level, func(t *testing.T) {
			t.Parallel()
			err := Initialize(tt.level)
			tt.expectedErr(t, err)
			assert.NotNil(t, Log)

		})
	}

}

func TestLog_DefaultNotNil(t *testing.T) {
	// Log must never be nil
	assert.NotNil(t, Log)
}
