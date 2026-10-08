package queries

import (
	"context"
	"time"

	common "github.com/quangbach27/golang-common"

	"backend/internal/trainers/domain"
)

type TrainerHoursReadModel interface {
	GetTrainerHours(ctx context.Context, query GetTrainerHoursQuery) ([]Date, error)
}

type GetTrainerHoursQuery struct {
	TrainerUUID string
	DateFrom    time.Time
	DateTo      time.Time
	// Status keeps only the hours with this status (and drops the days left without hours).
	// Nil or zero keeps every hour.
	Status *domain.HourStatus
}

func (q GetTrainerHoursQuery) Validate() error {
	errDetails := []common.ErrorDetails{}

	if q.TrainerUUID == "" {
		errDetails = append(errDetails, common.ErrorDetails{
			EntityType: "GetTrainerHoursQuery",
			ErrorSlug:  "invalid-trainer-uuid",
			Message:    "trainerUUID can't be empty",
		})
	}

	if q.DateFrom.IsZero() {
		errDetails = append(errDetails, common.ErrorDetails{
			EntityType: "GetTrainerHoursQuery",
			ErrorSlug:  "invalid-date-from",
			Message:    "dateFrom can't be empty",
		})
	}

	if q.DateTo.IsZero() {
		errDetails = append(errDetails, common.ErrorDetails{
			EntityType: "GetTrainerHoursQuery",
			ErrorSlug:  "invalid-date-to",
			Message:    "dateTo can't be empty",
		})
	}

	if !q.DateFrom.IsZero() && !q.DateTo.IsZero() && q.DateTo.Before(q.DateFrom) {
		errDetails = append(errDetails, common.ErrorDetails{
			EntityType: "GetTrainerHoursQuery",
			ErrorSlug:  "invalid-date-range",
			Message:    "dateTo can't be before dateFrom",
		})
	}

	if len(errDetails) != 0 {
		return common.NewInvalidInputError(
			"invalid-get-trainer-hours-query",
			"get trainer hours query is not valid",
		).WithDetails(errDetails)
	}

	return nil
}

func (h *Handler) GetTrainerHours(ctx context.Context, query GetTrainerHoursQuery) ([]Date, error) {
	if err := query.Validate(); err != nil {
		return nil, err
	}

	return h.trainerHoursReadModel.GetTrainerHours(ctx, query)
}
