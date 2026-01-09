// Package router является оберткой над вспомогательным пакетом wbf/ginext.
package router

import (
	"github.com/wb-go/wbf/ginext"
)

type Router = ginext.Engine
type Context = ginext.Context
type H = ginext.H

type Handler interface {
	Register(router *Router)
}

type Config struct {
	GinMode string `mapstructure:"gin_mode"`
}

func NewRouter(cfg Config) *Router {
	return ginext.New(cfg.GinMode)
}
