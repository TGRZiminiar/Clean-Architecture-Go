package authrepository

import (
	"context"
	"errors"
	"fmt"

	"github.com/TGRZiminiar/Clean-Architecture-Go/modules/auth"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type (
	AuthRepositoryService interface {
		CreateUser(pctx context.Context, user *auth.CreateUser) (*uuid.UUID, error)
		FindUserById(pctx context.Context, userId *uuid.UUID) (*auth.User, error)
	}

	authRepository struct {
		db *sqlx.DB
	}
)

func NewAuthRepository(db *sqlx.DB) AuthRepositoryService {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) CreateUser(pctx context.Context, user *auth.CreateUser) (*uuid.UUID, error) {
	query := `
    INSERT INTO "users" (
        "username",
        "email",
        "password"
    )
    VALUES ($1, $2, $3)
	RETURNING "id";
    `
	auth := new(auth.User)
	if err := r.db.QueryRowxContext(pctx, query, user.Username, user.Email, user.Password).Scan(&auth.Id); err != nil {
		// r.db.MustBegin().Tx.Rollback()
		// return nil, fmt.Errorf("insert user failed: %v", err.Error())
		return nil, errors.New("insert user repo failed")
	}
	fmt.Println(auth.Id)
	return &auth.Id, nil
}

func (r *authRepository) FindUserById(pctx context.Context, userId *uuid.UUID) (*auth.User, error) {
	query := `
	SELECT 
		"id",
		"username",
		"email",
		"created_at",
		"updated_at"
	FROM "users"
	WHERE "id" = $1;
	`

	userData := new(auth.User)
	if err := r.db.Get(userData, query, userId); err != nil {
		return nil, errors.New("find user by id failed")
	}

	return userData, nil

}
