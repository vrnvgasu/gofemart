// Package middleware содержит gin-middleware для HTTP-обработчиков.
package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/vrnvgasu/gofemart/pkg/jwt"
)

const userIDKey = "userID"

// Auth возвращает middleware для проверки JWT-токена.
// Токен читается из заголовка Authorization (Bearer) или cookie "token".
// Если токен отсутствует или недействителен, возвращает 401.
func Auth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := ""

		if auth := c.GetHeader("Authorization"); strings.HasPrefix(auth, "Bearer ") {
			tokenStr = strings.TrimPrefix(auth, "Bearer ")
		}

		if tokenStr == "" {
			tokenStr, _ = c.Cookie("token")
		}

		if tokenStr == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		userID, err := jwt.Parse(tokenStr, jwtSecret)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set(userIDKey, userID)
		c.Next()
	}
}

// UserIDFromCtx возвращает ID пользователя из контекста gin.
// Возвращает 0, если ID не найден.
func UserIDFromCtx(c *gin.Context) int64 {
	v, _ := c.Get(userIDKey)
	id, _ := v.(int64)
	return id
}
