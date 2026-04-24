package otel

import (
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	semconv "go.opentelemetry.io/otel/semconv/v1.30.0"

	"github.com/wb-go/wbf/ginext"
	"github.com/adexcell/delayed-notifier/pkg/otel/tracer"
)

func Middleware() ginext.HandlerFunc {
	return func(c *ginext.Context) {
		// Извлекаем контекст из заголовков запроса
		ctx := otel.GetTextMapPropagator().Extract(c, propagation.HeaderCarrier(c.Request.Header))

		// Создаем корневой span
		ctx, span := tracer.Start(ctx, "", trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		// Получаем статус-код
		statusCode := c.Writer.Status()

		// Вызываем следующий обработчик (или сам handler)
		c.Next()

		// Записываем полезные атрибуты
		span.SetAttributes(
			semconv.HTTPResponseStatusCode(statusCode),
			semconv.HTTPRequestMethodKey.String(c.Request.Method),
			semconv.HTTPRoute(c.Request.URL.Path),
		)

		// Помечаем span как ошибочный для 4xx и 5xx статусов
		if statusCode >= http.StatusBadRequest {
			span.SetStatus(codes.Error, http.StatusText(statusCode))
			span.AddEvent("error", trace.WithAttributes(
				attribute.String("error.message", http.StatusText(statusCode)),
			))
		}
	}

}
