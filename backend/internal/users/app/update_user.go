package app

import (
	"context"

	common "github.com/quangbach27/golang-common"

	"backend/internal/shared"
	"backend/internal/users/app/models"
)

// UpdateUserCmd carries the fields to change; nil leaves a field unchanged.
type UpdateUserCmd struct {
	UUID        string
	DisplayName *string
	Role        *shared.Role
}

// UpdateUser changes the display name and/or role. It returns a 400 common.Error for invalid input
// or an empty update, and a 404 for an unknown user.
func (s *Service) UpdateUser(ctx context.Context, cmd UpdateUserCmd) error {
	if cmd.UUID == "" {
		return errEmptyUserUUID()
	}
	if cmd.DisplayName == nil && cmd.Role == nil {
		return common.NewInvalidInputError("invalid-user-update", "nothing to update")
	}

	return s.userRepo.UpdateUser(ctx, cmd.UUID, func(user *models.User) error {
		return user.Update(cmd.DisplayName, cmd.Role)
	})
}
