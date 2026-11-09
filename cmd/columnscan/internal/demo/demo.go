// Package demo illustrates a vectorized filter and aggregate over aligned columns.
package demo

import (
	"example.com/learn/btrees/cmd/columnscan/internal/scan"
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
	table, err := scan.NewTable(
		[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		[]int64{10, 20, 30, 40, 50, 60, 70, 80, 90, 100},
		[]string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"},
	)
	if err != nil {
		return err
	}
	scanner, err := scan.NewScanner(table, 4)
	if err != nil {
		return err
	}
	fmt.Fprintln(d.out, "Column scan: process aligned columns in batches (no index).")
	fmt.Fprintln(d.out, "Filter id BETWEEN 3 AND 7; project id and amount; aggregate SUM(amount).")
	fmt.Fprintln(d.out, "Selection vectors contain positions local to each batch.")
	result := scanner.Scan(3, 7, d.out)
	fmt.Fprintf(d.out, "IDs: %v; SUM(amount)=%d\n", result.IDs, result.Sum)
	fmt.Fprintf(d.out, "Batches=%d; filter values read=%d; amount values read=%d\n", result.Batches, result.FilterValuesRead, result.AmountValuesRead)
	fmt.Fprintln(d.out, "The unused notes column is not read by the scanner.")
	return nil
}
