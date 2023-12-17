package authusecase

import (
	"context"
	"fmt"

	"github.com/TGRZiminiar/Clean-Architecture-Go/config"
	"github.com/TGRZiminiar/Clean-Architecture-Go/modules/auth"
	authrepository "github.com/TGRZiminiar/Clean-Architecture-Go/modules/auth/authRepository"
	"github.com/TGRZiminiar/Clean-Architecture-Go/pkg/jwtauth"
)

type (
	AuthUsecaseService interface {
		CreateUser(cfg *config.Config, pctx context.Context, req *auth.CreateUser) (*auth.UserJson, error)
	}

	authUsecase struct {
		authRepo authrepository.AuthRepositoryService
	}
)

func NewAuthUsecase(authRepo authrepository.AuthRepositoryService) AuthUsecaseService {
	return &authUsecase{
		authRepo: authRepo,
	}
}

func (u *authUsecase) CreateUser(cfg *config.Config, pctx context.Context, req *auth.CreateUser) (*auth.UserJson, error) {

	userId, err := u.authRepo.CreateUser(pctx, req)
	if err != nil {
		return nil, err
	}

	userData, err := u.authRepo.FindUserById(pctx, userId)
	if err != nil {
		return nil, err
	}

	accessToken := jwtauth.NewAccessToken(cfg.Jwt.AccessSecretKey, cfg.Jwt.AccessDuration, &jwtauth.Claims{
		UserId:   userData.Id,
		Email:    userData.Email,
		Username: userData.Username,
	}).SignToken()

	fmt.Println(accessToken)

	return &auth.UserJson{
		Id:        userData.Id,
		Username:  userData.Username,
		Email:     userData.Email,
		Image:     userData.Image,
		CreatedAt: userData.CreatedAt,
		UpdatedAt: userData.UpdatedAt,
	}, nil
}
