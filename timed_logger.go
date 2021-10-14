package mandelgif

import (
	"log"
	"time"
)

// A tlogger wraps the stdlib logger, prefixing each log messages with the
// duration elapsed since its creation.
type tlogger int64

// timedLogger creates a tlogger, using time.Now() as base time.
func timedLogger() tlogger {
	return tlogger(time.Now().UnixNano())
}

// dur returns the time elapsed since the logger creation.
func (l tlogger) dur() string {
	return time.Duration(time.Now().UnixNano() - int64(l)).String()
}

// logf prints a log message with the standard logger.
// Arguments are handled in the manner of fmt.Printf.
func (l tlogger) logf(format string, args ...interface{}) {
	log.Printf(l.dur()+" "+format, args...)
}
