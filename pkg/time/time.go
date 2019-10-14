package time

import (
	"time"
)

// GetLocalTime return local time same as the machine time
func GetLocalTime() (finalDate time.Time, err error) {
	date := time.Now()
	newTmpDate := date.Format(time.RFC3339)

	finalDate, err = time.Parse(time.RFC3339, newTmpDate)
	if err != nil {
		return
	}

	return
}

// GetUTCTime get time in utc formar utc format
func GetUTCTime() (finalDate time.Time, err error) {
	date := time.Now()
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		return
	}

	finalDate, err = ConvertTimeWithTimeStamp(date.In(loc))
	if err != nil {
		return
	}

	return
}

// ConvertToUTCTime return date with utc timezone
func ConvertToUTCTime(date time.Time) (finalDate time.Time, err error) {
	loc, err := time.LoadLocation("UTC")
	if err != nil {
		return
	}
	finalDate, err = ConvertTimeToDifferentTimezone(date, loc)
	if err != nil {
		return
	}
	return
}

// ConvertTimeToDifferentTimezone convert time to same as location
func ConvertTimeToDifferentTimezone(date time.Time, loc *time.Location) (finalDate time.Time, err error) {
	_, offset := date.Zone()

	tmpDate := time.Date(date.Year(), date.Month(), date.Day(), date.Hour(), date.Minute(), date.Second(), date.Nanosecond(), loc)
	finalDate, err = ConvertTimeWithTimeStamp(tmpDate)
	if err != nil {
		return
	}

	if offset > 0 {
		finalDate = finalDate.Add(time.Second * time.Duration(offset) * -1)
	} else if offset < 0 {
		finalDate = finalDate.Add(time.Second * time.Duration(offset))
	}

	return
}

// ConvertTimeWithTimeStamp will convert time format with time offset
func ConvertTimeWithTimeStamp(date time.Time) (finalDate time.Time, err error) {
	newTmpDate := date.Format(time.RFC3339)

	finalDate, err = time.Parse(time.RFC3339, newTmpDate)
	if err != nil {
		return
	}
	return
}
