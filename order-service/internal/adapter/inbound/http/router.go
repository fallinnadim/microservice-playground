package http

import (
	"time"

	"github.com/gin-gonic/gin"
)

func NewRouter(h *Handler) *gin.Engine {
	r := gin.New()
	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			public := v1.Group("/")
			public.POST("/login", h.Login)
			public.POST("/register", h.Register)

			private := v1.Group("/")
			private.Use(
				Tracer(h.log),
				Timeout(5*time.Second),
				AuthRequired(h.authUC.JWTAdapter, h.log),
				RateLimit(h.rateLimitUC, h.log),
			)
			private.GET("/ping", h.Ping)
		}
	}

	return r
}
