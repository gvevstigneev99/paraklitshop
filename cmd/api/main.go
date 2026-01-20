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

func main() {

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	//Наcтраиваем логгер
	l := logger.New(cfg.Env)

	// Создаем сервер
	srv := server.NewServer(cfg, l)

	// Запускаем сервер в отдельной горутине
	go func() {
		if err := srv.Start(); err != nil {
			l.Error("failed to start server", "error", err)
			os.Exit(1)
		}
	}()

	// Ожидаем сигнал завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	l.Info("shutting down server...")

	// Создаем контекст с таймаутом для завершения работы сервера
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		l.Error("failed to shutdown server", "error", err)
		os.Exit(1)
	}
	l.Info("server stopped gracefully")
}
