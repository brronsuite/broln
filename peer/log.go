package peer

import (
	"github.com/brronsuite/broln/build"
	"github.com/brronsuite/bronlog"
)

// peerLog is a logger that is initialized with the bronlog.Disabled logger.
var peerLog bronlog.Logger

// The default amount of logging is none.
func init() {
	UseLogger(build.NewSubLogger("PEER", nil))
}

// DisableLog disables all logging output.
func DisableLog() {
	UseLogger(bronlog.Disabled)
}

// UseLogger uses a specified Logger to output package logging info.
func UseLogger(logger bronlog.Logger) {
	peerLog = logger
}

// logClosure is used to provide a closure over expensive logging operations
// so they aren't performed when the logging level doesn't warrant it.
type logClosure func() string

// String invokes the underlying function and returns the result.
func (c logClosure) String() string {
	return c()
}

// newLogClosure returns a new closure over a function that returns a string
// which itself provides a Stringer interface so that it can be used with the
// logging system.
func newLogClosure(c func() string) logClosure {
	return logClosure(c)
}
