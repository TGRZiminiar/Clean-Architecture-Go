package authhandler

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/TGRZiminiar/Clean-Architecture-Go/config"
	"github.com/TGRZiminiar/Clean-Architecture-Go/modules/auth"
	authusecase "github.com/TGRZiminiar/Clean-Architecture-Go/modules/auth/authUsecase"
	"github.com/TGRZiminiar/Clean-Architecture-Go/pkg/request"
	"github.com/TGRZiminiar/Clean-Architecture-Go/pkg/response"
	"github.com/gofiber/fiber/v2"
)

type (
	AuthHttpHandlerService interface {
		CreateUser(c *fiber.Ctx) error
	}

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

	ctx := context.Background()

	wrapper := request.NewContextWrapper(c)

	req := new(auth.CreateUser)
	if err := wrapper.ParseJson(req); err != nil {
		return response.ErrorRes(c, http.StatusBadRequest, err.Error())
	}

	if errs := wrapper.Validate(req); len(errs) > 0 && errs[0].Error {
		errMsgs := make([]string, 0)

		for _, err := range errs {
			errMsgs = append(errMsgs, fmt.Sprintf(
				"[%s]: '%v' | Needs to implement '%s'",
				err.FailedField,
				err.Value,
				err.Tag,
			))
		}
		return response.ErrorRes(c, http.StatusBadRequest, strings.Join(errMsgs, " and "))
	}

	token, user, err := h.authUsecase.CreateUser(h.cfg, ctx, req)
	if err != nil {
		return response.ErrorRes(c, http.StatusBadRequest, err.Error())
	}

	oneWeek := time.Now().Add(7 * 24 * time.Hour)

	c.Cookie(&fiber.Cookie{
		Name:     "accessToken",
		Value:    token,
		HTTPOnly: true,
		Secure:   true,
		Expires:  oneWeek,
	})
	c.Cookie(&fiber.Cookie{
		Name:     "accessChecker",
		Value:    user.Id.String(),
		HTTPOnly: false,
		Secure:   false,
		Expires:  oneWeek,
	})

	return response.SuccessRes(c, 201, user)
}
