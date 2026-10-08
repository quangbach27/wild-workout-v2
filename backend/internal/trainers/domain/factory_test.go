package domain_test

import (
	"net/http"
	"testing"
	"time"

	common "github.com/quangbach27/golang-common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"backend/internal/trainers/domain"
)

const testTrainerUUID = "trainer-1"

// utcHourInDays returns a full hour at the given UTC hour, days from now.
func utcHourInDays(days, hour int) time.Time {
	d := time.Now().UTC().AddDate(0, 0, days) // AddDate, not Add: not every day has 24h

	return time.Date(d.Year(), d.Month(), d.Day(), hour, 0, 0, 0, time.UTC)
}

func newDefaultFactory(t *testing.T) *domain.HourFactory {
	t.Helper()

	f, err := domain.NewHourFactory()
	require.NoError(t, err)

	return f
}

func TestNewHourFactory_Config(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		fn       func(c *domain.HourFactoryConfig)
		wantErrs int // 0 = valid config, N = number of joined validation errors
	}{
		{
			name: "default config",
		},
		{
			name: "custom utc hours",
			fn: func(c *domain.HourFactoryConfig) {
				c.MinUtcHour = 10
				c.MaxUtcHour = 12
			},
		},
		{
			name: "custom max weeks in the future",
			fn: func(c *domain.HourFactoryConfig) {
				c.MaxWeeksInTheFutureToSet = 1
			},
		},
		{
			name: "full day boundaries",
			fn: func(c *domain.HourFactoryConfig) {
				c.MinUtcHour = 0
				c.MaxUtcHour = 24
			},
		},
		{
			name: "min equals max",
			fn: func(c *domain.HourFactoryConfig) {
				c.MinUtcHour = 10
				c.MaxUtcHour = 10
			},
		},

		// invalid configs
		{
			name:     "zero max weeks",
			fn:       func(c *domain.HourFactoryConfig) { c.MaxWeeksInTheFutureToSet = 0 },
			wantErrs: 1,
		},
		{
			name:     "negative max weeks",
			fn:       func(c *domain.HourFactoryConfig) { c.MaxWeeksInTheFutureToSet = -1 },
			wantErrs: 1,
		},
		{
			name:     "negative min hour",
			fn:       func(c *domain.HourFactoryConfig) { c.MinUtcHour = -1 },
			wantErrs: 1,
		},
		{
			// also fails the min-after-max rule, since max can't exceed 24 without failing too
			name: "min hour above 24",
			fn: func(c *domain.HourFactoryConfig) {
				c.MinUtcHour = 25
				c.MaxUtcHour = 24
			},
			wantErrs: 2,
		},
		{
			// also fails the min-after-max rule against the default min hour
			name:     "negative max hour",
			fn:       func(c *domain.HourFactoryConfig) { c.MaxUtcHour = -1 },
			wantErrs: 2,
		},
		{
			name:     "max hour above 24",
			fn:       func(c *domain.HourFactoryConfig) { c.MaxUtcHour = 25 },
			wantErrs: 1,
		},
		{
			name: "min after max",
			fn: func(c *domain.HourFactoryConfig) {
				c.MinUtcHour = 15
				c.MaxUtcHour = 10
			},
			wantErrs: 1,
		},
		{
			name: "multiple problems are joined",
			fn: func(c *domain.HourFactoryConfig) {
				c.MaxWeeksInTheFutureToSet = 0
				c.MinUtcHour = 15
				c.MaxUtcHour = 10
			},
			wantErrs: 2,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var fns []func(c *domain.HourFactoryConfig)
			if tc.fn != nil {
				fns = append(fns, tc.fn)
			}

			f, err := domain.NewHourFactory(fns...)

			if tc.wantErrs == 0 {
				require.NoError(t, err)
				assert.NotNil(t, f)

				return
			}

			require.Error(t, err)
			assert.Nil(t, f)

			var joined interface{ Unwrap() []error }
			require.ErrorAs(t, err, &joined)
			assert.Len(t, joined.Unwrap(), tc.wantErrs)
		})
	}
}

func TestHourFactory_NewHour(t *testing.T) {
	t.Parallel()
	f := newDefaultFactory(t)

	var (
		notFullHour = utcHourInDays(1, 10).Add(30 * time.Minute)
		pastHour    = utcHourInDays(-1, 10)
		currentHour = time.Now().Truncate(time.Hour)
		tooDistant  = utcHourInDays(f.Config().MaxWeeksInTheFutureToSet*7+1, 10)
		tooEarly    = utcHourInDays(1, f.Config().MinUtcHour-1)
		tooLate     = utcHourInDays(1, f.Config().MaxUtcHour+1)
	)

	const (
		slugTrainerUUID = "invalid-trainer-uuid"
		slugHour        = "invalid-training-hour"
		msgEmptyUUID    = "trainerUUID can't be empty"
	)

	type wantDetail struct {
		slug    string
		message string
	}

	testCases := []struct {
		name        string
		trainerUUID string
		hour        time.Time
		wantDetails []wantDetail // empty = valid hour, otherwise in order
	}{
		// valid hours
		{
			name:        "min utc hour",
			trainerUUID: testTrainerUUID,
			hour:        utcHourInDays(1, f.Config().MinUtcHour),
		},
		{
			name:        "max utc hour",
			trainerUUID: testTrainerUUID,
			hour:        utcHourInDays(1, f.Config().MaxUtcHour),
		},

		// invalid hours
		{
			name:        "not a full hour",
			trainerUUID: testTrainerUUID,
			hour:        notFullHour,
			wantDetails: []wantDetail{{slugHour, domain.ErrNotFullHour.Error()}},
		},
		{
			name:        "past hour",
			trainerUUID: testTrainerUUID,
			hour:        pastHour,
			wantDetails: []wantDetail{{slugHour, domain.ErrPastHour.Error()}},
		},
		{
			name:        "current hour",
			trainerUUID: testTrainerUUID,
			hour:        currentHour,
			wantDetails: []wantDetail{{slugHour, domain.ErrPastHour.Error()}},
		},
		{
			name:        "too distant date",
			trainerUUID: testTrainerUUID,
			hour:        tooDistant,
			wantDetails: []wantDetail{{slugHour, domain.TooDistantDateError{
				MaxWeeksInTheFutureToSet: f.Config().MaxWeeksInTheFutureToSet,
				ProvidedDate:             tooDistant,
			}.Error()}},
		},
		{
			name:        "too early hour",
			trainerUUID: testTrainerUUID,
			hour:        tooEarly,
			wantDetails: []wantDetail{{slugHour, domain.TooEarlyHourError{
				MinUtcHour:   f.Config().MinUtcHour,
				ProvidedTime: tooEarly,
			}.Error()}},
		},
		{
			name:        "too late hour",
			trainerUUID: testTrainerUUID,
			hour:        tooLate,
			wantDetails: []wantDetail{{slugHour, domain.TooLateHourError{
				MaxUtcHour:   f.Config().MaxUtcHour,
				ProvidedTime: tooLate,
			}.Error()}},
		},

		// invalid trainer uuid
		{
			name:        "empty trainer uuid",
			hour:        utcHourInDays(1, f.Config().MinUtcHour),
			wantDetails: []wantDetail{{slugTrainerUUID, msgEmptyUUID}},
		},
		{
			name:        "empty trainer uuid and past hour",
			hour:        pastHour,
			wantDetails: []wantDetail{{slugTrainerUUID, msgEmptyUUID}, {slugHour, domain.ErrPastHour.Error()}},
		},
	}

	constructors := []struct {
		name       string
		create     func(f *domain.HourFactory, trainerUUID string, h time.Time) (*domain.Hour, error)
		wantStatus domain.HourStatus
	}{
		{
			name:       "available",
			create:     (*domain.HourFactory).NewAvailableHour,
			wantStatus: domain.Available,
		},
		{
			name:       "not available",
			create:     (*domain.HourFactory).NewNotAvailableHour,
			wantStatus: domain.NotAvailable,
		},
	}

	for _, constructor := range constructors {
		t.Run(constructor.name, func(t *testing.T) {
			t.Parallel()

			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
					t.Parallel()

					h, err := constructor.create(f, tc.trainerUUID, tc.hour)

					if len(tc.wantDetails) > 0 {
						var commonErr common.Error
						require.ErrorAs(t, err, &commonErr)
						assert.Equal(t, http.StatusBadRequest, commonErr.HttpErrorCode)
						require.Len(t, commonErr.Details, len(tc.wantDetails))
						for i, want := range tc.wantDetails {
							assert.Equal(t, "Hour", commonErr.Details[i].EntityType)
							assert.Equal(t, want.slug, commonErr.Details[i].ErrorSlug)
							assert.Equal(t, want.message, commonErr.Details[i].Message)
						}
						assert.Nil(t, h)

						return
					}

					require.NoError(t, err)
					assert.True(t, tc.hour.Equal(h.Hour()), "expected hour %s, got %s", tc.hour, h.Hour())
					assert.Equal(t, constructor.wantStatus, h.Status())
					assert.Equal(t, tc.trainerUUID, h.TrainerUUID())
				})
			}
		})
	}
}
