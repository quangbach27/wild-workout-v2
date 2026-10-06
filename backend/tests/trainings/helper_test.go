//go:build component

package trainings_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	common "github.com/quangbach27/golang-common"
	"github.com/stretchr/testify/require"

	"backend/internal/shared"
	"backend/internal/trainers/domain"
	trainersclient "backend/internal/trainers/ports/http/clients"
	trainingsclient "backend/internal/trainings/ports/http/clients"
	"backend/tests"
)

type testUser struct {
	UUID  string
	Token string
}

func newUser(role shared.Role) testUser {
	userUUID := common.NewUUIDv7().String()
	return testUser{
		UUID:  userUUID,
		Token: tests.NewSession(userUUID, role.String()),
	}
}

func hourStatus(
	ctx context.Context,
	t *testing.T,
	clients tests.TestClients,
	trainer testUser,
	hour time.Time,
) domain.HourStatus {
	t.Helper()

	resp, err := clients.Trainers.GetTrainerHoursWithResponse(
		ctx,
		&trainersclient.GetTrainerHoursParams{From: hour, To: hour},
		tests.WithAuth(trainer.Token),
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode(), string(resp.Body))
	require.NotNil(t, resp.JSON200)

	for _, date := range *resp.JSON200 {
		for _, h := range date.Hours {
			if h.Hour.Equal(hour) {
				return h.Status
			}
		}
	}

	require.FailNow(t, "hour not found", hour.String())

	return domain.HourStatus{}
}

func scheduleTraining(
	ctx context.Context,
	t *testing.T,
	clients tests.TestClients,
	token string,
	body trainingsclient.ScheduleTrainingRequest,
) *trainingsclient.ScheduleTrainingClientResponse {
	t.Helper()

	resp, err := clients.Trainings.ScheduleTrainingWithResponse(ctx, body, tests.WithAuth(token))
	require.NoError(t, err)

	return resp
}
