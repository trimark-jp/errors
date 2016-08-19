package errors

import (
	"encoding/json"
)

type errMarshal struct {
	err error
}

func (e *errMarshal) MarshalJSON() ([]byte, error) {
	if m, ok := e.err.(json.Marshaler); ok {
		return m.MarshalJSON()
	}
	if e.err == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(e.err.Error())
}
