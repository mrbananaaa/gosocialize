package validator

import (
	"errors"
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

var (
	v    *validator.Validate
	once sync.Once
)

func get() *validator.Validate {
	once.Do(func() {
		v = validator.New()

		// register custom tag name (json instead of struct field)
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := fld.Tag.Get("json")
			if name == "-" {
				return ""
			}

			// remove ",omitempty"
			if idx := strings.Index(name, ","); idx != -1 {
				name = name[:idx]
			}

			return name
		})
	})

	return v
}

func ValidateStruct(s any) error {
	err := get().Struct(s)
	if err == nil {
		return nil
	}

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		return NewValidationError(validationErrors)
	}

	return err
}
