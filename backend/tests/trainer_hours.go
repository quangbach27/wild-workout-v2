//go:build component

package tests

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"backend/internal/trainers/domain"
	trainersclient "backend/internal/trainers/ports/http/clients"
)

// TomorrowHours returns n consecutive full hours, starting tomorrow at the min UTC hour.
func TomorrowHours(t *testing.T, n int) []time.Time {
	t.Helper()

	factory, err := domain.NewHourFactory()
	require.NoError(t, err)

	d := time.Now().UTC().AddDate(0, 0, 1)
	first := time.Date(d.Year(), d.Month(), d.Day(), factory.Config().MinUtcHour, 0, 0, 0, time.UTC)

	hours := make([]time.Time, 0, n)
	for i := range n {
		hours = append(hours, first.Add(time.Duration(i)*time.Hour))
	}

	return hours
}

// MakeHoursAvailable makes the hours available for the trainer that owns trainerToken, so other
// modules' tests can book them.
func MakeHoursAvailable(
	ctx context.Context,
	t *testing.T,
	clients TestClients,
	trainerToken string,
	hours []time.Time,
) {
	t.Helper()

	resp, err := clients.Trainers.MakeHoursAvailableWithResponse(
		ctx,
		trainersclient.HoursRequest{Hours: hours},
		WithAuth(trainerToken),
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode(), string(resp.Body))
}
