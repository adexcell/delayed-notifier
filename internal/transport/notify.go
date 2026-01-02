package transport

import "github.com/wb-go/wbf/ginext"

type notifyHandler struct {

}

func NewNotifyHandler() Handler{
	return &notifyHandler{}
}

func (h *notifyHandler) Register(router *ginext.Engine) {
	router.POST("/notify", h.create)
	router.GET("/notify/{id}", h.get)
	router.DELETE("/notify/{id}", h.delete)
}

func (h *notifyHandler) create(c *ginext.Context) {

}

func (h *notifyHandler) get(c *ginext.Context) {

}

func (h *notifyHandler) delete(c *ginext.Context) {
	
}

