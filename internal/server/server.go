package server

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gvevstigneev99/myProject/internal/handler"
	"github.com/gvevstigneev99/myProject/internal/middleware"
	"golang.org/x/exp/slog"

	"paraklitshop/internal/config"
)

type Server struct {
	app    *fiber.App
	cfg    *config.Config
	logger *slog.Logger
}

func NewServer(cfg *config.Config, logger *slog.Logger) *Server {
	app := fiber.New(fiber.Config{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	})
	s := &Server{
		app:    app,
		cfg:    cfg,
		logger: logger,
	}
	s.registerRoutes()
	return s
}

func (s *Server) registerRoutes() {
	//public routes
	s.app.Get("/health", handler.Health())
	//global middleware
	s.app.Use(middleware.TimingMiddleware(s.logger))
	s.app.Use(middleware.LoggingMiddleware(s.logger))
	//protected routes
	protected := s.app.Group("/api")
	protected.Use(middleware.JWTAuthMiddleware())
	protected.Get("/secret", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"msg": "you have access"})
	})
}

func (s *Server) Start() error {
	s.logger.Info("starting server", slog.String("port", s.cfg.ServerPort))
	return s.app.Listen(":" + s.cfg.ServerPort)
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("shutting down server")
	shutdown := make(chan error, 1)
	go func() {
		if err := s.app.Shutdown(); err != nil {
			shutdown <- err
		}
	}()
	select {
	case err := <-shutdown:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
