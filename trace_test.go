package errors

import (
	"errors"
	"fmt"
	"io"
	"os"
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

	t.Log(JSONAll(e))
}

func TestStringWithInner(t *testing.T) {
	inner := New("Inner Error")
	outer := Wrap(inner, "Outer Error")
	t.Log(StringWithInner(outer))
}

func TestMultipleError(t *testing.T) {
	f, _ := os.Open("unexistent_file")
	t.Log(writeDivision(f, 2, 0))
}

func writeDivision(w io.Writer, a int, b int) error {
	if b == 0 {
		msg := "tried to divide by zero"
		err := errors.New(msg)
		_, writeError := w.Write([]byte(msg))
		return Merge(err, writeError)
	}

	ret := a / b
	_, writeError := fmt.Fprintln(w, ret)
	return writeError
}
