//go:build component

package trainers_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	common "github.com/quangbach27/golang-common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"backend/internal/trainers/domain"
	trainersclient "backend/internal/trainers/ports/http/clients"
	"backend/tests"
)

func TestMakeHoursAvailable(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	clients := tests.NewTestClients(t)

	t.Run("not stored hours can be made available", func(t *testing.T) {
		t.Parallel()

		trainer := newTrainer(t)
		hours := tests.TomorrowHours(t, 2)

		tests.MakeHoursAvailable(ctx, t, clients, trainer.Token, hours)

		dates := getTrainerHours(ctx, t, clients, trainer, hours[0], hours[0])
		require.Len(t, dates, 1)
		assert.True(t, dates[0].HasFreeHours)
		assert.Equal(t, domain.Available, statusAt(t, dates, hours[0]))
		assert.Equal(t, domain.Available, statusAt(t, dates, hours[1]))
	})

	t.Run("already available hours conflict", func(t *testing.T) {
		t.Parallel()

		trainer := newTrainer(t)
		hours := tests.TomorrowHours(t, 1)
		tests.MakeHoursAvailable(ctx, t, clients, trainer.Token, hours)

		resp, err := clients.Trainers.MakeHoursAvailableWithResponse(
			ctx,
			makeHoursBody(hours),
			tests.WithAuth(trainer.Token),
		)

		require.NoError(t, err)
		require.Equal(t, http.StatusConflict, resp.StatusCode())
		require.NotNil(t, resp.JSON409)
		assert.Equal(t, "hours-incorrect-status", resp.JSON409.Slug)
	})

	t.Run("without token is unauthorized", func(t *testing.T) {
		t.Parallel()

		resp, err := clients.Trainers.MakeHoursAvailableWithResponse(ctx, makeHoursBody(tests.TomorrowHours(t, 1)))

		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode())
	})
}

func TestMakeHoursNotAvailable(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	clients := tests.NewTestClients(t)

	t.Run("available hours can be made not available", func(t *testing.T) {
		t.Parallel()

		trainer := newTrainer(t)
		hours := tests.TomorrowHours(t, 2)
		tests.MakeHoursAvailable(ctx, t, clients, trainer.Token, hours)

		makeHoursNotAvailable(ctx, t, clients, trainer, hours)

		dates := getTrainerHours(ctx, t, clients, trainer, hours[0], hours[0])
		require.Len(t, dates, 1)
		assert.False(t, dates[0].HasFreeHours)
		assert.Equal(t, domain.NotAvailable, statusAt(t, dates, hours[0]))
		assert.Equal(t, domain.NotAvailable, statusAt(t, dates, hours[1]))
	})

	t.Run("not stored hours are already not available and conflict", func(t *testing.T) {
		t.Parallel()

		trainer := newTrainer(t)

		resp, err := clients.Trainers.MakeHoursNotAvailableWithResponse(
			ctx,
			makeHoursBody(tests.TomorrowHours(t, 1)),
			tests.WithAuth(trainer.Token),
		)

		require.NoError(t, err)
		require.Equal(t, http.StatusConflict, resp.StatusCode())
		require.NotNil(t, resp.JSON409)
		assert.Equal(t, "hours-incorrect-status", resp.JSON409.Slug)
	})

	t.Run("without token is unauthorized", func(t *testing.T) {
		t.Parallel()

		resp, err := clients.Trainers.MakeHoursNotAvailableWithResponse(ctx, makeHoursBody(tests.TomorrowHours(t, 1)))

		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode())
	})
}

func TestGetTrainerHours(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	clients := tests.NewTestClients(t)
	cfg := hourFactoryConfig(t)
	hoursPerDay := cfg.MaxUtcHour - cfg.MinUtcHour + 1

	t.Run("trainer without data has every hour not available", func(t *testing.T) {
		t.Parallel()

		trainer := newTrainer(t)
		from := tests.TomorrowHours(t, 1)[0]

		dates := getTrainerHours(ctx, t, clients, trainer, from, from.AddDate(0, 0, 1))

		require.Len(t, dates, 2)
		for _, date := range dates {
			assert.False(t, date.HasFreeHours)
			require.Len(t, date.Hours, hoursPerDay)
			for _, h := range date.Hours {
				assert.Equal(t, domain.NotAvailable, h.Status)
			}
		}
	})

	t.Run("stored hours keep their status", func(t *testing.T) {
		t.Parallel()

		trainer := newTrainer(t)
		hours := tests.TomorrowHours(t, 1)
		nextDay := hours[0].AddDate(0, 0, 1)
		tests.MakeHoursAvailable(ctx, t, clients, trainer.Token, append(hours, nextDay))

		dates := getTrainerHours(ctx, t, clients, trainer, hours[0], nextDay)

		require.Len(t, dates, 2)
		assert.True(t, dates[0].HasFreeHours)
		assert.True(t, dates[1].HasFreeHours)
		assert.Equal(t, domain.Available, statusAt(t, dates, hours[0]))
		assert.Equal(t, domain.Available, statusAt(t, dates, nextDay))
		assert.Equal(t, domain.NotAvailable, statusAt(t, dates, hours[0].Add(time.Hour)))
	})

	t.Run("to before from is invalid", func(t *testing.T) {
		t.Parallel()

		trainer := newTrainer(t)
		from := tests.TomorrowHours(t, 1)[0]

		resp, err := clients.Trainers.GetTrainerHoursWithResponse(
			ctx,
			&trainersclient.GetTrainerHoursParams{From: from, To: from.AddDate(0, 0, -1)},
			tests.WithAuth(trainer.Token),
		)

		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode())
		require.NotNil(t, resp.JSON400)
		assert.Equal(t, "invalid-get-trainer-hours-query", resp.JSON400.Slug)
	})

	t.Run("without token is unauthorized", func(t *testing.T) {
		t.Parallel()

		from := tests.TomorrowHours(t, 1)[0]

		resp, err := clients.Trainers.GetTrainerHoursWithResponse(
			ctx,
			&trainersclient.GetTrainerHoursParams{From: from, To: from},
		)

		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode())
	})
}

func TestGetTrainerHoursByUUID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	clients := tests.NewTestClients(t)

	t.Run("another user reads only the available hours", func(t *testing.T) {
		t.Parallel()

		trainer := newTrainer(t)
		attendeeToken := tests.NewSession(common.NewUUIDv7().String())
		hours := tests.TomorrowHours(t, 2)
		tests.MakeHoursAvailable(ctx, t, clients, trainer.Token, hours)
		makeHoursNotAvailable(ctx, t, clients, trainer, hours[1:])

		status := trainersclient.Availability
		dates := getTrainerHoursByUUID(ctx, t, clients, attendeeToken, trainer.UUID,
			&trainersclient.GetTrainerHoursByUuidParams{From: hours[0], To: hours[0], Status: &status})

		require.Len(t, dates, 1)
		assert.True(t, dates[0].HasFreeHours)
		require.Len(t, dates[0].Hours, 1)
		assert.True(t, dates[0].Hours[0].Hour.Equal(hours[0]))
		assert.Equal(t, domain.Available, dates[0].Hours[0].Status)
	})

	t.Run("without status returns every hour", func(t *testing.T) {
		t.Parallel()

		trainer := newTrainer(t)
		attendeeToken := tests.NewSession(common.NewUUIDv7().String())
		hours := tests.TomorrowHours(t, 1)
		tests.MakeHoursAvailable(ctx, t, clients, trainer.Token, hours)

		dates := getTrainerHoursByUUID(ctx, t, clients, attendeeToken, trainer.UUID,
			&trainersclient.GetTrainerHoursByUuidParams{From: hours[0], To: hours[0]})

		require.Len(t, dates, 1)
		cfg := hourFactoryConfig(t)
		assert.Len(t, dates[0].Hours, cfg.MaxUtcHour-cfg.MinUtcHour+1)
	})

	t.Run("trainer without available hours gets an empty list", func(t *testing.T) {
		t.Parallel()

		trainer := newTrainer(t)
		attendeeToken := tests.NewSession(common.NewUUIDv7().String())
		hours := tests.TomorrowHours(t, 1)

		status := trainersclient.Availability
		dates := getTrainerHoursByUUID(ctx, t, clients, attendeeToken, trainer.UUID,
			&trainersclient.GetTrainerHoursByUuidParams{From: hours[0], To: hours[0], Status: &status})

		assert.Empty(t, dates)
	})

	t.Run("empty status returns every hour", func(t *testing.T) {
		t.Parallel()

		trainer := newTrainer(t)
		attendeeToken := tests.NewSession(common.NewUUIDv7().String())
		hours := tests.TomorrowHours(t, 1)

		status := trainersclient.GetTrainerHoursByUuidParamsStatus("")
		dates := getTrainerHoursByUUID(ctx, t, clients, attendeeToken, trainer.UUID,
			&trainersclient.GetTrainerHoursByUuidParams{From: hours[0], To: hours[0], Status: &status})

		require.Len(t, dates, 1)
		cfg := hourFactoryConfig(t)
		assert.Len(t, dates[0].Hours, cfg.MaxUtcHour-cfg.MinUtcHour+1)
	})

	t.Run("unknown status is invalid input", func(t *testing.T) {
		t.Parallel()

		hours := tests.TomorrowHours(t, 1)
		status := trainersclient.GetTrainerHoursByUuidParamsStatus("foo")
		resp, err := clients.Trainers.GetTrainerHoursByUuidWithResponse(
			ctx,
			newTrainer(t).UUID,
			&trainersclient.GetTrainerHoursByUuidParams{From: hours[0], To: hours[0], Status: &status},
			tests.WithAuth(tests.NewSession(common.NewUUIDv7().String())),
		)

		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode())
	})

	t.Run("without token is unauthorized", func(t *testing.T) {
		t.Parallel()

		hours := tests.TomorrowHours(t, 1)
		resp, err := clients.Trainers.GetTrainerHoursByUuidWithResponse(
			ctx,
			newTrainer(t).UUID,
			&trainersclient.GetTrainerHoursByUuidParams{From: hours[0], To: hours[0]},
		)

		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode())
	})
}
