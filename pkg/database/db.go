package database

import (
	"log"

	"github.com/TGRZiminiar/Clean-Architecture-Go/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func DbConnect(cfg *config.Config) *sqlx.DB {
	// Connect
	// host=%s port=%d user=%s password=%s dbname=%s sslmode=%s
	db, err := sqlx.Connect("pgx", "host=127.0.0.1 port=5432 user=mix dbname=test password=secret sslmode=disable")
	if err != nil {
		log.Fatalf("connect to db failed: %v\n", err)
	}
	db.DB.SetMaxOpenConns(10)
	return db
}
