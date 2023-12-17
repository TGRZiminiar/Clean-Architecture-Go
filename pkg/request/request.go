package request

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type (
	contextWrapperService interface {
		Validate(data interface{}) []ErrorResponse
		ParseJson(data interface{}) error
	}

	contextWrapper struct {
		Context   *fiber.Ctx
		validator *validator.Validate
	}

	ErrorResponse struct {
		Error       bool
		FailedField string
		Tag         string
		Value       interface{}
	}
)

func NewContextWrapper(ctx *fiber.Ctx) contextWrapperService {
	return &contextWrapper{
		Context:   ctx,
		validator: validator.New(),
	}
}

func (c *contextWrapper) ParseJson(data interface{}) error {
	if err := c.Context.BodyParser(&data); err != nil {
		return err
	}
	return nil
}

func (c *contextWrapper) Validate(data interface{}) []ErrorResponse {
	validationErrors := []ErrorResponse{}

	errs := c.validator.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			// In this case data object is actually holding the User struct
			var elem ErrorResponse

			elem.FailedField = err.Field() // Export struct field name
			elem.Tag = err.Tag()           // Export struct tag
			elem.Value = err.Value()       // Export field value
			elem.Error = true

			validationErrors = append(validationErrors, elem)
		}
	}

	return validationErrors
}
