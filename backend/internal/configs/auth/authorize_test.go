package auth_test

import (
	"context"
	"net/http"
	"testing"

	common "github.com/quangbach27/golang-common"
	commonAuth "github.com/quangbach27/golang-common/http/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"backend/internal/configs/auth"
)

func TestAuthorize(t *testing.T) {
	t.Parallel()

	policy := auth.NewPolicy()
	withSession := func(role string) context.Context {
		return commonAuth.WithSession(context.Background(), &commonAuth.Session{Subject: "u1", Role: role})
	}

	t.Run("allowed returns the session", func(t *testing.T) {
		t.Parallel()

		session, err := auth.Authorize(withSession("trainer"), policy, auth.HourMakeAvailable)

		require.NoError(t, err)
		assert.Equal(t, "u1", session.Subject)
	})

	t.Run("denied is forbidden", func(t *testing.T) {
		t.Parallel()

		_, err := auth.Authorize(withSession("attendee"), policy, auth.HourMakeAvailable)

		assertHTTPStatus(t, err, http.StatusForbidden)
	})

	t.Run("empty or unknown role is forbidden", func(t *testing.T) {
		t.Parallel()

		for _, role := range []string{"", "admin"} {
			_, err := auth.Authorize(withSession(role), policy, auth.HourView)

			assertHTTPStatus(t, err, http.StatusForbidden)
		}
	})

	t.Run("no session is unauthorized", func(t *testing.T) {
		t.Parallel()

		_, err := auth.Authorize(context.Background(), policy, auth.HourView)

		assertHTTPStatus(t, err, http.StatusUnauthorized)
	})
}

func assertHTTPStatus(t *testing.T, err error, status int) {
	t.Helper()

	var e common.Error
	require.ErrorAs(t, err, &e)
	assert.Equal(t, status, e.HttpErrorCode)
}
