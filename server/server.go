package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/TGRZiminiar/Clean-Architecture-Go/config"
	"github.com/TGRZiminiar/Clean-Architecture-Go/pkg/jwtauth"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type (
	server struct {
		app *fiber.App
		db  *sqlx.DB
		cfg *config.Config
	}
)

// func newMiddleware(cfg *config.Config) middlewarehandler.MiddlewareHandlerService {
// 	repo := middlewarerepository.NewMiddlewareRepository()
// 	usecase := middlewareusecase.NewMiddlewareUsecase(repo)
// 	return middlewarehandler.NewMiddlewareHandler(cfg, usecase)
// }

func (s *server) gracefulShutdown(pctx context.Context, quit <-chan os.Signal) {

	log.Printf("Starting service: %s", s.cfg.App.Name)

	<-quit

	log.Printf("Shutting down service: %s", s.cfg.App.Name)

	if err := s.app.Shutdown(); err != nil {
		log.Fatalf("Error: %v", err)
	}

}

func (s *server) httpListening() {
	if err := s.app.Listen(s.cfg.App.Url); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Error: %v", err)
	}
}

func Start(pctx context.Context, cfg *config.Config, db *sqlx.DB) {
	s := &server{
		db:  db,
		cfg: cfg,
		app: fiber.New(fiber.Config{
			AppName:      "test",
			BodyLimit:    10,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 20 * time.Second,
			JSONEncoder:  json.Marshal,
			JSONDecoder:  json.Unmarshal,
		}),
	}

	jwtauth.SetApiKey(cfg.Jwt.ApiSecretKey)

	// Body Limit

	switch s.cfg.App.Name {
	// case "auth":
	// 	s.authService()
	}

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go s.gracefulShutdown(pctx, quit)

	// Listening
	s.httpListening()

}
