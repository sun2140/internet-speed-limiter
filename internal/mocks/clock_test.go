package mocks

import (
	"encoding/json"
	"fmt"
	"goraj/limited-network-driver/pkg/file"
	"log"
	"testing"
	"time"
)

const everyPeriodDataset = "./clock_dataset_every_period.txt"

type everyPeriodTestCase struct {
	CurrentDate string `json:"currentDate"`
	FirstDay    int    `json:"FirstDay"`
}

func everyPeriodDatasetGenerator() file.Generator[[]byte] {
	return EveryPeriodGenerator(func(fakeClock *FakeClock, firstDay int) []byte {
		testCase := everyPeriodTestCase{
			CurrentDate: fakeClock.Now().Format("01/02/2006"),
			FirstDay:    firstDay,
		}
		stringResult, err := json.Marshal(testCase)
		if err != nil {
			panic(err)
		}
		return stringResult
	})
}

func TestFakeClock(t *testing.T) {

	t.Run("Assigned date is retrieved as is", func(t *testing.T) {
		fakeClock := NewFakeClock(2024, time.December, 1)

		expected := time.Date(2024, time.December, 1, 0, 0, 0, 0, time.UTC)

		actual := fakeClock.Now()

		if expected != actual {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	})

	if !file.Exists(everyPeriodDataset) {
		err := file.WriteStringAsLine(everyPeriodDataset, everyPeriodDatasetGenerator())
		if err != nil {
			t.Error(err)
		}
	}

	t.Run("Every period generator work properly", func(t *testing.T) {
		var testCase = everyPeriodTestCase{}
		var actual, expected []everyPeriodTestCase

		for tc := range everyPeriodDatasetGenerator() {
			err := json.Unmarshal(tc, &testCase)
			if err != nil {
				t.Errorf("unmarshal failed")
			}
			actual = append(actual, testCase)
		}

		for tc := range file.ReadJsonLineAsStruct[everyPeriodTestCase](everyPeriodDataset, &testCase) {
			expected = append(expected, *tc)
		}

		actualSize, expectedSize := len(actual), len(expected)

		log.Printf("expected has %v elements, actual got %v elements", expectedSize, actualSize)
		if actualSize != expectedSize {
			t.Errorf("expected has %v elements, actual got %v elements", expectedSize, actualSize)
		}

		for index := 0; index < expectedSize; index++ {
			actualElement := actual[index]
			expectedElement := expected[index]

			t.Run(fmt.Sprintf("date is %v, firstDay is %v", expectedElement.CurrentDate, expectedElement.FirstDay), func(t *testing.T) {
				if expectedElement.CurrentDate != actualElement.CurrentDate {
					t.Errorf("expected %v, got %v", expectedElement.CurrentDate, actualElement.CurrentDate)
				}

				if expectedElement.FirstDay != actualElement.FirstDay {
					t.Errorf("expected %v, got %v", expectedElement.CurrentDate, actualElement.CurrentDate)
				}
			})
		}
	})
}
