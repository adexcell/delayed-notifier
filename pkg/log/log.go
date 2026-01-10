// Package log является оберткой над вспомогательным пакетом wbf/zlog.
package log

import "github.com/wb-go/wbf/zlog"

type Log = zlog.Zerolog

func New() Log {
	zlog.InitConsole()
	return zlog.Logger
}
