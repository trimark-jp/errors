package errors

import "encoding/json"

type (
	errMarshal struct {
		err         error
		callerCount int
	}

	callerMarshal interface {
		json.Marshaler
		setCallerCount(int)
	}
)

func (e *errMarshal) setCallerCount(n int) {
	if m, ok := e.err.(callerMarshal); ok {
		m.setCallerCount(n)
	} else {
	}
	e.callerCount = n
}

func (e *errMarshal) MarshalJSON() ([]byte, error) {
	if m, ok := e.err.(json.Marshaler); ok {
		return m.MarshalJSON()
	}
	if e.err == nil {
		return json.Marshal(nil)
	}
	obj := struct {
		Message string `json:"message"`
	}{
		Message: e.err.Error(),
	}
	return json.Marshal(&obj)
}
