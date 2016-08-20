package errors

import (
	"encoding/json"
	"fmt"
)

var (
	// StringWithLocationFormat is the format for StringWithLocation.
	StringWithLocationFormat = "%s\t%s"
)

// Trace returns a json string which has error trace.
func Trace(err error) (string, error) {
	return TraceWithStack(err, 1)
}

// TraceAll returns a json string which has error trace.
func TraceAll(err error) (string, error) {
	return TraceWithStack(err, CallerInfoMaxStack)
}

// TraceWithStack returns a json string which has error trace.
func TraceWithStack(err error, stackCount int) (string, error) {
	em := &errMarshal{
		err: err,
	}
	em.setCallerCount(stackCount)
	b, e := json.Marshal(em)
	return string(b), e
}

// StringWithLocation returns error location and error message as string.
func StringWithLocation(err error) string {
	if err == nil {
		return ""
	}

	if e, ok := err.(*errorType); ok {
		return fmt.Sprintf(StringWithLocationFormat, e.info.String(), e.Error())
	}
	if c, ok := err.(*collection); ok {
		return c.StringWithLocation()
	}
	return err.Error()
}
