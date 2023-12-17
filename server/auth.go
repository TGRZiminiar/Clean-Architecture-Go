package server

import (
	authhandler "github.com/TGRZiminiar/Clean-Architecture-Go/modules/auth/authHandler"
	authrepository "github.com/TGRZiminiar/Clean-Architecture-Go/modules/auth/authRepository"
	authusecase "github.com/TGRZiminiar/Clean-Architecture-Go/modules/auth/authUsecase"
)

func (s *server) authService() {
	authRepo := authrepository.NewAuthRepository(s.db)
	authUsecase := authusecase.NewAuthUsecase(authRepo)
	authHttpHandler := authhandler.NewAuthHttpHandler(s.cfg, authUsecase)

	auth := s.app.Group("/auth_v1")

	auth.Post("/register", authHttpHandler.CreateUser)

}
