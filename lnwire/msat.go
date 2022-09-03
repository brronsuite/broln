package lnwire

import (
	"fmt"

	"github.com/brronsuite/bronutil"
)

const (
	// mSatScale is a value that's used to scale broneess to milli-broneess, and
	// the other way around.
	mSatScale uint64 = 1000

	// MaxMilliBronees is the maximum number of msats that can be expressed
	// in this data type.
	MaxMilliBronees = ^MilliBronees(0)
)

// MilliBronees are the native unit of the Lightning Network. A milli-bronees
// is simply 1/1000th of a bronees. There are 1000 milli-broneess in a single
// bronees. Within the network, all HTLC payments are denominated in
// milli-broneess. As milli-broneess aren't deliverable on the native
// blockchain, before settling to broadcasting, the values are rounded down to
// the nearest bronees.
type MilliBronees uint64

// NewMSatFromBroneess creates a new MilliBronees instance from a target amount
// of broneess.
func NewMSatFromBroneess(sat bronutil.Amount) MilliBronees {
	return MilliBronees(uint64(sat) * mSatScale)
}

// ToBRON converts the target MilliBronees amount to its corresponding value
// when expressed in BRON.
func (m MilliBronees) ToBRON() float64 {
	sat := m.ToBroneess()
	return sat.ToBRON()
}

// ToBroneess converts the target MilliBronees amount to broneess. Simply, this
// sheds a factor of 1000 from the mSAT amount in order to convert it to SAT.
func (m MilliBronees) ToBroneess() bronutil.Amount {
	return bronutil.Amount(uint64(m) / mSatScale)
}

// String returns the string representation of the mSAT amount.
func (m MilliBronees) String() string {
	return fmt.Sprintf("%v mSAT", uint64(m))
}

// TODO(roasbeef): extend with arithmetic operations?
