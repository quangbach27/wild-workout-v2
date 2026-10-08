package commands

import (
	"context"
	"time"

	trainersClient "backend/internal/trainers/ports/module/client"
	"backend/internal/trainings/domain"
	usersClient "backend/internal/users/ports/module/client"
)

type CancelTrainingCmd struct {
	TrainingUUID domain.TrainingUUID
	User         domain.User
}

func (h *Handler) CancelTraining(ctx context.Context, cmd CancelTrainingCmd) error {
	canBeCancelForFree := false
	var attendeeUUID string
	var trainerUUID string
	var hour time.Time

	if err := h.trainingRepo.UpdateTraining(ctx, cmd.TrainingUUID, cmd.User, func(ctx context.Context, training *domain.Training) error {
		canBeCancelForFree = training.CanBeCanceledForFree()
		attendeeUUID = training.Attendee().ID()
		trainerUUID = training.Trainer().ID()
		hour = training.Hour()

		return training.Cancel(cmd.User)
	}); err != nil {
		return err
	}

	if _, err := h.modules.CancelHourSchedule(ctx, trainersClient.CancelHourScheduleRequest{
		TrainerUUID: trainerUUID,
		Hour:        hour,
	}); err != nil {
		return err
	}

	if canBeCancelForFree {
		_, err := h.modules.UpdateBalance(ctx, usersClient.UpdateBalanceRequest{
			UserUUID:     attendeeUUID,
			AmountChange: 1,
		})
		return err
	}

	return nil
}
