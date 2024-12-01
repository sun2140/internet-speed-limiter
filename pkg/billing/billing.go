package billing

import (
	"goraj/limited-network-driver/internal/mocks"
	"goraj/limited-network-driver/pkg/clock"
	"goraj/limited-network-driver/pkg/file"
	"time"
)

const DateLayout = "01/02/2006"

type Billing struct {
	clock    clock.Clock
	FirstDay int
}

func NewBilling(clock clock.Clock, firstDay int) *Billing {
	return &Billing{
		clock:    clock,
		FirstDay: firstDay,
	}
}

func (billing *Billing) getDayOfMonth() int {
	return billing.clock.Now().Day()
}

func (billing *Billing) getDaysInMonth() int {
	date := billing.clock.Now()
	nextMonth := date.AddDate(0, 1, -date.Day()+1)
	currentMonthLastDay := nextMonth.AddDate(0, 0, -1)
	return currentMonthLastDay.Day()
}

func (billing *Billing) getDaysInPastMonth() int {
	date := billing.clock.Now()
	lastDay := date.AddDate(0, 0, -date.Day())
	return lastDay.Day()
}

func (billing *Billing) getDaysInNextMonth() int {
	date := billing.clock.Now()
	nextNextMonth := date.AddDate(0, 2, -date.Day()+1)
	nextMonthLastDay := nextNextMonth.AddDate(0, 0, -1)
	return nextMonthLastDay.Day()
}

func (billing *Billing) GetDaysInCurrentBillingPeriod() int {
	firstDay := billing.FirstDay
	currentDay := billing.getDayOfMonth()

	daysInPastMonth := billing.getDaysInPastMonth()
	daysInCurrentMonth := billing.getDaysInMonth()
	daysInNextMonth := billing.getDaysInNextMonth()

	periodDaysPastMonth := 0
	periodDaysCurrentMonth := 0
	periodDaysNextMonth := 0

	periodStartedLastMonth := firstDay > currentDay

	if periodStartedLastMonth {
		// From first day (included) to end of the last month
		if firstDay <= daysInPastMonth {
			periodDaysPastMonth = daysInPastMonth - firstDay + 1
		}

		// From begin of the month to first day (excluded)
		if firstDay > daysInCurrentMonth {
			periodDaysCurrentMonth = daysInCurrentMonth
		} else {
			periodDaysCurrentMonth = firstDay - 1
		}
	} else { // period started this month and end next month

		// From first day (included) to the end of the month
		periodDaysCurrentMonth = daysInCurrentMonth - firstDay + 1

		// From begin of the next month  to first day (excluded)
		if firstDay > daysInNextMonth {
			periodDaysNextMonth = daysInNextMonth
		} else {
			periodDaysNextMonth = firstDay - 1
		}
	}

	daysInCurrentBillingPeriod := periodDaysPastMonth + periodDaysCurrentMonth + periodDaysNextMonth
	return daysInCurrentBillingPeriod
}

func (billing *Billing) GetBillingPeriodCurrentDay() int {
	firstDay := billing.FirstDay
	currentDay := billing.getDayOfMonth()

	daysInPastMonth := billing.getDaysInPastMonth()

	currentPeriodStartedThisMonth := currentDay >= firstDay

	if currentPeriodStartedThisMonth {
		return currentDay - firstDay + 1
	}

	// currentPeriodStartedLastMonth
	periodFirstDayExistInPastMonth := firstDay <= daysInPastMonth
	if periodFirstDayExistInPastMonth {
		return daysInPastMonth - firstDay + 1 + currentDay
	}

	// last month is too short
	return currentDay
}

type PeriodTestCaseBuildFn func(time time.Time, billing *Billing) []byte

func EveryPeriodGenerator(buildFn PeriodTestCaseBuildFn) file.Generator[[]byte] {
	return func(yield func([]byte) bool) {
		for year := 2023; year <= 2024; year++ {
			for month := 1; month <= 12; month++ {
				daysInMonth := time.Date(year, time.Month((month+1)%12), 1, 0, 0, 0, 0, time.UTC).AddDate(0, 0, -1).Day()
				for day := 1; day <= daysInMonth; day++ {
					for firstDay := 1; firstDay <= 31; firstDay++ {
						currentDate := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)

						myClock := mocks.NewFakeClock(year, time.Month(month), day)
						myBilling := NewBilling(myClock, firstDay)

						if !yield(buildFn(currentDate, myBilling)) {
							return
						}
					}
				}
			}
		}
	}
}
