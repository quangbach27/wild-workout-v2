package users

import (
	"context"

	"backend/internal/modules/contracts"
	"backend/internal/users/ports/module"

	"github.com/quangbach27/golang-common/http"
)

type Module struct{}

func NewModule() *Module {
	return &Module{}
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
