package logger

import (
	"github.com/rs/zerolog/log"
	"github.com/wb-go/wbf/ginext"
)

func Middleware() ginext.HandlerFunc {
	return func(c *ginext.Context) {
		c.Next()

		statusCode := c.Writer.Status()

		method := c.Request.Method + " " + c.FullPath()

		log.Info().
			Int("code", statusCode).
			Str("method", method).
			Send()
	}
}
