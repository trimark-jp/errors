package errors

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestMarhsalNil(t *testing.T) {
	e := &errMarshal{
		err: nil,
	}

	b, _ := json.Marshal(e)
	s := string(b)
	if s != `{"message":""}` {
		t.Fatal("marshalled nil wrong")
	}
}
func TestMarshalNormalErr(t *testing.T) {
	const (
		msg = "normal error"
	)
	e := &errMarshal{
		err: errors.New(msg),
	}

	b, _ := json.Marshal(e)
	s := string(b)
	if s != `{"message":"`+msg+`"}` {
		t.Fatal("marshalled nil wrong", s)
	}
}
