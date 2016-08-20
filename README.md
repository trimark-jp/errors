# errors

This package provides interface for error handling.

## Purpose

* Error Stack Trace
* Trace multiple errors
* Mark source error


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

### With Explicit Source of error

```go
func TestExplicitSource(t *testing.T) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "https://localhost/login?user=username&password=secret", nil)
	if err != nil {
		t.Fatal(err)
	}
	handle(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatal("invalid status", w.Code)
	}
	buf := &bytes.Buffer{}
	io.Copy(buf, w.Body)
	response := buf.String()

	if response != errMessageInvalidUserIDFormat {
		t.Fatal("invalid response", response)
	}
}

const (
	errMessageUserIsRequired      = "user is required"
	errMessagePasswordIsRequired  = "password is required"
	errMessageUserNotFound        = "user not found"
	errMessageInvalidUserIDFormat = "user id must be integer"
	errMessagePasswordMismatch    = "password wrong"
	errMessageLoginFailed         = "login failed"
	errMessageURLNotFound         = "url not found"
)

type httperror struct {
	StatusCode int
	Message    string
}

func (e *httperror) Error() string {
	return fmt.Sprintf("%d %s", e.StatusCode, e.Message)
}

func handle(w http.ResponseWriter, req *http.Request) {
	var err error

	if req.URL.Path == "/login" {
		err = loginByRequest(req)
	} else {
		err = AsSource(&httperror{
			StatusCode: http.StatusNotFound,
			Message:    errMessageURLNotFound,
		})
	}
	if e, ok := SourceOf(err).(*httperror); ok {
		w.WriteHeader(e.StatusCode)
		fmt.Fprint(w, e.Message)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "unknown err")
	}
}

func loginByRequest(req *http.Request) error {
	query := req.URL.Query()
	user := query.Get("user")
	password := query.Get("password")
	err := login(user, password)
	return Wrap(err, errMessageLoginFailed)
}

func login(user string, password string) error {
	if user == "" {
		return AsSource(&httperror{
			StatusCode: http.StatusBadRequest,
			Message:    errMessageUserIsRequired,
		})
	}

	userID, err := strconv.Atoi(user)
	if err != nil {
		return WrapBySourceError(err,
			&httperror{
				StatusCode: http.StatusBadRequest,
				Message:    errMessageInvalidUserIDFormat,
			})
	}
	if password == "" {
		return AsSource(&httperror{
			StatusCode: http.StatusBadRequest,
			Message:    errMessagePasswordIsRequired,
		})
	}

	if userID != 10 {
		return AsSource(&httperror{
			StatusCode: http.StatusForbidden,
			Message:    errMessageURLNotFound,
		})
	}
	if password != "secret" {
		return AsSource(&httperror{
			StatusCode: http.StatusForbidden,
			Message:    errMessagePasswordMismatch,
		})
	}
	return nil
}
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