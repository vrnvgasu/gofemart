package luhn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValid(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		number string
		want   bool
	}{
		{
			name:   "valid number from task",
			number: "12345678903",
			want:   true,
		},
		{
			name:   "valid number 9278923470",
			number: "9278923470",
			want:   true,
		},
		{
			name:   "valid number 2377225624",
			number: "2377225624",
			want:   true,
		},
		{
			name:   "invalid number",
			number: "12345678901",
			want:   false,
		},
		{
			name:   "single digit zero",
			number: "0",
			want:   true,
		},
		{
			name:   "empty string",
			number: "",
			want:   false,
		},
		{
			name:   "non-digit characters",
			number: "1234abc",
			want:   false,
		},
		{
			name:   "valid 346436439",
			number: "346436439",
			want:   true,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, Valid(tt.number))
		})
	}
}
