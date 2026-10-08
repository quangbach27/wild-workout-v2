package http

import (
	"context"
	"errors"
	"slices"

	common "github.com/quangbach27/golang-common"
	commonHttp "github.com/quangbach27/golang-common/http"
	"github.com/quangbach27/golang-common/http/auth"

	"backend/internal/shared"
	"backend/internal/trainings/app/commands"
	"backend/internal/trainings/app/queries"
	"backend/internal/trainings/domain"
)

type Handler struct {
	commands *commands.Handler
	queries  *queries.Handler
}

var _ StrictServerInterface = (*Handler)(nil)

func NewHandler(
	commands *commands.Handler,
	queries *queries.Handler,
) *Handler {
	var errs []error

	if commands == nil {
		errs = append(errs, errors.New("commands can't be nil"))
	}
	if queries == nil {
		errs = append(errs, errors.New("queries can't be nil"))
	}

	if len(errs) != 0 {
		panic(errors.Join(errs...))
	}

	return &Handler{
		commands: commands,
		queries:  queries,
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
	session, err := sessionFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if !slices.Contains(session.Roles, shared.RoleAttendee.String()) {
		return nil, common.NewForbiddenError("forbidden-role", "only attendees can schedule a training")
	}

	username, _ := session.Extra["username"].(string)
	if username == "" {
		return nil, common.NewUnauthorizedError("missing-username", "session has no username")
	}

	trainingUUID, err := h.commands.ScheduleTraining(ctx, commands.ScheduleTrainingCmd{
		AttendeeUUID:     session.UserID,
		AttendeeUsername: username,
		TrainerUUID:      request.Body.TrainerUuid,
		TrainerUsername:  request.Body.TrainerUsername,
		Hour:             request.Body.Hour,
		Notes:            request.Body.Notes,
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
	session, err := sessionFromContext(ctx)
	if err != nil {
		return nil, err
	}

	username, _ := session.Extra["username"].(string)
	if username == "" {
		return nil, common.NewUnauthorizedError("missing-username", "session has no username")
	}

	role := shared.RoleAttendee
	if slices.Contains(session.Roles, shared.RoleTrainer.String()) {
		role = shared.RoleTrainer
	}

	user, err := domain.NewUser(session.UserID, username, role)
	if err != nil {
		return nil, err
	}

	var trainingUUID domain.TrainingUUID
	if err := trainingUUID.UnmarshalText([]byte(request.TrainingUuid)); err != nil {
		return nil, common.NewInvalidInputError("invalid-training-uuid", "training uuid is not valid")
	}

	if err := h.commands.CancelTraining(ctx, commands.CancelTrainingCmd{
		TrainingUUID: trainingUUID,
		User:         user,
	}); err != nil {
		return nil, err
	}

	return CancelTraining204Response{}, nil
}

func (h *Handler) GetUserTrainings(
	ctx context.Context,
	request GetUserTrainingsRequestObject,
) (GetUserTrainingsResponseObject, error) {
	session, err := sessionFromContext(ctx)
	if err != nil {
		return nil, err
	}

	query := queries.GetUserTrainingsQuery{
		UserUUID: session.UserID,
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

func sessionFromContext(ctx context.Context) (*auth.Session, error) {
	session, ok := auth.SessionFromContext(ctx)
	if !ok || session.UserID == "" {
		return nil, common.NewUnauthorizedError("missing-session", "missing session")
	}

	return session, nil
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
