package commands

import (
	"context"
	"errors"

	trainersClient "backend/internal/trainers/ports/module/client"
	"backend/internal/trainings/domain"
)

// ModulesContract lists what the trainings module needs from other modules. The global
// contracts.Contracts satisfies it.
type ModulesContract interface {
	ScheduleHour(
		ctx context.Context,
		req trainersClient.ScheduleHourRequest,
	) (trainersClient.ScheduleHourResponse, error)
	CancelHourSchedule(
		ctx context.Context,
		req trainersClient.CancelHourScheduleRequest,
	) (trainersClient.CancelHourScheduleResponse, error)
}

type Handler struct {
	trainingRepo domain.TrainingRepository
	modules      ModulesContract
}

func NewHandler(
	trainingRepo domain.TrainingRepository,
	modules ModulesContract,
) *Handler {
	var errs []error
	if trainingRepo == nil {
		errs = append(errs, errors.New("trainingRepo can't be nil"))
	}
	if modules == nil {
		errs = append(errs, errors.New("modules can't be nil"))
	}

	if len(errs) != 0 {
		panic(errors.Join(errs...))
	}

	return &Handler{
		trainingRepo: trainingRepo,
		modules:      modules,
	}
}
