package http

import (
	"context"
	"errors"

	common "github.com/quangbach27/golang-common"
	commonHttp "github.com/quangbach27/golang-common/http"
	"github.com/quangbach27/golang-common/http/auth"

	"backend/internal/trainers/app/commands"
	"backend/internal/trainers/app/queries"
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

// Register registers the trainer hours routes. The router should be the protected one, because
// the handlers read the trainer from the auth session.
func Register(router commonHttp.EchoRouter, h *Handler) {
	RegisterHandlers(router, NewStrictHandler(h, nil))
}

func (h *Handler) GetTrainerHours(
	ctx context.Context,
	request GetTrainerHoursRequestObject,
) (GetTrainerHoursResponseObject, error) {
	trainerUUID, err := trainerUUIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	dates, err := h.queries.GetTrainerHours(ctx, queries.GetTrainerHoursQuery{
		TrainerUUID: trainerUUID,
		DateFrom:    request.Params.From,
		DateTo:      request.Params.To,
	})
	if err != nil {
		return nil, err
	}

	return GetTrainerHours200JSONResponse(datesToResponse(dates)), nil
}

func (h *Handler) MakeHoursAvailable(
	ctx context.Context,
	request MakeHoursAvailableRequestObject,
) (MakeHoursAvailableResponseObject, error) {
	trainerUUID, err := trainerUUIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := h.commands.MakeHourAvailable(ctx, commands.MakeHourAvailableCmd{
		TrainerUUID: trainerUUID,
		Hours:       request.Body.Hours,
	}); err != nil {
		return nil, err
	}

	return MakeHoursAvailable204Response{}, nil
}

func (h *Handler) MakeHoursNotAvailable(
	ctx context.Context,
	request MakeHoursNotAvailableRequestObject,
) (MakeHoursNotAvailableResponseObject, error) {
	trainerUUID, err := trainerUUIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if err := h.commands.MakeHourNotAvailable(ctx, commands.MakeHourNotAvailableCmd{
		TrainerUUID: trainerUUID,
		Hours:       request.Body.Hours,
	}); err != nil {
		return nil, err
	}

	return MakeHoursNotAvailable204Response{}, nil
}

func trainerUUIDFromContext(ctx context.Context) (string, error) {
	// TODO: check roles for this session. Only trainer can access this
	session, ok := auth.SessionFromContext(ctx)
	if !ok || session.UserID == "" {
		return "", common.NewUnauthorizedError("missing-session", "missing session")
	}

	return session.UserID, nil
}

func datesToResponse(dates []queries.Date) []Date {
	response := make([]Date, 0, len(dates))
	for _, date := range dates {
		hours := make([]Hour, 0, len(date.Hours))
		for _, hour := range date.Hours {
			hours = append(hours, Hour{Hour: hour.Hour, Status: hour.Status})
		}

		response = append(response, Date{
			Date:         date.Date,
			HasFreeHours: date.HasFreeHours,
			Hours:        hours,
		})
	}

	return response
}
