package domain

import "context"

type TrainingRepository interface {
	AddTraining(ctx context.Context, training *Training) error

	// GetTraining returns a 404 common.Error if the training doesn't exist and a 403 one if
	// the user is neither its trainer nor its attendee.
	GetTraining(ctx context.Context, id TrainingUUID, user User) (*Training, error)

	// UpdateTraining loads the training under a row lock and runs updateFn on it inside a
	// transaction; an error from updateFn rolls everything back. Errors for a missing training
	// and for a user who can't see it are the same as GetTraining.
	UpdateTraining(
		ctx context.Context,
		id TrainingUUID,
		user User,
		updateFn func(ctx context.Context, training *Training) error,
	) error
}
