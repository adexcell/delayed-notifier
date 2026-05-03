package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Явная карта сообщений: struct field -> tag -> custom message.
var customMessages = map[string]map[string]string{
	"Channel": {"oneof": "Укажите email или telegram"},
	"Subject": {
		"required": "Поле обязательно для заполнения",
		"max=255":  "Не более 255 символов",
	},
	"ScheduledAt": {"gt": "Дата должна быть в будущем"},
}

// ExtractErrors преобразует валидационные ошибки в map[field]message.
func ExtractErrors(errs validator.ValidationErrors) map[string]string {
	out := make(map[string]string, len(errs))

	for _, e := range errs {
		field := e.Field() // Имя поля в структуре (например, "Name")
		tag := e.Tag()     // Название тега (например, "min")

		// 1. Ищем в явной карте
		if fieldMsgs, ok := customMessages[field]; ok {
			if msg, ok := fieldMsgs[tag]; ok {
				out[field] = msg

				continue
			}
		}

		// 2. Fallback: генерируем сообщение динамически или берём дефолтное
		out[field] = fallbackMessage(e)
	}

	return out
}

// fallbackMessage создаёт адекватное сообщение, если для тега нет кастомного.
func fallbackMessage(e validator.FieldError) string {
	param := e.Param() // Значение из тега (например, "3" для min=3)

	switch e.Tag() {
	case "required":
		return "Поле обязательно для заполнения"
	case "min":
		return "Значение должно быть не менее " + param
	case "max":
		return "Значение не должно превышать " + param
	case "gte", "lte":
		return fmt.Sprintf("Значение должно быть в допустимом диапазоне (параметр: %s)", param)
	case "oneof":
		return "Допустимые значения: " + strings.Join(strings.Split(param, " "), ", ")
	default:
		return "Ошибка валидации поля"
	}
}
