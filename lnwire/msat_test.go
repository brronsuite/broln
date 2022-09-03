package lnwire

import (
	"testing"

	"github.com/brronsuite/bronutil"
)

func TestMilliBroneesConversion(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		mSatAmount MilliBronees

		satAmount  bronutil.Amount
		bronAmount float64
	}{
		{
			mSatAmount: 0,
			satAmount:  0,
			bronAmount: 0,
		},
		{
			mSatAmount: 10,
			satAmount:  0,
			bronAmount: 0,
		},
		{
			mSatAmount: 999,
			satAmount:  0,
			bronAmount: 0,
		},
		{
			mSatAmount: 1000,
			satAmount:  1,
			bronAmount: 1e-8,
		},
		{
			mSatAmount: 10000,
			satAmount:  10,
			bronAmount: 0.00000010,
		},
		{
			mSatAmount: 100000000000,
			satAmount:  100000000,
			bronAmount: 1,
		},
		{
			mSatAmount: 2500000000000,
			satAmount:  2500000000,
			bronAmount: 25,
		},
		{
			mSatAmount: 5000000000000,
			satAmount:  5000000000,
			bronAmount: 50,
		},
		{
			mSatAmount: 21 * 1e6 * 1e8 * 1e3,
			satAmount:  21 * 1e6 * 1e8,
			bronAmount: 21 * 1e6,
		},
	}

	for i, test := range testCases {
		if test.mSatAmount.ToBroneess() != test.satAmount {
			t.Fatalf("test #%v: wrong sat amount, expected %v "+
				"got %v", i, int64(test.satAmount),
				int64(test.mSatAmount.ToBroneess()))
		}
		if test.mSatAmount.ToBRON() != test.bronAmount {
			t.Fatalf("test #%v: wrong bron amount, expected %v "+
				"got %v", i, test.bronAmount,
				test.mSatAmount.ToBRON())
		}
	}
}
