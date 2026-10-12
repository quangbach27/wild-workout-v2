package auth

import "backend/internal/shared"

// Policy maps each role to the permissions it grants.
type Policy struct {
	grants map[shared.Role]map[Permission]struct{}
}

func NewPolicy() *Policy {
	return &Policy{grants: map[shared.Role]map[Permission]struct{}{
		shared.RoleTrainer: permissionSet(
			mustPermission(ResourceHour, ActionMakeAvailable),
			mustPermission(ResourceHour, ActionMakeUnavailable),
			mustPermission(ResourceHour, ActionView),
			mustPermission(ResourceTraining, ActionRescheduleRequest),
			mustPermission(ResourceTraining, ActionRescheduleResponse),
			mustPermission(ResourceTraining, ActionCancel),
			mustPermission(ResourceTraining, ActionView),
		),
		shared.RoleAttendee: permissionSet(
			mustPermission(ResourceHour, ActionSchedule),
			mustPermission(ResourceHour, ActionCancel),
			mustPermission(ResourceHour, ActionView),
			mustPermission(ResourceTraining, ActionSchedule),
			mustPermission(ResourceTraining, ActionRescheduleRequest),
			mustPermission(ResourceTraining, ActionRescheduleResponse),
			mustPermission(ResourceTraining, ActionCancel),
			mustPermission(ResourceTraining, ActionView),
		),
	}}
}

// Allows reports whether any of the roles grants the permission.
func (p *Policy) Allows(roles []shared.Role, permission Permission) bool {
	for _, role := range roles {
		if _, ok := p.grants[role][permission]; ok {
			return true
		}
	}

	return false
}

func permissionSet(permissions ...Permission) map[Permission]struct{} {
	set := make(map[Permission]struct{}, len(permissions))
	for _, p := range permissions {
		set[p] = struct{}{}
	}

	return set
}
