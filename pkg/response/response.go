package response

import (
	"github.com/gofiber/fiber/v2"
)

type ()

func ErrorRes(c *fiber.Ctx, status int, msg string) error {
	return c.Status(status).JSON(map[string]string{
		"msg": msg,
	})
}

func SuccessRes(c *fiber.Ctx, status int, data any) error {
	return c.Status(status).JSON(data)
}
