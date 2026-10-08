//go:build integration

package db_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"backend/internal/trainers/adapters/db"
	"backend/internal/trainers/adapters/db/dbmodels"
	"backend/internal/trainers/domain"
)

const concurrentUpdates = 5

// newTrainerUUID returns a unique trainer, so tests never collide with each other or with dev
// data. Its hours are deleted when the test ends.
func newTrainerUUID(t *testing.T) string {
	t.Helper()

	trainerUUID := uuid.NewV7().String()
	t.Cleanup(func() {
		_, err := pgxDB.Exec(
			context.Background(),
			"DELETE FROM trainers.hours WHERE trainer_uuid = $1",
			trainerUUID,
		)
		assert.NoError(t, err)
	})

	return trainerUUID
}

func newRepo(t *testing.T) (domain.HourRepository, *domain.HourFactory) {
	t.Helper()

	factory, err := domain.NewHourFactory()
	require.NoError(t, err)

	return db.NewHourRepository(pgxDB, factory), factory
}

// tomorrowHours returns one hour per status: consecutive full hours, starting tomorrow at the
// factory's min UTC hour.
func tomorrowHours(
	t *testing.T,
	factory *domain.HourFactory,
	trainerUUID string,
	statuses ...domain.HourStatus,
) []*domain.Hour {
	t.Helper()

	d := time.Now().UTC().AddDate(0, 0, 1)
	first := time.Date(d.Year(), d.Month(), d.Day(), factory.Config().MinUtcHour, 0, 0, 0, time.UTC)

	hours := make([]*domain.Hour, 0, len(statuses))
	for i, status := range statuses {
		hour := first.Add(time.Duration(i) * time.Hour)

		var (
			h   *domain.Hour
			err error
		)
		switch status {
		case domain.Available:
			h, err = factory.NewAvailableHour(trainerUUID, hour)
		case domain.NotAvailable:
			h, err = factory.NewNotAvailableHour(trainerUUID, hour)
		case domain.TrainingScheduled:
			h, err = factory.NewAvailableHour(trainerUUID, hour)
			require.NoError(t, err)
			err = h.ScheduleTraining()
		default:
			require.FailNow(t, "unsupported status", status.String())
		}
		require.NoError(t, err)

		hours = append(hours, h)
	}

	return hours
}

func hourTimes(hours []*domain.Hour) []time.Time {
	times := make([]time.Time, 0, len(hours))
	for _, h := range hours {
		times = append(times, h.Hour())
	}

	return times
}

func seedHours(t *testing.T, trainerUUID string, hours []*domain.Hour) {
	t.Helper()

	statuses := make([]string, 0, len(hours))
	for _, h := range hours {
		statuses = append(statuses, h.Status().String())
	}

	err := dbmodels.New(pgxDB).UpsertHours(context.Background(), dbmodels.UpsertHoursParams{
		TrainerUuid: trainerUUID,
		Hours:       hourTimes(hours),
		Status:      statuses,
	})
	require.NoError(t, err)
}

// requireSameHours compares hours field by field: whole hours can't be compared with assert.Equal
// because times loaded from the database carry another time zone.
func requireSameHours(t *testing.T, want, got []*domain.Hour) {
	t.Helper()

	if !assert.Len(t, got, len(want)) {
		return
	}
	for i := range want {
		assert.Equal(t, want[i].TrainerUUID(), got[i].TrainerUUID(), "trainer of hour %d", i)
		assert.Equal(t, want[i].Status(), got[i].Status(), "status of hour %d", i)
		assert.True(t, want[i].Hour().Equal(got[i].Hour()), "time of hour %d", i)
	}
}

// getStoredStatuses returns the stored status of each hour, keyed by Unix time so time zones
// don't matter.
func getStoredStatuses(t *testing.T, trainerUUID string, hours []*domain.Hour) map[int64]domain.HourStatus {
	t.Helper()

	models, err := dbmodels.New(pgxDB).GetHours(context.Background(), dbmodels.GetHoursParams{
		TrainerUuid: trainerUUID,
		Hours:       hourTimes(hours),
	})
	require.NoError(t, err)

	statuses := make(map[int64]domain.HourStatus, len(models))
	for _, m := range models {
		statuses[m.Hour.Unix()] = m.Status
	}

	return statuses
}

// requireStoredStatuses asserts that hours[i] is stored with want[i], and nothing else is stored.
func requireStoredStatuses(t *testing.T, trainerUUID string, hours []*domain.Hour, want ...domain.HourStatus) {
	t.Helper()

	stored := getStoredStatuses(t, trainerUUID, hours)
	require.Len(t, stored, len(want))
	for i, h := range hours {
		assert.Equal(t, want[i], stored[h.Hour().Unix()], "stored status of hour %d", i)
	}
}

func TestHourRepository_UpsertHours_HoursExist(t *testing.T) {
	t.Parallel()

	repo, factory := newRepo(t)
	trainerUUID := newTrainerUUID(t)
	hours := tomorrowHours(t, factory, trainerUUID, domain.Available, domain.NotAvailable, domain.Available)
	seedHours(t, trainerUUID, hours)

	err := repo.UpsertHours(context.Background(), trainerUUID, hourTimes(hours), func(got []*domain.Hour) error {
		// the hours are loaded from the database, not the not-available default
		requireSameHours(t, hours, got)

		if err := got[0].MakeNotAvailable(); err != nil {
			return err
		}

		return got[1].MakeAvailable()
	})
	require.NoError(t, err)

	// on conflict the status is updated: no new rows, the untouched hour keeps its status
	requireStoredStatuses(t, trainerUUID, hours, domain.NotAvailable, domain.Available, domain.Available)
}

func TestHourRepository_UpsertHours_HoursNotExist(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		statuses   []domain.HourStatus // what upsertFn should receive
		seeded     int                 // how many leading hours are already in the database
		change     func(got []*domain.Hour) error
		wantStored []domain.HourStatus
	}{
		{
			name:       "none stored, left untouched",
			statuses:   []domain.HourStatus{domain.NotAvailable, domain.NotAvailable, domain.NotAvailable},
			change:     func([]*domain.Hour) error { return nil },
			wantStored: []domain.HourStatus{domain.NotAvailable, domain.NotAvailable, domain.NotAvailable},
		},
		{
			name:     "none stored, made available",
			statuses: []domain.HourStatus{domain.NotAvailable, domain.NotAvailable, domain.NotAvailable},
			change: func(got []*domain.Hour) error {
				for _, h := range got {
					if err := h.MakeAvailable(); err != nil {
						return err
					}
				}

				return nil
			},
			wantStored: []domain.HourStatus{domain.Available, domain.Available, domain.Available},
		},
		{
			name:     "some stored, some missing",
			statuses: []domain.HourStatus{domain.Available, domain.NotAvailable},
			seeded:   1,
			change: func(got []*domain.Hour) error {
				if err := got[0].MakeNotAvailable(); err != nil {
					return err
				}

				return got[1].MakeAvailable()
			},
			wantStored: []domain.HourStatus{domain.NotAvailable, domain.Available},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			repo, factory := newRepo(t)
			trainerUUID := newTrainerUUID(t)
			hours := tomorrowHours(t, factory, trainerUUID, tc.statuses...)
			seedHours(t, trainerUUID, hours[:tc.seeded])

			err := repo.UpsertHours(context.Background(), trainerUUID, hourTimes(hours), func(got []*domain.Hour) error {
				// a missing hour is handed over as not available
				requireSameHours(t, hours, got)

				return tc.change(got)
			})
			require.NoError(t, err)

			requireStoredStatuses(t, trainerUUID, hours, tc.wantStored...)
		})
	}
}

func TestHourRepository_UpsertHours_UpsertFnError(t *testing.T) {
	t.Parallel()

	repo, factory := newRepo(t)
	trainerUUID := newTrainerUUID(t)
	hours := tomorrowHours(t, factory, trainerUUID, domain.NotAvailable, domain.NotAvailable, domain.NotAvailable)
	errUpsertFn := errors.New("upsert fn failed")

	err := repo.UpsertHours(context.Background(), trainerUUID, hourTimes(hours), func(got []*domain.Hour) error {
		if err := got[0].MakeAvailable(); err != nil {
			return err
		}

		return errUpsertFn
	})

	require.ErrorIs(t, err, errUpsertFn)
	assert.Empty(t, getStoredStatuses(t, trainerUUID, hours), "nothing should be saved when upsertFn fails")
}

func TestHourRepository_UpsertHours_ConcurrentOnlyOneWins(t *testing.T) {
	t.Parallel()

	repo, factory := newRepo(t)
	trainerUUID := newTrainerUUID(t)
	hours := tomorrowHours(t, factory, trainerUUID, domain.Available, domain.Available, domain.Available)
	seedHours(t, trainerUUID, hours)

	errCh := make(chan error, concurrentUpdates) // buffered, so failing goroutines never block
	start := make(chan struct{})

	var wg sync.WaitGroup
	for range concurrentUpdates {
		wg.Add(1)
		go func() {
			defer wg.Done()

			<-start
			err := repo.UpsertHours(context.Background(), trainerUUID, hourTimes(hours), func(got []*domain.Hour) error {
				for _, h := range got {
					if err := h.ScheduleTraining(); err != nil {
						return err
					}
				}

				return nil
			})
			if err != nil {
				errCh <- err
			}
		}()
	}
	close(start)
	wg.Wait()
	close(errCh)

	// only one update wins, every other one fails
	assert.Len(t, errCh, concurrentUpdates-1)

	requireStoredStatuses(
		t, trainerUUID, hours,
		domain.TrainingScheduled, domain.TrainingScheduled, domain.TrainingScheduled,
	)
}
