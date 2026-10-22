package app

import (
	"context"

	"backend/internal/shared"
)

type UserRepository interface {
	OnboardUser(ctx context.Context, user *shared.User) error
	UpdateUser(
		ctx context.Context,
		userID string,
		updateFn func(user *shared.User) error,
	) error
	UpdateBalance(ctx context.Context, userID string, amountChange int) (int, error)
	DeleteUser(ctx context.Context, userID string) error
}

type Service struct {
	userRepo UserRepository
}

func NewService(userRepo UserRepository) *Service {
	if userRepo == nil {
		panic("userRepo can't be nil")
	}
	return &Service{
		userRepo: userRepo,
	}
}
