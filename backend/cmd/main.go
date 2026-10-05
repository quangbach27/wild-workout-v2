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

	// TODO: Replace with Real Identity Provider for production
	tokenVerifier := commonAuth.NewStubTokenVerifier()
	registerDevTokens(tokenVerifier)

	svc, err := internal.New(
		ctx,
		config,
		dbPgx,
		internal.ExternalServices{
			TokenVerifier: tokenVerifier,
		},
	)
	if err != nil {
		panic(err)
	}

	if err := svc.Run(ctx); err != nil {
		panic(err)
	}
}

// registerDevTokens registers fixed mock sessions so the web app can call
// protected endpoints locally. The web app sends the same tokens (web/src/lib/auth.ts).
func registerDevTokens(verifier *commonAuth.StubTokenVerifier) {
	verifier.
		Add("mock-trainer-token", &commonAuth.Session{
			UserID: "11111111-1111-4111-8111-111111111111",
			Roles:  []string{"trainer"},
		}).
		Add("mock-attendee-token", &commonAuth.Session{
			UserID: "22222222-2222-4222-8222-222222222222",
			Roles:  []string{"attendee"},
		})
}
