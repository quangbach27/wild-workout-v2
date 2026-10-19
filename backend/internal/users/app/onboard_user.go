package app

import (
	"context"

	"backend/internal/shared"
	"backend/internal/users/app/models"
)

// attendeeStartBalance is the credit a new attendee gets when onboarded.
const attendeeStartBalance = 20

type OnboardUserCmd struct {
	UUID        string
	DisplayName string
	Role        shared.Role
}

func (s *Service) OnboardUser(ctx context.Context, cmd OnboardUserCmd) error {
	startBalance := 0
	if cmd.Role == shared.RoleAttendee {
		startBalance = attendeeStartBalance
	}

	user, err := models.NewUser(cmd.UUID, cmd.DisplayName, cmd.Role, startBalance)
	if err != nil {
		return err
	}

	return s.userRepo.OnboardUser(ctx, user)
}
