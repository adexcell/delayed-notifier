// Package uuid является оберткой над вспомогательным пакетом wbf/helpers.
package uuid

import (
	"github.com/wb-go/wbf/helpers"
)

// New создает новый случайный UUID.
func New() string {
	return helpers.CreateUUID()
}

// Parse проверяет, является ли строка валидным UUID.
func Parse(s string) error {
	return helpers.ParseUUID(s)
}
