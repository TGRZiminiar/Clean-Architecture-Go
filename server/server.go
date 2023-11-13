package server

// type (
// 	server struct {
// 		app        *echo.Echo
// 		db         *mongo.Client
// 		cfg        *config.Config
// 		middleware middlewarehandler.MiddlewareHandlerService
// 	}
// )

// func newMiddleware(cfg *config.Config) middlewarehandler.MiddlewareHandlerService {
// 	repo := middlewarerepository.NewMiddlewareRepository()
// 	usecase := middlewareusecase.NewMiddlewareUsecase(repo)
// 	return middlewarehandler.NewMiddlewareHandler(cfg, usecase)
// }

// func (s *server) gracefulShutdown(pctx context.Context, quit <-chan os.Signal) {

// 	log.Printf("Starting service: %s", s.cfg.App.Name)

// 	<-quit

// 	log.Printf("Shutting down service: %s", s.cfg.App.Name)

// 	ctx, cancel := context.WithTimeout(pctx, 10*time.Second)
// 	defer cancel()

// 	if err := s.app.Shutdown(ctx); err != nil {
// 		log.Fatalf("Error: %v", err)
// 	}

// }

// func (s *server) httpListening() {
// 	if err := s.app.Start(s.cfg.App.Url); err != nil && err != http.ErrServerClosed {
// 		log.Fatalf("Error: %v", err)
// 	}
// }

// func Start(pctx context.Context, cfg *config.Config, db *mongo.Client) {
// 	s := &server{
// 		app:        echo.New(),
// 		db:         db,
// 		cfg:        cfg,
// 		middleware: newMiddleware(cfg),
// 	}

// 	jwtauth.SetApiKey(cfg.Jwt.ApiSecretKey)

// 	// Basic Middleware

// 	// Request Timeout
// 	s.app.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
// 		Skipper:      middleware.DefaultSkipper,
// 		Timeout:      30 * time.Second,
// 		ErrorMessage: "Error: Request Timeout",
// 	}))

// 	// Cors
// 	s.app.Use(middleware.CORSWithConfig(middleware.CORSConfig{
// 		Skipper:      middleware.DefaultSkipper,
// 		AllowOrigins: []string{"*"},
// 		AllowMethods: []string{echo.GET, echo.PATCH, echo.POST, echo.PUT, echo.DELETE},
// 	}))

// 	// Body Limit
// 	s.app.Use(middleware.BodyLimit("10M"))

// 	switch s.cfg.App.Name {
// 	case "auth":
// 		s.authService()
// 	}

// 	// Graceful Shutdown
// 	quit := make(chan os.Signal, 1)
// 	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

// 	s.app.Use(middleware.Logger())

// 	go s.gracefulShutdown(pctx, quit)

// 	// Listening
// 	s.httpListening()

// }