package errors

import (
	"encoding/json"
	"fmt"
)

type (
	errorSource struct {
		*errorType
		source error
	}
)

// Error implements error interface.
func (e *errorSource) Error() string {
	return e.source.Error()
}

// NewAsSource returns a new error which is the source.
func NewAsSource(msg string) error {
	return newSource(nil, new(nil, msg, 1), 1)
}

// NewAsSourcef returns a new error which is the source.
func NewAsSourcef(format string, a ...interface{}) error {
	return newSource(nil, new(nil, fmt.Sprintf(format, a...), 1), 1)
}

// AsSource returns a new source error.
func AsSource(err error) error {
	return newSource(nil, err, 1)
}

// WrapBySourceError returns a new error.
// If the error is passed to errors.SourceOf function,
// returns the source.
func WrapBySourceError(inner error, source error) error {
	if inner == nil {
		return nil
	}
	return newSource(inner, source, 1)
}

// WrapBySourceMsg returns a new error.
func WrapBySourceMsg(inner error, msg string) error {
	if inner == nil {
		return nil
	}
	return newSource(inner, new(nil, msg, 1), 1)
}

// WrapBySourceMsgf returns a new error.
func WrapBySourceMsgf(inner error, format string, a ...interface{}) error {
	if inner == nil {
		return nil
	}
	return newSource(inner, new(nil, fmt.Sprintf(format, a...), 1), 1)
}

// SourceOf returns the source error of the err.
func SourceOf(err error) error {
	if e, ok := err.(*errorSource); ok {
		return e.source
	}
	if e, ok := err.(*collection); ok {
		return e.source()
	}
	if e, ok := err.(*errorType); ok {
		if e.inner != nil {
			return SourceOf(e.inner)
		}
		return e
	}
	return err
}

// ExplicitSourceOf returns an error if the error is
// explicitly specified Source by WrapBySourceXxx.
func ExplicitSourceOf(err error) error {
	if e, ok := err.(*errorSource); ok {
		return e.source
	}
	if e, ok := err.(*collection); ok {
		return e.explicitSource()
	}
	if e, ok := err.(*errorType); ok {
		if e.inner != nil {
			return ExplicitSourceOf(e.inner)
		}
		return nil
	}
	return nil
}

// MarshalJSON implements json.Marshaler interface.
func (e *errorSource) MarshalJSON() ([]byte, error) {
	obj := struct {
		Inner    *errMarshal `json:"inner"`
		Callers  *callerInfo `json:"callers"`
		Message  string      `json:"message"`
		IsSource bool        `json:"isSource"`
	}{
		Inner: &errMarshal{
			err:         e.inner,
			callerCount: e.callerCount,
		},
		Callers:  e.info,
		Message:  e.source.Error(),
		IsSource: true,
	}
	return json.Marshal(&obj)
}

func newSource(inner error, source error, skip int) error {
	return &errorSource{
		errorType: new(inner, "", skip+1).(*errorType),
		source:    source,
	}
}
