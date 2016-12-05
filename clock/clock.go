package clock

import "time"

type Time interface {
	Now() time.Time
	After(time.Duration) <-chan time.Time
}

type SystemTime struct{}

func (st SystemTime) Now() time.Time {
	return time.Now()
}

func (st SystemTime) After(d time.Duration) <-chan time.Time {
	return time.After(d)
}
