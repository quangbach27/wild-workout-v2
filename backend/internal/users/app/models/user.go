package models

import (
	"context"

	common "github.com/quangbach27/golang-common"

	"backend/internal/shared"
)

type UserRepository interface {
	OnboardUser(ctx context.Context, user *User) error
	// GetUserByUUID returns a 404 common.Error if the user doesn't exist.
	GetUserByUUID(ctx context.Context, userID string) (*User, error)
	// UpdateDisplayName sets the display name. It returns a 404 common.Error if the user doesn't exist.
	UpdateDisplayName(ctx context.Context, userID, displayName string) error
	// UpdateUser persists the display name and role. It returns a 404 common.Error if the user doesn't exist.
	UpdateUser(ctx context.Context, user *User) error
	// UpdateBalance adds amountChange (negative to spend) to the balance and returns the new one.
	// It returns a 404 common.Error if the user doesn't exist and a 409 if the balance would
	// become negative.
	UpdateBalance(ctx context.Context, userID string, amountChange int) (int, error)
}

type User struct {
	uuid        string
	balance     int
	displayName string
	role        shared.Role
}

func (u User) UUID() string {
	return u.uuid
}

func (u User) Balance() int {
	return u.balance
}

func (u User) DisplayName() string {
	return u.displayName
}

func (u User) Role() shared.Role {
	return u.role
}

// NewUser validates the input and returns a new user with a zero balance. On failure it returns a
// common.Error (400) whose details list each failure (invalid-user-uuid, invalid-user-display-name,
// invalid-user-role).
func NewUser(uuid, displayName string, role shared.Role) (*User, error) {
	errDetails := []common.ErrorDetails{}

	if uuid == "" {
		errDetails = append(errDetails, common.ErrorDetails{
			EntityType: "User",
			ErrorSlug:  "invalid-user-uuid",
			Message:    "uuid can't be empty",
		})
	}

	if displayName == "" {
		errDetails = append(errDetails, common.ErrorDetails{
			EntityType: "User",
			ErrorSlug:  "invalid-user-display-name",
			Message:    "display name can't be empty",
		})
	}

	if role.IsZero() {
		errDetails = append(errDetails, common.ErrorDetails{
			EntityType: "User",
			ErrorSlug:  "invalid-user-role",
			Message:    "role can't be empty",
		})
	}

	if len(errDetails) != 0 {
		return nil, common.NewInvalidInputError("invalid-user", "user is not valid").WithDetails(errDetails)
	}

	return &User{uuid: uuid, displayName: displayName, role: role}, nil
}

// Update changes the display name and/or role; nil leaves a field unchanged. It returns a
// common.Error (400) whose details list each failure (invalid-user-display-name, invalid-user-role)
// and leaves the user untouched on failure.
func (u *User) Update(displayName *string, role *shared.Role) error {
	errDetails := []common.ErrorDetails{}

	if displayName != nil && *displayName == "" {
		errDetails = append(errDetails, common.ErrorDetails{
			EntityType: "User",
			ErrorSlug:  "invalid-user-display-name",
			Message:    "display name can't be empty",
		})
	}

	if role != nil && role.IsZero() {
		errDetails = append(errDetails, common.ErrorDetails{
			EntityType: "User",
			ErrorSlug:  "invalid-user-role",
			Message:    "role can't be empty",
		})
	}

	if len(errDetails) != 0 {
		return common.NewInvalidInputError("invalid-user", "user is not valid").WithDetails(errDetails)
	}

	if displayName != nil {
		u.displayName = *displayName
	}
	if role != nil {
		u.role = *role
	}

	return nil
}

// UnmarshalUser rebuilds a user from persisted data, bypassing validation.
func UnmarshalUser(uuid, displayName string, balance int, role shared.Role) *User {
	return &User{uuid: uuid, balance: balance, displayName: displayName, role: role}
}
