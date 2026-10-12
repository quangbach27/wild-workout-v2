package app

import (
	"context"

	"backend/internal/users/app/models"
)

// GetUser returns the user. It returns a 404 common.Error for an unknown user.
func (s *Service) GetUser(ctx context.Context, userUUID string) (*models.User, error) {
	if userUUID == "" {
		return nil, errEmptyUserUUID()
	}

	return s.userRepo.GetUserByUUID(ctx, userUUID)
}
