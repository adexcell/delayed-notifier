package logger

import (
	"github.com/rs/zerolog/log"
	"github.com/wb-go/wbf/ginext"
)

func Middleware() ginext.HandlerFunc {
	return func(c *ginext.Context) {
		statusCode := c.Writer.Status()

		c.Next()

		method := c.Request.Method + " " + c.FullPath()

		log.Info().
			Str("proto", "http").
			Int("code", statusCode).
			Str("method", method).
			Send()
	}
}
