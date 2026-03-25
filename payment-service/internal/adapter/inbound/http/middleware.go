package http

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

func Tracer(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		span := trace.SpanFromContext(c.Request.Context())
		traceID := span.SpanContext().TraceID().String()
		start := time.Now()

		c.Next()

		log.Info("finish request",
			slog.String("service_name", "payment-service"),
			slog.String("trace_id", traceID),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", c.Writer.Status()),
			slog.String("duration", time.Since(start).String()),
		)
	}
}

func TimeoutPropagation(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		timeoutStr := c.GetHeader("X-Timeout-Remaining")

		if timeoutStr != "" {
			ms, err := strconv.Atoi(timeoutStr)
			if err == nil && ms > 0 {
				var cancel context.CancelFunc

				ctx, cancel = context.WithTimeout(ctx, time.Duration(ms)*time.Millisecond)
				defer cancel()

				c.Request = c.Request.WithContext(ctx)

				log.Info("timeout propagated",
					"timeout_ms", ms,
				)
			}
		}

		c.Next()
	}
}
