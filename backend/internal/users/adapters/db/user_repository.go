package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	common "github.com/quangbach27/golang-common"

	"backend/internal/shared"
	"backend/internal/users/adapters/db/dbmodels"
	"backend/internal/users/app/models"
)

type userRepo struct {
	db *pgxpool.Pool
}

var _ models.UserRepository = (*userRepo)(nil)

func NewUserRepository(pgxDb *pgxpool.Pool) *userRepo {
	if pgxDb == nil {
		panic(errors.New("pgxDb can't be nil"))
	}

	return &userRepo{db: pgxDb}
}

func (r *userRepo) OnboardUser(ctx context.Context, user *models.User) error {
	err := dbmodels.New(r.db).InsertUser(ctx, dbmodels.InsertUserParams{
		ID:          user.UUID(),
		DisplayName: user.DisplayName(),
		Balance:     int32(user.Balance()),
		Role:        user.Role().String(),
	})
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
		return common.NewConflictError("user-already-onboarded", "user %s is already onboarded", user.UUID())
	}
	if err != nil {
		return fmt.Errorf("error inserting user to database: %w", err)
	}

	return nil
}

func (r *userRepo) GetUserByUUID(ctx context.Context, userID string) (*models.User, error) {
	row, err := dbmodels.New(r.db).GetUserByUUID(ctx, userID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errUserNotFound(userID)
	}
	if err != nil {
		return nil, fmt.Errorf("error retrieving user from database: %w", err)
	}

	role := shared.Role{Enum: common.MustEnumFromString[shared.RoleType](row.Role)}

	return models.UnmarshalUser(row.ID, row.DisplayName, int(row.Balance), role), nil
}

func (r *userRepo) UpdateBalance(ctx context.Context, userID string, amountChange int) (int, error) {
	balance, err := dbmodels.New(r.db).UpdateBalance(ctx, dbmodels.UpdateBalanceParams{
		ID:           userID,
		AmountChange: int32(amountChange),
	})
	if errors.Is(err, pgx.ErrNoRows) {
		// no row updated: the user doesn't exist, or the balance would become negative
		if _, err = r.GetUserByUUID(ctx, userID); err != nil {
			return 0, err
		}

		return 0, common.NewConflictError("insufficient-balance", "balance can't become negative")
	}
	if err != nil {
		return 0, fmt.Errorf("error updating balance in database: %w", err)
	}

	return int(balance), nil
}

func errUserNotFound(userID string) error {
	return common.NewNotFoundError("user-not-found", "user %s not found", userID)
}

// uniqueViolationCode is the Postgres error code for a unique constraint violation.
const uniqueViolationCode = "23505"
