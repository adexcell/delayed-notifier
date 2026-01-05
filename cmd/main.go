package main

import (
	"net/http"

	transport "github.com/adexcell/delayed-notifier/internal/transport/http"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

func main() {
	zlog.Init()

	zlog.Logger.Info().Msg("create router")
	httprouter := ginext.New("debug")

	zlog.Logger.Info().Msg("register notify handler")
	notifyHandler := transport.NewNotifyHandler()
	notifyHandler.Register(httprouter)

	httprouter.GET("/", Hello)

	httprouter.Run()
}

func Hello(c *ginext.Context) {
	c.JSON(http.StatusOK, ginext.H{
			"msg": "hello_world",
		})
}
