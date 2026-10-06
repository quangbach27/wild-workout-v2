package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	common "github.com/quangbach27/golang-common"
	commonDb "github.com/quangbach27/golang-common/db"

	"backend/internal/shared"
	"backend/internal/trainings/adapters/db/dbmodels"
	"backend/internal/trainings/domain"
)

type trainingRepo struct {
	db *pgxpool.Pool
}

var _ domain.TrainingRepository = (*trainingRepo)(nil)

func NewTrainingRepository(pgxDb *pgxpool.Pool) *trainingRepo {
	if pgxDb == nil {
		panic(errors.New("pgxDb can't be nil"))
	}

	return &trainingRepo{db: pgxDb}
}

func (r *trainingRepo) AddTraining(ctx context.Context, training *domain.Training) error {
	if err := dbmodels.New(r.db).UpsertTraining(ctx, upsertParams(training)); err != nil {
		return fmt.Errorf("error upserting training to database: %w", err)
	}

	return nil
}

func (r *trainingRepo) GetTraining(
	ctx context.Context,
	id domain.TrainingUUID,
	user domain.User,
) (*domain.Training, error) {
	model, err := dbmodels.New(r.db).GetTraining(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error retrieiving training from database: %w", err)
	}

	return toVisibleTraining(model, user)
}

func (r *trainingRepo) UpdateTraining(
	ctx context.Context,
	id domain.TrainingUUID,
	user domain.User,
	updateFn func(ctx context.Context, training *domain.Training) error,
) error {
	return commonDb.UpdateInTx(ctx, r.db, func(ctx context.Context, tx pgx.Tx) error {
		queries := dbmodels.New(tx)

		model, err := queries.GetTraining(ctx, id)
		if err != nil {
			return fmt.Errorf("error retrieiving training from database: %w", err)
		}

		training, err := toVisibleTraining(model, user)
		if err != nil {
			return err
		}

		if err := updateFn(ctx, training); err != nil {
			return err
		}

		if err := queries.UpsertTraining(ctx, upsertParams(training)); err != nil {
			return fmt.Errorf("error upserting training to database: %w", err)
		}

		return nil
	})
}

// toVisibleTraining maps the model to the domain and returns a 403 common.Error when the user
// is neither the trainer nor the attendee.
func toVisibleTraining(model dbmodels.TrainingsTraining, user domain.User) (*domain.Training, error) {
	training := toDomain(model)

	if err := domain.CanUserSeeTraining(user, training); err != nil {
		return nil, common.NewForbiddenError("training-forbidden", "%s", err.Error()).WithInternalError(err)
	}

	return training, nil
}

func toDomain(m dbmodels.TrainingsTraining) *domain.Training {
	attendee := domain.UnmarshalUser(m.AttendeeID, m.AttendeeUsername, shared.RoleAttendee)
	trainer := domain.UnmarshalUser(m.TrainerID, m.TrainerUsername, shared.RoleTrainer)

	// the proposer is always one of the two participants
	var moveProposedBy *domain.User
	if m.MoveProposedByID != nil {
		proposer := attendee
		if *m.MoveProposedByID == trainer.ID() {
			proposer = trainer
		}
		moveProposedBy = &proposer
	}

	// the database returns times in the local zone, the domain works in UTC
	hour := m.Hour.UTC()
	proposedNewTime := m.ProposedNewTime
	if proposedNewTime != nil {
		utc := proposedNewTime.UTC()
		proposedNewTime = &utc
	}

	return domain.UnmarshalTraining(
		m.ID,
		hour,
		m.Notes,
		attendee,
		trainer,
		proposedNewTime,
		moveProposedBy,
		m.Canceled,
	)
}

func upsertParams(t *domain.Training) dbmodels.UpsertTrainingParams {
	proposedBy := t.MoveProposedBy()

	params := dbmodels.UpsertTrainingParams{
		ID:               t.ID(),
		Hour:             t.Hour(),
		Notes:            t.Notes(),
		AttendeeID:       t.Attendee().ID(),
		AttendeeUsername: t.Attendee().Username(),
		TrainerID:        t.Trainer().ID(),
		TrainerUsername:  t.Trainer().Username(),
		ProposedNewTime:  t.ProposedNewTime(),
		Canceled:         t.IsCanceled(),
	}
	if proposedBy != nil {
		id := proposedBy.ID()
		params.MoveProposedByID = &id
	}

	return params
}
