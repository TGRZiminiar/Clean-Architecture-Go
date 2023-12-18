package authusecase

import (
	"context"

	"github.com/TGRZiminiar/Clean-Architecture-Go/config"
	"github.com/TGRZiminiar/Clean-Architecture-Go/modules/auth"
	authrepository "github.com/TGRZiminiar/Clean-Architecture-Go/modules/auth/authRepository"
	"github.com/TGRZiminiar/Clean-Architecture-Go/pkg/jwtauth"
)

type (
	AuthUsecaseService interface {
		CreateUser(cfg *config.Config, pctx context.Context, req *auth.CreateUser) (string, *auth.UserJson, error)
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

func (u *authUsecase) CreateUser(cfg *config.Config, pctx context.Context, req *auth.CreateUser) (string, *auth.UserJson, error) {

	userId, err := u.authRepo.CreateUser(pctx, req)
	if err != nil {
		return "", nil, err
	}

	userData, err := u.authRepo.FindUserById(pctx, userId)
	if err != nil {
		return "", nil, err
	}

	accessToken, err := jwtauth.NewAccessToken(
		&cfg.Jwt,
		&jwtauth.Claims{
			UserId:   userData.Id,
			Email:    userData.Email,
			Username: userData.Username,
		}, cfg.Jwt.AccessDuration, "accessToken").SignToken()
	if err != nil {
		return "", nil, err
	}

	// parse, _ := jwtauth.ParseToken(accessToken, &cfg.Jwt)
	// fmt.Println(parse)

	return accessToken,
		&auth.UserJson{
			Id:        userData.Id,
			Username:  userData.Username,
			Email:     userData.Email,
			Image:     userData.Image,
			UpdatedAt: userData.UpdatedAt,
			CreatedAt: userData.CreatedAt,
		}, nil

}
