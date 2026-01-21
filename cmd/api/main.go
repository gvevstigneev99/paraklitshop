package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"paraklitshop/internal/config"
	"paraklitshop/internal/logger"
	"paraklitshop/internal/server"
)

// @title Paraklit Shop API
// @version 1.0
// @description This is a backend API for Paraklit Shop.
// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	//Наcтраиваем логгер
	log := logger.New(cfg.App.Env)

	log.Info(
		"starting server",
		"app", cfg.App.Name,
		"env", cfg.App.Env,
		"port", cfg.HTTP.Port,
	)

	// Создаем сервер
	srv := server.NewServer(cfg, log)

	// Запускаем сервер в отдельной горутине
	go func() {
		if err := srv.Start(); err != nil {
			log.Error("failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	// Ожидаем сигнал завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Info("shutting down server...")

	// Создаем контекст с таймаутом для завершения работы сервера
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to shutdown server", "error", err)
		os.Exit(1)
	}
	log.Info("server stopped gracefully")
}
