package main

import (
	"context"
	"net/http"
	"os"
	"time"

	"example.com/authservice/internal/config"
	"example.com/authservice/internal/db"
	"example.com/authservice/internal/handlers"
	"example.com/authservice/internal/middleware"
	"example.com/authservice/internal/repository"
	"example.com/authservice/internal/service"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	zerolog.TimeFieldFormat = time.RFC3339
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	cfg := config.LoadConfig()

	ctx := context.Background()
	conn, err := db.Connect(ctx, cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to connect to the database")
	}
	defer conn.Close(ctx)

	userRepo := repository.NewUserRepository(conn)
	tokenRepo := repository.NewTokenRepository(conn)

	emailService := service.NewEmailService(cfg.EmailSender)
	authService := service.NewAuthService(tokenRepo, userRepo, cfg.JWTSecret, emailService)
	authHandler := handlers.NewAuthHandler(authService)

	r := mux.NewRouter()
	r.Use(middleware.LoggingMiddleware)

	r.HandleFunc("/auth/token", authHandler.IssueTokens).Methods("GET")
	r.HandleFunc("/auth/refresh", authHandler.RefreshTokens).Methods("POST")

	log.Info().Msgf("Starting server on %s", cfg.ServerAddr)
	if err := http.ListenAndServe(cfg.ServerAddr, r); err != nil {
		log.Fatal().Err(err).Msg("server encountered an error")
	}
}
