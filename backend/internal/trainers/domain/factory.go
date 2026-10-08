package domain

import (
	"errors"
	"fmt"
	"time"

	common "github.com/quangbach27/golang-common"
)

type HourFactoryConfig struct {
	MinUtcHour               int
	MaxUtcHour               int
	MaxWeeksInTheFutureToSet int
	Location                 *time.Location
}

func (c *HourFactoryConfig) setDefault() {
	c.MinUtcHour = 1
	c.MaxUtcHour = 10
	c.MaxWeeksInTheFutureToSet = 6
}

func (c *HourFactoryConfig) validate() error {
	var errs []error

	if c.MaxWeeksInTheFutureToSet < 1 {
		errs = append(errs, fmt.Errorf(
			"MaxWeeksInTheFutureToSet should be greater than 1, but is %d",
			c.MaxWeeksInTheFutureToSet,
		))
	}
	if c.MinUtcHour < 0 || c.MinUtcHour > 24 {
		errs = append(errs, fmt.Errorf(

			"MinUtcHour should be value between 0 and 24, but is %d",
			c.MinUtcHour,
		))
	}
	if c.MaxUtcHour < 0 || c.MaxUtcHour > 24 {
		errs = append(errs, fmt.Errorf(
			"MinUtcHour should be value between 0 and 24, but is %d",
			c.MaxUtcHour,
		))
	}

	if c.MinUtcHour > c.MaxUtcHour {
		errs = append(errs, fmt.Errorf(
			"MinUtcHour (%d) can't be after MaxUtcHour (%d)",
			c.MinUtcHour, c.MaxUtcHour,
		))
	}

	if len(errs) != 0 {
		return errors.Join(errs...)
	}
	return nil
}

type HourFactory struct {
	config *HourFactoryConfig
}

func NewHourFactory(updateConfigFns ...func(config *HourFactoryConfig)) (*HourFactory, error) {
	config := &HourFactoryConfig{}
	config.setDefault()
	for _, updateFn := range updateConfigFns {
		updateFn(config)
	}
	if err := config.validate(); err != nil {
		return nil, err
	}

	return &HourFactory{
		config: config,
	}, nil
}

func (f *HourFactory) Config() *HourFactoryConfig {
	return f.config
}

func (f *HourFactory) NewAvailableHour(trainerUUID string, h time.Time) (*Hour, error) {
	errDetails := []common.ErrorDetails{}

	if trainerUUID == "" {
		errDetails = append(errDetails, common.ErrorDetails{
			EntityType: "Hour",
			ErrorSlug:  "invalid-trainer-uuid",
			Message:    "trainerUUID can't be empty",
		})
	}

	if err := f.validateTime(h); err != nil {
		errDetails = append(errDetails, common.ErrorDetails{
			EntityType: "Hour",
			ErrorSlug:  "invalid-training-hour",
			Message:    err.Error(),
		})
	}

	if len(errDetails) != 0 {
		return nil, common.NewInvalidInputError("invalid-available-hour", "hour available is not valid").WithDetails(errDetails)
	}

	return &Hour{
		hour:        h,
		status:      Available,
		trainerUUID: trainerUUID,
	}, nil
}

func (f *HourFactory) NewNotAvailableHour(trainerUUID string, h time.Time) (*Hour, error) {
	errDetails := []common.ErrorDetails{}

	if trainerUUID == "" {
		errDetails = append(errDetails, common.ErrorDetails{
			EntityType: "Hour",
			ErrorSlug:  "invalid-trainer-uuid",
			Message:    "trainerUUID can't be empty",
		})
	}

	if err := f.validateTime(h); err != nil {
		errDetails = append(errDetails, common.ErrorDetails{
			EntityType: "Hour",
			ErrorSlug:  "invalid-training-hour",
			Message:    err.Error(),
		})
	}

	if len(errDetails) != 0 {
		return nil, common.NewInvalidInputError("invalid-available-hour", "hour available is not valid").WithDetails(errDetails)
	}

	return &Hour{
		hour:        h,
		status:      NotAvailable,
		trainerUUID: trainerUUID,
	}, nil
}

var (
	ErrNotFullHour = errors.New("hour should be a full hour")
	ErrPastHour    = errors.New("cannot create hour from past")
)

type TooDistantDateError struct {
	MaxWeeksInTheFutureToSet int
	ProvidedDate             time.Time
}

func (e TooDistantDateError) Error() string {
	return fmt.Sprintf(
		"schedule can be only set for next %d weeks, provided date: %s",
		e.MaxWeeksInTheFutureToSet,
		e.ProvidedDate,
	)
}

type TooEarlyHourError struct {
	MinUtcHour   int
	ProvidedTime time.Time
}

func (e TooEarlyHourError) Error() string {
	return fmt.Sprintf(
		"too early hour, min UTC hour: %d, provided time: %s",
		e.MinUtcHour,
		e.ProvidedTime,
	)
}

type TooLateHourError struct {
	MaxUtcHour   int
	ProvidedTime time.Time
}

func (e TooLateHourError) Error() string {
	return fmt.Sprintf(
		"too late hour, min UTC hour: %d, provided time: %s",
		e.MaxUtcHour,
		e.ProvidedTime,
	)
}

func (f *HourFactory) validateTime(hour time.Time) error {
	if !hour.Round(time.Hour).Equal(hour) {
		return ErrNotFullHour
	}

	// AddDate is better than Add for adding days, because not every day have 24h!
	if hour.After(time.Now().AddDate(0, 0, f.config.MaxWeeksInTheFutureToSet*7)) {
		return TooDistantDateError{
			MaxWeeksInTheFutureToSet: f.config.MaxWeeksInTheFutureToSet,
			ProvidedDate:             hour,
		}
	}

	currentHour := time.Now().Truncate(time.Hour)
	if hour.Before(currentHour) || hour.Equal(currentHour) {
		return ErrPastHour
	}
	if hour.UTC().Hour() > f.config.MaxUtcHour {
		return TooLateHourError{
			MaxUtcHour:   f.config.MaxUtcHour,
			ProvidedTime: hour,
		}
	}
	if hour.UTC().Hour() < f.config.MinUtcHour {
		return TooEarlyHourError{
			MinUtcHour:   f.config.MinUtcHour,
			ProvidedTime: hour,
		}
	}

	return nil
}
