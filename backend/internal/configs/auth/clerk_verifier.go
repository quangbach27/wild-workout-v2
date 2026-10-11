package auth

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	commonAuth "github.com/quangbach27/golang-common/http/auth"
)

type customClaims struct {
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
}

type verifyFn func(ctx context.Context, params *jwt.VerifyParams) (*clerk.SessionClaims, error)

type clerkVerifier struct {
	verify            verifyFn
	authorizedParties []string
}

var _ commonAuth.TokenVerifier = (*clerkVerifier)(nil)

func NewClerkVerifier(secretKey string, authrorizedParties []string) *clerkVerifier {
	if secretKey == "" {
		panic("clerkauth: secret key is required")
	}
	clerk.SetKey(secretKey)

	return &clerkVerifier{
		authorizedParties: authrorizedParties,
		verify:            jwt.Verify,
	}
}

func (v *clerkVerifier) Verify(ctx context.Context, token string) (*commonAuth.Session, error) {
	params := v.buildVerifyParams(token)

	claims, err := v.verify(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to verify clerk token: %w", err)
	}

	return v.sessionFromClaims(claims)
}

func (v *clerkVerifier) buildVerifyParams(token string) *jwt.VerifyParams {
	params := &jwt.VerifyParams{
		Token: token,
		CustomClaimsConstructor: func(context.Context) any {
			return &customClaims{}
		},
	}

	if len(v.authorizedParties) != 0 {
		params.AuthorizedPartyHandler = func(azp string) bool {
			return slices.Contains(v.authorizedParties, azp)
		}
	}

	return params
}

func (v *clerkVerifier) sessionFromClaims(claims *clerk.SessionClaims) (*commonAuth.Session, error) {
	if claims == nil || claims.Subject == "" {
		return nil, errors.New("token has no subject")
	}

	custom, _ := claims.Custom.(*customClaims)
	if custom == nil {
		custom = &customClaims{}
	}

	session := &commonAuth.Session{
		UserID: claims.Subject,
		// Clone so callers can mutate Roles without writing into the decoded claims.
		Roles: slices.Clone(custom.Roles),
	}
	if custom.Username != "" {
		session.Extra = map[string]any{"username": custom.Username}
	}

	return session, nil
}
