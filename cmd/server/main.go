package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jorzel/booking-service/internal/app"
	"github.com/jorzel/booking-service/internal/infrastructure"
	"github.com/jorzel/booking-service/internal/transport"
	"github.com/rs/zerolog"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	config := infrastructure.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     5432,
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		Database: getEnv("DB_NAME", "booking_service"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	db, err := infrastructure.NewPostgresDB(config)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to database")
	}
	defer db.Close()

	eventRepo := infrastructure.NewPostgresEventRepository(db)
	bookingRepo := infrastructure.NewPostgresBookingRepository(db)

	eventService := app.NewEventService(eventRepo, logger)
	bookingService := app.NewBookingService(bookingRepo, eventRepo, db, logger)

	router := transport.NewRouter(eventService, bookingService, db, logger)

	port := getEnv("PORT", "8080")
	addr := fmt.Sprintf(":%s", port)

	go func() {
		logger.Info().Str("address", addr).Msg("starting server")
		if err := router.Start(addr); err != nil {
			logger.Fatal().Err(err).Msg("server failed to start")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Info().Msg("shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := router.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msg("server forced to shutdown")
	}

	logger.Info().Msg("server exited")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
