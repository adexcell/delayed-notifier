package v1

import (
	"encoding/json"
	"net/http"

	"github.com/adexcell/delayed-notifier/internal/notify/dto"
	"github.com/wb-go/wbf/ginext"
)

func (h *Handler) CreateNotify(c *ginext.Context) {
	input := dto.CreateNotifyInput{}

	if err := json.NewDecoder(c.Request.Body).Decode(&input); err != nil {
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
