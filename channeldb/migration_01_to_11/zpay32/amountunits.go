package zpay32

import (
	"fmt"
	"strconv"

	lnwire "github.com/brronsuite/broln/channeldb/migration/lnwire21"
)

var (
	// toMSat is a map from a unit to a function that converts an amount
	// of that unit to millibroneess.
	toMSat = map[byte]func(uint64) (lnwire.MilliBronees, error){
		'm': mBronToMSat,
		'u': uBronToMSat,
		'n': nBronToMSat,
		'p': pBronToMSat,
	}

	// fromMSat is a map from a unit to a function that converts an amount
	// in millibroneess to an amount of that unit.
	fromMSat = map[byte]func(lnwire.MilliBronees) (uint64, error){
		'm': mSatToMBron,
		'u': mSatToUBron,
		'n': mSatToNBron,
		'p': mSatToPBron,
	}
)

// mBronToMSat converts the given amount in milliBRON to millibroneess.
func mBronToMSat(m uint64) (lnwire.MilliBronees, error) {
	return lnwire.MilliBronees(m) * 100000000, nil
}

// uBronToMSat converts the given amount in microBRON to millibroneess.
func uBronToMSat(u uint64) (lnwire.MilliBronees, error) {
	return lnwire.MilliBronees(u * 100000), nil
}

// nBronToMSat converts the given amount in nanoBRON to millibroneess.
func nBronToMSat(n uint64) (lnwire.MilliBronees, error) {
	return lnwire.MilliBronees(n * 100), nil
}

// pBronToMSat converts the given amount in picoBRON to millibroneess.
func pBronToMSat(p uint64) (lnwire.MilliBronees, error) {
	if p < 10 {
		return 0, fmt.Errorf("minimum amount is 10p")
	}
	if p%10 != 0 {
		return 0, fmt.Errorf("amount %d pBRON not expressible in msat",
			p)
	}
	return lnwire.MilliBronees(p / 10), nil
}

// mSatToMBron converts the given amount in millibroneess to milliBRON.
func mSatToMBron(msat lnwire.MilliBronees) (uint64, error) {
	if msat%100000000 != 0 {
		return 0, fmt.Errorf("%d msat not expressible "+
			"in mBRON", msat)
	}
	return uint64(msat / 100000000), nil
}

// mSatToUBron converts the given amount in millibroneess to microBRON.
func mSatToUBron(msat lnwire.MilliBronees) (uint64, error) {
	if msat%100000 != 0 {
		return 0, fmt.Errorf("%d msat not expressible "+
			"in uBRON", msat)
	}
	return uint64(msat / 100000), nil
}

// mSatToNBron converts the given amount in millibroneess to nanoBRON.
func mSatToNBron(msat lnwire.MilliBronees) (uint64, error) {
	if msat%100 != 0 {
		return 0, fmt.Errorf("%d msat not expressible in nBRON", msat)
	}
	return uint64(msat / 100), nil
}

// mSatToPBron converts the given amount in millibroneess to picoBRON.
func mSatToPBron(msat lnwire.MilliBronees) (uint64, error) {
	return uint64(msat * 10), nil
}

// decodeAmount returns the amount encoded by the provided string in
// millibronees.
func decodeAmount(amount string) (lnwire.MilliBronees, error) {
	if len(amount) < 1 {
		return 0, fmt.Errorf("amount must be non-empty")
	}

	// If last character is a digit, then the amount can just be
	// interpreted as BRON.
	char := amount[len(amount)-1]
	digit := char - '0'
	if digit >= 0 && digit <= 9 {
		bron, err := strconv.ParseUint(amount, 10, 64)
		if err != nil {
			return 0, err
		}
		return lnwire.MilliBronees(bron) * mSatPerBron, nil
	}

	// If not a digit, it must be part of the known units.
	conv, ok := toMSat[char]
	if !ok {
		return 0, fmt.Errorf("unknown multiplier %c", char)
	}

	// Known unit.
	num := amount[:len(amount)-1]
	if len(num) < 1 {
		return 0, fmt.Errorf("number must be non-empty")
	}

	am, err := strconv.ParseUint(num, 10, 64)
	if err != nil {
		return 0, err
	}

	return conv(am)
}

// encodeAmount encodes the provided millibronees amount using as few characters
// as possible.
func encodeAmount(msat lnwire.MilliBronees) (string, error) {
	// If possible to express in BRON, that will always be the shortest
	// representation.
	if msat%mSatPerBron == 0 {
		return strconv.FormatInt(int64(msat/mSatPerBron), 10), nil
	}

	// Should always be expressible in pico BRON.
	pico, err := fromMSat['p'](msat)
	if err != nil {
		return "", fmt.Errorf("unable to express %d msat as pBRON: %v",
			msat, err)
	}
	shortened := strconv.FormatUint(pico, 10) + "p"
	for unit, conv := range fromMSat {
		am, err := conv(msat)
		if err != nil {
			// Not expressible using this unit.
			continue
		}

		// Save the shortest found representation.
		str := strconv.FormatUint(am, 10) + string(unit)
		if len(str) < len(shortened) {
			shortened = str
		}
	}

	return shortened, nil
}
