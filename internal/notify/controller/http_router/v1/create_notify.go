package v1

import (
	"errors"
	"net/http"

	"github.com/adexcell/delayed-notifier/internal/notify/controller/http_router/response"
	"github.com/adexcell/delayed-notifier/internal/notify/domain/validation"
	"github.com/adexcell/delayed-notifier/internal/notify/dto"
	"github.com/go-playground/validator/v10"
	"github.com/wb-go/wbf/ginext"
)

// CreateNotify POST /notify — создание уведомлений с датой и временем отправки.
func (h *Handler) CreateNotify(c *ginext.Context) {
	input := dto.CreateNotifyInput{}

	err := c.ShouldBindJSON(&input)
	if err != nil {
		var verrs validator.ValidationErrors
		if errors.As(err, &verrs) {
			errMap := validation.ExtractErrors(verrs)
			response.ValidationsError(c, errMap)

			return
		}

		response.BadRequest(c, err)
		return
	}

	output, err := h.usecase.CreateNotify(c.Request.Context(), input)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	c.JSON(http.StatusCreated, output)
}
