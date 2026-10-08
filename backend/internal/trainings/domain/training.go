package domain

import (
	"context"
	"time"

	common "github.com/quangbach27/golang-common"

	"backend/internal/shared"
)

type TrainingsRepository interface {
	AddTraining(ctx context.Context, tr *Training) error

	GetTraining(ctx context.Context, trainingUUID string, user User) (*Training, error)

	UpdateTraining(
		ctx context.Context,
		trainingUUID string,
		user User,
		updateFn func(ctx context.Context, tr *Training) (*Training, error),
	) error
}

// FreeCancellationPeriod is how long before the training it can still be canceled or
// rescheduled without losing credits.
const FreeCancellationPeriod = time.Hour * 24

type TrainingUUID struct {
	common.UUID
}

type Training struct {
	id TrainingUUID

	hour  time.Time
	notes string

	attendee User
	trainer  User

	proposedNewTime *time.Time
	moveProposedBy  *User

	canceled bool
}

func (t *Training) ID() TrainingUUID {
	return t.id
}

func (t *Training) Hour() time.Time {
	return t.hour
}

func (t *Training) Attendee() User {
	return t.attendee
}

func (t *Training) Trainer() User {
	return t.trainer
}

func (t *Training) Notes() string {
	return t.notes
}

func (t *Training) IsCanceled() bool {
	return t.canceled
}

func (t *Training) ProposedNewTime() *time.Time {
	return t.proposedNewTime
}

func (t *Training) MoveProposedBy() *User {
	return t.moveProposedBy
}

// NewTraining validates the input and returns a common.Error (400) whose details list each failure
// (invalid-training-hour, invalid-training-attendee, invalid-training-trainer, invalid-training-notes).
func NewTraining(attendee, trainer User, hour time.Time, notes string) (*Training, error) {
	errDetails := []common.ErrorDetails{}

	if hour.IsZero() || hour.Before(time.Now()) {
		errDetails = append(errDetails, common.ErrorDetails{
			EntityType: "Training",
			ErrorSlug:  "invalid-training-hour",
			Message:    "hour must be set and can't be in the past",
		})
	}

	if !attendee.Role().Equal(shared.RoleAttendee.Enum) {
		errDetails = append(errDetails, common.ErrorDetails{
			EntityType: "Training",
			ErrorSlug:  "invalid-training-attendee",
			Message:    "attendee must have the attendee role",
		})
	}

	if !trainer.Role().Equal(shared.RoleTrainer.Enum) {
		errDetails = append(errDetails, common.ErrorDetails{
			EntityType: "Training",
			ErrorSlug:  "invalid-training-trainer",
			Message:    "trainer must have the trainer role",
		})
	}

	if len(errDetails) != 0 {
		return nil, common.NewInvalidInputError("invalid-training", "training is not valid").WithDetails(errDetails)
	}

	return &Training{
		id:       TrainingUUID{UUID: common.NewUUIDv7()},
		hour:     hour,
		notes:    notes,
		attendee: attendee,
		trainer:  trainer,
	}, nil
}

func (t *Training) CanBeCanceledForFree() bool {
	return time.Until(t.hour) >= FreeCancellationPeriod
}

func (t *Training) Cancel(user User) error {
	if err := CanUserSeeTraining(user, t); err != nil {
		return common.NewForbiddenError("failed-to-cancel", "%s", err.Error()).WithInternalError(err)
	}

	if t.canceled {
		return common.NewConflictError("failed-to-cancel", "training is already canceled")
	}

	t.canceled = true
	return nil
}

func (t *Training) ProposeReschedule(newTime time.Time, proposedBy User) error {
	if err := CanUserSeeTraining(proposedBy, t); err != nil {
		return common.NewForbiddenError("failed-to-propose-reschedule", "%s", err.Error()).WithInternalError(err)
	}

	if t.canceled {
		return common.NewConflictError("failed-to-propose-reschedule", "training is canceled")
	}

	if newTime.IsZero() || newTime.Before(time.Now()) {
		return common.NewInvalidInputError("failed-to-propose-reschedule", "new time must be set and can't be in the past")
	}

	t.moveProposedBy = &proposedBy
	t.proposedNewTime = &newTime
	return nil
}

func (t *Training) IsRescheduleProposed() bool {
	return t.moveProposedBy != nil && t.proposedNewTime != nil
}

// ApproveReschedule moves the training to the proposed time. Only the party that did not propose
// the move can approve it.
func (t *Training) ApproveReschedule(user User) error {
	if err := CanUserSeeTraining(user, t); err != nil {
		return common.NewForbiddenError("failed-to-approve-reschedule", "%s", err.Error()).WithInternalError(err)
	}

	if t.canceled {
		return common.NewConflictError("failed-to-approve-reschedule", "training is canceled")
	}

	if !t.IsRescheduleProposed() {
		return common.NewConflictError("failed-to-approve-reschedule", "no training reschedule was requested yet")
	}

	if t.moveProposedBy.Equal(user) {
		return common.NewForbiddenError("failed-to-approve-reschedule", "the proposer can't approve their own reschedule")
	}

	t.hour = *t.proposedNewTime

	t.proposedNewTime = nil
	t.moveProposedBy = nil

	return nil
}

func (t *Training) RejectReschedule(user User) error {
	if err := CanUserSeeTraining(user, t); err != nil {
		return common.NewForbiddenError("failed-to-reject-reschedule", "%s", err.Error()).WithInternalError(err)
	}

	if !t.IsRescheduleProposed() {
		return common.NewConflictError("failed-to-reject-reschedule", "no training reschedule was requested yet")
	}

	t.proposedNewTime = nil
	t.moveProposedBy = nil

	return nil
}
