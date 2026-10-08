package http

import (
	"context"
	"errors"
	"slices"

	common "github.com/quangbach27/golang-common"
	commonHttp "github.com/quangbach27/golang-common/http"
	"github.com/quangbach27/golang-common/http/auth"

	"backend/internal/shared"
	"backend/internal/users/app"
)

const (
	defaultPage     = 1
	defaultPageSize = 10
	maxPageSize     = 50
	maxTopUp        = 1_000_000
)

type TrainersReadModel interface {
	// GetTrainers returns the page of trainers ordered by (display name, id), with the total
	// across all pages.
	GetTrainers(ctx context.Context, page, pageSize int) (TrainersPage, error)
}

type Handler struct {
	service           *app.Service
	trainersReadModel TrainersReadModel
}

var _ StrictServerInterface = (*Handler)(nil)

func NewHandler(service *app.Service, trainersReadModel TrainersReadModel) *Handler {
	var errs []error

	if service == nil {
		errs = append(errs, errors.New("service can't be nil"))
	}
	if trainersReadModel == nil {
		errs = append(errs, errors.New("trainersReadModel can't be nil"))
	}

	if len(errs) != 0 {
		panic(errors.Join(errs...))
	}

	return &Handler{service: service, trainersReadModel: trainersReadModel}
}

// OnboardUser creates the session's user. The uuid and the role come from the session, never from
// the request.
func (h *Handler) OnboardUser(
	ctx context.Context,
	request OnboardUserRequestObject,
) (OnboardUserResponseObject, error) {
	session, ok := auth.SessionFromContext(ctx)
	if !ok || session.UserID == "" {
		return nil, common.NewUnauthorizedError("missing-session", "missing session")
	}

	var role shared.Role
	switch {
	case slices.Contains(session.Roles, shared.RoleTrainer.String()):
		role = shared.RoleTrainer
	case slices.Contains(session.Roles, shared.RoleAttendee.String()):
		role = shared.RoleAttendee
	}

	err := h.service.OnboardUser(ctx, app.OnboardUserCmd{
		UUID:        session.UserID,
		DisplayName: request.Body.DisplayName,
		Role:        role,
	})
	if err != nil {
		return nil, err
	}

	return OnboardUser204Response{}, nil
}

// Register registers the users routes on the protected router.
func Register(router commonHttp.EchoRouter, h *Handler) {
	RegisterHandlers(router, NewStrictHandler(h, nil))
}

// TopUpBalance adds the amount to the balance of the session user. There is no payment step yet.
func (h *Handler) TopUpBalance(
	ctx context.Context,
	request TopUpBalanceRequestObject,
) (TopUpBalanceResponseObject, error) {
	session, ok := auth.SessionFromContext(ctx)
	if !ok || session.UserID == "" {
		return nil, common.NewUnauthorizedError("missing-session", "missing session")
	}

	if request.Body.Amount < 1 || request.Body.Amount > maxTopUp {
		return nil, common.NewInvalidInputError(
			"invalid-top-up-balance-request",
			"top up balance request is not valid",
		).WithDetails([]common.ErrorDetails{{
			EntityType: "TopUpBalanceRequest",
			ErrorSlug:  "invalid-amount",
			Message:    "amount must be between 1 and 1000000",
		}})
	}

	balance, err := h.service.UpdateBalance(ctx, app.UpdateBalanceCmd{
		UserUUID:     session.UserID,
		AmountChange: request.Body.Amount,
	})
	if err != nil {
		return nil, err
	}

	return TopUpBalance200JSONResponse{Balance: balance}, nil
}

func (h *Handler) GetTrainers(
	ctx context.Context,
	request GetTrainersRequestObject,
) (GetTrainersResponseObject, error) {
	if session, ok := auth.SessionFromContext(ctx); !ok || session.UserID == "" {
		return nil, common.NewUnauthorizedError("missing-session", "missing session")
	}

	page, pageSize := defaultPage, defaultPageSize
	if request.Params.Page != nil {
		page = *request.Params.Page
	}
	if request.Params.PageSize != nil {
		pageSize = *request.Params.PageSize
	}

	if err := validatePagination(page, pageSize); err != nil {
		return nil, err
	}

	trainers, err := h.trainersReadModel.GetTrainers(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	return GetTrainers200JSONResponse(trainers), nil
}

func validatePagination(page, pageSize int) error {
	errDetails := []common.ErrorDetails{}

	if page < 1 {
		errDetails = append(errDetails, common.ErrorDetails{
			EntityType: "GetTrainersQuery",
			ErrorSlug:  "invalid-page",
			Message:    "page must be at least 1",
		})
	}

	if pageSize < 1 || pageSize > maxPageSize {
		errDetails = append(errDetails, common.ErrorDetails{
			EntityType: "GetTrainersQuery",
			ErrorSlug:  "invalid-page-size",
			Message:    "pageSize must be between 1 and 50",
		})
	}

	if len(errDetails) != 0 {
		return common.NewInvalidInputError(
			"invalid-get-trainers-query",
			"get trainers query is not valid",
		).WithDetails(errDetails)
	}

	return nil
}
