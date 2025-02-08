package clock

import "time"

type Clock interface {
	Now() time.Time
}

type clock struct {
}

func NewClock() Clock {
	return &clock{}
}

func (s *clock) Now() time.Time {
	return time.Now().UTC()
}
