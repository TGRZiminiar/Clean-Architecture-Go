package authrepository

import (
	"context"
	"log"

	"github.com/TGRZiminiar/Clean-Architecture-Go/modules/auth"
	"github.com/jmoiron/sqlx"
)

type (
	AuthRepositoryService interface {
		CreateUser(pctx context.Context) (*auth.User, error)
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

func (r *authRepository) CreateUser(pctx context.Context) (*auth.User, error) {
	query := `
	INSERT INTO "users" (
		"id",
		"username",
		"email",
		"password",
		"created_at",
		"updated_at",
	)
	VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING "id", "username", "email", "created_at", "updated_at";`
	auth := new(auth.User)
	_, err := r.db.NamedExecContext(pctx, query, auth)
	if err != nil {
		log.Printf("Error: CreateUser Failed: %s", err.Error())
		return nil, err
	}
	return auth, nil
}
