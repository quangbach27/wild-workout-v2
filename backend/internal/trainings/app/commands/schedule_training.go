package commands

import (
	"context"
	"fmt"
	"time"

	"backend/internal/shared"
	"backend/internal/trainings/domain"
	"backend/internal/users/ports/module/client"

	trainersClient "backend/internal/trainers/ports/module/client"
)

type ScheduleTrainingCmd struct {
	AttendeeUUID string
	TrainerUUID  string
	Hour         time.Time
	Notes        string
}

// ScheduleTraining books the hour in the trainers module and then stores the training. If storing
// fails, the hour is released again.
func (h *Handler) ScheduleTraining(ctx context.Context, cmd ScheduleTrainingCmd) (domain.TrainingUUID, error) {
	attendeeUser, err := h.modules.GetUser(ctx, client.GetUserRequest{UserID: cmd.AttendeeUUID})
	if err != nil {
		return domain.TrainingUUID{}, err
	}

	trainerUser, err := h.modules.GetUser(ctx, client.GetUserRequest{UserID: cmd.TrainerUUID})
	if err != nil {
		return domain.TrainingUUID{}, err
	}

	attendee, err := domain.NewUser(cmd.AttendeeUUID, attendeeUser.DisplayName, shared.RoleAttendee)
	if err != nil {
		return domain.TrainingUUID{}, err
	}

	trainer, err := domain.NewUser(cmd.TrainerUUID, trainerUser.DisplayName, shared.RoleTrainer)
	if err != nil {
		return domain.TrainingUUID{}, err
	}

	training, err := domain.NewTraining(attendee, trainer, cmd.Hour, cmd.Notes)
	if err != nil {
		return domain.TrainingUUID{}, err
	}

	if _, err := h.modules.ScheduleHour(ctx, trainersClient.ScheduleHourRequest{
		TrainerUUID: cmd.TrainerUUID,
		Hour:        cmd.Hour,
	}); err != nil {
		return domain.TrainingUUID{}, err
	}

	if err := h.trainingRepo.AddTraining(ctx, training); err != nil {
		return domain.TrainingUUID{}, err
	}

	if _, err := h.modules.UpdateBalance(ctx, client.UpdateBalanceRequest{
		UserUUID:     cmd.AttendeeUUID,
		AmountChange: -1,
	}); err != nil {
		return domain.TrainingUUID{}, fmt.Errorf("failed to update balance: %w", err)
	}

	return training.ID(), nil
}
