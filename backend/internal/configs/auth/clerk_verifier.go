package auth

import (
	"backend/internal/shared"
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	clerk "github.com/clerk/clerk-sdk-go/v2"
	"github.com/clerk/clerk-sdk-go/v2/jwt"
	commonAuth "github.com/quangbach27/golang-common/http/auth"
	"github.com/quangbach27/golang-common/log"
)

type verifyFn func(ctx context.Context, params *jwt.VerifyParams) (*clerk.SessionClaims, error)

type UserReadModel interface {
	GetUser(ctx context.Context, userID string) (*shared.User, error)
}

type clerkVerifier struct {
	verify            verifyFn
	authorizedParties []string
	userReadModel     UserReadModel
}

var _ commonAuth.TokenVerifier = (*clerkVerifier)(nil)

func NewClerkVerifier(
	secretKey string,
	authorizedParties []string,
	userReadModel UserReadModel,
) *clerkVerifier {
	if secretKey == "" {
		panic("clerkauth: secret key is required")
	}
	if userReadModel == nil {
		panic("userReadModel can't be nil")
	}
	clerk.SetKey(secretKey)

	return &clerkVerifier{
		authorizedParties: authorizedParties,
		verify:            jwt.Verify,
		userReadModel:     userReadModel,
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

	return v.sessionFromClaims(ctx, claims)
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

func (v *clerkVerifier) sessionFromClaims(ctx context.Context, claims *clerk.SessionClaims) (*commonAuth.Session, error) {
	if claims == nil || claims.Subject == "" {
		return nil, errors.New("token has no subject")
	}

	// OrgRole has format: org:{role}
	var role string
	roles := strings.Split(claims.ActiveOrganizationRole, ":")
	if len(roles) == 2 {
		role = roles[1]
	}

	if role == "" {
		log.FromContext(ctx).Debug("Roles empty in claims. Retrieving role from user database")
		user, err := v.userReadModel.GetUser(ctx, claims.Subject)
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve user info: %w", err)
		}

		role = user.Role().String()
	}

	session := &commonAuth.Session{
		Subject:   claims.Subject,
		SessionID: claims.SessionID,
		Role:      role,
	}

	return session, nil
}
