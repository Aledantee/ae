package ae

import "time"

// ErrorTimestamp defines an interface for errors that can provide a timestamp.
type ErrorTimestamp interface {
	// ErrorTimestamp returns the timestamp of the error.
	// Returns zero time if no timestamp is set.
	ErrorTimestamp() time.Time
}

func Timestamp(err error) time.Time {
	if ae, ok := err.(ErrorTimestamp); ok {
		return ae.ErrorTimestamp()
	}

	return time.Time{}
}
