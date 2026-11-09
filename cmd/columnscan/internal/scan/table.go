// Package scan demonstrates column projection and batch-at-a-time filtering.
package scan

import "fmt"

// Table owns aligned columns. Position i refers to the same row in each column.
type Table struct {
	ids     []int
	amounts []int64
	notes   []string
}

// NewTable copies the input, so external changes cannot break column alignment.
func NewTable(ids []int, amounts []int64, notes []string) (*Table, error) {
	if len(ids) != len(amounts) || len(ids) != len(notes) {
		return nil, fmt.Errorf("columns must have equal lengths")
	}
	return &Table{
		ids:     append([]int(nil), ids...),
		amounts: append([]int64(nil), amounts...),
		notes:   append([]string(nil), notes...),
	}, nil
}
