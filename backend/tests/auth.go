//go:build component

package tests

import (
	"context"
	"net/http"

	common "github.com/quangbach27/golang-common"
	commonAuth "github.com/quangbach27/golang-common/http/auth"
)

// Username is the username NewSession gives to userID.
func Username(userID string) string {
	return "name-" + userID
}

// NewSession registers a bearer token that authenticates as userID (with username Username(userID))
// and the given role (only the first of roles is used), and returns the token.
func NewSession(userID string, roles ...string) string {
	token := common.NewUUIDv7().String()

	var role string
	if len(roles) > 0 {
		role = roles[0]
	}

	verifier.Add(token, &commonAuth.Session{
		Subject: userID,
		Role:    role,
		Extra:   map[string]any{"username": Username(userID)},
	})

	return token
}

// WithAuth is a request editor that authenticates the request with the bearer token.
func WithAuth(token string) func(ctx context.Context, req *http.Request) error {
	return func(_ context.Context, req *http.Request) error {
		req.Header.Set("Authorization", "Bearer "+token)
		return nil
	}
}
