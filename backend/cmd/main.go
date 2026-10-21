package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	commonLog "github.com/quangbach27/golang-common/log"

	"backend/internal"
	"backend/internal/configs"
	"backend/internal/configs/auth"
	userDB "backend/internal/users/adapters/db"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	commonLog.Init(slog.LevelInfo)

	config := configs.NewConfig()

	dbPgx, err := pgxpool.New(ctx, config.DB.Dsn())
	if err != nil {
		panic(err)
	}

	userReadModel := userDB.NewUserReadModel(dbPgx)

	svc, err := internal.New(
		ctx,
		config,
		dbPgx,
		internal.ExternalServices{
			TokenVerifier: auth.NewClerkVerifier(
				config.Clerk.ClerkSecretKey,
				config.App.CorsAllowedOrigins,
				userReadModel,
			),
		},
	)
	if err != nil {
		panic(err)
	}

	if err := svc.Run(ctx); err != nil {
		panic(err)
	}
}
