package module

import (
	"context"
	"errors"

	"backend/internal/users/app"
	"backend/internal/users/ports/module/client"
)

type Users struct {
	service *app.Service
}

var _ client.Users = (*Users)(nil)

func New(service *app.Service) *Users {
	if service == nil {
		panic(errors.New("service can't be nil"))
	}

	return &Users{service: service}
}

func (u *Users) GetBalance(ctx context.Context, req client.GetBalanceRequest) (client.GetBalanceResponse, error) {
	balance, err := u.service.GetBalance(ctx, req.UserUUID)
	if err != nil {
		return client.GetBalanceResponse{}, err
	}

	return client.GetBalanceResponse{Balance: balance}, nil
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
