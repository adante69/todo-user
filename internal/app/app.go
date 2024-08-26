package app

import (
	"log/slog"
	"time"
	"todo-sso/internal/app/grpc"
	"todo-sso/internal/services/auth"
	"todo-sso/internal/storage/postgres"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(logger *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTLL time.Duration,
) *App {
	storage, err := postgres.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(logger, storage, storage, storage, tokenTLL)

	grpcAPP := grpcapp.New(logger, authService, grpcPort)

	return &App{
		GRPCServer: grpcAPP,
	}
}
