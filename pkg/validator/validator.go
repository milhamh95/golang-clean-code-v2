package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

// Validate will validate struct based on tag
func Validate(data interface{}) error {
	validate := validator.New()

	if err := validate.Struct(data); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			newError := fmt.Errorf("error field validation for %s failed on the '%s' tag", e.Field(), e.Tag())
			return newError
		}
	}

	return nil
}
