package app

import (
	"context"
	"errors"

	common "github.com/quangbach27/golang-common"
)

const slugUserAlreadyOnboarded = "user-already-onboarded"

// EnsureUser onboards the user unless it already exists, so it is safe to call more than once
// (e.g. for redelivered webhooks). It reports whether the user was created.
func (s *Service) EnsureUser(ctx context.Context, cmd OnboardUserCmd) (bool, error) {
	err := s.OnboardUser(ctx, cmd)

	var commonErr common.Error
	if errors.As(err, &commonErr) && commonErr.ErrorSlug == slugUserAlreadyOnboarded {
		return false, nil
	}

	return err == nil, err
}

// RenameUser sets the display name. It returns a 404 common.Error for an unknown user.
func (s *Service) RenameUser(ctx context.Context, userUUID, displayName string) error {
	if displayName == "" {
		return common.NewInvalidInputError("invalid-user-display-name", "display name can't be empty")
	}

	return s.userRepo.UpdateDisplayName(ctx, userUUID, displayName)
}
