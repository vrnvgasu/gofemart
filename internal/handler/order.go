package handler

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	serviceerrors "github.com/vrnvgasu/gofemart/internal/app"
	"github.com/vrnvgasu/gofemart/internal/handler/middleware"
	"github.com/vrnvgasu/gofemart/internal/handler/response"
)

// UploadOrder handles POST /api/user/orders.
func (h *Handler) UploadOrder(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil || len(body) == 0 {
		response.ResponseError(c, serviceerrors.BadRequestError())
		return
	}

	userID := middleware.UserIDFromCtx(c)
	isCreated, err := h.app.CreateOrder(c.Request.Context(), userID, string(body))
	if err != nil {
		response.ResponseError(c, err)
		return
	}
	if isCreated {
		c.Status(http.StatusAccepted)
		return
	} else {
		c.Status(http.StatusOK)
		return
	}
}

// GetOrders handles GET /api/user/orders.
func (h *Handler) GetOrders(c *gin.Context) {
	userID := middleware.UserIDFromCtx(c)

	resp, err := h.app.GetUserOrders(c.Request.Context(), userID)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}
