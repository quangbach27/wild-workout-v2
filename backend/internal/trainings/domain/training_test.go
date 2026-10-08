package domain_test

import (
	"net/http"
	"testing"
	"time"

	common "github.com/quangbach27/golang-common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"backend/internal/shared"
	"backend/internal/trainings/domain"
)

func newUser(t *testing.T, id string, role shared.Role) domain.User {
	t.Helper()

	user, err := domain.NewUser(id, "name-"+id, role)
	require.NoError(t, err)

	return user
}

func futureHour(days int) time.Time {
	return time.Now().UTC().Truncate(time.Hour).AddDate(0, 0, days)
}

type trainingUsers struct {
	attendee, trainer, stranger domain.User
}

func newTrainingUsers(t *testing.T) trainingUsers {
	t.Helper()

	return trainingUsers{
		attendee: newUser(t, "attendee-1", shared.RoleAttendee),
		trainer:  newUser(t, "trainer-1", shared.RoleTrainer),
		stranger: newUser(t, "attendee-2", shared.RoleAttendee),
	}
}

func newTraining(t *testing.T, users trainingUsers) *domain.Training {
	t.Helper()

	training, err := domain.NewTraining(users.attendee, users.trainer, futureHour(3), "legs day")
	require.NoError(t, err)

	return training
}

func requireCommonError(t *testing.T, err error, wantStatus int, wantSlug string) common.Error {
	t.Helper()

	var commonErr common.Error
	require.ErrorAs(t, err, &commonErr)
	assert.Equal(t, wantStatus, commonErr.HttpErrorCode)
	assert.Equal(t, wantSlug, commonErr.ErrorSlug)

	return commonErr
}

func TestNewTraining(t *testing.T) {
	t.Parallel()

	users := newTrainingUsers(t)

	testCases := []struct {
		name      string
		attendee  domain.User
		trainer   domain.User
		hour      time.Time
		notes     string
		wantSlugs []string // empty = valid training, otherwise in order
	}{
		{name: "valid", attendee: users.attendee, trainer: users.trainer, hour: futureHour(1), notes: "notes"},
		{
			name: "past hour", attendee: users.attendee, trainer: users.trainer,
			hour: futureHour(-1), notes: "notes", wantSlugs: []string{"invalid-training-hour"},
		},
		{
			name: "zero hour", attendee: users.attendee, trainer: users.trainer,
			notes: "notes", wantSlugs: []string{"invalid-training-hour"},
		},
		{
			name: "attendee with trainer role", attendee: users.trainer, trainer: users.trainer,
			hour: futureHour(1), notes: "notes", wantSlugs: []string{"invalid-training-attendee"},
		},
		{
			name: "trainer with attendee role", attendee: users.attendee, trainer: users.attendee,
			hour: futureHour(1), notes: "notes", wantSlugs: []string{"invalid-training-trainer"},
		},
		{
			name: "empty notes", attendee: users.attendee, trainer: users.trainer,
			hour: futureHour(1),
		},
		{
			name: "multiple failures", wantSlugs: []string{
				"invalid-training-hour",
				"invalid-training-attendee",
				"invalid-training-trainer",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			training, err := domain.NewTraining(tc.attendee, tc.trainer, tc.hour, tc.notes)

			if len(tc.wantSlugs) > 0 {
				commonErr := requireCommonError(t, err, http.StatusBadRequest, "invalid-training")
				require.Len(t, commonErr.Details, len(tc.wantSlugs))
				for i, slug := range tc.wantSlugs {
					assert.Equal(t, "Training", commonErr.Details[i].EntityType)
					assert.Equal(t, slug, commonErr.Details[i].ErrorSlug)
				}
				assert.Nil(t, training)

				return
			}

			require.NoError(t, err)
			assert.False(t, training.ID().IsZero())
			assert.Equal(t, tc.hour, training.Hour())
			assert.Equal(t, tc.notes, training.Notes())
			assert.Equal(t, tc.attendee, training.Attendee())
			assert.Equal(t, tc.trainer, training.Trainer())
			assert.False(t, training.IsCanceled())
			assert.False(t, training.IsRescheduleProposed())
		})
	}
}

func TestTraining_Cancel(t *testing.T) {
	t.Parallel()

	users := newTrainingUsers(t)

	testCases := []struct {
		name        string
		user        domain.User
		cancelFirst bool
		wantStatus  int // 0 = success
	}{
		{name: "attendee cancels", user: users.attendee},
		{name: "trainer cancels", user: users.trainer},
		{name: "stranger", user: users.stranger, wantStatus: http.StatusForbidden},
		{name: "already canceled", user: users.attendee, cancelFirst: true, wantStatus: http.StatusConflict},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			training := newTraining(t, users)
			if tc.cancelFirst {
				require.NoError(t, training.Cancel(users.trainer))
			}

			err := training.Cancel(tc.user)

			if tc.wantStatus != 0 {
				_ = requireCommonError(t, err, tc.wantStatus, "failed-to-cancel")
				return
			}

			require.NoError(t, err)
			assert.True(t, training.IsCanceled())
		})
	}
}

func TestTraining_ProposeReschedule(t *testing.T) {
	t.Parallel()

	users := newTrainingUsers(t)

	testCases := []struct {
		name       string
		user       domain.User
		newTime    time.Time
		canceled   bool
		wantStatus int // 0 = success
	}{
		{name: "attendee proposes", user: users.attendee, newTime: futureHour(5)},
		{name: "trainer proposes", user: users.trainer, newTime: futureHour(5)},
		{name: "stranger", user: users.stranger, newTime: futureHour(5), wantStatus: http.StatusForbidden},
		{name: "canceled", user: users.attendee, newTime: futureHour(5), canceled: true, wantStatus: http.StatusConflict},
		{name: "past time", user: users.attendee, newTime: futureHour(-1), wantStatus: http.StatusBadRequest},
		{name: "zero time", user: users.attendee, wantStatus: http.StatusBadRequest},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			training := newTraining(t, users)
			originalHour := training.Hour()
			if tc.canceled {
				require.NoError(t, training.Cancel(users.attendee))
			}

			err := training.ProposeReschedule(tc.newTime, tc.user)

			if tc.wantStatus != 0 {
				_ = requireCommonError(t, err, tc.wantStatus, "failed-to-propose-reschedule")
				assert.False(t, training.IsRescheduleProposed())
				return
			}

			require.NoError(t, err)
			assert.True(t, training.IsRescheduleProposed())
			assert.Equal(t, tc.newTime, *training.ProposedNewTime())
			assert.Equal(t, tc.user, *training.MoveProposedBy())
			assert.Equal(t, originalHour, training.Hour(), "hour must not change until approved")
		})
	}
}

func TestTraining_ApproveReschedule(t *testing.T) {
	t.Parallel()

	users := newTrainingUsers(t)
	newTime := futureHour(5)

	testCases := []struct {
		name       string
		proposer   *domain.User // nil = nothing proposed
		approver   domain.User
		canceled   bool
		wantStatus int // 0 = success
	}{
		{name: "trainer approves attendee proposal", proposer: &users.attendee, approver: users.trainer},
		{name: "attendee approves trainer proposal", proposer: &users.trainer, approver: users.attendee},
		{name: "proposer approves own", proposer: &users.attendee, approver: users.attendee, wantStatus: http.StatusForbidden},
		{name: "stranger", proposer: &users.attendee, approver: users.stranger, wantStatus: http.StatusForbidden},
		{name: "nothing proposed", approver: users.trainer, wantStatus: http.StatusConflict},
		{
			name: "canceled", proposer: &users.attendee, approver: users.trainer, canceled: true,
			wantStatus: http.StatusConflict,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			training := newTraining(t, users)
			originalHour := training.Hour()
			if tc.proposer != nil {
				require.NoError(t, training.ProposeReschedule(newTime, *tc.proposer))
			}
			if tc.canceled {
				require.NoError(t, training.Cancel(users.attendee))
			}

			err := training.ApproveReschedule(tc.approver)

			if tc.wantStatus != 0 {
				_ = requireCommonError(t, err, tc.wantStatus, "failed-to-approve-reschedule")
				assert.Equal(t, originalHour, training.Hour())
				return
			}

			require.NoError(t, err)
			assert.Equal(t, newTime, training.Hour())
			assert.False(t, training.IsRescheduleProposed())
			assert.Nil(t, training.ProposedNewTime())
			assert.Nil(t, training.MoveProposedBy())
		})
	}
}

func TestTraining_RejectReschedule(t *testing.T) {
	t.Parallel()

	users := newTrainingUsers(t)

	testCases := []struct {
		name       string
		proposed   bool
		user       domain.User
		wantStatus int // 0 = success
	}{
		{name: "other party rejects", proposed: true, user: users.trainer},
		{name: "stranger", proposed: true, user: users.stranger, wantStatus: http.StatusForbidden},
		{name: "nothing proposed", user: users.trainer, wantStatus: http.StatusConflict},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			training := newTraining(t, users)
			originalHour := training.Hour()
			if tc.proposed {
				require.NoError(t, training.ProposeReschedule(futureHour(5), users.attendee))
			}

			err := training.RejectReschedule(tc.user)

			if tc.wantStatus != 0 {
				_ = requireCommonError(t, err, tc.wantStatus, "failed-to-reject-reschedule")
				return
			}

			require.NoError(t, err)
			assert.Equal(t, originalHour, training.Hour())
			assert.False(t, training.IsRescheduleProposed())
		})
	}
}

func TestTraining_CanBeCanceledForFree(t *testing.T) {
	t.Parallel()

	users := newTrainingUsers(t)

	testCases := []struct {
		name string
		hour time.Time
		want bool
	}{
		{name: "more than 24h away", hour: time.Now().Add(48 * time.Hour), want: true},
		{name: "less than 24h away", hour: time.Now().Add(2 * time.Hour), want: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			training, err := domain.NewTraining(users.attendee, users.trainer, tc.hour, "notes")
			require.NoError(t, err)

			assert.Equal(t, tc.want, training.CanBeCanceledForFree())
		})
	}
}
