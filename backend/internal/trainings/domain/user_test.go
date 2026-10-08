package domain_test

import (
	"net/http"
	"testing"

	common "github.com/quangbach27/golang-common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"backend/internal/shared"
	"backend/internal/trainings/domain"
)

func TestNewUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		id        string
		username  string
		role      shared.Role
		wantSlugs []string // empty = valid user, otherwise in order
	}{
		{name: "trainer", id: "user-1", username: "alice", role: shared.RoleTrainer},
		{name: "attendee", id: "user-1", username: "bob", role: shared.RoleAttendee},
		{name: "empty id", username: "alice", role: shared.RoleTrainer, wantSlugs: []string{"invalid-user-id"}},
		{name: "empty username", id: "user-1", role: shared.RoleTrainer, wantSlugs: []string{"invalid-user-username"}},
		{name: "empty role", id: "user-1", username: "alice", wantSlugs: []string{"invalid-user-role"}},
		{
			name:      "everything empty",
			wantSlugs: []string{"invalid-user-id", "invalid-user-username", "invalid-user-role"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			user, err := domain.NewUser(tc.id, tc.username, tc.role)

			if len(tc.wantSlugs) > 0 {
				var commonErr common.Error
				require.ErrorAs(t, err, &commonErr)
				assert.Equal(t, http.StatusBadRequest, commonErr.HttpErrorCode)
				assert.Equal(t, "invalid-user", commonErr.ErrorSlug)
				require.Len(t, commonErr.Details, len(tc.wantSlugs))
				for i, slug := range tc.wantSlugs {
					assert.Equal(t, "User", commonErr.Details[i].EntityType)
					assert.Equal(t, slug, commonErr.Details[i].ErrorSlug)
				}
				assert.Zero(t, user)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.id, user.ID())
			assert.Equal(t, tc.username, user.Username())
			assert.Equal(t, tc.role, user.Role())
		})
	}
}
