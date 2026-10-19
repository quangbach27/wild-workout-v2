package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"backend/internal/users/adapters/db/dbmodels"
	"backend/internal/users/ports/module"
	"backend/internal/users/ports/module/client"
)

type UserReadModel struct {
	db *pgxpool.Pool
}

var _ module.UserReadModel = (*UserReadModel)(nil)

func NewUserReadModel(db *pgxpool.Pool) *UserReadModel {
	if db == nil {
		panic(errors.New("db can't be nil"))
	}

	return &UserReadModel{db: db}
}

func (r *UserReadModel) GetUser(ctx context.Context, userUUID string) (client.GetUserResponse, error) {
	row, err := dbmodels.New(r.db).GetUserByUUID(ctx, userUUID)
	if errors.Is(err, pgx.ErrNoRows) {
		return client.GetUserResponse{}, errUserNotFound(userUUID)
	}
	if err != nil {
		return client.GetUserResponse{}, fmt.Errorf("error retrieving user from database: %w", err)
	}

	return client.GetUserResponse{
		UUID:        row.ID,
		DisplayName: row.DisplayName,
		Role:        row.Role.String(),
		Balance:     int(row.Balance),
	}, nil
}
