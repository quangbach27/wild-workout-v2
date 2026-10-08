package app

import "context"

// GetBalance returns the user's current balance.
func (s *Service) GetBalance(ctx context.Context, userUUID string) (int, error) {
	if userUUID == "" {
		return 0, errEmptyUserUUID()
	}

	user, err := s.userRepo.GetUserByUUID(ctx, userUUID)
	if err != nil {
		return 0, err
	}

	return user.Balance(), nil
}
