package errors

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"testing"
)

func TestCollection(t *testing.T) {
	const (
		leftMessage = "left"
	)
	left := New(leftMessage)
	right := io.ErrUnexpectedEOF
	c := Merge(left, right)

	expected := fmt.Sprintf(collectionStringFormat,
		strings.Join([]string{left.Error(), right.Error()}, collectionSeparator))
	if expected != c.Error() {
		t.Fatal("invalid message", c.Error())
	}

}

func TestSourceOfCollection(t *testing.T) {
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

	b, _ := json.Marshal(e)
	t.Log(string(b))
}
