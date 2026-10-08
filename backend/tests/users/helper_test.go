//go:build component

package users_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	common "github.com/quangbach27/golang-common"
	"github.com/stretchr/testify/require"

	"backend/internal/configs"
	"backend/internal/shared"
	"backend/internal/users/adapters/db"
	"backend/internal/users/app"
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

// onboardUser onboards the user through the HTTP API.
func onboardUser(
	ctx context.Context,
	t *testing.T,
	clients tests.TestClients,
	user testUser,
	displayName string,
) {
	t.Helper()

	resp, err := clients.Users.OnboardUserWithResponse(
		ctx,
		usersclient.OnboardUserJSONRequestBody{DisplayName: displayName},
		tests.WithAuth(user.Token),
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode(), string(resp.Body))
}

// onboardTrainers seeds trainers directly through the service, as the HTTP onboarding only creates
// the session's own user.
func onboardTrainers(ctx context.Context, t *testing.T, count int) {
	t.Helper()

	pool, err := pgxpool.New(ctx, configs.NewConfig().DB.Dsn())
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	service := app.NewService(db.NewUserRepository(pool))
	for range count {
		require.NoError(t, service.OnboardUser(ctx, app.OnboardUserCmd{
			UUID:        common.NewUUIDv7().String(),
			DisplayName: "trainer-" + common.NewUUIDv7().String(),
			Role:        shared.RoleTrainer,
		}))
	}
}
