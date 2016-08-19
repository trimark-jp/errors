package errors

import "encoding/json"

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
