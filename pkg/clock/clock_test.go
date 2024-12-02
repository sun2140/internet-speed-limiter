package clock

import (
	"testing"
	"time"
)

func TestClock(t *testing.T) {

	t.Run("The clock return the current time", func(t *testing.T) {
		realClock := NewRealClock()

		actual := realClock.Now().UnixNano()
		expected := time.Now().UnixNano()

		diff := expected - actual

		if diff > 1000 {
			t.Errorf("Difference between timestamp is too big: %v", diff)
		}

	})

}
