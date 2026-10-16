package users

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log/slog"

	"backend/internal/configs"
	"backend/internal/configs/modules/contracts"
	"backend/internal/users/adapters/db"
	"backend/internal/users/app"
	"backend/internal/users/ports/module"

	"github.com/jackc/pgx/v5/pgxpool"
	commonDb "github.com/quangbach27/golang-common/db"
	"github.com/quangbach27/golang-common/http"
	svix "github.com/svix/svix-webhooks/go"

	portHttp "backend/internal/users/ports/http"
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

	if err := m.runSeed(ctx); err != nil {
		return fmt.Errorf("error seeding users in module %s: %w", m.Name(), err)
	}

	m.service = app.NewService(db.NewUserRepository(m.dbPgx))
	m.trainersReadModel = db.NewTrainersReadModel(m.dbPgx)

	return nil
}

func (m *Module) RegisterHttp(ctx context.Context, router http.EchoRouter) error {
	portHttp.Register(router, portHttp.NewHandler(m.service, m.trainersReadModel))
	if err := m.registerwebhook(router); err != nil {
		return err
	}
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

func (m *Module) registerwebhook(router http.EchoRouter) error {
	if m.config.Clerk.WebhookSigningSecret == "" {
		slog.Warn("CLERK_WEBHOOK_SIGNING_SECRET is not set, Clerk webhook route is not registered")

		return nil
	}

	webhookVerifier, err := svix.NewWebhook(m.config.Clerk.WebhookSigningSecret)
	if err != nil {
		return fmt.Errorf("failed to init webhook verifier: %w", err)
	}

	webhookHandler, err := portHttp.NewClerkWebhook(m.service, webhookVerifier)
	if err != nil {
		return fmt.Errorf("failed to create webhook handler: %w", err)
	}

	if err := portHttp.RegisterWebhook(router, webhookHandler); err != nil {
		return fmt.Errorf("failed to register webhook: %w", err)
	}

	return nil
}

// runSeed inserts the initial users. Existing rows are left untouched.
// Temporary: delete it (and its call in Init) after the app is deployed.
func (m *Module) runSeed(ctx context.Context) error {
	_, err := m.dbPgx.Exec(ctx, `
		INSERT INTO users.users (id, display_name, balance, role) VALUES
			('user_3KaKzdEqG0PUdKwSy9w4MSOZ6kh', 'attendee27', 20, 'attendee'),
			('user_3KaKiiqacUY8hqxq5IOFgqm5gsI', 'trainer27', 20, 'trainer')
		ON CONFLICT (id) DO NOTHING`)

	return err
}
