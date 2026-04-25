package v1

import (
	"net/http"

	"github.com/wb-go/wbf/ginext"

	"github.com/adexcell/delayed-notifier/internal/notify/dto"
)

// CreateNotify POST /notify — создание уведомлений с датой и временем отправки.
func (h *Handler) CreateNotify(c *ginext.Context) {
	input := dto.CreateNotifyInput{}

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "json decode error"})
		return
	}

	output, err := h.usecase.CreateNotify(c, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "request failed"})
		return
	}

	c.JSON(http.StatusOK, output)
}
