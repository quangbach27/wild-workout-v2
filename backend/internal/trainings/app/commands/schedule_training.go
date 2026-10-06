package commands

import (
	"context"
	"log/slog"
	"time"

	"backend/internal/shared"
	trainersClient "backend/internal/trainers/ports/module/client"
	"backend/internal/trainings/domain"
)

type ScheduleTrainingCmd struct {
	AttendeeUUID     string
	AttendeeUsername string
	TrainerUUID      string
	TrainerUsername  string
	Hour             time.Time
	Notes            string
}

// ScheduleTraining books the hour in the trainers module and then stores the training. If storing
// fails, the hour is released again.
func (h *Handler) ScheduleTraining(ctx context.Context, cmd ScheduleTrainingCmd) (domain.TrainingUUID, error) {
	attendee, err := domain.NewUser(cmd.AttendeeUUID, cmd.AttendeeUsername, shared.RoleAttendee)
	if err != nil {
		return domain.TrainingUUID{}, err
	}

	trainer, err := domain.NewUser(cmd.TrainerUUID, cmd.TrainerUsername, shared.RoleTrainer)
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
		if _, cancelErr := h.modules.CancelHourSchedule(ctx, trainersClient.CancelHourScheduleRequest{
			TrainerUUID: cmd.TrainerUUID,
			Hour:        cmd.Hour,
		}); cancelErr != nil {
			slog.ErrorContext(ctx, "failed to release hour after training could not be stored",
				"trainer_uuid", cmd.TrainerUUID,
				"hour", cmd.Hour,
				"error", cancelErr,
			)
		}

		return domain.TrainingUUID{}, err
	}

	return training.ID(), nil
}
