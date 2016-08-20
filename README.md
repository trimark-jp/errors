# errors

This package provides interface for error handling.

## install

```bash
go get github.com/trimark-jp/errors
```

## Usage

### Error message with location

```go
e := errors.New("Error Message")
StringWithLocation(e)
```

Outputs:

```txt
github.com/trimark-jp/errors/example_test.go:9:github.com/trimark-jp/errors.Example	Error Message
```

### With internal errors

```go
inner := errors.New("Inner Error")
outer := errors.Wrap(inner, "Outer Error")
StringWithInner(outer)
```

Outputs:

```txt
	github.com/trimark-jp/errors/trace_test.go:74:github.com/trimark-jp/errors.TestStringWithInner	Outer Error
			github.com/trimark-jp/errors/trace_test.go:73:github.com/trimark-jp/errors.TestStringWithInner	Inner Error
```

### With multiple errors

```go
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
```

Outputs:

```txt
 Errors: [tried to divide by zero, invalid argument]
```

### With All Stack Trace

```go
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
```

Outputs:

```
{
    "inner": {
        "errors": [
            {
                "inner": {
                    "inner": {
                        "inner": {
                            "message": "io: read/write on closed pipe"
                        },
                        "callers": [
                            {
                                "file": "github.com/trimark-jp/errors/trace_test.go",
                                "line": 53,
                                "function": "github.com/trimark-jp/errors.TestStringWithLocation"
                            },
                            {
                                "file": "testing/testing.go",
                                "line": 611,
                                "function": "testing.tRunner"
                            },
                            {
                                "file": "runtime/asm_amd64.s",
                                "line": 2087,
                                "function": "runtime.goexit"
                            }
                        ],
                        "message": "left inner"
                    },
                    "callers": [
                        {
                            "file": "github.com/trimark-jp/errors/trace_test.go",
                            "line": 54,
                            "function": "github.com/trimark-jp/errors.TestStringWithLocation"
                        },
                        {
                            "file": "testing/testing.go",
                            "line": 611,
                            "function": "testing.tRunner"
                        },
                        {
                            "file": "runtime/asm_amd64.s",
                            "line": 2087,
                            "function": "runtime.goexit"
                        }
                    ],
                    "message": "left middle"
                },
                "callers": [
                    {
                        "file": "github.com/trimark-jp/errors/trace_test.go",
                        "line": 55,
                        "function": "github.com/trimark-jp/errors.TestStringWithLocation"
                    },
                    {
                        "file": "testing/testing.go",
                        "line": 611,
                        "function": "testing.tRunner"
                    },
                    {
                        "file": "runtime/asm_amd64.s",
                        "line": 2087,
                        "function": "runtime.goexit"
                    }
                ],
                "message": "left outer"
            },
            {
                "inner": {
                    "inner": {
                        "inner": null,
                        "callers": [
                            {
                                "file": "github.com/trimark-jp/errors/trace_test.go",
                                "line": 57,
                                "function": "github.com/trimark-jp/errors.TestStringWithLocation"
                            },
                            {
                                "file": "testing/testing.go",
                                "line": 611,
                                "function": "testing.tRunner"
                            },
                            {
                                "file": "runtime/asm_amd64.s",
                                "line": 2087,
                                "function": "runtime.goexit"
                            }
                        ],
                        "message": "right inner"
                    },
                    "callers": [
                        {
                            "file": "github.com/trimark-jp/errors/trace_test.go",
                            "line": 58,
                            "function": "github.com/trimark-jp/errors.TestStringWithLocation"
                        },
                        {
                            "file": "testing/testing.go",
                            "line": 611,
                            "function": "testing.tRunner"
                        },
                        {
                            "file": "runtime/asm_amd64.s",
                            "line": 2087,
                            "function": "runtime.goexit"
                        }
                    ],
                    "message": "right middle",
                    "isSource": true
                },
                "callers": [
                    {
                        "file": "github.com/trimark-jp/errors/trace_test.go",
                        "line": 59,
                        "function": "github.com/trimark-jp/errors.TestStringWithLocation"
                    },
                    {
                        "file": "testing/testing.go",
                        "line": 611,
                        "function": "testing.tRunner"
                    },
                    {
                        "file": "runtime/asm_amd64.s",
                        "line": 2087,
                        "function": "runtime.goexit"
                    }
                ],
                "message": "right outer"
            }
        ]
    },
    "callers": [
        {
            "file": "github.com/trimark-jp/errors/trace_test.go",
            "line": 68,
            "function": "github.com/trimark-jp/errors.TestStringWithLocation"
        },
        {
            "file": "testing/testing.go",
            "line": 611,
            "function": "testing.tRunner"
        },
        {
            "file": "runtime/asm_amd64.s",
            "line": 2087,
            "function": "runtime.goexit"
        }
    ],
    "message": "most out"
}
```