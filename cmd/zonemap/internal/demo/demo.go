// Package demo compares pruning for the same values in different storage orders.
package demo

import (
	"example.com/learn/btrees/cmd/zonemap/internal/column"
	"fmt"
	"io"
)

type Demo struct{ out io.Writer }

func New(out io.Writer) *Demo {
	if out == nil {
		out = io.Discard
	}
	return &Demo{out: out}
}

func (d *Demo) Run() error {
	fmt.Fprintln(d.out, "Zone maps: skip impossible segments; overlapping segments still need a scan.")
	datasets := []struct {
		name   string
		values []int
	}{
		{"ordered", []int{1, 2, 3, 4, 10, 11, 12, 13, 20, 21, 22, 23}},
		{"mixed", []int{1, 10, 20, 2, 11, 21, 3, 12, 22, 4, 13, 23}},
	}
	for _, dataset := range datasets {
		values, err := column.New(dataset.values, 4)
		if err != nil {
			return err
		}
		fmt.Fprintf(d.out, "\n%s values: %v\nQuery [10,12]\n", dataset.name, dataset.values)
		result := values.Scan(10, 12, d.out)
		fmt.Fprintf(d.out, "Matches: %v; skipped=%d scanned=%d values examined=%d\n", result.Values, result.SegmentsSkipped, result.SegmentsScanned, result.RowsExamined)
	}
	values, err := column.New([]int{1, 100}, 2)
	if err != nil {
		return err
	}
	fmt.Fprintln(d.out, "\nFalse positive: [1,100] has an overlapping summary for [50,50].")
	result := values.Scan(50, 50, d.out)
	fmt.Fprintf(d.out, "Matches: %v; values examined=%d\n", result.Values, result.RowsExamined)
	return nil
}
