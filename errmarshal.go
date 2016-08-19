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
		obj := struct {
			Message string `json:"message"`
		}{
			Message: "",
		}
		return json.Marshal(&obj)
	}
	obj := struct {
		Message string `json:"message"`
	}{
		Message: e.err.Error(),
	}
	return json.Marshal(&obj)
}
