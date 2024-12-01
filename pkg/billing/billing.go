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

func (billing *Billing) GetBillingPeriodCurrentDay() int {
	firstDay := billing.FirstDay
	currentDay := billing.getDayOfMonth()

	if hasPeriodStartedThisMonth(currentDay, firstDay) {
		return currentDay - firstDay + 1
	}

	daysInPastMonth := billing.getDaysInPastMonth()
	if firstDayExistInMonth(firstDay, daysInPastMonth) {
		return daysInPastMonth - firstDay + 1 + currentDay
	}

	return currentDay
}

func (billing *Billing) GetDaysInCurrentBillingPeriod() int {
	firstDay := billing.FirstDay
	currentDay := billing.getDayOfMonth()

	periodDaysCurrentMonth := computeCurrentMonthPeriodDays(currentDay, firstDay, billing.getDaysInMonth())

	if hasPeriodStartedPastMonth(currentDay, firstDay) {
		return computePastMonthPeriodDays(currentDay, firstDay, billing.getDaysInPastMonth()) + periodDaysCurrentMonth
	}
	return computeNextMonthPeriodDays(currentDay, firstDay, billing.getDaysInNextMonth()) + periodDaysCurrentMonth
}

func computeCurrentMonthPeriodDays(currentDay, firstDay, daysInCurrentMonth int) int {
	if hasPeriodStartedPastMonth(currentDay, firstDay) {
		if firstDayExistInMonth(firstDay, daysInCurrentMonth) {
			return firstDay - 1
		}
		return daysInCurrentMonth
	}
	return daysInCurrentMonth - firstDay + 1

}

func computePastMonthPeriodDays(currentDay, firstDay, daysInPastMonth int) int {
	if hasPeriodStartedPastMonth(currentDay, firstDay) && firstDayExistInMonth(firstDay, daysInPastMonth) {
		return daysInPastMonth - firstDay + 1
	}
	return 0
}

func computeNextMonthPeriodDays(currentDay, firstDay, daysInNextMonth int) int {
	if hasPeriodStartedPastMonth(currentDay, firstDay) {
		return 0
	}
	if firstDay > daysInNextMonth {
		return daysInNextMonth
	}
	return firstDay - 1
}

func hasPeriodStartedThisMonth(currentDay, firstDay int) bool {
	return !hasPeriodStartedPastMonth(currentDay, firstDay)
}

func hasPeriodStartedPastMonth(currentDay, firstDay int) bool {
	return firstDay > currentDay
}

func firstDayExistInMonth(firstDay, monthDurationInDays int) bool {
	return firstDay <= monthDurationInDays
}
