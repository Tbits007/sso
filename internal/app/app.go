package app

import (
	"log/slog"
	"sso/internal/app/grpcapp"
	"sso/internal/services/auth"
	"sso/internal/storage/postgres"
	"time"
)

type App struct {
    GRPCServer *grpcapp.App
}

func New(
    log *slog.Logger,
    grpcPort int,
    storagePath string,
    tokenTTL time.Duration,
) *App {
    storage, err := postgres.New(storagePath)
    if err != nil {
        panic(err)
    }

    authService := auth.New(log, storage, storage, storage, tokenTTL)

    grpcApp := grpcapp.New(log, authService, grpcPort)

    return &App{
        GRPCServer: grpcApp,
    }
}