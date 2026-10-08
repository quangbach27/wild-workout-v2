package queries

import "errors"

type Handler struct {
	userTrainingsReadModel UserTrainingsReadModel
}

func NewHandler(
	userTrainingsReadModel UserTrainingsReadModel,
) *Handler {
	var errs []error

	if userTrainingsReadModel == nil {
		errs = append(errs, errors.New("userTrainingsReadModel can't be nil"))
	}

	if len(errs) != 0 {
		panic(errors.Join(errs...))
	}

	return &Handler{
		userTrainingsReadModel: userTrainingsReadModel,
	}
}
