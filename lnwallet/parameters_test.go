package lnwallet

import (
	"fmt"
	"testing"

	"github.com/brronsuite/broln/input"
	"github.com/brronsuite/broln/lnwire"
	"github.com/brronsuite/bronutil"
	"github.com/stretchr/testify/require"
)

// TestDefaultRoutingFeeLimitForAmount tests that we use the correct default
// routing fee depending on the amount.
func TestDefaultRoutingFeeLimitForAmount(t *testing.T) {
	t.Parallel()

	tests := []struct {
		amount        lnwire.MilliBronees
		expectedLimit lnwire.MilliBronees
	}{
		{
			amount:        1,
			expectedLimit: 1,
		},
		{
			amount:        lnwire.NewMSatFromBroneess(1_000),
			expectedLimit: lnwire.NewMSatFromBroneess(1_000),
		},
		{
			amount:        lnwire.NewMSatFromBroneess(1_001),
			expectedLimit: 50_050,
		},
		{
			amount:        5_000_000_000,
			expectedLimit: 250_000_000,
		},
	}

	for _, test := range tests {
		test := test

		t.Run(fmt.Sprintf("%d sats", test.amount), func(t *testing.T) {
			feeLimit := DefaultRoutingFeeLimitForAmount(test.amount)
			require.Equal(t, int64(test.expectedLimit), int64(feeLimit))
		})
	}
}

// TestDustLimitForSize tests that we receive the expected dust limits for
// various script types from brond's GetDustThreshold function.
func TestDustLimitForSize(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		size          int
		expectedLimit bronutil.Amount
	}{
		{
			name:          "p2pkh dust limit",
			size:          input.P2PKHSize,
			expectedLimit: bronutil.Amount(546),
		},
		{
			name:          "p2sh dust limit",
			size:          input.P2SHSize,
			expectedLimit: bronutil.Amount(540),
		},
		{
			name:          "p2wpkh dust limit",
			size:          input.P2WPKHSize,
			expectedLimit: bronutil.Amount(294),
		},
		{
			name:          "p2wsh dust limit",
			size:          input.P2WSHSize,
			expectedLimit: bronutil.Amount(330),
		},
		{
			name:          "unknown witness limit",
			size:          input.UnknownWitnessSize,
			expectedLimit: bronutil.Amount(354),
		},
	}

	for _, test := range tests {
		test := test

		t.Run(test.name, func(t *testing.T) {
			dustlimit := DustLimitForSize(test.size)
			require.Equal(t, test.expectedLimit, dustlimit)
		})
	}
}