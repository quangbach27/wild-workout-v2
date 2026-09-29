package internal

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"backend/internal/configs"
	"backend/internal/modules"
	"backend/internal/modules/contracts"
	"backend/internal/trainers"
	"backend/internal/trainings"
	"backend/internal/users"

	"github.com/jackc/pgx/v5/pgxpool"
	commonHttp "github.com/quangbach27/golang-common/http"
	commonAuth "github.com/quangbach27/golang-common/http/auth"
	"golang.org/x/sync/errgroup"
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
		users.NewModule(),
		trainers.NewModule(),
		trainings.NewModule(),
	}

	if err = initAndRegisterModuleContracts(ctx, modules, moduleContracts); err != nil {
		return Svc{}, err
	}

	for _, module := range modules {
		if err := module.RegisterHttp(
			ctx,
			echoServer.GlobalRouter(),
			echoServer.ProtectedRouter(),
		); err != nil {
			return Svc{}, fmt.Errorf("error register http in %s module: %w", module.Name(), err)
		}
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

func initAndRegisterModuleContracts(
	ctx context.Context,
	modules []modules.Module,
	moduleContracts *contracts.Contracts,
) error {
	for _, module := range modules {
		start := time.Now()

		if err := module.Init(ctx); err != nil {
			return fmt.Errorf("error init %s module: %w", module.Name(), err)
		}

		if err := module.RegisterContracts(ctx, moduleContracts); err != nil {
			return fmt.Errorf("error register contract for %s module: %w", module.Name(), err)
		}

		slog.With(
			"duration", time.Since(start),
			"module", module.Name(),
		).Debug("Initialized module")
	}

	return moduleContracts.Verify()
}
