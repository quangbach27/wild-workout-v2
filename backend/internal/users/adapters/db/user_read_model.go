package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"backend/internal/shared"
	"backend/internal/users/adapters/db/dbmodels"
	portHttp "backend/internal/users/ports/http"
	"backend/internal/users/ports/module"
)

type UserReadModel struct {
	db *pgxpool.Pool
}

var (
	_ module.UserReadModel   = (*UserReadModel)(nil)
	_ portHttp.UserReadModel = (*UserReadModel)(nil)
)

func NewUserReadModel(db *pgxpool.Pool) *UserReadModel {
	if db == nil {
		panic(errors.New("db can't be nil"))
	}

	return &UserReadModel{db: db}
}

func (r *UserReadModel) GetUser(ctx context.Context, userID string) (*shared.User, error) {
	row, err := dbmodels.New(r.db).GetUserByUUID(ctx, userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errUserNotFound(userID)
	}
	if err != nil {
		return nil, fmt.Errorf("error retrieving user from database: %w", err)
	}

	return shared.UnmarshalUser(
		userID,
		row.DisplayName,
		int(row.Balance),
		row.Role,
	), nil
}
