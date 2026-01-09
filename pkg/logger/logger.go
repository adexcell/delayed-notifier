// Package logger является оберткой над вспомогательным пакетом wbf/zlog.
package logger

import "github.com/wb-go/wbf/zlog"

type Log = zlog.Zerolog

func NewLogger() Log {
	zlog.InitConsole()
	return zlog.Logger
}
