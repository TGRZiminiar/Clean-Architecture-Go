package authusecase

import (
	"context"

	"github.com/TGRZiminiar/Clean-Architecture-Go/modules/auth"
	authrepository "github.com/TGRZiminiar/Clean-Architecture-Go/modules/auth/authRepository"
)

type (
	AuthUsecaseService interface {
		CreateUser(pctx context.Context, req *auth.CreateUser) (*auth.UserJson, error)
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

func (u *authUsecase) CreateUser(pctx context.Context, req *auth.CreateUser) (*auth.UserJson, error) {

	return nil, nil
}
