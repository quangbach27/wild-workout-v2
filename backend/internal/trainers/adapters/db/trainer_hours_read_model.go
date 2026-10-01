package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"backend/internal/trainers/adapters/db/dbmodels"
	"backend/internal/trainers/app/queries"
	"backend/internal/trainers/domain"
)

type TrainerHoursReadModel struct {
	db                *pgxpool.Pool
	hourFactoryConfig *domain.HourFactoryConfig
}

func NewTrainerHoursReadModel(
	db *pgxpool.Pool,
	hourFactoryConfig *domain.HourFactoryConfig,
) *TrainerHoursReadModel {
	errs := []error{}
	if db == nil {
		errs = append(errs, errors.New("db can't be nil"))
	}
	if hourFactoryConfig == nil {
		errs = append(errs, errors.New("hourFactoryConfig can't be nil"))
	}
	if len(errs) != 0 {
		panic(errors.Join(errs...))
	}

	return &TrainerHoursReadModel{
		db:                db,
		hourFactoryConfig: hourFactoryConfig,
	}
}

var _ queries.TrainerHoursReadModel = (*TrainerHoursReadModel)(nil)

// GetTrainerHours returns every day from query.DateFrom to query.DateTo (both inclusive, in UTC),
// each with all the hours between MinUtcHour and MaxUtcHour. Hours missing in the database are
// returned as not available.
func (rm *TrainerHoursReadModel) GetTrainerHours(
	ctx context.Context,
	query queries.GetTrainerHoursQuery,
) ([]queries.Date, error) {
	firstDay := utcDay(query.DateFrom)
	endExclusive := utcDay(query.DateTo).AddDate(0, 0, 1)

	models, err := dbmodels.New(rm.db).GetHoursInRange(ctx, dbmodels.GetHoursInRangeParams{
		TrainerUuid: query.TrainerUUID,
		DateFrom:    firstDay,
		DateTo:      endExclusive,
	})
	if err != nil {
		return nil, fmt.Errorf("error retrieving hours in range from database: %w", err)
	}

	return buildDates(firstDay, endExclusive, rm.hourFactoryConfig, models), nil
}

func utcDay(t time.Time) time.Time {
	y, m, d := t.UTC().Date()
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}

// buildDates returns one Date per day in [firstDay, endExclusive), each holding every hour from
// cfg.MinUtcHour to cfg.MaxUtcHour (inclusive). Hours without a row in models are not available.
func buildDates(
	firstDay time.Time,
	endExclusive time.Time,
	cfg *domain.HourFactoryConfig,
	models []dbmodels.TrainersHour,
) []queries.Date {
	// keyed by Unix time: a time.Time key also compares the time zone, and the database returns
	// times in the local zone
	statusByHour := make(map[int64]domain.HourStatus, len(models))
	for _, model := range models {
		statusByHour[model.Hour.Unix()] = model.Status
	}

	dates := []queries.Date{}
	for day := firstDay; day.Before(endExclusive); day = day.AddDate(0, 0, 1) {
		date := queries.Date{
			Date:  day,
			Hours: make([]queries.Hour, 0, cfg.MaxUtcHour-cfg.MinUtcHour+1),
		}

		for h := cfg.MinUtcHour; h <= cfg.MaxUtcHour; h++ {
			hour := time.Date(day.Year(), day.Month(), day.Day(), h, 0, 0, 0, time.UTC)

			status, exist := statusByHour[hour.Unix()]
			if !exist {
				status = domain.NotAvailable
			}
			if status == domain.Available {
				date.HasFreeHours = true
			}

			date.Hours = append(date.Hours, queries.Hour{Hour: hour, Status: status})
		}

		dates = append(dates, date)
	}

	return dates
}
