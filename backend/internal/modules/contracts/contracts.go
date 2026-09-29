package contracts

import (
	"errors"

	trainersClient "backend/internal/trainers/ports/module/client"
	trainingsClient "backend/internal/trainings/ports/module/client"
	usersClient "backend/internal/users/ports/module/client"
)

type Contracts struct {
	trainersClient.Trainers
	trainingsClient.Trainings
	usersClient.Users
}

// Verify checks that every module has registered its contract.
func (c *Contracts) Verify() error {
	var errs []error

	if c.Trainers == nil {
		errs = append(errs, errors.New("trainers contract is not registered"))
	}
	if c.Trainings == nil {
		errs = append(errs, errors.New("trainings contract is not registered"))
	}
	if c.Users == nil {
		errs = append(errs, errors.New("users contract is not registered"))
	}

	return errors.Join(errs...)
}
