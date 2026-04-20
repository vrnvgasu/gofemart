package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	serviceerrors "github.com/vrnvgasu/gofemart/internal/app"
	"github.com/vrnvgasu/gofemart/internal/handler/response"
)

type loginRequest struct {
	Login    string `json:"login"    binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login handles POST /api/user/login.
func (h *Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ResponseError(c, serviceerrors.BadRequestError())
		return
	}

	token, err := h.app.Login(c, req.Login, req.Password)
	if err != nil {
		response.ResponseError(c, err)
		return
	}

	c.Header("Authorization", "Bearer "+token)
	c.Status(http.StatusOK)
}
