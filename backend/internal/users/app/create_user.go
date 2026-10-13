package app

import (
	"context"

	"backend/internal/shared"
	"backend/internal/users/app/models"
)

type CreateUserCmd struct {
	UUID        string
	DisplayName string
	Role        shared.Role
}

// CreateUser validates and stores a new user with a zero balance. It returns a 400 common.Error
// for invalid input and a 409 if the user already exists.
func (s *Service) CreateUser(ctx context.Context, cmd CreateUserCmd) error {
	user, err := models.NewUser(cmd.UUID, cmd.DisplayName, cmd.Role)
	if err != nil {
		return err
	}

	return s.userRepo.OnboardUser(ctx, user)
}
