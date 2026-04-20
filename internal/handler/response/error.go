package response

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	serviceerrors "github.com/vrnvgasu/gofemart/internal/app"
	"github.com/vrnvgasu/gofemart/internal/logger"
)

type Error struct {
	Code        string `json:"code"`
	HTTPCode    int    `json:"http_code"`
	UserMessage string `json:"user_message"`
}

func ResponseError(c *gin.Context, err error) {
	if err != nil {
		_ = c.Error(err)
	}

	var serviceError *serviceerrors.ServiceError

	switch {
	case errors.As(err, &serviceError):
		parseServiceError(c, serviceError)
	default:
		logger.Log.Errorw("http request",
			"uri", c.Request.RequestURI,
			"method", c.Request.Method,
			"error", err,
		)
		c.AbortWithStatusJSON(http.StatusInternalServerError, Error{
			HTTPCode:    http.StatusInternalServerError,
			UserMessage: "Unhandled error",
		})
	}
}

func parseServiceError(c *gin.Context, err *serviceerrors.ServiceError) {
	c.AbortWithStatusJSON(err.HTTPCode, Error{
		Code:        string(err.Type),
		HTTPCode:    err.HTTPCode,
		UserMessage: err.Message,
	})
}
