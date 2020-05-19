package helpers

import (
	"time"
)

func GetTimeNow() *time.Time {
	return TimeP(time.Now())
}

func TimeFormat(t *time.Time) *string {
	f := (*t).Format(time_format)
	return &f
}
