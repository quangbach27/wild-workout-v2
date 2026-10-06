//go:build component

package tests

import (
	"context"
	"net/http"
	"testing"
	"time"

	commonHTTP "github.com/quangbach27/golang-common/http"
	"github.com/quangbach27/golang-common/log"
	"github.com/stretchr/testify/require"

	trainersclient "backend/internal/trainers/ports/http/clients"
	trainingsclient "backend/internal/trainings/ports/http/clients"
)

// TestClients holds the API client of every module.
type TestClients struct {
	Trainers  *trainersclient.ClientWithResponses
	Trainings *trainingsclient.ClientWithResponses
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

	trainings, err := trainingsclient.NewClientWithResponses(
		BaseURL,
		trainingsclient.WithHTTPClient(httpClient),
		trainingsclient.WithRequestEditorFn(editorFn),
	)
	require.NoError(t, err)

	return TestClients{
		Trainers:  trainers,
		Trainings: trainings,
	}
}
