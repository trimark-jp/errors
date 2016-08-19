package errors

import (
	"encoding/json"
	"testing"
)

func TestSourceOf(t *testing.T) {
	const (
		innerMessage  = "inner"
		middleMessage = "middle"
		outerMessage  = "outer"
	)
	outer := Wrap(Wrap(New(innerMessage), middleMessage), outerMessage).(*errorType)
	source := SourceOf(outer)

	if source.Error() != innerMessage {
		t.Fatal("invalid source", source)
	}
}

func TestWrapAsSource(t *testing.T) {
	const (
		innerMessage  = "inner"
		middleMessage = "middle"
		outerMessage  = "outer"
	)
	inner := New(innerMessage)
	middle := WrapBySourceMsg(inner, middleMessage)
	outer := Wrap(middle, outerMessage)
	source := SourceOf(outer)

	if source.Error() != middleMessage {
		t.Fatal("invalid source", source)
	}

	outer2 := WrapBySourceMsg(middle, outerMessage)
	source2 := SourceOf(outer2)
	if source2.Error() != outerMessage {
		t.Fatal("invalid source", source2)
	}

	b, _ := json.Marshal(outer2)
	t.Log(string(b))
}
