package controller

import (
	"errors"
	"net/http"
	"time"

	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/adexcell/delayed-notifier/pkg/log"
	"github.com/adexcell/delayed-notifier/pkg/router"
	"github.com/adexcell/delayed-notifier/pkg/utils/uuid"
)

const (
	postNotify        = "/notify"
	getOrDeleteNotify = "/notify/:id"
)

type notifyHandler struct {
	usecase domain.NotifyUsecase
	log     log.Log
}

func NewNotifyHandler(u domain.NotifyUsecase, l log.Log) router.Handler {
	return &notifyHandler{usecase: u, log: l}
}

func (h *notifyHandler) Register(router *router.Router) {
	router.POST(postNotify, h.Create)
	router.GET(getOrDeleteNotify, h.Get)
	router.DELETE(getOrDeleteNotify, h.Delete)
}

// TODO: добавить логирование
func (h *notifyHandler) Create(c *router.Context) {
	var dto NotifyControllerDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, router.H{
			"error": "invalid json",
		})
		return
	}

	dto.ID = uuid.New()

	if dto.ScheduledAt.Before(time.Now()) {
		c.JSON(http.StatusUnprocessableEntity, router.H{
			"error": "scheduled_at in the past",
		})
		return
	}

	n := toDomain(dto)
	id, err := h.usecase.Save(c, n)
	if err != nil {
		if errors.Is(err, domain.ErrNotifyAlreadyExisis) {
			c.JSON(http.StatusConflict, router.H{
				"error": "notify already exists",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, router.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, router.H{"id": id})
}

func (h *notifyHandler) Get(c *router.Context) {
	id := c.Param("id")
	if err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, router.H{
			"error": "cannot parse ID",
		})
		return
	}

	notify, err := h.usecase.GetByID(c, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusNotFound, router.H{
				"error": "not found notify",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, router.H{
			"error": err.Error(),
		})
		return
	}

	res := toResponse(notify)

	c.JSON(http.StatusOK, res)
}

func (h *notifyHandler) Delete(c *router.Context) {
	id := c.Param("id")
	if err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, router.H{
			"error": "cannot parse ID",
		})
		return
	}

	if err := h.usecase.Delete(c, id); err != nil {
		if !errors.Is(err, domain.ErrNotFound) {
			c.JSON(http.StatusInternalServerError, router.H{
				"error": err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusNoContent, router.H{
		"success": "successfully deleted",
	})
}
