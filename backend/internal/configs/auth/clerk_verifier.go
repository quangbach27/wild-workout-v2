package auth

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	clerk "github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	commonAuth "github.com/quangbach27/golang-common/http/auth"
)

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
	claims, err := v.verify(
		ctx,
		v.buildVerifyParams(token),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to verify clerk token: %w", err)
	}

	return v.sessionFromClaims(claims)
}

func (v *clerkVerifier) buildVerifyParams(token string) *jwt.VerifyParams {
	params := &jwt.VerifyParams{
		Token: token,
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

	// OrgRole has format: org:{role}
	var role string
	roles := strings.Split(claims.ActiveOrganizationRole, ":")
	if len(roles) == 2 {
		role = roles[1]
	}

	session := &commonAuth.Session{
		Subject:   claims.Subject,
		SessionID: claims.SessionID,
		Role:      role,
	}

	return session, nil
}
