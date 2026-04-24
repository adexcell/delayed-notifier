package v1

import (
	"net/http"

	"github.com/wb-go/wbf/ginext"
)

type UseCase interface{}

type Handler struct {
	uc UseCase
}

func New(uc UseCase) *Handler {
	return &Handler{uc: uc}
}

func(h *Handler) Ping(c *ginext.Context) {
	c.JSON(http.StatusOK, ginext.H{"hello":"pong"})
}
