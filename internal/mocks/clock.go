package mocks

import (
	"time"
)

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

func EveryPeriodGenerator(buildFn func(fakeClock *FakeClock, firstDay int) []byte) func(yield func([]byte) bool) {
	return func(yield func([]byte) bool) {
		for year := 2023; year <= 2024; year++ {
			for month := 1; month <= 12; month++ {
				daysInMonth := time.Date(year, time.Month((month+1)%12), 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1).Day()
				for day := 1; day <= daysInMonth; day++ {
					for firstDay := 1; firstDay <= 31; firstDay++ {
						fakeClock := NewFakeClock(year, time.Month(month), day)
						if !yield(buildFn(fakeClock, firstDay)) {
							return
						}
					}
				}
			}
		}
	}
}
