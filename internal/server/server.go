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

type Dependencies struct {
	UserRepository    repository.UserRepository
	CartRepository    repository.CartRepository
	ProductRepository repository.ProductRepository
	OrderRepository   repository.OrderRepository
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

	//global middleware
	s.app.Use(middleware.TimingMiddleware(s.logger))
	s.app.Use(middleware.LoggingMiddleware(s.logger))

	return s
}

func (s *Server) SetupDependencies() (Dependencies, error) {
	// Postgres connection
	pgDSN := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		s.cfg.Postgres.Host,
		s.cfg.Postgres.Port,
		s.cfg.Postgres.User,
		s.cfg.Postgres.Password,
		s.cfg.Postgres.DBName,
		s.cfg.Postgres.SSLMode,
	)
	db, err := sqlx.Connect("postgres", pgDSN)
	if err != nil {
		s.logger.Error("failed to connect to postgres", slog.Any("error", err))
		return Dependencies{}, err
	}

	//Redis connection
	cartRepo, err := redis.NewCartRepository(s.cfg.Redis.Addr, s.cfg.Redis.Password, s.cfg.Redis.DB)
	if err != nil {
		s.logger.Error("failed to connect to redis", slog.Any("error", err))
		return Dependencies{}, err
	}

	deps := Dependencies{
		UserRepository:    postgres.NewUserRepository(db),
		CartRepository:    cartRepo,
		ProductRepository: postgres.NewProductRepository(db),
		OrderRepository:   postgres.NewOrderRepository(db),
	}
	return deps, nil
}

func (s *Server) RegisterRoutes(deps Dependencies) {

	s.app.Get("/health", handler.Health())
	s.app.Get("/swagger/*", swagger.HandlerDefault) // default swagger UI

	productHandler := handler.NewProductHandler(service.NewProductService(deps.ProductRepository))
	s.app.Get("/products", productHandler.List)

	authService := service.NewAuthService(deps.UserRepository, s.cfg.JWT.Secret, s.cfg.JWT.TTL)
	authHandler := handler.NewAuthHandler(authService)
	auth := s.app.Group("/auth")
	auth.Post("/login", authHandler.Login)

	//protected routes
	protected := s.app.Group("/api")
	protected.Use(middleware.JWTMiddleware(s.cfg.JWT.Secret))

	cartService := service.NewCartService(deps.CartRepository, deps.ProductRepository)
	cartHandler := handler.NewCartHandler(cartService)

	orderService := service.NewOrderService(deps.OrderRepository, deps.CartRepository, deps.ProductRepository)
	orderHandler := handler.NewOrderHandler(orderService)

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
