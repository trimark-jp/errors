package errors

import (
	"io"
	"testing"
)

func TestTrace(t *testing.T) {
	const (
		leftInnerMessage   = "left inner"
		leftMiddleMessage  = "left middle"
		leftOuterMessage   = "left outer"
		rightInnerMessage  = "right inner"
		rightMiddleMessage = "right middle"
		rightOuterMessage  = "right outer"
	)

	li := New(leftInnerMessage)
	lm := Wrap(li, leftMiddleMessage)
	lo := Wrap(lm, leftOuterMessage)

	ri := New(rightInnerMessage)
	rm := WrapBySourceMsg(ri, rightMiddleMessage)
	ro := Wrap(rm, rightOuterMessage)

	c := Merge(lo, ro)

	s := SourceOf(c)
	if s.Error() != rightMiddleMessage {
		t.Error("invalid source", s)
	}

	e := Wrap(c, "most out")

	t.Log(JSON(e))
	t.Log(JSONWithStack(e, 1))
}

func TestStringWithLocation(t *testing.T) {
	const (
		leftInnerMessage   = "left inner"
		leftMiddleMessage  = "left middle"
		leftOuterMessage   = "left outer"
		rightInnerMessage  = "right inner"
		rightMiddleMessage = "right middle"
		rightOuterMessage  = "right outer"
	)

	other := io.ErrClosedPipe
	li := Wrap(other, leftInnerMessage)
	lm := Wrap(li, leftMiddleMessage)
	lo := Wrap(lm, leftOuterMessage)

	ri := New(rightInnerMessage)
	rm := WrapBySourceMsg(ri, rightMiddleMessage)
	ro := Wrap(rm, rightOuterMessage)

	c := Merge(lo, ro)

	s := SourceOf(c)
	if s.Error() != rightMiddleMessage {
		t.Error("invalid source", s)
	}

	e := Wrap(c, "most out")

	// t.Log(StringWithLocation(e))
	// t.Log(StringWithLocation(c))
	t.Log(StringWithInner(e))
}
