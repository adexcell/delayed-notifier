package httprouter

import (
	ver1 "github.com/adexcell/delayed-notifier/internal/notify/controller/http_router/v1"
	"github.com/adexcell/delayed-notifier/pkg/logger"
	"github.com/adexcell/delayed-notifier/pkg/metrics"
	"github.com/adexcell/delayed-notifier/pkg/otel"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/wb-go/wbf/ginext"
)

// NewRouter creates and configures the HTTP router with middleware and routes.
func NotifyRouter(r *ginext.Engine, uc ver1.Usecase, m *metrics.HTTPServer) {
	v1 := ver1.New(uc)

	// Expose metrics endpoint (separate from application routes)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	api := r.Group("api/v1/notifies")
	{
		api.Use(logger.Middleware())
		api.Use(metrics.Middleware(m))
		api.Use(otel.Middleware())

		api.POST("", v1.CreateNotify)
		api.GET("/:id", v1.GetNotify)
		api.DELETE("/:id", v1.DeleteNotify)
		api.PUT("/:id", v1.UpdateNotify)

	}

}
