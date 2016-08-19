package errors

import (
	"encoding/json"
	"fmt"
	"runtime"
)

type (
	callerInfoItem struct {
		File     string `json:"file"`
		Line     int    `json:"line"`
		Function string `json:"function"`
	}
	callerInfo struct {
		Items       []*callerInfoItem `json:"items"`
		outputCount int
	}
)

var (
	// CallerInfoMaxStack is maximum stack count for caller info.
	CallerInfoMaxStack = 10
)

func (c *callerInfo) String() string {
	if len(c.Items) <= 0 {
		return ""
	}

	first := c.Items[0]
	return fmt.Sprintf("%s:%d:%s", first.File, first.Line, first.Function)
}

func (c *callerInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.Items[:c.outputCount])
}

func caller(skip int) *callerInfo {
	pcs := make([]uintptr, CallerInfoMaxStack)
	n := runtime.Callers(skip+2, pcs)
	pcs = pcs[:n]

	result := &callerInfo{
		Items: make([]*callerInfoItem, len(pcs)),
	}

	for index, pc := range pcs {
		f := runtime.FuncForPC(pc)
		function := f.Name()
		file, line := f.FileLine(pc)
		item := &callerInfoItem{
			Function: function,
			File:     trimGOPATHProbably(file, function),
			Line:     line,
		}
		result.Items[index] = item
	}
	return result
}

func (c *callerInfo) setCount(n int) {
	count := n
	if len(c.Items) < count {
		count = len(c.Items)
	}
	if count < 0 {
		count = 0
	}
	c.outputCount = count
}
