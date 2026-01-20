package server

import (
	"context"
	"time"

	"paraklitshop/internal/config"
	"paraklitshop/internal/handler"
	"paraklitshop/internal/middleware"
	"paraklitshop/internal/repository/postgres"
	"paraklitshop/internal/repository/redis"
	"paraklitshop/internal/service"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/exp/slog"
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

	cartRepository := redis.NewCartRepository()
	_ = cartRepository // to avoid unused variable error for now
	productRepository := postgres.NewProductRepository()
	_ = productRepository // to avoid unused variable error for now
	//public routes
	cartService := service.NewCartService(cartRepository, productRepository)
	_ = cartService // to avoid unused variable error for now
	cartHandler := handler.NewCartHandler(cartService)
	_ = cartHandler // to avoid unused variable error for now
	orderRepository := postgres.NewOrderRepository()
	_ = orderRepository // to avoid unused variable error for now
	orderService := service.NewOrderService(orderRepository, cartRepository, productRepository)
	_ = orderService // to avoid unused variable error for now
	orderHandler := handler.NewOrderHandler(orderService)
	_ = orderHandler // to avoid unused variable error for now

	s.app.Get("/health", handler.Health())

	//global middleware
	s.app.Use(middleware.TimingMiddleware(s.logger))
	s.app.Use(middleware.LoggingMiddleware(s.logger))

	//products routes
	productHandler := handler.NewProductHandler(service.NewProductService(postgres.NewProductRepository()))
	s.app.Get("/products", productHandler.List)

	//protected routes
	protected := s.app.Group("/api")
	protected.Post("/orders", orderHandler.CreateOrder)
	protected.Use(middleware.JWTAuthMiddleware())
	protected.Post("/cart/add/:productId/:qty", cartHandler.AddToCart)
	protected.Get("/cart", cartHandler.ViewCart)
	protected.Delete("/cart/clear", cartHandler.ClearCart)
	protected.Delete("/cart/remove/:productId", cartHandler.RemoveFromCart)
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
