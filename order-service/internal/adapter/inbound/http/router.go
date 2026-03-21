package http

import "github.com/gin-gonic/gin"

func NewRouter(h *Handler) *gin.Engine {
	r := gin.New()

	r.GET("/ping",
		AuthRequired(h.authUC, h.log),
		RateLimit(h.rateLimitUC, h.log),
		h.Ping)

	return r
}
