package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func LogMiddleware(c *gin.Context) {
	start := time.Now()


	// ALWAYS log (even if aborted)
	defer func() {
	// Duration
	duration := time.Since(start)
	var formattedDuration string
	switch {
	case duration < time.Millisecond:
		formattedDuration = duration.String()
	case duration < time.Second:
		formattedDuration = duration.Round(time.Microsecond).String()
	default:
		formattedDuration = duration.Round(time.Second).String()
	}

		log.Info().
			Str("method", c.Request.Method).
			Str("path", c.FullPath()).
			Int("status", c.Writer.Status()).
			Str("ip", c.ClientIP()).
			Str("duration", formattedDuration).
			Int("size_bytes", c.Writer.Size()).
			Msg("request finished")
	}()

	c.Next()
}