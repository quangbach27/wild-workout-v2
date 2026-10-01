//go:build component

package trainers_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	common "github.com/quangbach27/golang-common"
	"github.com/stretchr/testify/require"

	"backend/internal/trainers/domain"
	trainersclient "backend/internal/trainers/ports/http/clients"
	"backend/tests"
)

type testTrainer struct {
	UUID  string
	Token string
}

// newTrainer returns a trainer with a fresh UUID and its own session, so tests never share hours.
func newTrainer(t *testing.T) testTrainer {
	t.Helper()

	trainerUUID := common.NewUUIDv7().String()
	return testTrainer{
		UUID:  trainerUUID,
		Token: tests.NewSession(trainerUUID),
	}
}

func hourFactoryConfig(t *testing.T) *domain.HourFactoryConfig {
	t.Helper()

	factory, err := domain.NewHourFactory()
	require.NoError(t, err)

	return factory.Config()
}

// tomorrowHours returns n consecutive full hours, starting tomorrow at the min UTC hour.
func tomorrowHours(t *testing.T, n int) []time.Time {
	t.Helper()

	d := time.Now().UTC().AddDate(0, 0, 1)
	first := time.Date(d.Year(), d.Month(), d.Day(), hourFactoryConfig(t).MinUtcHour, 0, 0, 0, time.UTC)

	hours := make([]time.Time, 0, n)
	for i := range n {
		hours = append(hours, first.Add(time.Duration(i)*time.Hour))
	}

	return hours
}

func makeHoursAvailable(
	ctx context.Context,
	t *testing.T,
	clients tests.TestClients,
	trainer testTrainer,
	hours []time.Time,
) {
	t.Helper()

	resp, err := clients.Trainers.MakeHoursAvailableWithResponse(
		ctx,
		makeHoursBody(hours),
		tests.WithAuth(trainer.Token),
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode(), string(resp.Body))
}

func makeHoursNotAvailable(
	ctx context.Context,
	t *testing.T,
	clients tests.TestClients,
	trainer testTrainer,
	hours []time.Time,
) {
	t.Helper()

	resp, err := clients.Trainers.MakeHoursNotAvailableWithResponse(
		ctx,
		makeHoursBody(hours),
		tests.WithAuth(trainer.Token),
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode(), string(resp.Body))
}

func getTrainerHours(
	ctx context.Context,
	t *testing.T,
	clients tests.TestClients,
	trainer testTrainer,
	from time.Time,
	to time.Time,
) []trainersclient.Date {
	t.Helper()

	resp, err := clients.Trainers.GetTrainerHoursWithResponse(
		ctx,
		&trainersclient.GetTrainerHoursParams{From: from, To: to},
		tests.WithAuth(trainer.Token),
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode(), string(resp.Body))
	require.NotNil(t, resp.JSON200)

	return *resp.JSON200
}

func makeHoursBody(hours []time.Time) trainersclient.HoursRequest {
	return trainersclient.HoursRequest{Hours: hours}
}

// statusAt returns the status of the hour in the dates, failing the test if the hour is missing.
func statusAt(t *testing.T, dates []trainersclient.Date, hour time.Time) domain.HourStatus {
	t.Helper()

	for _, date := range dates {
		for _, h := range date.Hours {
			if h.Hour.Equal(hour) {
				return h.Status
			}
		}
	}

	require.FailNow(t, "hour not found in dates", hour.String())

	return domain.HourStatus{}
}
