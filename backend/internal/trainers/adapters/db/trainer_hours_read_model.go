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

// GetTrainerHours returns the hours of query.TrainerUUID from query.DateFrom to query.DateTo (both
// inclusive, in UTC).
//
//   - No status (nil or zero, meaning all): every day with all the hours between MinUtcHour and
//     MaxUtcHour. Hours missing in the database are returned as not available.
//   - NotAvailable: the same grid, keeping only the not available hours. Most of them have no row,
//     so the status can't be filtered in SQL.
//   - Any other status: filtered in SQL, so only the stored hours with that status are returned,
//     grouped by day (days without hours are omitted).
func (rm *TrainerHoursReadModel) GetTrainerHours(
	ctx context.Context,
	query queries.GetTrainerHoursQuery,
) ([]queries.Date, error) {
	firstDay := utcDay(query.DateFrom)
	endExclusive := utcDay(query.DateTo).AddDate(0, 0, 1)

	status := statusFilter(query)
	// not available hours are mostly not stored, so they are filtered after building the grid
	buildGrid := status == nil || *status == domain.NotAvailable
	sqlStatus := status
	if buildGrid {
		sqlStatus = nil
	}

	models, err := dbmodels.New(rm.db).GetHoursInRange(ctx, dbmodels.GetHoursInRangeParams{
		TrainerUuid: query.TrainerUUID,
		DateFrom:    firstDay,
		DateTo:      endExclusive,
		Status:      sqlStatus,
	})
	if err != nil {
		return nil, fmt.Errorf("error retrieving hours in range from database: %w", err)
	}

	if !buildGrid {
		return groupByDay(models, *status), nil
	}

	dates := buildDates(firstDay, endExclusive, rm.hourFactoryConfig, models)
	if status != nil {
		dates = filterByStatus(dates, *status)
	}

	return dates, nil
}

// statusFilter returns the status to filter by, or nil for all the statuses (nil or zero status).
func statusFilter(query queries.GetTrainerHoursQuery) *domain.HourStatus {
	if query.Status == nil || query.Status.IsZero() {
		return nil
	}

	return query.Status
}

// groupByDay groups the stored hours (already filtered by status and ordered by hour) by UTC day.
// HasFreeHours is only known to be true for available hours.
func groupByDay(models []dbmodels.TrainersHour, status domain.HourStatus) []queries.Date {
	dates := []queries.Date{}
	for _, model := range models {
		hour := model.Hour.UTC()
		day := utcDay(hour)

		if len(dates) == 0 || !dates[len(dates)-1].Date.Equal(day) {
			dates = append(dates, queries.Date{
				Date:         day,
				HasFreeHours: status == domain.Available,
			})
		}

		last := &dates[len(dates)-1]
		last.Hours = append(last.Hours, queries.Hour{Hour: hour, Status: model.Status})
	}

	return dates
}

// filterByStatus keeps only the hours with the given status and drops the days left without
// hours. HasFreeHours is left as computed from all the hours of the day.
func filterByStatus(dates []queries.Date, status domain.HourStatus) []queries.Date {
	filtered := []queries.Date{}
	for _, date := range dates {
		hours := make([]queries.Hour, 0, len(date.Hours))
		for _, hour := range date.Hours {
			if hour.Status.Equal(status.Enum) {
				hours = append(hours, hour)
			}
		}
		if len(hours) == 0 {
			continue
		}

		date.Hours = hours
		filtered = append(filtered, date)
	}

	return filtered
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
