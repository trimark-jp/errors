package errors

import "testing"

func TestCaller(t *testing.T) {
	info := caller(0).Items[0]
	if info.File != "github.com/trimark-jp/errors/caller_test.go" {
		t.Fatal("can't get caller info file ", info.File)
	}
	if info.Function != "github.com/trimark-jp/errors.TestCaller" {
		t.Fatal("can't get caller info Function", info.Function)
	}
	if info.Line != 6 {
		t.Fatal("can't get caller info Line", info.Line)
	}
}
