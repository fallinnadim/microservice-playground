package http

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func NewRouter(h *Handler) *gin.Engine {
	r := gin.New()
	r.Use(otelgin.Middleware("payment-service"))
	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			payment := v1.Group("/payment")
			payment.Use(Tracer(h.Log))
			payment.Use(TimeoutPropagation(h.Log))
			payment.POST("/", h.Payment)
		}
	}

	return r
}
