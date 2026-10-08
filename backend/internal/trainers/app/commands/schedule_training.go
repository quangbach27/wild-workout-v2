package commands

import (
	"context"
	"fmt"
	"time"

	common "github.com/quangbach27/golang-common"

	"backend/internal/trainers/domain"
)

type ScheduleTrainingCmd struct {
	TrainerUUID string
	Hours       []time.Time
}

func (h *Handler) ScheduleTraining(ctx context.Context, cmd ScheduleTrainingCmd) error {
	return h.hourRepo.UpsertHours(ctx, cmd.TrainerUUID, cmd.Hours, func(hours []*domain.Hour) error {
		errDetails := []common.ErrorDetails{}

		for _, hour := range hours {
			if err := hour.ScheduleTraining(); err != nil {
				errDetails = append(errDetails, common.ErrorDetails{
					EntityType: "Hour",
					EntityID:   hour.Hour().String(),
					ErrorSlug:  fmt.Sprintf("status-%s-not-valid", hour.Status().String()),
					Message:    err.Error(),
				})
			}
		}

		if len(errDetails) != 0 {
			return common.NewConflictError(
				"hours-incorrect-status",
				"cannot schedule training because some of the hours have an incorrect status",
			).WithDetails(errDetails)
		}

		return nil
	})
}
