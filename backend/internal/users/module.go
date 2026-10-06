package users

import (
	"context"
	"errors"

	"backend/internal/modules/contracts"
	"backend/internal/users/ports/module"

	"github.com/quangbach27/golang-common/http"
)

type Module struct {
	contracts *contracts.Contracts
}

func NewModule(contracts *contracts.Contracts) *Module {
	if contracts == nil {
		panic(errors.New("contracts can't be nil"))
	}

	return &Module{contracts: contracts}
}

func (m *Module) Name() string {
	return "users"
}

func (m *Module) Init(ctx context.Context) error {
	return nil
}

func (m *Module) RegisterHttp(ctx context.Context, publicRouter http.EchoRouter, protectedRouter http.EchoRouter) error {
	return nil
}

func (m *Module) RegisterContracts(ctx context.Context, c *contracts.Contracts) error {
	c.Users = module.New()

	return nil
}
