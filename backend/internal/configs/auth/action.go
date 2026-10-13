package auth

import common "github.com/quangbach27/golang-common"

type Action struct {
	common.Enum[ActionType]
}

type ActionType string

func (a ActionType) Values() []string {
	return []string{
		"make-available",
		"make-unavailable",
		"schedule",
		"cancel",
		"view",
		"reschedule-request",
		"reschedule-response",
	}
}

var (
	ActionMakeAvailable      = common.MustEnum[Action]("make-available")
	ActionMakeUnavailable    = common.MustEnum[Action]("make-unavailable")
	ActionSchedule           = common.MustEnum[Action]("schedule")
	ActionCancel             = common.MustEnum[Action]("cancel")
	ActionView               = common.MustEnum[Action]("view")
	ActionRescheduleRequest  = common.MustEnum[Action]("reschedule-request")
	ActionRescheduleResponse = common.MustEnum[Action]("reschedule-response")
)
