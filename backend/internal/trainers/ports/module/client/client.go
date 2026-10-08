package client

import (
	"context"
	"time"
)

type Trainers interface {
	// ScheduleHour marks the trainer's hour as training-scheduled. It returns a 409 common.Error
	// if the hour is not available.
	ScheduleHour(ctx context.Context, req ScheduleHourRequest) (ScheduleHourResponse, error)
	// CancelHourSchedule frees an hour previously scheduled with ScheduleHour.
	CancelHourSchedule(ctx context.Context, req CancelHourScheduleRequest) (CancelHourScheduleResponse, error)
}

type ScheduleHourRequest struct {
	TrainerUUID string
	Hour        time.Time
}

type ScheduleHourResponse struct{}

type CancelHourScheduleRequest struct {
	TrainerUUID string
	Hour        time.Time
}

type CancelHourScheduleResponse struct{}
