package v1

import (
	"net/http"

	"github.com/adexcell/delayed-notifier/internal/notify/controller/http_router/response"
	"github.com/adexcell/delayed-notifier/internal/notify/dto"
	"github.com/wb-go/wbf/ginext"
)

// GetNotifyStatus handles the GET request to retrieve the status of a specific notification.
func (h *Handler) GetNotifyStatus(c *ginext.Context) {
	input := dto.GetStatusInput{
		ID: c.Param("id"),
	}

	output, err := h.usecase.GetStatus(c.Request.Context(), input)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, output)
}
