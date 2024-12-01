package billing

import (
	"goraj/limited-network-driver/pkg/clock"
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
