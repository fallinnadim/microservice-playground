package http

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/fallinnadim/order-service/internal/adapter/inbound/http/response"
	"github.com/fallinnadim/order-service/internal/port/outbound"
	"github.com/fallinnadim/order-service/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Tracer(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := uuid.NewString()
		c.Set("trace_id", traceID)
		c.Writer.Header().Set("traceparent", traceID)
		c.Next()
		start := time.Now()
		c.Next()
		log.Info("incoming request",
			slog.String("trace_id", traceID),
			slog.String("path", c.Request.URL.Path),
			slog.Int("status", c.Writer.Status()),
			slog.Duration("duration", time.Since(start)),
		)
	}
}

func Timeout(d time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), d)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func AuthRequired(authUsecase outbound.JWTAuthPort, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			response.ErrorMsg(c, http.StatusUnauthorized, "missing authorization header")
			return
		}

		const prefix = "Bearer "
		if len(authHeader) < len(prefix) || authHeader[:len(prefix)] != prefix {
			response.ErrorMsg(c, http.StatusUnauthorized, "invalid authorization format")
			return
		}

		tokenStr := authHeader[len(prefix):]

		claims, err := authUsecase.ValidateToken(tokenStr)
		if err != nil {
			log.Warn("invalid token", "error", err)
			response.ErrorMsg(c, http.StatusUnauthorized, "invalid or expired token")
			return
		}

		c.Set("user_id", claims.UserID)

		c.Next()
	}
}

func RateLimit(uc *usecase.RateLimitUsecase, log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDVal, exists := c.Get("user_id")
		if !exists {
			log.Warn("missing user_id in context")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		userID, ok := userIDVal.(string)
		if !ok {
			log.Warn("invalid user_id type")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "invalid user context",
			})
			return
		}

		allowed, err := uc.Allow(c.Request.Context(), userID)
		if err != nil {
			log.Error("rate limit error", "err", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "internal error",
			})
			return
		}

		if !allowed {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			return
		}

		c.Next()
	}
}
