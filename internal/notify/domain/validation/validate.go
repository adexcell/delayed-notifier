package validation

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Validator wraps the standard validator.Validate.
type Validator struct {
	inner *validator.Validate
}

var validate = validator.New(validator.WithRequiredStructEnabled()) //nolint:gochecknoglobals

// Validate performs struct validation using the global validator instance.
func Validate(structure any) error {
	if err := validate.Struct(structure); err != nil {
		return fmt.Errorf("validate.Struct: %w", err)
	}

	return nil
}
