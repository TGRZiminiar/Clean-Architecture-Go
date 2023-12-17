package main

import (
	"context"
	"log"
	"os"

	"github.com/TGRZiminiar/Clean-Architecture-Go/config"
	"github.com/TGRZiminiar/Clean-Architecture-Go/pkg/database"
	"github.com/TGRZiminiar/Clean-Architecture-Go/server"
)

func main() {
	ctx := context.Background()

	cfg := config.LoadConfig(func() string {
		if len(os.Args) < 2 {
			log.Fatal("Error: .env path is invalid")
		}
		return os.Args[1]
	}())

	db := database.DbConnect(&cfg)
	defer db.Close()

	server.Start(ctx, &cfg, db)
}
