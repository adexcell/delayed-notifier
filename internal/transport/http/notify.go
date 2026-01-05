package http

import (
	"errors"
	"net/http"
	"time"

	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/google/uuid"
	"github.com/wb-go/wbf/ginext"
)

type notifyHandler struct {
	NotifyUsecase domain.NotifyUsecase
}

func NewNotifyHandler() Handler {
	return &notifyHandler{}
}

func (h *notifyHandler) Register(router *ginext.Engine) {
	router.POST("/notify", h.Create)
	router.GET("/notify/{id}", h.Get)
	router.DELETE("/notify/{id}", h.Delete)
}

// TODO: добавить логирование
func (h *notifyHandler) Create(c *ginext.Context) {
	var dto NotifyTransportDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "invalid json",
		})
		return
	}

	dto.ID = uuid.New()

	if dto.ScheduledAt.Before(time.Now()) {
		c.JSON(http.StatusUnprocessableEntity, ginext.H{
			"error": "scheduled_at in the past",
		})
		return
	}

	n := toDomain(dto)
	id, err := h.NotifyUsecase.Save(c, n)
	if err != nil {
		if errors.Is(err, domain.ErrNotifyAlreadyExisis) {
			c.JSON(http.StatusConflict, ginext.H{
				"error": "notify already exists",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ginext.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, ginext.H{"id": id})
}

func (h *notifyHandler) Get(c *ginext.Context) {
	stringID := c.Param("id")
	id, err := uuid.Parse(stringID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "cannot parse ID",
		})
		return
	}

	notify, err := h.NotifyUsecase.GetByID(c, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFoundNotify) {
			c.JSON(http.StatusNotFound, ginext.H{
				"error": "not found notify",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ginext.H{
			"error": err.Error(),
		})
		return
	}

	res := toResponse(notify)

	c.JSON(http.StatusOK, res)
}

func (h *notifyHandler) Delete(c *ginext.Context) {
	stringID := c.Param("id")
	id, err := uuid.Parse(stringID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginext.H{
			"error": "cannot parse ID",
		})
		return
	}

	if err := h.NotifyUsecase.Delete(c, id); err != nil {
		if !errors.Is(err, domain.ErrNotFoundNotify) {
			c.JSON(http.StatusInternalServerError, ginext.H{
				"error": err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusNoContent, ginext.H{
		"success": "successfully deleted",
	})
}
