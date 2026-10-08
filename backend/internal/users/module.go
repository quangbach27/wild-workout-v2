package users

import (
	"context"
	"embed"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	commonDb "github.com/quangbach27/golang-common/db"
	"github.com/quangbach27/golang-common/http"

	"backend/internal/configs"
	"backend/internal/modules/contracts"
	"backend/internal/users/adapters/db"
	"backend/internal/users/app"
	portHttp "backend/internal/users/ports/http"
	"backend/internal/users/ports/module"
)

type Module struct {
	dbPgx *pgxpool.Pool

	config    *configs.Config
	contracts *contracts.Contracts

	service           *app.Service
	trainersReadModel portHttp.TrainersReadModel
}

func NewModule(config *configs.Config, dbPgx *pgxpool.Pool, contracts *contracts.Contracts) *Module {
	errs := []error{}
	if config == nil {
		errs = append(errs, errors.New("config can't be nil"))
	}
	if dbPgx == nil {
		errs = append(errs, errors.New("dbPgx can't be nil"))
	}
	if contracts == nil {
		errs = append(errs, errors.New("contracts can't be nil"))
	}
	if len(errs) != 0 {
		panic(errors.Join(errs...))
	}

	return &Module{config: config, dbPgx: dbPgx, contracts: contracts}
}

func (m *Module) Name() string {
	return "users"
}

func (m *Module) Init(ctx context.Context) error {
	if err := m.runMigration(ctx); err != nil {
		return fmt.Errorf("error running migration in module %s: %w", m.Name(), err)
	}

	m.service = app.NewService(db.NewUserRepository(m.dbPgx))
	m.trainersReadModel = db.NewTrainersReadModel(m.dbPgx)

	return nil
}

func (m *Module) RegisterHttp(ctx context.Context, publicRouter http.EchoRouter, protectedRouter http.EchoRouter) error {
	portHttp.Register(protectedRouter, portHttp.NewHandler(m.service, m.trainersReadModel))

	return nil
}

func (m *Module) RegisterContracts(ctx context.Context, c *contracts.Contracts) error {
	c.Users = module.New(m.service)

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
