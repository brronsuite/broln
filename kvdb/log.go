package kvdb

import (
	"github.com/brronsuite/broln/kvdb/postgres"
	"github.com/brronsuite/bronlog"
)

// log is a logger that is initialized as disabled.  This means the package will
// not perform any logging by default until a logger is set.
var log = bronlog.Disabled

// UseLogger uses a specified Logger to output package logging info.
func UseLogger(logger bronlog.Logger) {
	log = logger

	postgres.UseLogger(log)
}
