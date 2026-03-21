package response

import (
	"time"

	"github.com/gin-gonic/gin"
)

func ErrorMsg(c *gin.Context, status int, message string) {
	c.AbortWithStatusJSON(status, ErrorResponse{
		Timestamp: time.Now(),
		Error:     message,
	})
}

func OK[T any](c *gin.Context, message string, data T) {
	c.JSON(200, SuccessResponse[T]{
		Timestamp: time.Now(),
		Data:      data,
	})
}

func Created[T any](c *gin.Context, message string, data T) {
	c.JSON(201, SuccessResponse[T]{
		Timestamp: time.Now(),
		Data:      data,
	})
}
