package funding

import (
	"github.com/brronsuite/broln/build"
	"github.com/brronsuite/bronlog"
)

// Subsystem defines the logging code for this subsystem.
const Subsystem = "FNDG"

// log is a logger that is initialized with the bronlog.Disabled logger.
var log bronlog.Logger

// The default amount of logging is none.
func init() {
	UseLogger(build.NewSubLogger(Subsystem, nil))
}

// DisableLog disables all logging output.
func DisableLog() {
	UseLogger(bronlog.Disabled)
}

// UseLogger uses a specified Logger to output package logging info.
func UseLogger(logger bronlog.Logger) {
	log = logger
}
