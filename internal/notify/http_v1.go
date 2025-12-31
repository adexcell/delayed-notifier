package notify

import (
	"net/http"

	"github.com/adexcell/delayed-notifier/internal/controllers"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

const (
	notifiesURL = "/notify"
	notifyURL   = "/notify/:id"
)

type handler struct {
}

func NewNotifyHandler() controllers.Handler {
	return &handler{}
}

func (h *handler) Register(router *ginext.Engine) {
	router.POST(notifiesURL, h.PostNotify)
	router.GET(notifyURL, h.GetNotify)
	router.DELETE(notifyURL, h.DeleteNotify)
}

func (h *handler) PostNotify(c *ginext.Context) {
	zlog.Logger.Info().Msg("post notify")
	c.JSON(http.StatusOK, ginext.H{
		"msg": "post notify",
	})
}

func (h *handler) GetNotify(c *ginext.Context) {
	var id int
	c.ShouldBindJSON(&id)

	zlog.Logger.Info().Msg("get notify")
	c.JSON(http.StatusOK, ginext.H{
		"msg": "get notify",
		"id":  id,
	})
}

func (h *handler) DeleteNotify(c *ginext.Context) {
	var id int
	c.ShouldBindJSON(&id)

	zlog.Logger.Info().Msg("delete notify")
	c.JSON(http.StatusOK, ginext.H{
		"msg": "delete notify",
		"id":  id,
	})
}
