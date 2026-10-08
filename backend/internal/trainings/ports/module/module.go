package module

import "backend/internal/trainings/ports/module/client"

type Trainings struct{}

var _ client.Trainings = (*Trainings)(nil)

func New() *Trainings {
	return &Trainings{}
}

func (t *Trainings) PingTrainings(request client.PingTrainingsRequest) error {
	return nil
}
