package commands

import (
	"errors"

	"backend/internal/trainers/domain"
)

type Handler struct {
	hourRepo domain.HourRepository
}

func NewHandler(
	hourRepo domain.HourRepository,
) *Handler {
	var errs []error
	if hourRepo == nil {
		errs = append(errs, errors.New("hourRepo can't be nil"))
	}

	if len(errs) != 0 {
		panic(errors.Join(errs...))
	}

	return &Handler{
		hourRepo: hourRepo,
	}
}
