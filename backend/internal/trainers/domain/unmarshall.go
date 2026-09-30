package domain

import "time"

// UnmarshalHour is only used in term of loading hour from database. We trust the data stored in database
func UnmarshalHour(
	h time.Time,
	status HourStatus,
	trainerUUID string,
) *Hour {
	return &Hour{
		hour:        h,
		status:      status,
		trainerUUID: trainerUUID,
	}
}
