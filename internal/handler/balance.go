package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/vrnvgasu/gofemart/internal/handler/middleware"
	"github.com/vrnvgasu/gofemart/internal/handler/response"
)

// GetBalance handles GET /api/user/balance.
func (h *Handler) GetBalance(c *gin.Context) {
	userID := middleware.UserIDFromCtx(c)

	balance, err := h.app.GetBalance(c, userID)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	c.JSON(http.StatusOK, balance)
}
