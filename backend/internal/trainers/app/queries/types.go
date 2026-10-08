package queries

import (
	"time"

	"backend/internal/trainers/domain"
)

type Date struct {
	Date         time.Time
	HasFreeHours bool
	Hours        []Hour
}

type Hour struct {
	Hour   time.Time
	Status domain.HourStatus
}
