package auth

import (
	"errors"
	"fmt"
)

var ErrInvalidPermission = errors.New("invalid permission")

// Permission is an action on a resource. It is comparable, so it can be a map key.
type Permission struct {
	Resource Resource
	Action   Action
}

var validActions = map[Resource][]Action{
	ResourceHour: {
		ActionMakeAvailable,
		ActionMakeUnavailable,
		ActionSchedule,
		ActionCancel,
		ActionView,
	},
	ResourceTraining: {
		ActionSchedule,
		ActionRescheduleRequest,
		ActionRescheduleResponse,
		ActionCancel,
		ActionView,
	},
}

// NewPermission rejects actions that don't exist for the resource.
func NewPermission(resource Resource, action Action) (Permission, error) {
	for _, valid := range validActions[resource] {
		if valid == action {
			return Permission{Resource: resource, Action: action}, nil
		}
	}

	return Permission{}, fmt.Errorf(
		"%w: action %q on resource %q", ErrInvalidPermission, action, resource,
	)
}

func mustPermission(resource Resource, action Action) Permission {
	p, err := NewPermission(resource, action)
	if err != nil {
		panic(err)
	}

	return p
}

// Permissions checked by the HTTP handlers.
var (
	HourMakeAvailable   = mustPermission(ResourceHour, ActionMakeAvailable)
	HourMakeUnavailable = mustPermission(ResourceHour, ActionMakeUnavailable)
	HourSchedule        = mustPermission(ResourceHour, ActionSchedule)
	HourCancel          = mustPermission(ResourceHour, ActionCancel)
	HourView            = mustPermission(ResourceHour, ActionView)

	TrainingSchedule           = mustPermission(ResourceTraining, ActionSchedule)
	TrainingRescheduleRequest  = mustPermission(ResourceTraining, ActionRescheduleRequest)
	TrainingRescheduleResponse = mustPermission(ResourceTraining, ActionRescheduleResponse)
	TrainingCancel             = mustPermission(ResourceTraining, ActionCancel)
	TrainingView               = mustPermission(ResourceTraining, ActionView)
)
