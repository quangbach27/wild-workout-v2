package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"backend/internal/trainings/adapters/db/dbmodels"
	"backend/internal/trainings/app/queries"
)

type UserTrainingsReadModel struct {
	db *pgxpool.Pool
}

var _ queries.UserTrainingsReadModel = (*UserTrainingsReadModel)(nil)

func NewUserTrainingsReadModel(db *pgxpool.Pool) *UserTrainingsReadModel {
	if db == nil {
		panic(errors.New("db can't be nil"))
	}

	return &UserTrainingsReadModel{db: db}
}

func (r *UserTrainingsReadModel) GetUpcomingTrainings(
	ctx context.Context,
	query queries.GetUserTrainingsQuery,
) (queries.TrainingsPage, error) {
	queriesDb := dbmodels.New(r.db)

	models, err := queriesDb.GetUpcomingTrainingsByUser(ctx, dbmodels.GetUpcomingTrainingsByUserParams{
		UserID:     query.UserUUID,
		After:      query.After,
		PageLimit:  int64(query.PageSize),
		PageOffset: int64(query.Page-1) * int64(query.PageSize),
	})
	if err != nil {
		return queries.TrainingsPage{}, fmt.Errorf("error retrieving upcoming trainings from database: %w", err)
	}

	total, err := queriesDb.CountUpcomingTrainingsByUser(ctx, dbmodels.CountUpcomingTrainingsByUserParams{
		UserID: query.UserUUID,
		After:  query.After,
	})
	if err != nil {
		return queries.TrainingsPage{}, fmt.Errorf("error counting upcoming trainings in database: %w", err)
	}

	trainings := make([]queries.Training, 0, len(models))
	for _, m := range models {
		training := queries.Training{
			UUID:               m.ID.String(),
			Hour:               m.Hour.UTC(),
			Notes:              m.Notes,
			AttendeeUUID:       m.AttendeeID,
			AttendeeUsername:   m.AttendeeUsername,
			TrainerUUID:        m.TrainerID,
			TrainerUsername:    m.TrainerUsername,
			ProposedBy:         m.MoveProposedByID,
			ProposedByUsername: proposedByUsername(m),
			Canceled:           m.Canceled,
		}
		if m.ProposedNewTime != nil {
			proposed := m.ProposedNewTime.UTC()
			training.ProposedNewTime = &proposed
		}

		trainings = append(trainings, training)
	}

	return queries.TrainingsPage{
		Trainings: trainings,
		Page:      query.Page,
		PageSize:  query.PageSize,
		Total:     int(total),
	}, nil
}

// proposedByUsername resolves the proposer's username, who is always the trainer or the attendee.
func proposedByUsername(m dbmodels.TrainingsTraining) *string {
	if m.MoveProposedByID == nil {
		return nil
	}

	username := m.AttendeeUsername
	if *m.MoveProposedByID == m.TrainerID {
		username = m.TrainerUsername
	}

	return &username
}
