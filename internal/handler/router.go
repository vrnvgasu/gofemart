package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/vrnvgasu/gofemart/internal/config"
	"github.com/vrnvgasu/gofemart/internal/handler/middleware"
)

func NewRouter(h *Handler, cnf *config.Config) *gin.Engine {
	r := gin.New()
	r.Use(middleware.Logger())
	r.Use(middleware.Gzip())

	api := r.Group("/api/user")
	{
		api.POST("/register", h.Register)
		api.POST("/login", h.Login)

		auth := api.Group("")
		auth.Use(middleware.Auth(cnf.JWTSecret))
		{
			auth.POST("/orders", h.UploadOrder)
			auth.GET("/orders", h.GetOrders)
			auth.GET("/balance", h.GetBalance)
			auth.POST("/balance/withdraw", h.SaveWithdraw)
			auth.GET("/withdrawals", h.GetWithdrawals)
		}
	}

	return r
}
