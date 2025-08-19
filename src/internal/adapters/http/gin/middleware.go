package ginadp

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.GetHeader("X-Request-ID")
		if rid == "" {
			rid = uuid.NewString()
		}
		c.Writer.Header().Set("X-Request-ID", rid)
		c.Set("request_id", rid)
		c.Next()
	}
}

func Logger(log zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		dur := time.Since(start)
		ev := log.Info().
			Str("rid", c.GetString("request_id")).
			Str("method", c.Request.Method).
			Str("path", c.FullPath()).
			Int("status", c.Writer.Status()).
			Dur("duration", dur)
		if len(c.Errors) > 0 {
			ev = log.Error().
				Str("rid", c.GetString("request_id")).
				Str("method", c.Request.Method).
				Str("path", c.FullPath()).
				Int("status", c.Writer.Status()).
				Dur("duration", dur).
				Err(c.Errors.Last())
		}
		ev.Msg("http")
	}
}
