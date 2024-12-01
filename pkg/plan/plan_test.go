package plan

import (
	"encoding/json"
	"fmt"
	"goraj/limited-network-driver/internal/mocks"
	"goraj/limited-network-driver/pkg/billing"
	"goraj/limited-network-driver/pkg/file"
	"testing"
	"time"
)

const currentThresholdDataset = "./plan_dataset_current_threshold.txt"

type currentThresholdTestCase struct {
	CurrentDate       string `json:"currentDate"`
	FirstDay          int    `json:"firstDay"`
	MonthData         int    `json:"monthData"`
	ExpectedThreshold int    `json:"expectedThreshold"`
}

func CurrentThresholdDatasetGenerator() file.Generator[[]byte] {
	return mocks.EveryPeriodGenerator(func(fakeClock *mocks.FakeClock, firstDay int) []byte {
		myBilling := billing.NewBilling(fakeClock, firstDay)
		billingPeriodDuration := myBilling.GetDaysInCurrentBillingPeriod()
		monthData := billingPeriodDuration
		expectedThreshold := (monthData / billingPeriodDuration) * myBilling.GetBillingPeriodCurrentDay()

		testCase := currentThresholdTestCase{
			CurrentDate:       fakeClock.Now().Format(billing.DateLayout),
			FirstDay:          firstDay,
			MonthData:         monthData,
			ExpectedThreshold: expectedThreshold,
		}

		stringResult, err := json.Marshal(testCase)
		if err != nil {
			panic(err)
		}
		return stringResult
	})
}

func TestPlan(t *testing.T) {

	t.Run("Get current threshold", func(t *testing.T) {

		type CurrentThresholdTestCase struct {
			CurrentDate       string `json:"currentDate"`
			FirstDay          int    `json:"firstDay"`
			MonthData         Mb     `json:"monthData"`
			ExpectedThreshold int    `json:"expectedThreshold"`
		}

		if !file.Exists(currentThresholdDataset) {
			err := file.WriteStringAsLine(currentThresholdDataset, CurrentThresholdDatasetGenerator())
			if err != nil {
				t.Error(err)
			}
		}

		for tc := range file.ReadJsonLineAsStruct(currentThresholdDataset, &CurrentThresholdTestCase{}) {

			t.Run(fmt.Sprintf("Today is %v, the first day is %v, I have %v Mb of data, my current threshold should be %v",
				tc.CurrentDate, tc.FirstDay, tc.MonthData, tc.ExpectedThreshold), func(t *testing.T) {

				date, err := time.Parse(billing.DateLayout, tc.CurrentDate)
				if err != nil {
					t.Errorf("for date %v, parsing error: %v", tc.CurrentDate, err)
				}

				fakeClock := mocks.NewFakeClock(date.Year(), date.Month(), date.Day())
				myBilling := billing.NewBilling(fakeClock, tc.FirstDay)

				myPlan := NewPlan(tc.MonthData, 50)

				actualThreshold := myPlan.getCurrentThreshold(*myBilling)
				expectedThreshold := tc.ExpectedThreshold

				if actualThreshold != expectedThreshold {
					t.Errorf("for date %v, expected: %v, got: %v", tc.CurrentDate, expectedThreshold, actualThreshold)
				}
			})
		}
	})
}
