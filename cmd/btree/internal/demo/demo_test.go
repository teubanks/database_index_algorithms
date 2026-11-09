package demo_test

import (
	"bytes"
	"testing"

	"example.com/learn/btrees/cmd/btree/internal/demo"
)

func TestRunStartsFresh(t *testing.T) {
	var out bytes.Buffer
	walkthrough := demo.New(&out)
	walkthrough.Run()
	first := out.String()
	out.Reset()
	walkthrough.Run()
	if first == "" || out.String() != first {
		t.Fatal("repeated walkthrough changed its output")
	}
}
