package billing

import (
	"encoding/json"
	"fmt"
	"goraj/limited-network-driver/internal/mocks"
	"goraj/limited-network-driver/pkg/clock"
	"goraj/limited-network-driver/pkg/file"
	"testing"
	"time"
)

const periodCurrentDayDataset = "./billing_dataset_current_day.txt"

type periodCurrentDayTestCase struct {
	CurrentDate               string `json:"currentDate"`
	FirstDay                  int    `json:"FirstDay"`
	CurrentDayOfBillingPeriod int    `json:"currentDayOfBillingPeriod"`
}

func periodCurrentDayDatasetGenerator() file.Generator[[]byte] {
	return EveryPeriodGenerator(func(currentDate time.Time, myBilling *Billing) []byte {
		currentDayOfBillingPeriod := myBilling.GetBillingPeriodCurrentDay()

		testCase := periodCurrentDayTestCase{
			CurrentDate:               currentDate.Format(DateLayout),
			FirstDay:                  myBilling.FirstDay,
			CurrentDayOfBillingPeriod: currentDayOfBillingPeriod,
		}

		stringResult, err := json.Marshal(testCase)
		if err != nil {
			panic(err)
		}
		return stringResult
	})
}

func TestBillingPeriod(t *testing.T) {

	t.Run("Get number of days in any month", func(t *testing.T) {
		tests := []struct {
			year             int
			month            time.Month
			pastMonthDays    int
			currentMonthDays int
			nextMonthDays    int
		}{
			{pastMonthDays: 31, currentMonthDays: 31, nextMonthDays: 28, year: 2023, month: time.January},
			{pastMonthDays: 31, currentMonthDays: 28, nextMonthDays: 31, year: 2023, month: time.February},
			{pastMonthDays: 28, currentMonthDays: 31, nextMonthDays: 30, year: 2023, month: time.March},
			{pastMonthDays: 31, currentMonthDays: 30, nextMonthDays: 31, year: 2023, month: time.April},
			{pastMonthDays: 30, currentMonthDays: 31, nextMonthDays: 30, year: 2023, month: time.May},
			{pastMonthDays: 31, currentMonthDays: 30, nextMonthDays: 31, year: 2023, month: time.June},
			{pastMonthDays: 30, currentMonthDays: 31, nextMonthDays: 31, year: 2023, month: time.July},
			{pastMonthDays: 31, currentMonthDays: 31, nextMonthDays: 30, year: 2023, month: time.August},
			{pastMonthDays: 31, currentMonthDays: 30, nextMonthDays: 31, year: 2023, month: time.September},
			{pastMonthDays: 30, currentMonthDays: 31, nextMonthDays: 30, year: 2023, month: time.October},
			{pastMonthDays: 31, currentMonthDays: 30, nextMonthDays: 31, year: 2023, month: time.November},
			{pastMonthDays: 30, currentMonthDays: 31, nextMonthDays: 31, year: 2023, month: time.December},
			{pastMonthDays: 31, currentMonthDays: 29, nextMonthDays: 31, year: 2024, month: time.February},
		}

		for _, tc := range tests {
			t.Run(fmt.Sprintf("In %v %v the expected number of days is computed correctly", tc.month, tc.year), func(t *testing.T) {
				fakeClock := mocks.NewFakeClock(tc.year, tc.month, 1)
				date := fakeClock.Now().Format(DateLayout)
				myBilling := Billing{clock: fakeClock}

				daysInPastMonth := myBilling.getDaysInPastMonth()
				if daysInPastMonth != tc.pastMonthDays {
					t.Errorf("For date %v, expected: %v, got: %v", date, tc.pastMonthDays, daysInPastMonth)
				}

				daysInMonth := myBilling.getDaysInMonth()
				if daysInMonth != tc.currentMonthDays {
					t.Errorf("For date %v, expected: %v, got: %v", date, tc.currentMonthDays, daysInMonth)
				}

				daysInNextMonth := myBilling.getDaysInNextMonth()
				if daysInNextMonth != tc.nextMonthDays {
					t.Errorf("For date %v, expected: %v, got: %v", date, tc.nextMonthDays, daysInNextMonth)
				}

			})
		}
	})

	t.Run("Get the current day in the current month", func(t *testing.T) {
		expectedDay := 10
		fakeClock := mocks.NewFakeClock(1994, time.May, expectedDay)
		date := fakeClock.Now().Format(DateLayout)
		myBilling := Billing{clock: fakeClock}

		actualDay := myBilling.getDayOfMonth()

		if actualDay != 10 {
			t.Errorf("For date %v, expected: %v, got: %v", date, expectedDay, actualDay)
		}
	})

	t.Run("Get billing period current day", func(t *testing.T) {

		t.Run("Today with first day 1 should return day of the month", func(t *testing.T) {
			realClock := clock.NewRealClock()
			myBilling := NewBilling(realClock, 1)

			today := realClock.Now().Day()
			billingPeriodCurrentDay := myBilling.GetBillingPeriodCurrentDay()

			if today != billingPeriodCurrentDay {
				t.Errorf("for date %v, expected: %v, got: %v", today, today, billingPeriodCurrentDay)
			}
		})

		if !file.Exists(periodCurrentDayDataset) {
			err := file.WriteStringAsLine(periodCurrentDayDataset, periodCurrentDayDatasetGenerator())
			if err != nil {
				t.Error(err)
			}
		}

		for tc := range file.ReadJsonLineAsStruct(periodCurrentDayDataset, &periodCurrentDayTestCase{}) {

			t.Run(fmt.Sprintf("%v period first day %v => current period day %v", tc.CurrentDate, tc.FirstDay, tc.CurrentDayOfBillingPeriod), func(t *testing.T) {
				date, err := time.Parse(DateLayout, tc.CurrentDate)
				if err != nil {
					t.Errorf("for date %v, parsing error: %v", tc.CurrentDate, err)
				}

				fakeClock := mocks.NewFakeClock(date.Year(), date.Month(), date.Day())
				myBilling := NewBilling(fakeClock, tc.FirstDay)

				actualDay := myBilling.GetBillingPeriodCurrentDay()
				expectedDay := tc.CurrentDayOfBillingPeriod

				if actualDay != expectedDay {
					t.Errorf("for date %v, expected: %v, got: %v", tc.CurrentDate, expectedDay, actualDay)
				}
			})

		}
	})

	t.Run("Get number of days in the current billing period", func(t *testing.T) {

		t.Run(
			"The period started last month"+
				"and the first day exists in the last month"+
				"and the first day exists in the current month", func(t *testing.T) {
				fakeClock := mocks.NewFakeClock(2023, time.March, 14)
				myBilling := NewBilling(fakeClock, 15)

				expected := 28
				actual := myBilling.GetDaysInCurrentBillingPeriod()

				if actual != expected {
					t.Errorf("expected %v, got %v", expected, actual)
				}
			})

		t.Run(
			"The period started last month"+
				"and the first day exists in the last month"+
				"and the first day does not exists in the current month", func(t *testing.T) {
				fakeClock := mocks.NewFakeClock(2023, time.June, 14)
				myBilling := NewBilling(fakeClock, 31)

				expected := 31
				actual := myBilling.GetDaysInCurrentBillingPeriod()

				if actual != expected {
					t.Errorf("expected %v, got %v", expected, actual)
				}
			})

		t.Run(
			"The period started last month"+
				"and the first day does not exists in the last month"+
				"and the first day exists in the current month", func(t *testing.T) {
				fakeClock := mocks.NewFakeClock(2023, time.July, 14)
				myBilling := NewBilling(fakeClock, 31)

				expected := 30
				actual := myBilling.GetDaysInCurrentBillingPeriod()

				if actual != expected {
					t.Errorf("expected %v, got %v", expected, actual)
				}
			})

		t.Run(
			"The period started this month"+
				"and the first day exists in the next month", func(t *testing.T) {
				fakeClock := mocks.NewFakeClock(2023, time.February, 14)
				myBilling := NewBilling(fakeClock, 10)

				expected := 28
				actual := myBilling.GetDaysInCurrentBillingPeriod()

				if actual != expected {
					t.Errorf("expected %v, got %v", expected, actual)
				}
			})

		t.Run(
			"The period started this month"+
				"and the first day does not exists in the next month", func(t *testing.T) {
				fakeClock := mocks.NewFakeClock(2023, time.January, 30)
				myBilling := NewBilling(fakeClock, 30)

				expected := 30
				actual := myBilling.GetDaysInCurrentBillingPeriod()

				if actual != expected {
					t.Errorf("expected %v, got %v", expected, actual)
				}
			})

	})
}
