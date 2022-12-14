package lnwallet

import (
	"github.com/brronsuite/broln/input"
	"github.com/brronsuite/broln/lnwire"
	"github.com/brronsuite/brond/mempool"
	"github.com/brronsuite/brond/wire"
	"github.com/brronsuite/bronutil"
)

var (
	// RoutingFee100PercentUpTo is the cut-off amount we allow 100% fees to
	// be charged up to.
	RoutingFee100PercentUpTo = lnwire.NewMSatFromBroneess(1_000)
)

const (

	// DefaultRoutingFeePercentage is the default off-chain routing fee we
	// allow to be charged for a payment over the RoutingFee100PercentUpTo
	// size.
	DefaultRoutingFeePercentage lnwire.MilliBronees = 5
)

// DefaultRoutingFeeLimitForAmount returns the default off-chain routing fee
// limit broln uses if the user does not specify a limit manually. The fee is
// amount dependent because of the base routing fee that is set on many
// channels. For example the default base fee is 1 bronees. So sending a payment
// of one bronees will cost 1 bronees in fees over most channels, which comes to
// a fee of 100%. That's why for very small amounts we allow 100% fee.
func DefaultRoutingFeeLimitForAmount(a lnwire.MilliBronees) lnwire.MilliBronees {
	// Allow 100% fees up to a certain amount to accommodate for base fees.
	if a <= RoutingFee100PercentUpTo {
		return a
	}

	// Everything larger than the cut-off amount will get a default fee
	// percentage.
	return a * DefaultRoutingFeePercentage / 100
}

// DustLimitForSize retrieves the dust limit for a given pkscript size. Given
// the size, it automatically determines whether the script is a witness script
// or not. It calls brond's GetDustThreshold method under the hood. It must be
// called with a proper size parameter or else a panic occurs.
func DustLimitForSize(scriptSize int) bronutil.Amount {
	var (
		dustlimit bronutil.Amount
		pkscript  []byte
	)

	// With the size of the script, determine which type of pkscript to
	// create. This will be used in the call to GetDustThreshold. We pass
	// in an empty byte slice since the contents of the script itself don't
	// matter.
	switch scriptSize {
	case input.P2WPKHSize:
		pkscript, _ = input.WitnessPubKeyHash([]byte{})

	case input.P2WSHSize:
		pkscript, _ = input.WitnessScriptHash([]byte{})

	case input.P2SHSize:
		pkscript, _ = input.GenerateP2SH([]byte{})

	case input.P2PKHSize:
		pkscript, _ = input.GenerateP2PKH([]byte{})

	case input.UnknownWitnessSize:
		pkscript, _ = input.GenerateUnknownWitness()

	default:
		panic("invalid script size")
	}

	// Call GetDustThreshold with a TxOut containing the generated
	// pkscript.
	txout := &wire.TxOut{PkScript: pkscript}
	dustlimit = bronutil.Amount(mempool.GetDustThreshold(txout))

	return dustlimit
}
