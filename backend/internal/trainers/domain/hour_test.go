package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"backend/internal/trainers/domain"
)

func newAvailableHour(t *testing.T, f *domain.HourFactory) *domain.Hour {
	t.Helper()

	h, err := f.NewAvailableHour(testTrainerUUID, utcHourInDays(1, f.Config().MinUtcHour))
	require.NoError(t, err)

	return h
}

func newNotAvailableHour(t *testing.T, f *domain.HourFactory) *domain.Hour {
	t.Helper()

	h, err := f.NewNotAvailableHour(testTrainerUUID, utcHourInDays(1, f.Config().MinUtcHour))
	require.NoError(t, err)

	return h
}

func newTrainingScheduledHour(t *testing.T, f *domain.HourFactory) *domain.Hour {
	t.Helper()

	h := newAvailableHour(t, f)
	require.NoError(t, h.ScheduleTraining())

	return h
}

func TestHour_ScheduleTraining(t *testing.T) {
	t.Parallel()

	f := newDefaultFactory(t)

	testCases := []struct {
		name       string
		hour       *domain.Hour
		wantErr    bool
		wantStatus domain.HourStatus
	}{
		{
			name:       "available hour",
			hour:       newAvailableHour(t, f),
			wantStatus: domain.TrainingScheduled,
		},
		{
			name:       "not available hour",
			hour:       newNotAvailableHour(t, f),
			wantErr:    true,
			wantStatus: domain.NotAvailable,
		},
		{
			name:       "already scheduled hour",
			hour:       newTrainingScheduledHour(t, f),
			wantErr:    true,
			wantStatus: domain.TrainingScheduled,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.hour.ScheduleTraining()

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.wantStatus, tc.hour.Status())
		})
	}
}

func TestHour_MakeAvailable(t *testing.T) {
	t.Parallel()

	f := newDefaultFactory(t)

	testCases := []struct {
		name       string
		hour       *domain.Hour
		wantErr    bool
		wantStatus domain.HourStatus
	}{
		{
			name:       "not available hour",
			hour:       newNotAvailableHour(t, f),
			wantStatus: domain.Available,
		},
		{
			name:       "already available hour",
			hour:       newAvailableHour(t, f),
			wantErr:    true,
			wantStatus: domain.Available,
		},
		{
			name:       "training scheduled hour",
			hour:       newTrainingScheduledHour(t, f),
			wantErr:    true,
			wantStatus: domain.TrainingScheduled,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.hour.MakeAvailable()

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.wantStatus, tc.hour.Status())
		})
	}
}

func TestHour_MakeNotAvailable(t *testing.T) {
	t.Parallel()

	f := newDefaultFactory(t)

	testCases := []struct {
		name       string
		hour       *domain.Hour
		wantErr    bool
		wantStatus domain.HourStatus
	}{
		{
			name:       "available hour",
			hour:       newAvailableHour(t, f),
			wantStatus: domain.NotAvailable,
		},
		{
			name:       "already not available hour",
			hour:       newNotAvailableHour(t, f),
			wantErr:    true,
			wantStatus: domain.NotAvailable,
		},
		{
			name:       "training scheduled hour",
			hour:       newTrainingScheduledHour(t, f),
			wantErr:    true,
			wantStatus: domain.TrainingScheduled,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.hour.MakeNotAvailable()

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.wantStatus, tc.hour.Status())
		})
	}
}

func TestHour_CancelTraining(t *testing.T) {
	t.Parallel()

	f := newDefaultFactory(t)

	testCases := []struct {
		name       string
		hour       *domain.Hour
		wantErr    bool
		wantStatus domain.HourStatus
	}{
		{
			name:       "training scheduled hour",
			hour:       newTrainingScheduledHour(t, f),
			wantStatus: domain.Available,
		},
		{
			name:       "available hour",
			hour:       newAvailableHour(t, f),
			wantErr:    true,
			wantStatus: domain.Available,
		},
		{
			name:       "not available hour",
			hour:       newNotAvailableHour(t, f),
			wantErr:    true,
			wantStatus: domain.NotAvailable,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.hour.CancelTraining()

			if tc.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, tc.wantStatus, tc.hour.Status())
		})
	}
}

func TestHour_StatusPredicates(t *testing.T) {
	t.Parallel()

	f := newDefaultFactory(t)

	testCases := []struct {
		name                     string
		hour                     *domain.Hour
		wantAvailable            bool
		wantNotAvailable         bool
		wantHasTrainingScheduled bool
	}{
		{
			name:          "available hour",
			hour:          newAvailableHour(t, f),
			wantAvailable: true,
		},
		{
			name:             "not available hour",
			hour:             newNotAvailableHour(t, f),
			wantNotAvailable: true,
		},
		{
			name:                     "training scheduled hour",
			hour:                     newTrainingScheduledHour(t, f),
			wantHasTrainingScheduled: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.wantAvailable, tc.hour.IsAvailable())
			assert.Equal(t, tc.wantNotAvailable, tc.hour.IsNotAvailable())
			assert.Equal(t, tc.wantHasTrainingScheduled, tc.hour.HasTrainingScheduled())
		})
	}
}
