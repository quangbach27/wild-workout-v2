package auth_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"backend/internal/configs/auth"
	"backend/internal/shared"
)

func TestPolicy_Allows(t *testing.T) {
	t.Parallel()

	type grant struct {
		resource auth.Resource
		actions  []auth.Action
	}

	hourAll := []auth.Action{
		auth.ActionMakeAvailable, auth.ActionMakeUnavailable, auth.ActionSchedule,
		auth.ActionCancel, auth.ActionView,
	}
	trainingAll := []auth.Action{
		auth.ActionSchedule, auth.ActionRescheduleRequest, auth.ActionRescheduleResponse,
		auth.ActionCancel, auth.ActionView,
	}

	tests := map[string]struct {
		role    shared.Role
		granted []grant
	}{
		"trainer": {shared.RoleTrainer, []grant{
			{auth.ResourceHour, []auth.Action{auth.ActionMakeAvailable, auth.ActionMakeUnavailable, auth.ActionView}},
			{auth.ResourceTraining, []auth.Action{auth.ActionRescheduleRequest, auth.ActionRescheduleResponse, auth.ActionCancel, auth.ActionView}},
		}},
		"attendee": {shared.RoleAttendee, []grant{
			{auth.ResourceHour, []auth.Action{auth.ActionSchedule, auth.ActionCancel, auth.ActionView}},
			{auth.ResourceTraining, trainingAll},
		}},
	}

	policy := auth.NewPolicy()

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			allowed := map[auth.Permission]bool{}
			for _, g := range tc.granted {
				for _, a := range g.actions {
					p, err := auth.NewPermission(g.resource, a)
					require.NoError(t, err)
					allowed[p] = true
				}
			}

			all := map[auth.Resource][]auth.Action{
				auth.ResourceHour:     hourAll,
				auth.ResourceTraining: trainingAll,
			}
			for resource, actions := range all {
				for _, a := range actions {
					p, err := auth.NewPermission(resource, a)
					require.NoError(t, err)
					assert.Equal(t, allowed[p], policy.Allows([]shared.Role{tc.role}, p), "%s %s", resource, a)
				}
			}
		})
	}
}

func TestPolicy_Allows_UnionOfRoles(t *testing.T) {
	t.Parallel()

	p, err := auth.NewPermission(auth.ResourceHour, auth.ActionMakeAvailable)
	require.NoError(t, err)

	policy := auth.NewPolicy()
	assert.False(t, policy.Allows([]shared.Role{shared.RoleAttendee}, p))
	assert.True(t, policy.Allows([]shared.Role{shared.RoleAttendee, shared.RoleTrainer}, p))
	assert.False(t, policy.Allows(nil, p))
}

func TestNewPermission_InvalidPair(t *testing.T) {
	t.Parallel()

	_, err := auth.NewPermission(auth.ResourceHour, auth.ActionRescheduleRequest)
	assert.ErrorIs(t, err, auth.ErrInvalidPermission)

	_, err = auth.NewPermission(auth.ResourceTraining, auth.ActionMakeAvailable)
	assert.ErrorIs(t, err, auth.ErrInvalidPermission)
}
