package autopilot_test

import (
	"testing"

	"github.com/brronsuite/bronutil"
	"github.com/brronsuite/broln/autopilot"
)

// TestMedian tests the Median method.
func TestMedian(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		values []bronutil.Amount
		median bronutil.Amount
	}{
		{
			values: []bronutil.Amount{},
			median: 0,
		},
		{
			values: []bronutil.Amount{10},
			median: 10,
		},
		{
			values: []bronutil.Amount{10, 20},
			median: 15,
		},
		{
			values: []bronutil.Amount{10, 20, 30},
			median: 20,
		},
		{
			values: []bronutil.Amount{30, 10, 20},
			median: 20,
		},
		{
			values: []bronutil.Amount{10, 10, 10, 10, 5000000},
			median: 10,
		},
	}

	for _, test := range testCases {
		res := autopilot.Median(test.values)
		if res != test.median {
			t.Fatalf("expected median %v, got %v", test.median, res)
		}
	}
}
