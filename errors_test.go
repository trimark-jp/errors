package errors

import (
	"encoding/json"
	"testing"
)

func TestErrors(t *testing.T) {
	const testMessage = "test message"
	e := New(testMessage)
	if e.Error() != testMessage {
		t.Fatal("invalid error message", e.Error())
	}
}

func TestWrap(t *testing.T) {
	const (
		innerMessage = "inner"
		outerMessage = "outer"
	)
	inner := Newf("%s", innerMessage)
	outer := Wrapf(inner, "%s", outerMessage)

	e := outer.(*errorType)
	if e.Error() != outerMessage {
		t.Fatal("invalid outer message", e)
	}
	if e.inner.Error() != innerMessage {
		t.Fatal("invalid inner message", e.inner)
	}
}

func TestMarshal(t *testing.T) {
	const (
		innerMessage = "inner"
		outerMessage = "outer"
	)
	inner := Newf("%s", innerMessage)
	outer := Wrapf(inner, "%s", outerMessage)

	b, _ := json.Marshal(outer)
	t.Log(string(b))
}
