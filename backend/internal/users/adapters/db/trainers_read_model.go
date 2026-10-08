package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"backend/internal/users/adapters/db/dbmodels"
	portHttp "backend/internal/users/ports/http"
)

type TrainersReadModel struct {
	db *pgxpool.Pool
}

var _ portHttp.TrainersReadModel = (*TrainersReadModel)(nil)

func NewTrainersReadModel(db *pgxpool.Pool) *TrainersReadModel {
	if db == nil {
		panic(errors.New("db can't be nil"))
	}

	return &TrainersReadModel{db: db}
}

func (r *TrainersReadModel) GetTrainers(ctx context.Context, page, pageSize int) (portHttp.TrainersPage, error) {
	queriesDb := dbmodels.New(r.db)

	rows, err := queriesDb.GetTrainers(ctx, dbmodels.GetTrainersParams{
		PageLimit:  int64(pageSize),
		PageOffset: int64(page-1) * int64(pageSize),
	})
	if err != nil {
		return portHttp.TrainersPage{}, fmt.Errorf("error retrieving trainers from database: %w", err)
	}

	total, err := queriesDb.CountTrainers(ctx)
	if err != nil {
		return portHttp.TrainersPage{}, fmt.Errorf("error counting trainers in database: %w", err)
	}

	items := make([]portHttp.Trainer, 0, len(rows))
	for _, row := range rows {
		items = append(items, portHttp.Trainer{
			Uuid:        row.ID,
			DisplayName: row.DisplayName,
			Balance:     int(row.Balance),
		})
	}

	return portHttp.TrainersPage{
		Items:      items,
		Pagination: portHttp.Pagination{Page: page, PageSize: pageSize, Total: int(total)},
	}, nil
}
