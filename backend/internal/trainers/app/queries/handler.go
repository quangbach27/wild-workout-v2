package queries

import "errors"

type Handler struct {
	trainerHoursReadModel TrainerHoursReadModel
}

func NewHandler(
	trainerHoursReadModel TrainerHoursReadModel,
) *Handler {
	var errs []error

	if trainerHoursReadModel == nil {
		errs = append(errs, errors.New("trainerHoursReadModel can't be nil"))
	}

	if len(errs) != 0 {
		panic(errors.Join(errs...))
	}

	return &Handler{
		trainerHoursReadModel: trainerHoursReadModel,
	}
}
