package client

type Trainings interface {
	PingTrainings(request PingTrainingsRequest) error
}

type PingTrainingsRequest struct{}
