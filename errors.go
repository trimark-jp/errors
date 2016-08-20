package errors

import (
	"encoding/json"
	"fmt"
)

type (
	errorType struct {
		inner       error
		msg         string
		info        *callerInfo
		callerCount int
	}
)

// Error implements error interface.
func (e *errorType) Error() string {
	return e.msg
}

// New returns a new error.
func New(msg string) error {
	return new(nil, msg, 1)
}

// Newf returns a new error.
func Newf(format string, a ...interface{}) error {
	return new(nil, fmt.Sprintf(format, a...), 1)
}

// Wrap returns the err by new error wich has msg.
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	return new(err, msg, 1)
}

// Wrapf returns the err by new error wich has msg.
func Wrapf(err error, format string, a ...interface{}) error {
	if err == nil {
		return nil
	}
	return new(err, fmt.Sprintf(format, a...), 1)
}

// MarshalJSON implements json.Marshaler interface.
func (e *errorType) MarshalJSON() ([]byte, error) {
	obj := struct {
		Inner   *errMarshal `json:"inner"`
		Callers *callerInfo `json:"callers"`
		Message string      `json:"message"`
	}{
		Inner: &errMarshal{
			err:         e.inner,
			callerCount: e.callerCount,
		},
		Callers: e.info,
		Message: e.msg,
	}
	return json.Marshal(&obj)
}

func (e *errorType) setCallerCount(n int) {
	if m, ok := e.inner.(callerMarshal); ok {
		m.setCallerCount(n)
	}
	e.callerCount = n
	e.info.setCount(n)
}

func new(inner error, msg string, skip int) error {
	return &errorType{
		inner: inner,
		msg:   msg,
		info:  caller(skip + 1),
	}
}
