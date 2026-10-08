//go:build component

package trainings_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	common "github.com/quangbach27/golang-common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"backend/internal/shared"
	"backend/internal/trainers/domain"
	trainingsclient "backend/internal/trainings/ports/http/clients"
	"backend/tests"
)

func TestScheduleTraining(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	clients := tests.NewTestClients(t)

	body := func(trainer testUser, hour time.Time) trainingsclient.ScheduleTrainingRequest {
		return trainingsclient.ScheduleTrainingRequest{
			TrainerUuid:     trainer.UUID,
			TrainerUsername: tests.Username(trainer.UUID),
			Hour:            hour,
			Notes:           "legs day",
		}
	}

	t.Run("attendee schedules an available hour", func(t *testing.T) {
		t.Parallel()

		trainer, attendee := newUser(shared.RoleTrainer), newOnboardedAttendee(ctx, t, clients)
		hour := tests.TomorrowHours(t, 1)[0]
		tests.MakeHoursAvailable(ctx, t, clients, trainer.Token, []time.Time{hour})

		resp := scheduleTraining(ctx, t, clients, attendee.Token, body(trainer, hour))

		require.Equal(t, http.StatusCreated, resp.StatusCode(), string(resp.Body))
		require.NotNil(t, resp.JSON201)
		assert.NotEmpty(t, resp.JSON201.Uuid)
		assert.Equal(t, domain.TrainingScheduled, hourStatus(ctx, t, clients, trainer, hour))
	})

	t.Run("the same hour can't be scheduled twice", func(t *testing.T) {
		t.Parallel()

		trainer := newUser(shared.RoleTrainer)
		hour := tests.TomorrowHours(t, 1)[0]
		tests.MakeHoursAvailable(ctx, t, clients, trainer.Token, []time.Time{hour})

		first := scheduleTraining(ctx, t, clients, newOnboardedAttendee(ctx, t, clients).Token, body(trainer, hour))
		require.Equal(t, http.StatusCreated, first.StatusCode(), string(first.Body))

		second := scheduleTraining(ctx, t, clients, newOnboardedAttendee(ctx, t, clients).Token, body(trainer, hour))
		assert.Equal(t, http.StatusConflict, second.StatusCode(), string(second.Body))
	})

	t.Run("an hour that was never made available conflicts", func(t *testing.T) {
		t.Parallel()

		trainer, attendee := newUser(shared.RoleTrainer), newOnboardedAttendee(ctx, t, clients)
		hour := tests.TomorrowHours(t, 1)[0]

		resp := scheduleTraining(ctx, t, clients, attendee.Token, body(trainer, hour))

		assert.Equal(t, http.StatusConflict, resp.StatusCode(), string(resp.Body))
	})

	t.Run("a trainer can't schedule a training", func(t *testing.T) {
		t.Parallel()

		trainer := newUser(shared.RoleTrainer)
		hour := tests.TomorrowHours(t, 1)[0]
		tests.MakeHoursAvailable(ctx, t, clients, trainer.Token, []time.Time{hour})

		resp := scheduleTraining(ctx, t, clients, trainer.Token, body(trainer, hour))

		assert.Equal(t, http.StatusForbidden, resp.StatusCode(), string(resp.Body))
		assert.Equal(t, domain.Available, hourStatus(ctx, t, clients, trainer, hour))
	})

	t.Run("missing token is unauthorized", func(t *testing.T) {
		t.Parallel()

		trainer := newUser(shared.RoleTrainer)

		resp, err := clients.Trainings.ScheduleTrainingWithResponse(ctx, body(trainer, tests.TomorrowHours(t, 1)[0]))
		require.NoError(t, err)

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode(), string(resp.Body))
	})

	t.Run("empty trainer username is rejected", func(t *testing.T) {
		t.Parallel()

		trainer, attendee := newUser(shared.RoleTrainer), newOnboardedAttendee(ctx, t, clients)
		hour := tests.TomorrowHours(t, 1)[0]
		tests.MakeHoursAvailable(ctx, t, clients, trainer.Token, []time.Time{hour})

		invalid := body(trainer, hour)
		invalid.TrainerUsername = ""

		resp := scheduleTraining(ctx, t, clients, attendee.Token, invalid)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode(), string(resp.Body))
		assert.Equal(t, domain.Available, hourStatus(ctx, t, clients, trainer, hour))
	})

	t.Run("invalid input is rejected and the hour is untouched", func(t *testing.T) {
		t.Parallel()

		trainer, attendee := newUser(shared.RoleTrainer), newOnboardedAttendee(ctx, t, clients)
		hour := tests.TomorrowHours(t, 1)[0]
		tests.MakeHoursAvailable(ctx, t, clients, trainer.Token, []time.Time{hour})

		invalid := body(trainer, hour)
		invalid.Hour = time.Now().Add(-time.Hour)

		resp := scheduleTraining(ctx, t, clients, attendee.Token, invalid)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode(), string(resp.Body))
		assert.Equal(t, domain.Available, hourStatus(ctx, t, clients, trainer, hour))
	})
}

func TestGetUserTrainings(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	clients := tests.NewTestClients(t)

	getPage := func(t *testing.T, token string, params *trainingsclient.GetUserTrainingsParams) trainingsclient.TrainingsPage {
		t.Helper()

		resp, err := clients.Trainings.GetUserTrainingsWithResponse(ctx, params, tests.WithAuth(token))
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode(), string(resp.Body))
		require.NotNil(t, resp.JSON200)

		return *resp.JSON200
	}

	getTrainings := func(t *testing.T, token string) []trainingsclient.Training {
		t.Helper()

		return getPage(t, token, &trainingsclient.GetUserTrainingsParams{}).Items
	}

	t.Run("trainer and attendee see their upcoming trainings, soonest first", func(t *testing.T) {
		t.Parallel()

		trainer, attendee := newUser(shared.RoleTrainer), newOnboardedAttendee(ctx, t, clients)
		hours := tests.TomorrowHours(t, 2)
		tests.MakeHoursAvailable(ctx, t, clients, trainer.Token, hours)

		// scheduled in reverse order to prove the result is sorted by hour
		for _, hour := range []time.Time{hours[1], hours[0]} {
			resp := scheduleTraining(ctx, t, clients, attendee.Token, trainingsclient.ScheduleTrainingRequest{
				TrainerUuid: trainer.UUID, TrainerUsername: tests.Username(trainer.UUID), Hour: hour, Notes: "legs day",
			})
			require.Equal(t, http.StatusCreated, resp.StatusCode(), string(resp.Body))
		}

		for _, user := range []testUser{attendee, trainer} {
			got := getTrainings(t, user.Token)

			require.Len(t, got, 2)
			assert.True(t, hours[0].Equal(got[0].Hour))
			assert.True(t, hours[1].Equal(got[1].Hour))
			assert.NotEmpty(t, got[0].Uuid)
			assert.Equal(t, attendee.UUID, got[0].AttendeeUuid)
			assert.Equal(t, trainer.UUID, got[0].TrainerUuid)
			assert.Equal(t, tests.Username(attendee.UUID), got[0].AttendeeUsername)
			assert.Equal(t, tests.Username(trainer.UUID), got[0].TrainerUsername)
			assert.Equal(t, "legs day", got[0].Notes)
			assert.Nil(t, got[0].ProposedNewTime)
			assert.Nil(t, got[0].ProposedBy)
			assert.False(t, got[0].Canceled)
		}
	})

	t.Run("an unrelated user has no trainings", func(t *testing.T) {
		t.Parallel()

		got := getTrainings(t, newOnboardedAttendee(ctx, t, clients).Token)

		assert.Empty(t, got)
	})

	t.Run("pages split the trainings and report the total", func(t *testing.T) {
		t.Parallel()

		trainer, attendee := newUser(shared.RoleTrainer), newOnboardedAttendee(ctx, t, clients)
		hours := tests.TomorrowHours(t, 3)
		tests.MakeHoursAvailable(ctx, t, clients, trainer.Token, hours)
		for _, hour := range hours {
			resp := scheduleTraining(ctx, t, clients, attendee.Token, trainingsclient.ScheduleTrainingRequest{
				TrainerUuid: trainer.UUID, TrainerUsername: tests.Username(trainer.UUID), Hour: hour, Notes: "n",
			})
			require.Equal(t, http.StatusCreated, resp.StatusCode(), string(resp.Body))
		}

		pageSize, first, second, beyond := 2, 1, 2, 3
		page1 := getPage(t, attendee.Token, &trainingsclient.GetUserTrainingsParams{Page: &first, PageSize: &pageSize})
		require.Len(t, page1.Items, 2)
		assert.Equal(t, 1, page1.Pagination.Page)
		assert.Equal(t, 2, page1.Pagination.PageSize)
		assert.Equal(t, 3, page1.Pagination.Total)
		assert.True(t, hours[0].Equal(page1.Items[0].Hour))
		assert.True(t, hours[1].Equal(page1.Items[1].Hour))

		page2 := getPage(t, attendee.Token, &trainingsclient.GetUserTrainingsParams{Page: &second, PageSize: &pageSize})
		require.Len(t, page2.Items, 1)
		assert.Equal(t, 3, page2.Pagination.Total)
		assert.True(t, hours[2].Equal(page2.Items[0].Hour))

		page3 := getPage(t, attendee.Token, &trainingsclient.GetUserTrainingsParams{Page: &beyond, PageSize: &pageSize})
		assert.Empty(t, page3.Items)
		assert.Equal(t, 3, page3.Pagination.Total)

		// defaults: first page of 10
		defaults := getPage(t, attendee.Token, &trainingsclient.GetUserTrainingsParams{})
		assert.Len(t, defaults.Items, 3)
		assert.Equal(t, 1, defaults.Pagination.Page)
		assert.Equal(t, 10, defaults.Pagination.PageSize)
	})

	t.Run("invalid page or page size is rejected", func(t *testing.T) {
		t.Parallel()

		token := newOnboardedAttendee(ctx, t, clients).Token
		zero, tooMany := 0, 51

		for name, params := range map[string]*trainingsclient.GetUserTrainingsParams{
			"zero page":         {Page: &zero},
			"zero page size":    {PageSize: &zero},
			"page size too big": {PageSize: &tooMany},
		} {
			resp, err := clients.Trainings.GetUserTrainingsWithResponse(ctx, params, tests.WithAuth(token))
			require.NoError(t, err, name)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode(), name+": "+string(resp.Body))
		}
	})

	t.Run("missing token is unauthorized", func(t *testing.T) {
		t.Parallel()

		resp, err := clients.Trainings.GetUserTrainingsWithResponse(ctx, &trainingsclient.GetUserTrainingsParams{})
		require.NoError(t, err)

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode(), string(resp.Body))
	})
}

func TestCancelTraining(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	clients := tests.NewTestClients(t)

	schedule := func(t *testing.T) (trainer, attendee testUser, hour time.Time, trainingUUID string) {
		t.Helper()

		trainer, attendee = newUser(shared.RoleTrainer), newOnboardedAttendee(ctx, t, clients)
		hour = tests.TomorrowHours(t, 1)[0]
		tests.MakeHoursAvailable(ctx, t, clients, trainer.Token, []time.Time{hour})

		resp := scheduleTraining(ctx, t, clients, attendee.Token, trainingsclient.ScheduleTrainingRequest{
			TrainerUuid:     trainer.UUID,
			TrainerUsername: tests.Username(trainer.UUID),
			Hour:            hour,
			Notes:           "legs day",
		})
		require.Equal(t, http.StatusCreated, resp.StatusCode(), string(resp.Body))
		require.NotNil(t, resp.JSON201)

		return trainer, attendee, hour, resp.JSON201.Uuid
	}

	t.Run("attendee cancels and the hour is available again", func(t *testing.T) {
		t.Parallel()

		trainer, attendee, hour, trainingUUID := schedule(t)

		resp, err := clients.Trainings.CancelTrainingWithResponse(ctx, trainingUUID, tests.WithAuth(attendee.Token))
		require.NoError(t, err)

		require.Equal(t, http.StatusNoContent, resp.StatusCode(), string(resp.Body))
		assert.Equal(t, domain.Available, hourStatus(ctx, t, clients, trainer, hour))
	})

	t.Run("trainer cancels", func(t *testing.T) {
		t.Parallel()

		trainer, _, _, trainingUUID := schedule(t)

		resp, err := clients.Trainings.CancelTrainingWithResponse(ctx, trainingUUID, tests.WithAuth(trainer.Token))
		require.NoError(t, err)

		assert.Equal(t, http.StatusNoContent, resp.StatusCode(), string(resp.Body))
	})

	t.Run("a training can't be canceled twice", func(t *testing.T) {
		t.Parallel()

		_, attendee, _, trainingUUID := schedule(t)

		first, err := clients.Trainings.CancelTrainingWithResponse(ctx, trainingUUID, tests.WithAuth(attendee.Token))
		require.NoError(t, err)
		require.Equal(t, http.StatusNoContent, first.StatusCode(), string(first.Body))

		second, err := clients.Trainings.CancelTrainingWithResponse(ctx, trainingUUID, tests.WithAuth(attendee.Token))
		require.NoError(t, err)
		assert.Equal(t, http.StatusConflict, second.StatusCode(), string(second.Body))
	})

	t.Run("a stranger can't cancel", func(t *testing.T) {
		t.Parallel()

		_, _, _, trainingUUID := schedule(t)

		resp, err := clients.Trainings.CancelTrainingWithResponse(
			ctx, trainingUUID, tests.WithAuth(newOnboardedAttendee(ctx, t, clients).Token),
		)
		require.NoError(t, err)

		assert.Equal(t, http.StatusForbidden, resp.StatusCode(), string(resp.Body))
	})

	t.Run("an unknown training is not found", func(t *testing.T) {
		t.Parallel()

		resp, err := clients.Trainings.CancelTrainingWithResponse(
			ctx, common.NewUUIDv7().String(), tests.WithAuth(newOnboardedAttendee(ctx, t, clients).Token),
		)
		require.NoError(t, err)

		assert.Equal(t, http.StatusNotFound, resp.StatusCode(), string(resp.Body))
	})

	t.Run("an invalid uuid is a bad request", func(t *testing.T) {
		t.Parallel()

		resp, err := clients.Trainings.CancelTrainingWithResponse(
			ctx, "not-a-uuid", tests.WithAuth(newOnboardedAttendee(ctx, t, clients).Token),
		)
		require.NoError(t, err)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode(), string(resp.Body))
	})
}
