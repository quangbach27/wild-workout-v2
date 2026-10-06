package domain

import (
	"errors"
	"fmt"

	common "github.com/quangbach27/golang-common"

	"backend/internal/shared"
)

type User struct {
	id       string
	username string
	role     shared.Role
}

func (u User) ID() string {
	return u.id
}

func (u User) Username() string {
	return u.username
}

func (u User) Role() shared.Role {
	return u.role
}

func (u User) Equal(o User) bool {
	return u.id == o.id && u.role.Equal(o.role.Enum)
}

// NewUser validates the input and returns a common.Error (400) whose details list each failure
// (invalid-user-id, invalid-user-username, invalid-user-role).
func NewUser(id, username string, role shared.Role) (User, error) {
	errDetails := []common.ErrorDetails{}

	if id == "" {
		errDetails = append(errDetails, common.ErrorDetails{
			EntityType: "User",
			ErrorSlug:  "invalid-user-id",
			Message:    "id can't be empty",
		})
	}

	if username == "" {
		errDetails = append(errDetails, common.ErrorDetails{
			EntityType: "User",
			ErrorSlug:  "invalid-user-username",
			Message:    "username can't be empty",
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
		return User{}, common.NewInvalidInputError("invalid-user", "user is not valid").WithDetails(errDetails)
	}

	return User{id: id, username: username, role: role}, nil
}

func CanUserSeeTraining(user User, training *Training) error {
	if training == nil {
		return errors.New("training is empty")
	}

	if !user.Equal(training.trainer) && !user.Equal(training.attendee) {
		return fmt.Errorf("user %s can't see training %s", user.id, training.id)
	}

	return nil
}
