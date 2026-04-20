package jwt_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vrnvgasu/gofemart/pkg/jwt"
)

const testSecret = "test-secret-key"

func TestGenerate(t *testing.T) {
	userID := int64(1)

	t.Run("success", func(t *testing.T) {
		token, err := jwt.Generate(userID, testSecret)
		require.NoError(t, err)
		assert.NotEmpty(t, token)

		got, err := jwt.Parse(token, testSecret)
		require.NoError(t, err)
		assert.Equal(t, userID, got)
	})

	t.Run("parse with wrong secret", func(t *testing.T) {
		token, err := jwt.Generate(userID, testSecret)
		require.NoError(t, err)

		_, err = jwt.Parse(token, "wrong-secret")
		assert.Error(t, err)
	})
}
