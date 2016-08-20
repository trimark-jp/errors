package errors

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
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
	log.Println(JSONAll(err))
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
