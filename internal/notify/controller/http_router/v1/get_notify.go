package v1

import (
	"net/http"

	"github.com/wb-go/wbf/ginext"

	"github.com/adexcell/delayed-notifier/internal/notify/dto"
)

// GetNotify GET /notify/{id} — получение статуса уведомления;.
func (h *Handler) GetNotify(c *ginext.Context) {
	input := dto.GetNotifyInput{
		ID: c.Param("id"),
	}

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid json"})
		return
	}

	output, err := h.usecase.GetNotify(c, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "request failed"})
		return
	}

	c.JSON(http.StatusOK, output)
}
