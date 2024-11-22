package web

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

// Инициализация валидатора
func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

// Выполняет валидацию полей структуры
func ValidateStruct(ctx context.Context, input any) error {
	err := validate.StructCtx(ctx, input)
	if err == nil {
		return nil
	}
	data := ""
	if allErrs, ok := err.(validator.ValidationErrors); ok {
		for _, fld := range allErrs {
			data += fmt.Sprintf("field: %#v\n", fld.Error())
		}
	}
	return fmt.Errorf("error: %s\n%s", err.Error(), data)
}
