package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jackc/pgx/v5/pgxpool"
	commonAuth "github.com/quangbach27/golang-common/http/auth"
	commonLog "github.com/quangbach27/golang-common/log"

	"backend/internal"
	"backend/internal/configs"
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

	svc, err := internal.New(
		ctx,
		config,
		dbPgx,
		internal.ExternalServices{
			TokenVerifier: commonAuth.NewStubTokenVerifier(),
		},
	)
	if err != nil {
		panic(err)
	}

	if err := svc.Run(ctx); err != nil {
		panic(err)
	}
}
