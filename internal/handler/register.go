package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	serviceerrors "github.com/vrnvgasu/gofemart/internal/app"
	"github.com/vrnvgasu/gofemart/internal/handler/response"
)

type registerRequest struct {
	Login    string `json:"login"    binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Register handles POST /api/user/register.
func (h *Handler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResponseError(c, serviceerrors.BadRequestError())
		return
	}

	token, err := h.app.Register(c, req.Login, req.Password)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	c.Header("Authorization", "Bearer "+token)
	c.Status(http.StatusOK)
}
