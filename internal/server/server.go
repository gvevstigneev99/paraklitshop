package server

import (
	"context"
	"fmt"
	"time"

	"paraklitshop/internal/config"
	"paraklitshop/internal/handler"
	"paraklitshop/internal/middleware"
	"paraklitshop/internal/repository"
	"paraklitshop/internal/repository/postgres"
	"paraklitshop/internal/repository/redis"
	"paraklitshop/internal/service"

	_ "paraklitshop/docs"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
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
	// Подключение к БД
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		s.cfg.Postgres.Host,
		s.cfg.Postgres.Port,
		s.cfg.Postgres.User,
		s.cfg.Postgres.Password,
		s.cfg.Postgres.DBName,
		s.cfg.Postgres.SSLMode,
	)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		s.logger.Error("failed to connect to postgres", slog.Any("error", err))
		// В учебном режиме продолжаем работу без БД
	}

	var userRepository repository.UserRepository
	if db != nil {
		userRepository = postgres.NewUserRepository(db)
	}
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

	authService := service.NewAuthService(userRepository, s.cfg.JWT.Secret, s.cfg.JWT.TTL)
	authHandler := handler.NewAuthHandler(authService)
	s.app.Post("/login", authHandler.Login)
	s.app.Get("/swagger/*", swagger.HandlerDefault) // default swagger UI

	//global middleware
	s.app.Use(middleware.TimingMiddleware(s.logger))
	s.app.Use(middleware.LoggingMiddleware(s.logger))

	//products routes
	productHandler := handler.NewProductHandler(service.NewProductService(postgres.NewProductRepository()))
	s.app.Get("/products", productHandler.List)

	//protected routes
	protected := s.app.Group("/api")
	protected.Use(middleware.JWTMiddleware(s.cfg.JWT.Secret))
	protected.Post("/orders", orderHandler.CreateOrder)
	protected.Post("/cart/add/:productId/:qty", cartHandler.AddToCart)
	protected.Get("/cart", cartHandler.ViewCart)
	protected.Delete("/cart/clear", cartHandler.ClearCart)
	protected.Delete("/cart/remove/:productId", cartHandler.RemoveFromCart)
	protected.Get("/secret", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"msg": "you have access"})
	})

	// Для покупателя
	buyer := protected.Group("/buyer")
	buyer.Use(middleware.RoleMiddleware("buyer"))
	buyer.Get("/cart", cartHandler.ViewCart)
	buyer.Post("/cart/add/:productId/:qty", cartHandler.AddToCart)
	buyer.Post("/orders", orderHandler.CreateOrder)

	// Для продавца
	seller := protected.Group("/seller")
	seller.Use(middleware.RoleMiddleware("seller"))
	seller.Get("/cart", cartHandler.ViewCart)
	seller.Post("/cart/add/:productId/:qty", cartHandler.AddToCart)
	seller.Post("/orders", orderHandler.CreateOrder)
	seller.Get("/products", productHandler.List)
}

func (s *Server) Start() error {
	s.logger.Info(
		"http server started",
		"port", s.cfg.HTTP.Port,
	)
	return s.app.Listen(fmt.Sprintf(":%d", s.cfg.HTTP.Port))
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
