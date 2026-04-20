package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	serviceerrors "github.com/vrnvgasu/gofemart/internal/app"
	"github.com/vrnvgasu/gofemart/internal/handler/middleware"
	"github.com/vrnvgasu/gofemart/internal/handler/response"
)

type withdrawRequest struct {
	Order string  `json:"order" binding:"required"`
	Sum   float64 `json:"sum"   binding:"required,gt=0"`
}

// SaveWithdraw handles POST /api/user/balance/withdraw.
func (h *Handler) SaveWithdraw(c *gin.Context) {
	var req withdrawRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResponseError(c, serviceerrors.BadRequestError())
		return
	}

	userID := middleware.UserIDFromCtx(c)

	err := h.app.CreateWithdrawal(c.Request.Context(), userID, req.Order, req.Sum)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

// GetWithdrawals handles GET /api/user/withdrawals.
func (h *Handler) GetWithdrawals(c *gin.Context) {
	userID := middleware.UserIDFromCtx(c)

	resp, err := h.app.GetWithdrawals(c, userID)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
