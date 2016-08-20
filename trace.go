package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
)

var (
	// StringWithLocationFormat is the format for StringWithLocation.
	StringWithLocationFormat = "%s\t%s"

	// StringWithInnerIndent is indent for inner errors.
	StringWithInnerIndent = "\t"
)

// JSON returns a json string which has error trace.
func JSON(err error) (string, error) {
	return JSONWithStack(err, 1)
}

// JSONAll returns a json string which has error trace.
func JSONAll(err error) (string, error) {
	return JSONWithStack(err, CallerInfoMaxStack)
}

// JSONWithStack returns a json string which has error trace.
func JSONWithStack(err error, stackCount int) (string, error) {
	if err == nil {
		return "", nil
	}

	em := &errMarshal{
		err: err,
	}
	em.setCallerCount(stackCount)
	b, e := json.Marshal(em)
	return string(b), e
}

// StringWithLocation returns error location and error message as string.
func StringWithLocation(err error) string {
	if err == nil {
		return ""
	}

	if e, ok := err.(*errorType); ok {
		return fmt.Sprintf(StringWithLocationFormat, e.info.String(), e.Error())
	}
	if e, ok := err.(*errorSource); ok {
		return fmt.Sprintf(StringWithLocationFormat, e.info.String(), e.Error())
	}
	if c, ok := err.(*collection); ok {
		return c.StringWithLocation()
	}
	return err.Error()
}

// StringWithInner inner returns string representation of the error and inner errors.
func StringWithInner(err error) string {
	if err == nil {
		return ""
	}
	return stringWithInner(err, "")
}

func stringWithInner(err error, indent string) string {
	buf := &bytes.Buffer{}

	if e, ok := err.(*errorType); ok {
		fmt.Fprintln(buf, indent+StringWithLocation(err))
		if e.inner != nil {
			fmt.Fprint(buf, stringWithInner(e.inner, indent+StringWithInnerIndent))
		}
		return buf.String()
	}
	if e, ok := err.(*errorSource); ok {
		fmt.Fprintln(buf, indent+StringWithLocation(err))
		if e.inner != nil {
			fmt.Fprint(buf, stringWithInner(e.inner, indent+StringWithInnerIndent))
		}
		return buf.String()
	}

	if c, ok := err.(*collection); ok {
		fmt.Fprint(buf, c.stringWithInner(indent+StringWithInnerIndent))
		return buf.String()
	}

	fmt.Fprintln(buf, indent+StringWithLocation(err))
	return buf.String()
}
