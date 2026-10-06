package domain

import (
	"time"

	"backend/internal/shared"
)

// UnmarshalTraining is only used in term of loading training from database.
func UnmarshalTraining(
	id TrainingUUID,
	hour time.Time,
	notes string,
	attendee, trainer User,
	proposedNewTime *time.Time,
	moveProposedBy *User,
	canceled bool,
) *Training {
	return &Training{
		id:              id,
		hour:            hour,
		notes:           notes,
		attendee:        attendee,
		trainer:         trainer,
		proposedNewTime: proposedNewTime,
		moveProposedBy:  moveProposedBy,
		canceled:        canceled,
	}
}

// UnmarshalUser is only used in term of loading user from database.
func UnmarshalUser(id, username string, role shared.Role) User {
	return User{id: id, username: username, role: role}
}
