package controller

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/adexcell/delayed-notifier/pkg/log"
	"github.com/adexcell/delayed-notifier/pkg/router"
	"github.com/adexcell/delayed-notifier/pkg/utils/uuid"
)

const (
	Notify   = "/notify"     // POST, GET
	NotifyID = "/notify/:id" // GET, DELETE
)

type notifyHandler struct {
	usecase domain.NotifyUsecase
	log     log.Log
}

func NewNotifyHandler(u domain.NotifyUsecase, l log.Log) router.Handler {
	return &notifyHandler{usecase: u, log: l}
}

func (h *notifyHandler) Register(router *router.Router) {
	router.POST(Notify, h.Create)
	router.GET(NotifyID, h.Get)
	router.DELETE(NotifyID, h.Delete)
	router.GET(Notify, h.List)
}

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
		h.log.Info().Time("scheduled_at", dto.ScheduledAt).Msg("wrong time in scheduled at")
		c.JSON(http.StatusUnprocessableEntity, router.H{
			"error": "scheduled_at in the past",
		})
		return
	}

	n := toDomain(dto)
	id, err := h.usecase.Save(c, n)
	if err != nil {
		if errors.Is(err, domain.ErrNotifyAlreadyExists) {
			h.log.Error().Err(err).Msg("notify already exists")
			c.JSON(http.StatusConflict, router.H{
				"error": "notify already exists",
			})
			return
		}
		h.log.Error().Err(err).Msg("internal server error")
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
		h.log.Error().Err(err).Msg("wrong ID format")
		c.JSON(http.StatusBadRequest, router.H{
			"error": "cannot parse ID",
		})
		return
	}

	notify, err := h.usecase.GetByID(c, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			h.log.Error().Err(err).Msg("not found notify")
			c.JSON(http.StatusNotFound, router.H{
				"error": "not found notify",
			})
			return
		}
		h.log.Error().Err(err).Msg("internal server error")
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
		h.log.Error().Err(err).Msg("wrong ID format")
		c.JSON(http.StatusBadRequest, router.H{
			"error": "cannot parse ID",
		})
		return
	}

	if err := h.usecase.Delete(c, id); err != nil {
		if !errors.Is(err, domain.ErrNotFound) {
			h.log.Error().Err(err).Msg("internal server error")
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

func (h *notifyHandler) List(c *router.Context) {
	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		limit = 50
	}
	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		offset = 0
	}
	notifies, err := h.usecase.List(c, limit, offset)
	if err != nil {
		h.log.Error().Err(err).Msg("internal server error")
		c.JSON(http.StatusInternalServerError, router.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, notifies)
}
