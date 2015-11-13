package lpxgen

import (
	"fmt"
	"time"
)

var timeFormats = []string{
	time.RFC3339,
	time.RFC3339Nano,
	time.Kitchen,
	time.Stamp,
}

type Clock interface {
	Now() time.Time
}

type RealClock time.Time

func (c *RealClock) Now() time.Time {
	return time.Now()
}

type MonotonicClock struct {
	Current time.Time
	Step    time.Duration
}

func NewMonotonicClock(start string, step string) (Clock, error) {
	for _, format := range timeFormats {
		if tm, err := time.Parse(format, start); err == nil {
			// set the day to now.
			if format == time.Kitchen {
				tx := time.Now()
				tm = time.Date(tx.Year(), tx.Month(), tx.Day(), tm.Hour(), tm.Minute(), tm.Second(), 0, tm.Location())
			}

			duration, err := time.ParseDuration(step)
			if err == nil {
				return &MonotonicClock{tm, duration}, nil
			} else {
				return nil, err
			}
		}
	}

	return nil, fmt.Errorf("Unable to create a clock from start and step")
}

func (c *MonotonicClock) Now() time.Time {
	old := c.Current

	c.Current = c.Current.Add(c.Step)
	return old
}

var DefaultClock Clock = &RealClock{}
