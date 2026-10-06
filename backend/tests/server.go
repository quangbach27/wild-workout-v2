//go:build component

// Package tests holds helpers shared by the component tests of every module.
package tests

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	commonAuth "github.com/quangbach27/golang-common/http/auth"

	"backend/internal"
	"backend/internal/configs"
)

var (
	BaseURL  string
	verifier = commonAuth.NewStubTokenVerifier()
)

// StartServer starts the whole application in the background, using the env config (.env.test),
// and returns immediately. Tests wait for it to answer on /healthz through NewTestClients.
// It stops with the test process.
func StartServer(ctx context.Context) {
	config := configs.NewConfig()
	BaseURL = fmt.Sprintf("http://localhost:%d", config.App.Port)

	db, err := pgxpool.New(ctx, config.DB.Dsn())
	if err != nil {
		panic(err)
	}

	svc, err := internal.New(ctx, config, db, internal.ExternalServices{TokenVerifier: verifier})
	if err != nil {
		panic(err)
	}

	go func() {
		if err := svc.Run(ctx); err != nil {
			panic(err)
		}
	}()

	waitForServer()
}

func waitForServer() {
	for range 100 {
		resp, err := http.Get(BaseURL + "/healthz")
		if err == nil && resp.StatusCode < 300 {
			_ = resp.Body.Close()
			return
		}
		if resp != nil {
			_ = resp.Body.Close()
		}
		time.Sleep(50 * time.Millisecond)
	}
	panic("server did not start within 5 seconds")
}
