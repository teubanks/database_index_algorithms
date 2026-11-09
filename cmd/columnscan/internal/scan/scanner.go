package scan

import (
	"fmt"
	"io"
)

// Scanner owns the execution configuration for a columnar table scan.
type Scanner struct {
	table     *Table
	batchSize int
}

// Result exposes matches, the aggregate, and the work actually performed.
type Result struct {
	IDs              []int
	Sum              int64
	Batches          int
	FilterValuesRead int
	AmountValuesRead int
}

func NewScanner(table *Table, batchSize int) (*Scanner, error) {
	if table == nil {
		return nil, fmt.Errorf("table must not be nil")
	}
	if batchSize <= 0 {
		return nil, fmt.Errorf("batch size must be positive")
	}
	return &Scanner{table: table, batchSize: batchSize}, nil
}

// Scan models SELECT id, amount WHERE id BETWEEN low AND high, plus SUM(amount).
// It reads IDs in batches, builds a selection vector of local row positions,
// then reads amounts only at those positions. The notes column is never read.
func (s *Scanner) Scan(low, high int, trace io.Writer) Result {
	if trace == nil {
		trace = io.Discard
	}
	var result Result
	if low > high {
		return result
	}
	for start := 0; start < len(s.table.ids); {
		end := start + min(s.batchSize, len(s.table.ids)-start)
		var selected []int
		for local, id := range s.table.ids[start:end] {
			result.FilterValuesRead++
			if id >= low && id <= high {
				selected = append(selected, local)
			}
		}
		fmt.Fprintf(trace, "batch [%d,%d): ids=%v selection=%v\n", start, end, s.table.ids[start:end], selected)
		for _, local := range selected {
			row := start + local
			result.IDs = append(result.IDs, s.table.ids[row])
			result.Sum += s.table.amounts[row]
			result.AmountValuesRead++
		}
		result.Batches++
		start = end
	}
	return result
}
