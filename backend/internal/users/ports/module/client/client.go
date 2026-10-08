package client

import "context"

type Users interface {
	// GetBalance returns the user's balance. It returns a 404 common.Error for an unknown user.
	GetBalance(ctx context.Context, req GetBalanceRequest) (GetBalanceResponse, error)
	// UpdateBalance adds AmountChange (negative to spend) to the user's balance. It returns a 404
	// common.Error for an unknown user and a 409 if the balance would become negative.
	UpdateBalance(ctx context.Context, req UpdateBalanceRequest) (UpdateBalanceResponse, error)
}

type GetBalanceRequest struct {
	UserUUID string
}

type GetBalanceResponse struct {
	Balance int
}

type UpdateBalanceRequest struct {
	UserUUID     string
	AmountChange int
}

type UpdateBalanceResponse struct {
	// Balance is the balance after the change.
	Balance int
}
