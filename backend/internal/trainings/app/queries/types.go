package queries

import "time"

type Training struct {
	UUID               string
	Hour               time.Time
	Notes              string
	AttendeeUUID       string
	AttendeeUsername   string
	TrainerUUID        string
	TrainerUsername    string
	ProposedNewTime    *time.Time
	ProposedBy         *string
	ProposedByUsername *string
	Canceled           bool
}
