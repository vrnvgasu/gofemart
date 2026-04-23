package postgres

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStorage_NotNil(t *testing.T) {
	s := NewStorage()
	assert.NotNil(t, s)
}
