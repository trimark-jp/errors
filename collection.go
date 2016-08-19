package errors

import (
	"encoding/json"
	"fmt"
	"strings"
)

type (
	// collection has a list of errors.
	collection struct {
		errs        []error
		callerCount int
	}
)

const (
	collectionStringFormat = "Errors: [%s]"
	collectionSeparator    = ", "
)

// Error implements error interface.
func (c *collection) Error() string {
	if len(c.errs) <= 0 {
		return ""
	}

	msgs := make([]string, len(c.errs))
	for index, err := range c.errs {
		msgs[index] = err.Error()
	}
	return fmt.Sprintf(collectionStringFormat,
		strings.Join(msgs, collectionSeparator))

}

// Merge returns an error which has l and r.
func Merge(l error, r error) error {
	if l == nil {
		return r
	}
	if r == nil {
		return l
	}

	if c, ok := l.(*collection); ok {
		c.append(r)
		return c
	}
	if c, ok := r.(*collection); ok {
		c.insertFront(l)
		return c
	}
	c := newCollection()
	c.append(l)
	c.append(r)
	return c
}

// MarshalJSON implements json.Marshaler interface.
func (c *collection) MarshalJSON() ([]byte, error) {
	obj := struct {
		Errs []*errMarshal `json:"errors"`
	}{
		Errs: make([]*errMarshal, len(c.errs)),
	}
	for index, e := range c.errs {
		em := &errMarshal{
			err:         e,
			callerCount: c.callerCount,
		}
		obj.Errs[index] = em
	}
	return json.Marshal(&obj)
}

// newCollection returns a new collection.
func newCollection() *collection {
	return &collection{
		errs: []error{},
	}
}

// append appends the error to the collection.
func (c *collection) append(err error) {
	if tail, ok := err.(*collection); ok {
		c.errs = append(c.errs, tail.errs...)
		return
	}
	c.errs = append(c.errs, err)
}

func (c *collection) insertFront(err error) {
	if head, ok := err.(*collection); ok {
		c.errs = append(head.errs, c.errs...)
		return
	}
	c.errs = append([]error{err}, c.errs...)
}

func (c *collection) source() error {
	if len(c.errs) <= 0 {
		return nil
	}

	s := c.explicitSource()
	if s != nil {
		return s
	}
	return SourceOf(c.errs[0])
}

func (c *collection) explicitSource() error {
	for _, err := range c.errs {
		s := ExplicitSourceOf(err)
		if s != nil {
			return s
		}
	}
	return nil
}

func (c *collection) setCallerCount(n int) {
	for _, e := range c.errs {
		if m, ok := e.(callerMarshal); ok {
			m.setCallerCount(n)
		}
	}
	c.callerCount = n
}
