package http

import (
	"net/http"
	"time"

	"github.com/wb-go/wbf/ginext"
)

type notifyHandler struct {
}

func NewNotifyHandler() Handler {
	return &notifyHandler{}
}

func (h *notifyHandler) Register(router *ginext.Engine) {
	router.POST("/notify", h.create)
	router.GET("/notify/{id}", h.get)
	router.DELETE("/notify/{id}", h.delete)
}

func (h *notifyHandler) create(c *ginext.Context) {
	var req CreateNotifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "invalid json",
		})
		return
	}

	if req.ScheduledAt.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "scheduled_at must be in the future",
		})
		return
	}
	
	// TODO: не дописано
}

func (h *notifyHandler) get(c *ginext.Context) {

}

func (h *notifyHandler) delete(c *ginext.Context) {

}
