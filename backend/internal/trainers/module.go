package trainers

import (
	"context"
	"embed"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	commonDb "github.com/quangbach27/golang-common/db"
	"github.com/quangbach27/golang-common/http"

	"backend/internal/modules/contracts"
	"backend/internal/trainers/ports/module"
)

type Module struct {
	dbPgx *pgxpool.Pool
}

func NewModule() *Module {
	return &Module{}
}

func (m *Module) Name() string {
	return "trainers"
}

func (m *Module) Init(ctx context.Context) error {
	if err := m.runMigration(ctx); err != nil {
		return fmt.Errorf("error running migration in moduel: %s", m.Name())
	}

	return nil
}

func (m *Module) RegisterHttp(ctx context.Context, publicRouter http.EchoRouter, protectedRouter http.EchoRouter) error {
	return nil
}

func (m *Module) RegisterContracts(ctx context.Context, c *contracts.Contracts) error {
	c.Trainers = module.New()

	return nil
}

//go:embed adapters/db/migrations/*.sql
var embedMigrations embed.FS

func (m *Module) runMigration(ctx context.Context) error {
	return commonDb.MigrateDatabaseUp(
		ctx,
		m.Name(),
		m.dbPgx,
		embedMigrations,
		"adapters/db/migrations",
	)
}
