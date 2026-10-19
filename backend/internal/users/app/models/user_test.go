package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"backend/internal/shared"
	"backend/internal/users/app/models"
)

func TestUser_Update(t *testing.T) {
	t.Parallel()

	newName := "New Name"
	emptyName := ""
	trainer := shared.RoleTrainer
	var zeroRole shared.Role

	testCases := []struct {
		name        string
		displayName *string
		role        *shared.Role
		wantErr     bool
		wantName    string
		wantRole    shared.Role
	}{
		{name: "display name only", displayName: &newName, wantName: newName, wantRole: shared.RoleAttendee},
		{name: "role only", role: &trainer, wantName: "Old Name", wantRole: shared.RoleTrainer},
		{name: "both", displayName: &newName, role: &trainer, wantName: newName, wantRole: shared.RoleTrainer},
		{name: "empty display name", displayName: &emptyName, wantErr: true, wantName: "Old Name", wantRole: shared.RoleAttendee},
		{name: "zero role", role: &zeroRole, wantErr: true, wantName: "Old Name", wantRole: shared.RoleAttendee},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			user, err := models.NewUser("uuid", "Old Name", shared.RoleAttendee, 0)
			require.NoError(t, err)

			err = user.Update(tc.displayName, tc.role)

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.wantName, user.DisplayName())
			assert.Equal(t, tc.wantRole, user.Role())
		})
	}
}

func TestNewUser_StartBalance(t *testing.T) {
	t.Parallel()

	user, err := models.NewUser("uuid", "Name", shared.RoleAttendee, 20)
	require.NoError(t, err)
	assert.Equal(t, 20, user.Balance())
}
