//go:build dev
// +build dev

package brocoind_test

import (
	"testing"

	chainntnfstest "github.com/brronsuite/broln/chainntnfs/test"
)

// TestInterfaces executes the generic notifier test suite against a brocoind
// powered chain notifier.
func TestInterfaces(t *testing.T) {
	chainntnfstest.TestInterfaces(t, "brocoind")
}
