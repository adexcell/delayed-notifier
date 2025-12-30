package router

import (
	"net/http"
	"strings"
	"time"

	"github.com/adexcell/delayed-notifier/pkg/auth"
	"github.com/adexcell/delayed-notifier/pkg/logger"

	"github.com/gin-gonic/gin"
)

func Logger(l *logger.Zerolog) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		l.Info().
			Str("method", c.Request.Method).
			Str("path", path).
			Str("raw", raw).
			Int("status", status).
			Dur("latency", latency).
			Msg("inbound request")
	}
}

func Auth(manager auth.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "empty auth header"})
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid auth header"})
			return
		}

		userID, err := manager.Parse(headerParts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		c.Set("userID", userID)

		c.Next()
	}
}
