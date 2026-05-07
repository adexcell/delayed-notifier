package httprouter

import (
	ver1 "github.com/adexcell/delayed-notifier/internal/notify/controller/http_router/v1"
	"github.com/adexcell/delayed-notifier/pkg/logger"
	"github.com/adexcell/delayed-notifier/pkg/metrics"
	"github.com/adexcell/delayed-notifier/pkg/otel"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/wb-go/wbf/ginext"
)

// NotifyRouter configures the HTTP routes for the notification service.
func NotifyRouter(r *ginext.Engine, uc ver1.Usecase, m *metrics.HTTPServer) {
	v1 := ver1.New(uc)

	// Expose metrics endpoint (separate from application routes)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// CORS configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:5500", // frontend
	}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	config.ExposeHeaders = []string{"Content-Length"}
	config.AllowCredentials = true // if cookies/auth are needed

	r.Use(cors.New(config))

	api := r.Group("api/v1/notifies")
	{
		api.Use(logger.Middleware())
		api.Use(metrics.Middleware(m))
		api.Use(otel.Middleware())

		api.POST("", v1.CreateNotify)
		api.GET("/status/:id", v1.GetNotifyStatus)
		api.GET("/:id", v1.GetNotify)
		api.DELETE("/:id", v1.DeleteNotify)
	}

}
