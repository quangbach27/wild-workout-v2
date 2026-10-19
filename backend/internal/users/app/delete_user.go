package app

import "context"

// DeleteUser removes the user. Deleting an unknown user is not an error, so a retried webhook is
// harmless. It returns a 400 common.Error for an empty uuid.
func (s *Service) DeleteUser(ctx context.Context, userUUID string) error {
	if userUUID == "" {
		return errEmptyUserUUID()
	}

	return s.userRepo.DeleteUser(ctx, userUUID)
}
