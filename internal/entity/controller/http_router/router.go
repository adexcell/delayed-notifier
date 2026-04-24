package httprouter

import (
	ver1 "github.com/adexcell/delayed-notifier/internal/entity/controller/http_router/v1"
	"github.com/adexcell/delayed-notifier/pkg/logger"
	"github.com/adexcell/delayed-notifier/pkg/metrics"
	"github.com/adexcell/delayed-notifier/pkg/otel"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/wb-go/wbf/ginext"
)

// NewRouter creates and configures the HTTP router with middleware and routes.
func EntityRouter(r *ginext.Engine, uc ver1.UseCase, m *metrics.HTTPServer) {
	v1 := ver1.New(uc)

	// Expose metrics endpoint (separate from application routes)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	api := r.Group("api/entity")
	{
		api.Use(logger.Middleware())
		api.Use(metrics.NewMiddleware(m))
		api.Use(otel.Middleware())
		
		api.GET("/ping", v1.Ping)
	}

}
