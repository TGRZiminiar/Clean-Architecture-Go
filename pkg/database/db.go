package database

import (
	"log"

	"github.com/TGRZiminiar/Clean-Architecture-Go/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func DbConnect(cfg *config.Config) *sqlx.DB {
	// Connect
	db, err := sqlx.Connect("pgx", cfg.Db.Url)
	if err != nil {
		log.Fatalf("connect to db failed: %v\n", err)
	}
	db.DB.SetMaxOpenConns(10)
	return db
}
