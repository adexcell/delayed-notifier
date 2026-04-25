package v1

import (
	"encoding/json"
	"net/http"

	"github.com/adexcell/delayed-notifier/internal/notify/dto"
	"github.com/wb-go/wbf/ginext"
)

func (h *Handler) UpdateNotify(c *ginext.Context) {
	var input dto.UpdateNotifyInput

	if err := json.NewDecoder(c.Request.Body).Decode(&input); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{"error": "invalid json"})
		return
	}

	err := h.usecase.UpdateNotify(c, input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ginext.H{"error": "request failed"})
		return
	}

	c.Status(http.StatusNoContent)
}
