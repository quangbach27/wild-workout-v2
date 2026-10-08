package app

import (
	"context"

	common "github.com/quangbach27/golang-common"
)

type UpdateBalanceCmd struct {
	UserUUID     string
	AmountChange int
}

// UpdateBalance adds the amount change to the user's balance and returns the new balance.
func (s *Service) UpdateBalance(ctx context.Context, cmd UpdateBalanceCmd) (int, error) {
	if cmd.UserUUID == "" {
		return 0, errEmptyUserUUID()
	}

	return s.userRepo.UpdateBalance(ctx, cmd.UserUUID, cmd.AmountChange)
}

func errEmptyUserUUID() error {
	return common.NewInvalidInputError("invalid-user-uuid", "user uuid can't be empty")
}
