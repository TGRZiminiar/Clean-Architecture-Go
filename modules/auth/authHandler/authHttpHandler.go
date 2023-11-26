package authhandler

import (
	"github.com/TGRZiminiar/Clean-Architecture-Go/config"
	"github.com/TGRZiminiar/Clean-Architecture-Go/modules/auth"
	authusecase "github.com/TGRZiminiar/Clean-Architecture-Go/modules/auth/authUsecase"
	"github.com/gofiber/fiber/v2"
)

type (
	AuthHttpHandlerService interface{}

	authHttpHandler struct {
		cfg         *config.Config
		authUsecase authusecase.AuthUsecaseService
	}
)

func NewAuthHttpHandler(cfg *config.Config, authUsecase authusecase.AuthUsecaseService) AuthHttpHandlerService {
	return &authHttpHandler{
		cfg:         cfg,
		authUsecase: authUsecase,
	}
}

func (h *authHttpHandler) CreateUser(c *fiber.Ctx) error {
	req := new(auth.CreateUser)
	if err := c.BodyParser(req); err != nil {
		return c.JSON(map[string]string{
			"": "",
		})
	}
}
