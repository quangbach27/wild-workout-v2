package modules

import (
	"context"

	"backend/internal/modules/contracts"

	"github.com/quangbach27/golang-common/http"
)

type Module interface {
	Name() string
	Init(ctx context.Context) error
	RegisterHttp(ctx context.Context, publicRouter http.EchoRouter, protectedRouter http.EchoRouter) error
	RegisterContracts(ctx context.Context, contracts *contracts.Contracts) error
}
