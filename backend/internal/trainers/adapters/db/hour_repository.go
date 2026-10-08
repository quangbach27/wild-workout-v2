package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	commonDb "github.com/quangbach27/golang-common/db"

	"backend/internal/trainers/adapters/db/dbmodels"
	"backend/internal/trainers/domain"
)

type hourRepo struct {
	hourFactory *domain.HourFactory
	db          *pgxpool.Pool
}

var _ domain.HourRepository = (*hourRepo)(nil)

func NewHourRepository(
	pgxDb *pgxpool.Pool,
	hourFactory *domain.HourFactory,
) *hourRepo {
	var errs []error

	if pgxDb == nil {
		errs = append(errs, errors.New("pgxDb can't be nil"))
	}
	if hourFactory == nil {
		errs = append(errs, errors.New("hourRepo can't be nil"))
	}
	if len(errs) != 0 {
		panic(errors.Join(errs...))
	}

	return &hourRepo{
		db:          pgxDb,
		hourFactory: hourFactory,
	}
}

func (r *hourRepo) UpsertHours(
	ctx context.Context,
	trainerUUID string,
	hours []time.Time,
	upsertFn func(h []*domain.Hour) error,
) error {
	return commonDb.UpdateInTx(ctx, r.db, func(ctx context.Context, tx pgx.Tx) error {
		queries := dbmodels.New(tx)

		hourDomains, err := r.getHours(ctx, queries, trainerUUID, hours)
		if err != nil {
			return err
		}

		if err := upsertFn(hourDomains); err != nil {
			return err
		}

		return r.upsertHours(
			ctx,
			queries,
			trainerUUID,
			hourDomains,
		)
	})
}

func (r *hourRepo) getHours(
	ctx context.Context,
	queries *dbmodels.Queries,
	trainerUUID string,
	hours []time.Time,
) ([]*domain.Hour, error) {
	models, err := queries.GetHours(ctx, dbmodels.GetHoursParams{
		TrainerUuid: trainerUUID,
		Hours:       hours,
	})
	if err != nil {
		return nil, fmt.Errorf("error retrieving hours from database: %w", err)
	}

	// keyed by Unix time: a time.Time key also compares the time zone, and the database returns
	// times in the local zone while callers pass UTC
	hourModelMap := make(map[int64]dbmodels.TrainersHour, len(models))
	for _, model := range models {
		hourModelMap[model.Hour.Unix()] = model
	}

	hourDomains := make([]*domain.Hour, 0, len(hours))
	for _, h := range hours {
		if hourModel, exist := hourModelMap[h.Unix()]; exist {
			hourDomains = append(
				hourDomains,
				domain.UnmarshalHour(
					hourModel.Hour,
					hourModel.Status,
					hourModel.TrainerUuid,
				),
			)
		} else {
			newHourDomain, err := r.hourFactory.NewNotAvailableHour(trainerUUID, h)
			if err != nil {
				return nil, err
			}

			hourDomains = append(hourDomains, newHourDomain)
		}
	}

	return hourDomains, nil
}

func (r *hourRepo) upsertHours(
	ctx context.Context,
	queries *dbmodels.Queries,
	trainerUUID string,
	hours []*domain.Hour,
) error {
	hoursParams := make([]time.Time, 0, len(hours))
	statusParams := make([]string, 0, len(hours))

	for _, hour := range hours {
		hoursParams = append(hoursParams, hour.Hour())
		statusParams = append(statusParams, hour.Status().String())
	}

	if err := queries.UpsertHours(ctx, dbmodels.UpsertHoursParams{
		TrainerUuid: trainerUUID,
		Hours:       hoursParams,
		Status:      statusParams,
	}); err != nil {
		return fmt.Errorf("error batch upsert hour to database: %w", err)
	}

	return nil
}
