package app

import (
	"context"

	"backend/internal/shared"
	"backend/internal/users/app/models"
)

type OnboardUserCmd struct {
	UUID        string
	DisplayName string
	Role        shared.Role
}

func (s *Service) OnboardUser(ctx context.Context, cmd OnboardUserCmd) error {
	user, err := models.NewUser(cmd.UUID, cmd.DisplayName, cmd.Role)
	if err != nil {
		return err
	}

	return s.userRepo.OnboardUser(ctx, user)
}
