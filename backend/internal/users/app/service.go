package app

import "backend/internal/users/app/models"

type Service struct {
	userRepo models.UserRepository
}

func NewService(userRepo models.UserRepository) *Service {
	if userRepo == nil {
		panic("userRepo can't be nil")
	}
	return &Service{
		userRepo: userRepo,
	}
}
