package v1

import (
	"net/http"

	"github.com/wb-go/wbf/ginext"

	"github.com/adexcell/delayed-notifier/internal/notify/dto"
)

// DeleteNotify DELETE /notify/{id} — отмена запланированного уведомления.
func (h *Handler) DeleteNotify(c *ginext.Context) {
	input := dto.DeleteNotifyInput{
		ID: c.Param("id"),
	}

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid json"})
		return
	}

	err = h.usecase.DeleteNotify(c, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "request failed"})
		return
	}

	c.Status(http.StatusNoContent)
}
