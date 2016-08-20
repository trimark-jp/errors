package errors

import (
	"fmt"
)

func Example() {
	msg := "Error Message"
	e := New(msg)
	s := StringWithLocation(e)
	fmt.Println(s)

	// Output:
	// github.com/trimark-jp/errors/example_test.go:9:github.com/trimark-jp/errors.Example	Error Message
}
