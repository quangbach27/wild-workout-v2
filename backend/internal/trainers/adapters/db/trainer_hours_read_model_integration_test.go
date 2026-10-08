//go:build integration

package db_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"backend/internal/trainers/adapters/db"
	"backend/internal/trainers/app/queries"
	"backend/internal/trainers/domain"
)

func TestTrainerHoursReadModel_GetTrainerHours(t *testing.T) {
	_, factory := newRepo(t)
	readModel := db.NewTrainerHoursReadModel(pgxDB, factory.Config())

	trainerUUID := newTrainerUUID(t)
	seeded := tomorrowHours(t, factory, trainerUUID, domain.Available, domain.TrainingScheduled)
	seedHours(t, trainerUUID, seeded)

	day := seeded[0].Hour()

	dates, err := readModel.GetTrainerHours(context.Background(), queries.GetTrainerHoursQuery{
		TrainerUUID: trainerUUID,
		DateFrom:    day,
		DateTo:      day.AddDate(0, 0, 1),
	})
	require.NoError(t, err)

	cfg := factory.Config()
	require.Len(t, dates, 2)
	for _, date := range dates {
		assert.Len(t, date.Hours, cfg.MaxUtcHour-cfg.MinUtcHour+1)
	}

	assert.True(t, dates[0].HasFreeHours)
	assert.Equal(t, domain.Available, dates[0].Hours[0].Status)
	assert.Equal(t, domain.TrainingScheduled, dates[0].Hours[1].Status)
	for _, h := range dates[0].Hours[2:] {
		assert.Equal(t, domain.NotAvailable, h.Status)
	}

	assert.False(t, dates[1].HasFreeHours)
	for _, h := range dates[1].Hours {
		assert.Equal(t, domain.NotAvailable, h.Status)
	}
}

func TestTrainerHoursReadModel_GetTrainerHours_StatusFilter(t *testing.T) {
	_, factory := newRepo(t)
	readModel := db.NewTrainerHoursReadModel(pgxDB, factory.Config())

	trainerUUID := newTrainerUUID(t)
	seeded := tomorrowHours(t, factory, trainerUUID, domain.Available, domain.TrainingScheduled)
	seedHours(t, trainerUUID, seeded)

	day := seeded[0].Hour()
	status := domain.Available

	dates, err := readModel.GetTrainerHours(context.Background(), queries.GetTrainerHoursQuery{
		TrainerUUID: trainerUUID,
		DateFrom:    day,
		DateTo:      day.AddDate(0, 0, 1),
		Status:      &status,
	})
	require.NoError(t, err)

	// the second day has no available hours, so it is omitted
	require.Len(t, dates, 1)
	assert.True(t, dates[0].HasFreeHours)
	require.Len(t, dates[0].Hours, 1)
	assert.Equal(t, seeded[0].Hour(), dates[0].Hours[0].Hour)
	assert.Equal(t, domain.Available, dates[0].Hours[0].Status)
}

func TestTrainerHoursReadModel_GetTrainerHours_OtherStatuses(t *testing.T) {
	_, factory := newRepo(t)
	readModel := db.NewTrainerHoursReadModel(pgxDB, factory.Config())
	cfg := factory.Config()

	trainerUUID := newTrainerUUID(t)
	seeded := tomorrowHours(t, factory, trainerUUID, domain.Available, domain.TrainingScheduled)
	seedHours(t, trainerUUID, seeded)

	day := seeded[0].Hour()
	get := func(t *testing.T, status *domain.HourStatus) []queries.Date {
		t.Helper()

		dates, err := readModel.GetTrainerHours(context.Background(), queries.GetTrainerHoursQuery{
			TrainerUUID: trainerUUID,
			DateFrom:    day,
			DateTo:      day.AddDate(0, 0, 1),
			Status:      status,
		})
		require.NoError(t, err)

		return dates
	}

	t.Run("training scheduled returns only the stored scheduled hours", func(t *testing.T) {
		status := domain.TrainingScheduled
		dates := get(t, &status)

		require.Len(t, dates, 1)
		assert.False(t, dates[0].HasFreeHours)
		require.Len(t, dates[0].Hours, 1)
		assert.Equal(t, seeded[1].Hour(), dates[0].Hours[0].Hour)
		assert.Equal(t, domain.TrainingScheduled, dates[0].Hours[0].Status)
	})

	t.Run("not available falls back to the hours without a row", func(t *testing.T) {
		status := domain.NotAvailable
		dates := get(t, &status)

		require.Len(t, dates, 2)
		// the first day lost the available and the scheduled hour
		assert.Len(t, dates[0].Hours, cfg.MaxUtcHour-cfg.MinUtcHour+1-2)
		assert.Len(t, dates[1].Hours, cfg.MaxUtcHour-cfg.MinUtcHour+1)
		for _, date := range dates {
			for _, h := range date.Hours {
				assert.Equal(t, domain.NotAvailable, h.Status)
			}
		}
	})

	t.Run("nil and zero status return every hour", func(t *testing.T) {
		for _, status := range []*domain.HourStatus{nil, {}} {
			dates := get(t, status)

			require.Len(t, dates, 2)
			for _, date := range dates {
				assert.Len(t, date.Hours, cfg.MaxUtcHour-cfg.MinUtcHour+1)
			}
		}
	})
}
