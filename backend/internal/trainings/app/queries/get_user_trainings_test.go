package queries_test

import (
	"net/http"
	"testing"

	common "github.com/quangbach27/golang-common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"backend/internal/trainings/app/queries"
)

func TestGetUserTrainingsQuery_Validate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		query     queries.GetUserTrainingsQuery
		wantSlugs []string // empty = valid query
	}{
		{name: "valid", query: queries.GetUserTrainingsQuery{UserUUID: "user-1", Page: 1, PageSize: 10}},
		{name: "later page", query: queries.GetUserTrainingsQuery{UserUUID: "user-1", Page: 7, PageSize: 10}},
		{
			name:  "max page size",
			query: queries.GetUserTrainingsQuery{UserUUID: "user-1", Page: 1, PageSize: queries.MaxPageSize},
		},
		{
			name:      "empty user uuid",
			query:     queries.GetUserTrainingsQuery{Page: 1, PageSize: 10},
			wantSlugs: []string{"invalid-user-uuid"},
		},
		{
			name:      "zero page",
			query:     queries.GetUserTrainingsQuery{UserUUID: "user-1", PageSize: 10},
			wantSlugs: []string{"invalid-page"},
		},
		{
			name:      "negative page",
			query:     queries.GetUserTrainingsQuery{UserUUID: "user-1", Page: -1, PageSize: 10},
			wantSlugs: []string{"invalid-page"},
		},
		{
			name:      "zero page size",
			query:     queries.GetUserTrainingsQuery{UserUUID: "user-1", Page: 1},
			wantSlugs: []string{"invalid-page-size"},
		},
		{
			name:      "page size above max",
			query:     queries.GetUserTrainingsQuery{UserUUID: "user-1", Page: 1, PageSize: queries.MaxPageSize + 1},
			wantSlugs: []string{"invalid-page-size"},
		},
		{
			name:      "everything invalid",
			query:     queries.GetUserTrainingsQuery{},
			wantSlugs: []string{"invalid-user-uuid", "invalid-page", "invalid-page-size"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := tc.query.Validate()

			if len(tc.wantSlugs) > 0 {
				var commonErr common.Error
				require.ErrorAs(t, err, &commonErr)
				assert.Equal(t, http.StatusBadRequest, commonErr.HttpErrorCode)
				assert.Equal(t, "invalid-get-user-trainings-query", commonErr.ErrorSlug)
				require.Len(t, commonErr.Details, len(tc.wantSlugs))
				for i, slug := range tc.wantSlugs {
					assert.Equal(t, "GetUserTrainingsQuery", commonErr.Details[i].EntityType)
					assert.Equal(t, slug, commonErr.Details[i].ErrorSlug)
				}

				return
			}

			require.NoError(t, err)
		})
	}
}
