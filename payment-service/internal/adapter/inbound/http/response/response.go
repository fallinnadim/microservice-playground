package response

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func ErrorMsg(c *gin.Context, err error, status int) {
	if errors.Is(err, context.DeadlineExceeded) {
		c.AbortWithStatusJSON(http.StatusGatewayTimeout, ErrorResponse{
			Timestamp: time.Now(),
			Error:     "Request Timeout",
		})
		return
	}
	c.AbortWithStatusJSON(status, ErrorResponse{
		Timestamp: time.Now(),
		Error:     err.Error(),
	})
}

func OK[T any](c *gin.Context, data T) {
	c.JSON(200, SuccessResponse[T]{
		Timestamp: time.Now(),
		Data:      data,
	})
}

func Created[T any](c *gin.Context, data T) {
	c.JSON(201, SuccessResponse[T]{
		Timestamp: time.Now(),
		Data:      data,
	})
}
