package http

import "github.com/wb-go/wbf/ginext"

type Handler interface {
	Register(router *ginext.Engine)
}
