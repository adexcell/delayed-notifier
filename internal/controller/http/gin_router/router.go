package ginrouter

import (
	"github.com/adexcell/delayed-notifier/internal/domain"
	"github.com/adexcell/delayed-notifier/internal/usecase/notify/schedule"
	"github.com/adexcell/delayed-notifier/internal/usecase/user/login"
	"github.com/adexcell/delayed-notifier/internal/usecase/user/register"
	"github.com/adexcell/delayed-notifier/pkg/auth"
	"github.com/adexcell/delayed-notifier/pkg/logger"
	"github.com/adexcell/delayed-notifier/pkg/router"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	user         domain.User
	notify       domain.Notify
	tokenManager auth.TokenManager
	l            *logger.Zerolog
}

func NewHandler(us domain.User, ns domain.Notify, tm auth.TokenManager, l *logger.Zerolog) *Handler {
	return &Handler{
		user:         us,
		notify:       ns,
		tokenManager: tm,
		l:            l,
	}
}

// InitRoutes — теперь это единственное место, где живут пути твоего API
func (h *Handler) InitRoutes() *gin.Engine {
	r := router.New(h.l)

	// Раздача статики
	r.Router.StaticFile("/", "./static/index.html")

	// Группировка роутов
	auth := r.Router.Group("/auth")
	{
		auth.POST("/register", register.HTTPv1)
		auth.POST("/login", login.HTTPv1)
	}

	api := r.Router.Group("/api", router.Auth(h.tokenManager))
	{
		api.POST("/notification", schedule.HTTPv1)
	}

	return r.Router

}
