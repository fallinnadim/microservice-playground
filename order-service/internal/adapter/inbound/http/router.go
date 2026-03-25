package http

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func NewRouter(h *Handler) *gin.Engine {
	r := gin.New()
	r.Use(otelgin.Middleware("order-service"))
	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			public := v1.Group("/")
			public.POST("/login", h.auth.Login)
			public.POST("/register", h.auth.Register)

			private := v1.Group("/")
			private.Use(
				Tracer(h.auth.Log),
				Timeout(5*time.Second),
				AuthRequired(h.auth.AuthUC.JWTAdapter, h.auth.Log),
				RateLimit(h.rateLimitUC, h.auth.Log),
			)
			private.POST("/order", h.order.Order)
		}
	}

	return r
}
