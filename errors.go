package flagx

import (
	"fmt"
)

// simpleflagError is an error object with message and underlying error.
type simpleflagError struct {
	msg string
	err error
}

func (e *simpleflagError) Unwrap() error { return e.err }
func (e *simpleflagError) Error() string { return e.msg }

// wrapNameError returns an error object with inner error and "name: err" message.
// If err is nil, wrapNameError returns nil.
func wrapNameError(err error, name string) error {
	if err == nil {
		return nil
	}
	return wrapErrorf(err, "%s: %s", name, err.Error())
}

// wrapNameErrorString returns an error object with inner error and "name: err str" message.
// If err is nil, wrapNameErrorString returns nil.
func wrapNameErrorString(err error, name, str string) error {
	if err == nil {
		return nil
	}
	return wrapErrorf(err, "%s: %s %q", name, err.Error(), str)
}

// wrapErrorf returns an error object with inner error and formatted message.
// If err is nil, wrapErrorf returns nil.
func wrapErrorf(err error, format string, a ...interface{}) error {
	if err == nil {
		return nil
	}
	return &simpleflagError{fmt.Sprintf(format, a...), err}
}
