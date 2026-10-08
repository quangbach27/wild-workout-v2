package internal

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	commonHttp "github.com/quangbach27/golang-common/http"
	commonAuth "github.com/quangbach27/golang-common/http/auth"
	"golang.org/x/sync/errgroup"

	"backend/internal/configs"
	"backend/internal/modules"
	"backend/internal/modules/contracts"
	"backend/internal/trainers"
	"backend/internal/trainings"
	"backend/internal/users"
)

type ExternalServices struct {
	TokenVerifier commonAuth.TokenVerifier
}

type Svc struct {
	echoServer *commonHttp.EchoServer

	config *configs.Config
	pgxDb  *pgxpool.Pool
}

func New(
	ctx context.Context,
	config *configs.Config,
	pgxDb *pgxpool.Pool,
	externalServices ExternalServices,
) (Svc, error) {
	echoServer, err := commonHttp.NewEchoServerWithTokenVerifier(
		externalServices.TokenVerifier,
		func(c *commonHttp.EchoServerConfig) {
			c.Port = config.App.Port
			c.AllowedOrigins = config.App.CorsAllowedOrigins
		},
	)
	if err != nil {
		return Svc{}, fmt.Errorf("error creating echo server: %w", err)
	}

	moduleContracts := &contracts.Contracts{}
	modules := []modules.Module{
		users.NewModule(config, pgxDb, moduleContracts),
		trainers.NewModule(config, pgxDb, moduleContracts),
		trainings.NewModule(config, pgxDb, moduleContracts),
	}

	if err = initModules(ctx, modules); err != nil {
		return Svc{}, err
	}

	if err = registerModuleContracts(ctx, modules, moduleContracts); err != nil {
		return Svc{}, err
	}

	if err = registerModuleHttp(ctx, modules, echoServer); err != nil {
		return Svc{}, err
	}

	return Svc{
		echoServer: echoServer,
		config:     config,
		pgxDb:      pgxDb,
	}, nil
}

func (s Svc) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		if err := s.echoServer.Start(ctx); err != nil {
			return err
		}

		return nil
	})

	g.Go(func() error {
		<-ctx.Done()
		defer s.pgxDb.Close()

		return nil
	})

	return g.Wait()
}

func initModules(ctx context.Context, modules []modules.Module) error {
	for _, module := range modules {
		if err := module.Init(ctx); err != nil {
			return fmt.Errorf("error init %s module: %w", module.Name(), err)
		}
	}

	return nil
}

func registerModuleContracts(
	ctx context.Context,
	modules []modules.Module,
	moduleContracts *contracts.Contracts,
) error {
	for _, module := range modules {
		if err := module.RegisterContracts(ctx, moduleContracts); err != nil {
			return fmt.Errorf("error register contract for %s module: %w", module.Name(), err)
		}
	}

	return moduleContracts.Verify()
}

func registerModuleHttp(
	ctx context.Context,
	modules []modules.Module,
	echoServer *commonHttp.EchoServer,
) error {
	for _, module := range modules {
		if err := module.RegisterHttp(
			ctx,
			echoServer.GlobalRouter(),
			echoServer.ProtectedRouter(),
		); err != nil {
			return fmt.Errorf("error register http in %s module: %w", module.Name(), err)
		}
	}

	return nil
}
