//go:build component

// Package tests holds helpers shared by the component tests of every module.
package tests

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	common "github.com/quangbach27/golang-common"
	commonHTTP "github.com/quangbach27/golang-common/http"
	commonAuth "github.com/quangbach27/golang-common/http/auth"
	"github.com/quangbach27/golang-common/log"
	"github.com/stretchr/testify/require"

	"backend/internal"
	"backend/internal/configs"
	trainersclient "backend/internal/trainers/ports/http/clients"
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

// NewSession registers a bearer token that authenticates as userID, and returns the token.
func NewSession(userID string) string {
	token := common.NewUUIDv7().String()
	verifier.Add(token, &commonAuth.Session{UserID: userID})

	return token
}

// WithAuth is a request editor that authenticates the request with the bearer token.
func WithAuth(token string) func(ctx context.Context, req *http.Request) error {
	return func(_ context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+token)
		return nil
	}
}

// TestClients holds the API client of every module.
type TestClients struct {
	Trainers *trainersclient.ClientWithResponses
}

func NewTestClients(t *testing.T) TestClients {
	t.Helper()

	httpClient := &http.Client{Timeout: 10 * time.Second}

	editorFn := func(ctx context.Context, req *http.Request) error {
		log.FromContext(ctx).
			With(
				"method", req.Method,
				"url", req.URL.String(),
				"test_name", t.Name(),
			).
			Info("Making component test API request")
		req.Header.Set(commonHTTP.TestNameHeader, t.Name())
		return nil
	}

	trainers, err := trainersclient.NewClientWithResponses(
		BaseURL,
		trainersclient.WithHTTPClient(httpClient),
		trainersclient.WithRequestEditorFn(editorFn),
	)
	require.NoError(t, err)

	return TestClients{
		Trainers: trainers,
	}
}
