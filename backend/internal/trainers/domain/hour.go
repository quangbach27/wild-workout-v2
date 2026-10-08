package domain

import (
	"context"
	"fmt"
	"time"

	common "github.com/quangbach27/golang-common"
)

type HourRepository interface {
	UpsertHours(
		ctx context.Context,
		trainerUUID string,
		hour []time.Time,
		upsertFn func(hours []*Hour) error,
	) error
}

type HourStatus struct {
	common.Enum[StatusTypes]
}

type StatusTypes string

func (t StatusTypes) Values() []string {
	return []string{"availability", "not-availability", "training-scheduled"}
}

var (
	Available         = common.MustEnum[HourStatus]("availability")
	NotAvailable      = common.MustEnum[HourStatus]("not-availability")
	TrainingScheduled = common.MustEnum[HourStatus]("training-scheduled")
)

type Hour struct {
	hour        time.Time
	status      HourStatus
	trainerUUID string
}

func (h *Hour) Hour() time.Time {
	return h.hour
}

func (h *Hour) Status() HourStatus {
	return h.status
}

func (h *Hour) TrainerUUID() string {
	return h.trainerUUID
}

func (h *Hour) IsAvailable() bool {
	return h.status == Available
}

func (h *Hour) IsNotAvailable() bool {
	return h.status == NotAvailable
}

func (h *Hour) HasTrainingScheduled() bool {
	return h.status == TrainingScheduled
}

func (h *Hour) ScheduleTraining() error {
	if !h.IsAvailable() {
		return fmt.Errorf("can not schedule training hour because hour is %s", h.status.String())
	}

	h.status = TrainingScheduled
	return nil
}

func (h *Hour) MakeAvailable() error {
	if !h.IsNotAvailable() {
		return fmt.Errorf("can not make hour available because hour is %s", h.status.String())
	}

	h.status = Available
	return nil
}

func (h *Hour) MakeNotAvailable() error {
	if !h.IsAvailable() {
		return fmt.Errorf("can not make hour not available because hour is %s", h.status.String())
	}

	h.status = NotAvailable
	return nil
}

func (h *Hour) CancelTraining() error {
	if !h.HasTrainingScheduled() {
		return fmt.Errorf("can not cancel training because hour isn't scheduled yet")
	}

	h.status = Available
	return nil
}
