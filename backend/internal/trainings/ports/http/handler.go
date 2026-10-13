package http

import (
	"context"
	"errors"

	common "github.com/quangbach27/golang-common"
	commonHttp "github.com/quangbach27/golang-common/http"

	"backend/internal/configs/auth"
	"backend/internal/shared"
	"backend/internal/trainings/app/commands"
	"backend/internal/trainings/app/queries"
	"backend/internal/trainings/domain"
)

type Handler struct {
	commands *commands.Handler
	queries  *queries.Handler
	policy   *auth.Policy
}

var _ StrictServerInterface = (*Handler)(nil)

func NewHandler(
	commands *commands.Handler,
	queries *queries.Handler,
	policy *auth.Policy,
) *Handler {
	var errs []error

	if commands == nil {
		errs = append(errs, errors.New("commands can't be nil"))
	}
	if queries == nil {
		errs = append(errs, errors.New("queries can't be nil"))
	}

	if policy == nil {
		errs = append(errs, errors.New("policy can't be nil"))
	}

	if len(errs) != 0 {
		panic(errors.Join(errs...))
	}

	return &Handler{
		commands: commands,
		queries:  queries,
		policy:   policy,
	}
}

// Register registers the trainings routes. The router should be the protected one, because the
// handlers read the user from the auth session.
func Register(router commonHttp.EchoRouter, h *Handler) {
	RegisterHandlers(router, NewStrictHandler(h, nil))
}

func (h *Handler) ScheduleTraining(
	ctx context.Context,
	request ScheduleTrainingRequestObject,
) (ScheduleTrainingResponseObject, error) {
	session, err := auth.Authorize(ctx, h.policy, auth.TrainingSchedule)
	if err != nil {
		return nil, err
	}

	trainingUUID, err := h.commands.ScheduleTraining(ctx, commands.ScheduleTrainingCmd{
		AttendeeUUID: session.Subject,
		TrainerUUID:  request.Body.TrainerUuid,
		Hour:         request.Body.Hour,
		Notes:        request.Body.Notes,
	})
	if err != nil {
		return nil, err
	}

	return ScheduleTraining201JSONResponse{Uuid: trainingUUID.String()}, nil
}

func (h *Handler) CancelTraining(
	ctx context.Context,
	request CancelTrainingRequestObject,
) (CancelTrainingResponseObject, error) {
	session, err := auth.Authorize(ctx, h.policy, auth.TrainingCancel)
	if err != nil {
		return nil, err
	}

	role := shared.RoleAttendee
	if session.Role == shared.RoleTrainer.String() {
		role = shared.RoleTrainer
	}

	var trainingUUID domain.TrainingUUID
	if err := trainingUUID.UnmarshalText([]byte(request.TrainingUuid)); err != nil {
		return nil, common.NewInvalidInputError("invalid-training-uuid", "training uuid is not valid")
	}

	if err := h.commands.CancelTraining(ctx, commands.CancelTrainingCmd{
		TrainingUUID: trainingUUID,
		UserUUID:     session.Subject,
		Role:         role,
	}); err != nil {
		return nil, err
	}

	return CancelTraining204Response{}, nil
}

func (h *Handler) GetUserTrainings(
	ctx context.Context,
	request GetUserTrainingsRequestObject,
) (GetUserTrainingsResponseObject, error) {
	session, err := auth.Authorize(ctx, h.policy, auth.TrainingView)
	if err != nil {
		return nil, err
	}

	query := queries.GetUserTrainingsQuery{
		UserUUID: session.Subject,
		Page:     queries.DefaultPage,
		PageSize: queries.DefaultPageSize,
	}
	if request.Params.Page != nil {
		query.Page = *request.Params.Page
	}
	if request.Params.PageSize != nil {
		query.PageSize = *request.Params.PageSize
	}

	page, err := h.queries.GetUserTrainings(ctx, query)
	if err != nil {
		return nil, err
	}

	return GetUserTrainings200JSONResponse{
		Items: trainingsToResponse(page.Trainings),
		Pagination: Pagination{
			Page:     page.Page,
			PageSize: page.PageSize,
			Total:    page.Total,
		},
	}, nil
}

func trainingsToResponse(trainings []queries.Training) []Training {
	response := make([]Training, 0, len(trainings))
	for _, training := range trainings {
		response = append(response, Training{
			Uuid:               training.UUID,
			Hour:               training.Hour,
			Notes:              training.Notes,
			AttendeeUuid:       training.AttendeeUUID,
			AttendeeUsername:   training.AttendeeUsername,
			TrainerUuid:        training.TrainerUUID,
			TrainerUsername:    training.TrainerUsername,
			ProposedNewTime:    training.ProposedNewTime,
			ProposedBy:         training.ProposedBy,
			ProposedByUsername: training.ProposedByUsername,
			Canceled:           training.Canceled,
		})
	}

	return response
}
