package errors

import "encoding/json"

// Trace returns a json string which has error trace.
func Trace(err error) (string, error) {
	em := &errMarshal{
		err: err,
	}
	b, e := json.Marshal(em)
	return string(b), e
}
