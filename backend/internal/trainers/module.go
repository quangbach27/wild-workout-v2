package trainers

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
	"backend/internal/trainers/adapters/db"
	"backend/internal/trainers/app/commands"
	"backend/internal/trainers/app/queries"
	"backend/internal/trainers/domain"
	portHttp "backend/internal/trainers/ports/http"
	"backend/internal/trainers/ports/module"
)

type Module struct {
	dbPgx *pgxpool.Pool

	config    *configs.Config
	contracts *contracts.Contracts

	commandsHandler *commands.Handler
	queriesHandler  *queries.Handler
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

	return &Module{
		config:    config,
		dbPgx:     dbPgx,
		contracts: contracts,
	}
}

func (m *Module) Name() string {
	return "trainers"
}

func (m *Module) Init(ctx context.Context) error {
	if err := m.runMigration(ctx); err != nil {
		return fmt.Errorf("error running migration in module %s: %w", m.Name(), err)
	}

	hourFactory, err := domain.NewHourFactory()
	if err != nil {
		return err
	}

	hourRepo := db.NewHourRepository(m.dbPgx, hourFactory)
	trainerHoursReadModel := db.NewTrainerHoursReadModel(m.dbPgx, hourFactory.Config())

	m.commandsHandler = commands.NewHandler(hourRepo)
	m.queriesHandler = queries.NewHandler(trainerHoursReadModel)

	return nil
}

func (m *Module) RegisterHttp(ctx context.Context, publicRouter http.EchoRouter, protectedRouter http.EchoRouter) error {
	portHttp.Register(
		protectedRouter,
		portHttp.NewHandler(m.commandsHandler, m.queriesHandler),
	)

	return nil
}

func (m *Module) RegisterContracts(ctx context.Context, c *contracts.Contracts) error {
	c.Trainers = module.New(m.commandsHandler)

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
