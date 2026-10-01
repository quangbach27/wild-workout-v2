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
