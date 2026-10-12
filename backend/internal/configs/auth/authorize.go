package auth

import (
	"context"

	common "github.com/quangbach27/golang-common"
	commonAuth "github.com/quangbach27/golang-common/http/auth"

	"backend/internal/shared"
)

// Authorize returns the session of ctx if its role is allowed the permission by the policy. It
// returns a 401 error when there is no session and a 403 error when the permission is denied or
// the session has no (or an unknown) role.
func Authorize(ctx context.Context, policy *Policy, permission Permission) (*commonAuth.Session, error) {
	session, ok := commonAuth.SessionFromContext(ctx)
	if !ok || session.Subject == "" {
		return nil, common.NewUnauthorizedError("missing-session", "missing session")
	}

	forbidden := common.NewForbiddenError(
		"forbidden", "you are not allowed to %s this %s", permission.Action, permission.Resource,
	)

	role, err := common.NewEnum[shared.Role](session.Role)
	if err != nil || role.IsZero() {
		return nil, forbidden
	}

	if !policy.Allows([]shared.Role{role}, permission) {
		return nil, forbidden
	}

	return session, nil
}
