package http

import (
	"context"
	"errors"

	common "github.com/quangbach27/golang-common"
	commonHttp "github.com/quangbach27/golang-common/http"

	"backend/internal/configs/auth"
	"backend/internal/trainers/app/commands"
	"backend/internal/trainers/app/queries"
	"backend/internal/trainers/domain"
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

// Register registers the trainer hours routes. The router should be the protected one, because
// the handlers read the trainer from the auth session.
func Register(router commonHttp.EchoRouter, h *Handler) {
	RegisterHandlers(router, NewStrictHandler(h, nil))
}

func (h *Handler) GetTrainerHours(
	ctx context.Context,
	request GetTrainerHoursRequestObject,
) (GetTrainerHoursResponseObject, error) {
	session, err := auth.Authorize(ctx, h.policy, auth.HourView)
	if err != nil {
		return nil, err
	}

	dates, err := h.queries.GetTrainerHours(ctx, queries.GetTrainerHoursQuery{
		TrainerUUID: session.Subject,
		DateFrom:    request.Params.From,
		DateTo:      request.Params.To,
	})
	if err != nil {
		return nil, err
	}

	return GetTrainerHours200JSONResponse(datesToResponse(dates)), nil
}

func (h *Handler) GetTrainerHoursByUuid(
	ctx context.Context,
	request GetTrainerHoursByUuidRequestObject,
) (GetTrainerHoursByUuidResponseObject, error) {
	if _, err := auth.Authorize(ctx, h.policy, auth.HourView); err != nil {
		return nil, err
	}

	status, err := toDomainHourStatus(request.Params.Status)
	if err != nil {
		return nil, err
	}

	dates, err := h.queries.GetTrainerHours(ctx, queries.GetTrainerHoursQuery{
		TrainerUUID: request.TrainerUuid,
		DateFrom:    request.Params.From,
		DateTo:      request.Params.To,
		Status:      status,
	})
	if err != nil {
		return nil, err
	}

	return GetTrainerHoursByUuid200JSONResponse(datesToResponse(dates)), nil
}

func (h *Handler) MakeHoursAvailable(
	ctx context.Context,
	request MakeHoursAvailableRequestObject,
) (MakeHoursAvailableResponseObject, error) {
	session, err := auth.Authorize(ctx, h.policy, auth.HourMakeAvailable)
	if err != nil {
		return nil, err
	}

	if err := h.commands.MakeHourAvailable(ctx, commands.MakeHourAvailableCmd{
		TrainerUUID: session.Subject,
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
	session, err := auth.Authorize(ctx, h.policy, auth.HourMakeUnavailable)
	if err != nil {
		return nil, err
	}

	if err := h.commands.MakeHourNotAvailable(ctx, commands.MakeHourNotAvailableCmd{
		TrainerUUID: session.Subject,
		Hours:       request.Body.Hours,
	}); err != nil {
		return nil, err
	}

	return MakeHoursNotAvailable204Response{}, nil
}

func toDomainHourStatus(status *GetTrainerHoursByUuidParamsStatus) (*domain.HourStatus, error) {
	if status == nil || *status == "" {
		return nil, nil
	}

	// MustEnum panics on an unknown value, so reject it first
	if !status.Valid() {
		return nil, common.NewInvalidInputError(
			"invalid-hour-status",
			"unknown hour status %q",
			*status,
		)
	}

	hourStatus := common.MustEnum[domain.HourStatus](string(*status))

	return &hourStatus, nil
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
