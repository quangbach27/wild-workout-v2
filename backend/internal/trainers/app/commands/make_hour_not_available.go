package commands

import (
	"context"
	"fmt"
	"time"

	common "github.com/quangbach27/golang-common"

	"backend/internal/trainers/domain"
)

type MakeHourNotAvailableCmd struct {
	TrainerUUID string
	Hours       []time.Time
}

func (h *Handler) MakeHourNotAvailable(ctx context.Context, cmd MakeHourNotAvailableCmd) error {
	return h.hourRepo.UpsertHours(ctx, cmd.TrainerUUID, cmd.Hours, func(hours []*domain.Hour) error {
		errDetails := []common.ErrorDetails{}

		for _, hour := range hours {
			if err := hour.MakeNotAvailable(); err != nil {
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
				"cannot make hours not available because some of them have an incorrect status",
			).WithDetails(errDetails)
		}

		return nil
	})
}
