//go:build component

package users_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"backend/internal/shared"
	usersclient "backend/internal/users/ports/http/clients"
	"backend/tests"
)

func TestGetTrainers(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	clients := tests.NewTestClients(t)
	attendee := newUser(shared.RoleAttendee)

	onboardTrainers(ctx, t, 3)

	t.Run("returns one page with the total", func(t *testing.T) {
		t.Parallel()

		pageSize := 2
		resp, err := clients.Users.GetTrainersWithResponse(
			ctx,
			&usersclient.GetTrainersParams{PageSize: &pageSize},
			tests.WithAuth(attendee.Token),
		)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode(), string(resp.Body))
		require.NotNil(t, resp.JSON200)

		assert.Len(t, resp.JSON200.Items, 2)
		assert.Equal(t, 1, resp.JSON200.Pagination.Page)
		assert.Equal(t, 2, resp.JSON200.Pagination.PageSize)
		assert.GreaterOrEqual(t, resp.JSON200.Pagination.Total, 3)
	})

	t.Run("invalid pagination", func(t *testing.T) {
		t.Parallel()

		page, pageSize := 0, 51
		resp, err := clients.Users.GetTrainersWithResponse(
			ctx,
			&usersclient.GetTrainersParams{Page: &page, PageSize: &pageSize},
			tests.WithAuth(attendee.Token),
		)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode(), string(resp.Body))
	})

	t.Run("requires authentication", func(t *testing.T) {
		t.Parallel()

		resp, err := clients.Users.GetTrainersWithResponse(ctx, &usersclient.GetTrainersParams{})
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode(), string(resp.Body))
	})
}

func TestOnboardUser(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	clients := tests.NewTestClients(t)

	body := usersclient.OnboardUserJSONRequestBody{DisplayName: "Jane"}

	t.Run("onboards the session user once", func(t *testing.T) {
		t.Parallel()

		trainer := newUser(shared.RoleTrainer)

		onboardUser(ctx, t, clients, trainer, body.DisplayName)

		resp, err := clients.Users.OnboardUserWithResponse(ctx, body, tests.WithAuth(trainer.Token))
		require.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.StatusCode(), string(resp.Body))
	})

	t.Run("rejects an empty display name", func(t *testing.T) {
		t.Parallel()

		attendee := newUser(shared.RoleAttendee)

		resp, err := clients.Users.OnboardUserWithResponse(
			ctx, usersclient.OnboardUserJSONRequestBody{}, tests.WithAuth(attendee.Token),
		)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode(), string(resp.Body))
	})

	t.Run("requires authentication", func(t *testing.T) {
		t.Parallel()

		resp, err := clients.Users.OnboardUserWithResponse(ctx, body)
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode(), string(resp.Body))
	})
}

func TestTopUpBalance(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	clients := tests.NewTestClients(t)

	topUp := func(t *testing.T, user testUser, amount int) *usersclient.TopUpBalanceClientResponse {
		t.Helper()

		resp, err := clients.Users.TopUpBalanceWithResponse(
			ctx, usersclient.TopUpBalanceJSONRequestBody{Amount: amount}, tests.WithAuth(user.Token),
		)
		require.NoError(t, err)

		return resp
	}

	t.Run("top ups accumulate", func(t *testing.T) {
		t.Parallel()

		user := newUser(shared.RoleAttendee)
		onboardUser(ctx, t, clients, user, "Jane")

		resp := topUp(t, user, 100)
		require.Equal(t, http.StatusOK, resp.StatusCode(), string(resp.Body))
		assert.Equal(t, 100, resp.JSON200.Balance)

		resp = topUp(t, user, 50)
		require.Equal(t, http.StatusOK, resp.StatusCode(), string(resp.Body))
		assert.Equal(t, 150, resp.JSON200.Balance)
	})

	t.Run("invalid amount", func(t *testing.T) {
		t.Parallel()

		user := newUser(shared.RoleAttendee)
		onboardUser(ctx, t, clients, user, "Jane")

		for _, amount := range []int{0, -5, 1_000_001} {
			assert.Equal(t, http.StatusBadRequest, topUp(t, user, amount).StatusCode(), "amount %d", amount)
		}
	})

	t.Run("user is not onboarded", func(t *testing.T) {
		t.Parallel()

		resp := topUp(t, newUser(shared.RoleAttendee), 100)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode(), string(resp.Body))
	})

	t.Run("requires authentication", func(t *testing.T) {
		t.Parallel()

		resp, err := clients.Users.TopUpBalanceWithResponse(ctx, usersclient.TopUpBalanceJSONRequestBody{Amount: 1})
		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode(), string(resp.Body))
	})
}
