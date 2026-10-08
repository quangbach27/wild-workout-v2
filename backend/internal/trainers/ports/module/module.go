package module

import (
	"context"
	"errors"
	"time"

	"backend/internal/trainers/app/commands"
	"backend/internal/trainers/ports/module/client"
)

type Trainers struct {
	commands *commands.Handler
}

var _ client.Trainers = (*Trainers)(nil)

func New(commands *commands.Handler) *Trainers {
	if commands == nil {
		panic(errors.New("commands can't be nil"))
	}

	return &Trainers{commands: commands}
}

func (t *Trainers) ScheduleHour(
	ctx context.Context,
	req client.ScheduleHourRequest,
) (client.ScheduleHourResponse, error) {
	err := t.commands.ScheduleTraining(ctx, commands.ScheduleTrainingCmd{
		TrainerUUID: req.TrainerUUID,
		Hours:       []time.Time{req.Hour},
	})

	return client.ScheduleHourResponse{}, err
}

func (t *Trainers) CancelHourSchedule(
	ctx context.Context,
	req client.CancelHourScheduleRequest,
) (client.CancelHourScheduleResponse, error) {
	err := t.commands.CancelTraining(ctx, commands.CancelTrainingCmd{
		TrainerUUID: req.TrainerUUID,
		Hours:       []time.Time{req.Hour},
	})

	return client.CancelHourScheduleResponse{}, err
}
