package chanfunding

import (
	"github.com/brronsuite/broln/build"
	"github.com/brronsuite/bronlog"
)

// log is a logger that is initialized with no output filters.  This
// means the package will not perform any logging by default until the caller
// requests it.
var log bronlog.Logger

// The default amount of logging is none.
func init() {
	UseLogger(build.NewSubLogger("CHFD", nil))
}

// DisableLog disables all library log output.  Logging output is disabled
// by default until UseLogger is called.
func DisableLog() {
	UseLogger(bronlog.Disabled)
}

// UseLogger uses a specified Logger to output package logging info.
// This should be used in preference to SetLogWriter if the caller is also
// using bronlog.
func UseLogger(logger bronlog.Logger) {
	log = logger
}
