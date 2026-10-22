package module

import (
	"context"
	"errors"

	common "github.com/quangbach27/golang-common"

	"backend/internal/shared"
	"backend/internal/users/app"
	"backend/internal/users/ports/module/client"
)

type UserReadModel interface {
	// GetUser returns a 404 common.Error for an unknown user.
	GetUser(ctx context.Context, userID string) (*shared.User, error)
}

type Users struct {
	service       *app.Service
	userReadModel UserReadModel
}

var _ client.Users = (*Users)(nil)

func New(service *app.Service, userReadModel UserReadModel) *Users {
	errs := []error{}
	if service == nil {
		errs = append(errs, errors.New("service can't be nil"))
	}
	if userReadModel == nil {
		errs = append(errs, errors.New("userReadModel can't be nil"))
	}
	if len(errs) != 0 {
		panic(errors.Join(errs...))
	}

	return &Users{service: service, userReadModel: userReadModel}
}

func (u *Users) GetUser(ctx context.Context, req client.GetUserRequest) (client.GetUserResponse, error) {
	if req.UserID == "" {
		return client.GetUserResponse{}, common.NewInvalidInputError("invalid-user-uuid", "user id can't be empty")
	}
	userModel, err := u.userReadModel.GetUser(ctx, req.UserID)
	if err != nil {
		return client.GetUserResponse{}, err
	}

	return client.GetUserResponse{
		ID:          userModel.UUID(),
		DisplayName: userModel.DisplayName(),
		Role:        userModel.Role(),
		Balance:     userModel.Balance(),
	}, nil
}

func (u *Users) GetBalance(ctx context.Context, req client.GetBalanceRequest) (client.GetBalanceResponse, error) {
	userInfo, err := u.userReadModel.GetUser(ctx, req.UserUUID)
	if err != nil {
		return client.GetBalanceResponse{}, err
	}

	return client.GetBalanceResponse{Balance: userInfo.Balance()}, nil
}

func (u *Users) UpdateBalance(
	ctx context.Context,
	req client.UpdateBalanceRequest,
) (client.UpdateBalanceResponse, error) {
	balance, err := u.service.UpdateBalance(ctx, app.UpdateBalanceCmd{
		UserUUID:     req.UserUUID,
		AmountChange: req.AmountChange,
	})
	if err != nil {
		return client.UpdateBalanceResponse{}, err
	}

	return client.UpdateBalanceResponse{Balance: balance}, nil
}
