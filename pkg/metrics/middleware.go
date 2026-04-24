package metrics

import (
	"time"

	"github.com/wb-go/wbf/ginext"
)

func NewMiddleware(metrics *HTTPServer) ginext.HandlerFunc {
	return func(c *ginext.Context) {
		// 1. Засекаем время
		start := time.Now()

		// 2. Обрабатываем запрос (вызываем следующие хендлеры в цепочке Gin)
		c.Next()

		// 3. После обработки (c.Next() вернул управление)

		// Получаем метод и путь
		method := c.Request.Method + " " + c.FullPath()

		// Получаем статус-код
		statusCode := c.Writer.Status()

		// 4. Записываем метрики
		metrics.Duration(method, start)
		metrics.TotalInc(method, statusCode)
	}
}
