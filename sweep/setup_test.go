package sweep

import (
	"testing"

	"github.com/brronsuite/broln/kvdb"
)

func TestMain(m *testing.M) {
	kvdb.RunTests(m)
}
