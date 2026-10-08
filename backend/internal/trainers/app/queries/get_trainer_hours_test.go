package queries_test

import (
	"testing"
	"time"

	common "github.com/quangbach27/golang-common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"backend/internal/trainers/app/queries"
)

func TestGetTrainerHoursQuery_Validate(t *testing.T) {
	t.Parallel()

	from := time.Date(2030, 1, 1, 0, 0, 0, 0, time.UTC)

	testCases := []struct {
		name         string
		query        queries.GetTrainerHoursQuery
		wantErrSlugs []string
	}{
		{
			name:  "valid query",
			query: queries.GetTrainerHoursQuery{TrainerUUID: "trainer", DateFrom: from, DateTo: from.AddDate(0, 0, 7)},
		},
		{
			name:  "same day",
			query: queries.GetTrainerHoursQuery{TrainerUUID: "trainer", DateFrom: from, DateTo: from},
		},
		{
			name:         "empty trainer uuid",
			query:        queries.GetTrainerHoursQuery{DateFrom: from, DateTo: from},
			wantErrSlugs: []string{"invalid-trainer-uuid"},
		},
		{
			name:         "zero dates",
			query:        queries.GetTrainerHoursQuery{TrainerUUID: "trainer"},
			wantErrSlugs: []string{"invalid-date-from", "invalid-date-to"},
		},
		{
			name:         "date to before date from",
			query:        queries.GetTrainerHoursQuery{TrainerUUID: "trainer", DateFrom: from, DateTo: from.AddDate(0, 0, -1)},
			wantErrSlugs: []string{"invalid-date-range"},
		},
		{
			name:         "everything invalid",
			query:        queries.GetTrainerHoursQuery{},
			wantErrSlugs: []string{"invalid-trainer-uuid", "invalid-date-from", "invalid-date-to"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.query.Validate()

			if len(tc.wantErrSlugs) == 0 {
				require.NoError(t, err)
				return
			}

			var commonErr common.Error
			require.ErrorAs(t, err, &commonErr)

			gotSlugs := make([]string, 0, len(commonErr.Details))
			for _, d := range commonErr.Details {
				gotSlugs = append(gotSlugs, d.ErrorSlug)
			}
			assert.ElementsMatch(t, tc.wantErrSlugs, gotSlugs)
		})
	}
}
