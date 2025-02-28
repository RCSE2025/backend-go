package logger

import (
	"time"

	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"log/slog"
)

const LoggerKey = "logger"

// New returns a Gin middleware function that logs HTTP requests
func New(log *slog.Logger) gin.HandlerFunc {
	log = log.With(
		slog.String("component", "middleware/logger"),
	)

	log.Info("Logger middleware enabled")

	return func(c *gin.Context) {
		start := time.Now()

		// Extract request ID if available
		reqID := requestid.Get(c)
		if reqID == "" {
			reqID = "unknown"
		}

		entry := log.With(
			slog.String("method", c.Request.Method),
			slog.String("path", c.Request.URL.Path),
			slog.String("remote_addr", c.ClientIP()),
			slog.String("user_agent", c.Request.UserAgent()),
			slog.String("request_id", reqID),
		)

		c.Set(LoggerKey, entry)

		// Process the request
		c.Next()

		// Log after response is sent
		entry.Info("request completed",
			slog.Int("status", c.Writer.Status()),
			slog.Int("bytes", c.Writer.Size()),
			slog.String("duration", time.Since(start).String()),
		)
	}
}

func FromContext(c *gin.Context) *slog.Logger {
	return c.MustGet(LoggerKey).(*slog.Logger)
}
