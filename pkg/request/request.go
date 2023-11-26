package request

import (
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type (
	contextWrapperService interface {
		Bind(data any) error
	}

	contextWrapper struct {
		Context   *fiber.Ctx
		validator *validator.Validate
	}
)

func NewContextWrapper(ctx *fiber.Ctx) contextWrapperService {
	return &contextWrapper{
		Context:   ctx,
		validator: validator.New(),
	}
}
func (c *contextWrapper) Bind(data interface{}) error {

	dataMap, ok := data.(fiber.Map)
	if !ok {
		return fmt.Errorf("error: data must be of type fiber.Map")
	}

	if err := c.Context.Bind(dataMap); err != nil {
		log.Printf("Error: Bind data failed: %s", err.Error())
		return fmt.Errorf("error: bind data failed: %s", err.Error())
	}

	if err := c.validator.Struct(data); err != nil {
		log.Printf("Error: Validate data failed: %s", err.Error())
		return fmt.Errorf("error: validate data failed: %s", err.Error())
	}

	return nil
}
