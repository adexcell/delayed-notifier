package v1

import (
	"net/http"

	"github.com/wb-go/wbf/ginext"

	"github.com/adexcell/delayed-notifier/internal/notify/controller/http_router/response"
	"github.com/adexcell/delayed-notifier/internal/notify/dto"
)

// DeleteNotify handles the DELETE request to cancel a scheduled notification.
func (h *Handler) DeleteNotify(c *ginext.Context) {
	input := dto.DeleteNotifyInput{
		ID: c.Param("id"),
	}

	err := h.usecase.DeleteNotify(c.Request.Context(), input)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}
