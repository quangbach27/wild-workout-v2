package shared

import common "github.com/quangbach27/golang-common"

type Role struct {
	common.Enum[RoleType]
}

type RoleType string

func (t RoleType) Values() []string {
	return []string{"trainer", "attendee"}
}

var (
	RoleTrainer  = common.MustEnum[Role]("trainer")
	RoleAttendee = common.MustEnum[Role]("attendee")
)
