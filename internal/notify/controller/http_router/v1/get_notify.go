package v1

import (
	"encoding/json"
	"net/http"

	"github.com/adexcell/delayed-notifier/internal/notify/dto"
	"github.com/wb-go/wbf/ginext"
)

func (h *Handler) GetNotify(c *ginext.Context) {
	var input dto.GetNotifyInput

	if err := json.NewDecoder(c.Request.Body).Decode(&input); err != nil {
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
