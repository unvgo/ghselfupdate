package ghselfupdate

import (
	"io"
	stdlog "log"
	"os"
)

var logger = stdlog.New(io.Discard, "", 0)
var logEnabled = false

// EnableLog enables to output logging messages in library
func EnableLog() {
	if logEnabled {
		return
	}
	logEnabled = true
	logger.SetOutput(os.Stderr)
	logger.SetFlags(stdlog.Ltime)
}

// DisableLog disables to output logging messages in library
func DisableLog() {
	if !logEnabled {
		return
	}
	logEnabled = false
	logger.SetOutput(io.Discard)
	logger.SetFlags(0)
}
