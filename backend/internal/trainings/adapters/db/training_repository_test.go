//go:build integration

package db_test

import (
	"context"
	"errors"
	"net/http"
	"sync"
	"testing"
	"time"

	common "github.com/quangbach27/golang-common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"backend/internal/shared"
	"backend/internal/trainings/adapters/db"
	"backend/internal/trainings/domain"
)

const concurrentUpdates = 5

type users struct {
	attendee, trainer, stranger domain.User
}

func newUsers(t *testing.T) users {
	t.Helper()

	newUser := func(id string, role shared.Role) domain.User {
		user, err := domain.NewUser(id+"-"+t.Name(), "name-"+id, role)
		require.NoError(t, err)
		return user
	}

	return users{
		attendee: newUser("attendee", shared.RoleAttendee),
		trainer:  newUser("trainer", shared.RoleTrainer),
		stranger: newUser("stranger", shared.RoleAttendee),
	}
}

// newTraining builds a training and deletes its row when the test ends.
func newTraining(t *testing.T, u users) *domain.Training {
	t.Helper()

	hour := time.Now().UTC().Truncate(time.Hour).AddDate(0, 0, 3)
	training, err := domain.NewTraining(u.attendee, u.trainer, hour, "legs day")
	require.NoError(t, err)

	t.Cleanup(func() {
		_, err := pgxDB.Exec(
			context.Background(),
			"DELETE FROM trainings.trainings WHERE id = $1",
			training.ID(),
		)
		assert.NoError(t, err)
	})

	return training
}

func requireCommonError(t *testing.T, err error, wantStatus int, wantSlug string) {
	t.Helper()

	var commonErr common.Error
	require.ErrorAs(t, err, &commonErr)
	assert.Equal(t, wantStatus, commonErr.HttpErrorCode)
	assert.Equal(t, wantSlug, commonErr.ErrorSlug)
}

func TestTrainingRepository_AddAndGet(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := db.NewTrainingRepository(pgxDB)
	u := newUsers(t)

	training := newTraining(t, u)
	require.NoError(t, training.ProposeReschedule(time.Now().UTC().Truncate(time.Hour).AddDate(0, 0, 5), u.attendee))
	require.NoError(t, repo.AddTraining(ctx, training))

	for _, user := range []domain.User{u.attendee, u.trainer} {
		got, err := repo.GetTraining(ctx, training.ID(), user)
		require.NoError(t, err)

		assert.Equal(t, training.ID(), got.ID())
		assert.True(t, training.Hour().Equal(got.Hour()))
		assert.Equal(t, training.Notes(), got.Notes())
		assert.Equal(t, training.Attendee(), got.Attendee())
		assert.Equal(t, training.Trainer(), got.Trainer())
		assert.False(t, got.IsCanceled())
		require.True(t, got.IsRescheduleProposed())
		assert.True(t, training.ProposedNewTime().Equal(*got.ProposedNewTime()))
		assert.Equal(t, u.attendee, *got.MoveProposedBy())
	}
}

func TestTrainingRepository_Get_Errors(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := db.NewTrainingRepository(pgxDB)
	u := newUsers(t)

	training := newTraining(t, u)
	require.NoError(t, repo.AddTraining(ctx, training))

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		_, err := repo.GetTraining(ctx, domain.TrainingUUID{UUID: common.NewUUIDv7()}, u.attendee)
		requireCommonError(t, err, http.StatusNotFound, "training-not-found")
	})

	t.Run("stranger", func(t *testing.T) {
		t.Parallel()

		_, err := repo.GetTraining(ctx, training.ID(), u.stranger)
		requireCommonError(t, err, http.StatusForbidden, "training-forbidden")
	})
}

func TestTrainingRepository_UpdateTraining(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := db.NewTrainingRepository(pgxDB)
	u := newUsers(t)

	t.Run("cancel is persisted", func(t *testing.T) {
		t.Parallel()

		training := newTraining(t, u)
		require.NoError(t, repo.AddTraining(ctx, training))

		err := repo.UpdateTraining(ctx, training.ID(), u.trainer, func(_ context.Context, tr *domain.Training) error {
			return tr.Cancel(u.trainer)
		})
		require.NoError(t, err)

		got, err := repo.GetTraining(ctx, training.ID(), u.trainer)
		require.NoError(t, err)
		assert.True(t, got.IsCanceled())
	})

	t.Run("propose then approve is persisted", func(t *testing.T) {
		t.Parallel()

		training := newTraining(t, u)
		require.NoError(t, repo.AddTraining(ctx, training))
		newHour := time.Now().UTC().Truncate(time.Hour).AddDate(0, 0, 6)

		err := repo.UpdateTraining(ctx, training.ID(), u.attendee, func(_ context.Context, tr *domain.Training) error {
			return tr.ProposeReschedule(newHour, u.attendee)
		})
		require.NoError(t, err)

		got, err := repo.GetTraining(ctx, training.ID(), u.trainer)
		require.NoError(t, err)
		require.True(t, got.IsRescheduleProposed())
		assert.True(t, newHour.Equal(*got.ProposedNewTime()))

		err = repo.UpdateTraining(ctx, training.ID(), u.trainer, func(_ context.Context, tr *domain.Training) error {
			return tr.ApproveReschedule(u.trainer)
		})
		require.NoError(t, err)

		got, err = repo.GetTraining(ctx, training.ID(), u.trainer)
		require.NoError(t, err)
		assert.True(t, newHour.Equal(got.Hour()))
		assert.False(t, got.IsRescheduleProposed())
	})

	t.Run("updateFn error rolls back", func(t *testing.T) {
		t.Parallel()

		training := newTraining(t, u)
		require.NoError(t, repo.AddTraining(ctx, training))
		errBoom := errors.New("boom")

		err := repo.UpdateTraining(ctx, training.ID(), u.trainer, func(_ context.Context, tr *domain.Training) error {
			require.NoError(t, tr.Cancel(u.trainer))
			return errBoom
		})
		require.ErrorIs(t, err, errBoom)

		got, err := repo.GetTraining(ctx, training.ID(), u.trainer)
		require.NoError(t, err)
		assert.False(t, got.IsCanceled())
	})

	t.Run("stranger", func(t *testing.T) {
		t.Parallel()

		training := newTraining(t, u)
		require.NoError(t, repo.AddTraining(ctx, training))

		err := repo.UpdateTraining(ctx, training.ID(), u.stranger, func(context.Context, *domain.Training) error {
			t.Error("updateFn must not run for a stranger")
			return nil
		})
		requireCommonError(t, err, http.StatusForbidden, "training-forbidden")
	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()

		err := repo.UpdateTraining(
			ctx,
			domain.TrainingUUID{UUID: common.NewUUIDv7()},
			u.trainer,
			func(context.Context, *domain.Training) error {
				t.Error("updateFn must not run for a missing training")
				return nil
			},
		)
		requireCommonError(t, err, http.StatusNotFound, "training-not-found")
	})
}

// Each update runs under a row lock, so exactly one of the concurrent cancels wins and the rest
// see an already canceled training.
func TestTrainingRepository_UpdateTraining_Concurrent(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := db.NewTrainingRepository(pgxDB)
	u := newUsers(t)

	training := newTraining(t, u)
	require.NoError(t, repo.AddTraining(ctx, training))

	errs := make([]error, concurrentUpdates)
	var wg sync.WaitGroup
	for i := range concurrentUpdates {
		wg.Go(func() {
			errs[i] = repo.UpdateTraining(ctx, training.ID(), u.trainer, func(_ context.Context, tr *domain.Training) error {
				return tr.Cancel(u.trainer)
			})
		})
	}
	wg.Wait()

	succeeded := 0
	for _, err := range errs {
		if err == nil {
			succeeded++
			continue
		}
		requireCommonError(t, err, http.StatusConflict, "failed-to-cancel")
	}
	assert.Equal(t, 1, succeeded)
}
