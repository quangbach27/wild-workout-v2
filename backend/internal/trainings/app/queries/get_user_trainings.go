package queries

import (
	"context"
	"time"

	common "github.com/quangbach27/golang-common"
)

const (
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 50
)

type UserTrainingsReadModel interface {
	// GetUpcomingTrainings returns the page the query asks for, ordered by (hour, id).
	GetUpcomingTrainings(ctx context.Context, query GetUserTrainingsQuery) (TrainingsPage, error)
}

type GetUserTrainingsQuery struct {
	// UserUUID matches the trainings where the user is the trainer or the attendee.
	UserUUID string
	// After keeps only the trainings later than this time. Defaults to now.
	After time.Time
	// Page is the 1-based page number.
	Page int
	// PageSize is the number of trainings per page, between 1 and MaxPageSize.
	PageSize int
}

type TrainingsPage struct {
	Trainings []Training
	Page      int
	PageSize  int
	// Total is the number of upcoming trainings across all pages.
	Total int
}

func (q GetUserTrainingsQuery) Validate() error {
	errDetails := []common.ErrorDetails{}

	if q.UserUUID == "" {
		errDetails = append(errDetails, common.ErrorDetails{
			EntityType: "GetUserTrainingsQuery",
			ErrorSlug:  "invalid-user-uuid",
			Message:    "userUUID can't be empty",
		})
	}

	if q.Page < 1 {
		errDetails = append(errDetails, common.ErrorDetails{
			EntityType: "GetUserTrainingsQuery",
			ErrorSlug:  "invalid-page",
			Message:    "page must be at least 1",
		})
	}

	if q.PageSize < 1 || q.PageSize > MaxPageSize {
		errDetails = append(errDetails, common.ErrorDetails{
			EntityType: "GetUserTrainingsQuery",
			ErrorSlug:  "invalid-page-size",
			Message:    "pageSize must be between 1 and 50",
		})
	}

	if len(errDetails) != 0 {
		return common.NewInvalidInputError(
			"invalid-get-user-trainings-query",
			"get user trainings query is not valid",
		).WithDetails(errDetails)
	}

	return nil
}

// GetUserTrainings returns one page of the user's upcoming trainings, canceled ones included,
// soonest first, with the total across all pages.
func (h *Handler) GetUserTrainings(ctx context.Context, query GetUserTrainingsQuery) (TrainingsPage, error) {
	if err := query.Validate(); err != nil {
		return TrainingsPage{}, err
	}

	if query.After.IsZero() {
		query.After = time.Now()
	}

	return h.userTrainingsReadModel.GetUpcomingTrainings(ctx, query)
}
