//go:build integration

package db_test

import (
	"context"
	"testing"
	"time"

	common "github.com/quangbach27/golang-common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"backend/internal/shared"
	"backend/internal/trainings/adapters/db"
	"backend/internal/trainings/app/queries"
	"backend/internal/trainings/domain"
)

// addTrainingAt stores a training at any hour, even a past one, and deletes it when the test ends.
func addTrainingAt(
	t *testing.T,
	repo domain.TrainingRepository,
	attendee, trainer domain.User,
	hour time.Time,
	canceled bool,
	prepare func(*domain.Training),
) *domain.Training {
	t.Helper()

	training := domain.UnmarshalTraining(
		domain.TrainingUUID{UUID: common.NewUUIDv7()},
		hour,
		"notes",
		attendee,
		trainer,
		nil,
		nil,
		canceled,
	)
	if prepare != nil {
		prepare(training)
	}
	require.NoError(t, repo.AddTraining(context.Background(), training))

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

func TestUserTrainingsReadModel_GetUpcomingTrainings(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := db.NewTrainingRepository(pgxDB)
	readModel := db.NewUserTrainingsReadModel(pgxDB)
	u := newUsers(t)
	otherAttendee, err := domain.NewUser("other-attendee-"+t.Name(), "other-attendee", shared.RoleAttendee)
	require.NoError(t, err)
	otherTrainer, err := domain.NewUser("other-trainer-"+t.Name(), "other-trainer", shared.RoleTrainer)
	require.NoError(t, err)
	now := time.Now().UTC().Truncate(time.Hour)

	proposedTime := now.AddDate(0, 0, 5)
	sooner := addTrainingAt(t, repo, u.attendee, u.trainer, now.AddDate(0, 0, 1), false, nil)
	canceled := addTrainingAt(t, repo, u.attendee, u.trainer, now.AddDate(0, 0, 2), true, nil)
	later := addTrainingAt(t, repo, u.attendee, u.trainer, now.AddDate(0, 0, 3), false, func(tr *domain.Training) {
		require.NoError(t, tr.ProposeReschedule(proposedTime, u.trainer))
	})
	addTrainingAt(t, repo, u.attendee, u.trainer, now.AddDate(0, 0, -1), false, nil) // past
	addTrainingAt(t, repo, otherAttendee, otherTrainer, now.AddDate(0, 0, 1), false, nil)

	for name, userUUID := range map[string]string{
		"attendee": u.attendee.ID(),
		"trainer":  u.trainer.ID(),
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			page, err := readModel.GetUpcomingTrainings(ctx, queries.GetUserTrainingsQuery{
				UserUUID: userUUID,
				After:    time.Now(),
				Page:     1,
				PageSize: 10,
			})
			require.NoError(t, err)

			got := page.Trainings
			assert.Equal(t, 3, page.Total)
			assert.Equal(t, 1, page.Page)
			assert.Equal(t, 10, page.PageSize)

			require.Len(t, got, 3)
			assert.Equal(t, sooner.ID().String(), got[0].UUID)
			assert.Equal(t, canceled.ID().String(), got[1].UUID)
			assert.Equal(t, later.ID().String(), got[2].UUID)
			assert.True(t, sooner.Hour().Equal(got[0].Hour))
			assert.Equal(t, u.attendee.ID(), got[0].AttendeeUUID)
			assert.Equal(t, u.trainer.ID(), got[0].TrainerUUID)
			assert.Equal(t, u.attendee.Username(), got[0].AttendeeUsername)
			assert.Equal(t, u.trainer.Username(), got[0].TrainerUsername)
			assert.Equal(t, "notes", got[0].Notes)
			assert.Nil(t, got[0].ProposedNewTime)
			assert.Nil(t, got[0].ProposedBy)
			assert.False(t, got[0].Canceled)

			assert.True(t, got[1].Canceled)

			require.NotNil(t, got[2].ProposedNewTime)
			assert.True(t, proposedTime.Equal(*got[2].ProposedNewTime))
			require.NotNil(t, got[2].ProposedBy)
			assert.Equal(t, u.trainer.ID(), *got[2].ProposedBy)
			require.NotNil(t, got[2].ProposedByUsername)
			assert.Equal(t, u.trainer.Username(), *got[2].ProposedByUsername)
			assert.False(t, got[2].Canceled)
		})
	}

	t.Run("user without trainings", func(t *testing.T) {
		t.Parallel()

		page, err := readModel.GetUpcomingTrainings(ctx, queries.GetUserTrainingsQuery{
			UserUUID: "nobody-" + t.Name(),
			After:    time.Now(),
			Page:     1,
			PageSize: 10,
		})
		require.NoError(t, err)
		assert.Empty(t, page.Trainings)
		assert.Zero(t, page.Total)
	})
}

func TestUserTrainingsReadModel_GetUpcomingTrainings_Pagination(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	repo := db.NewTrainingRepository(pgxDB)
	readModel := db.NewUserTrainingsReadModel(pgxDB)
	u := newUsers(t)
	hour := time.Now().UTC().Truncate(time.Hour).AddDate(0, 0, 1)

	// two trainings share the same hour, so the page boundary has to break the tie by id
	expected := []*domain.Training{
		addTrainingAt(t, repo, u.attendee, u.trainer, hour, false, nil),
		addTrainingAt(t, repo, u.attendee, u.trainer, hour, false, nil),
		addTrainingAt(t, repo, u.attendee, u.trainer, hour.Add(time.Hour), false, nil),
		addTrainingAt(t, repo, u.attendee, u.trainer, hour.Add(2*time.Hour), false, nil),
		addTrainingAt(t, repo, u.attendee, u.trainer, hour.Add(3*time.Hour), false, nil),
	}

	var got []string
	for pageNumber, wantLen := range []int{2, 2, 1} {
		page, err := readModel.GetUpcomingTrainings(ctx, queries.GetUserTrainingsQuery{
			UserUUID: u.attendee.ID(),
			After:    time.Now(),
			Page:     pageNumber + 1,
			PageSize: 2,
		})
		require.NoError(t, err)

		assert.Len(t, page.Trainings, wantLen)
		assert.Equal(t, pageNumber+1, page.Page)
		assert.Equal(t, 2, page.PageSize)
		assert.Equal(t, len(expected), page.Total)

		for _, training := range page.Trainings {
			got = append(got, training.UUID)
		}
	}

	require.Len(t, got, len(expected))
	assert.ElementsMatch(t, trainingIDs(expected), got)
	assert.Len(t, uniqueStrings(got), len(got), "a training appeared on two pages")

	// ordered by hour, so the two trainings of the first hour come first
	assert.ElementsMatch(t, trainingIDs(expected[:2]), got[:2])

	beyond, err := readModel.GetUpcomingTrainings(ctx, queries.GetUserTrainingsQuery{
		UserUUID: u.attendee.ID(),
		After:    time.Now(),
		Page:     4,
		PageSize: 2,
	})
	require.NoError(t, err)
	assert.Empty(t, beyond.Trainings)
	assert.Equal(t, len(expected), beyond.Total)
}

func trainingIDs(trainings []*domain.Training) []string {
	ids := make([]string, 0, len(trainings))
	for _, training := range trainings {
		ids = append(ids, training.ID().String())
	}

	return ids
}

func uniqueStrings(values []string) map[string]struct{} {
	set := make(map[string]struct{}, len(values))
	for _, v := range values {
		set[v] = struct{}{}
	}

	return set
}
