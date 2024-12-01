package mocks

import "time"

type FakeClock struct {
	year  int
	month time.Month
	day   int
}

func NewFakeClock(year int, month time.Month, day int) *FakeClock {
	return &FakeClock{
		year:  year,
		month: month,
		day:   day,
	}
}

func (clock FakeClock) Now() time.Time {
	return time.Date(clock.year, clock.month, clock.day, 0, 0, 0, 0, time.UTC)
}
