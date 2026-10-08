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
	usersclient "backend/internal/users/ports/http/clients"
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

// newOnboardedAttendee creates an onboarded attendee with credits, as scheduling debits one.
func newOnboardedAttendee(ctx context.Context, t *testing.T, clients tests.TestClients) testUser {
	t.Helper()

	attendee := newUser(shared.RoleAttendee)

	onboard, err := clients.Users.OnboardUserWithResponse(
		ctx,
		usersclient.OnboardUserJSONRequestBody{DisplayName: tests.Username(attendee.UUID)},
		tests.WithAuth(attendee.Token),
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, onboard.StatusCode(), string(onboard.Body))

	topUp, err := clients.Users.TopUpBalanceWithResponse(
		ctx, usersclient.TopUpBalanceJSONRequestBody{Amount: 5}, tests.WithAuth(attendee.Token),
	)
	require.NoError(t, err)
	require.Less(t, topUp.StatusCode(), http.StatusMultipleChoices, string(topUp.Body))

	return attendee
}
