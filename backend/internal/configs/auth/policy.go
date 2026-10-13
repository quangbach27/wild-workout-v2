package auth

import "backend/internal/shared"

// Policy maps each role to the permissions it grants.
type Policy struct {
	grants map[shared.Role]map[Permission]struct{}
}

func NewPolicy() *Policy {
	return &Policy{grants: map[shared.Role]map[Permission]struct{}{
		shared.RoleTrainer: permissionSet(
			HourMakeAvailable,
			HourMakeUnavailable,
			HourView,
			TrainingRescheduleRequest,
			TrainingRescheduleResponse,
			TrainingCancel,
			TrainingView,
		),
		shared.RoleAttendee: permissionSet(
			HourSchedule,
			HourCancel,
			HourView,
			TrainingSchedule,
			TrainingRescheduleRequest,
			TrainingRescheduleResponse,
			TrainingCancel,
			TrainingView,
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
