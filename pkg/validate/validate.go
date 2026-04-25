package validate

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	inner *validator.Validate
}

var validate = validator.New(validator.WithRequiredStructEnabled()) //nolint:gochecknoglobals

func Validate(structure any) error {
	if err := validate.Struct(structure); err != nil {
		return fmt.Errorf("validate.Struct: %w", err)
	}

	return nil
}
