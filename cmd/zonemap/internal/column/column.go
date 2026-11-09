// Package column demonstrates conservative min/max pruning of column segments.
package column

import (
	"fmt"
	"io"
)

// Column owns immutable segments and the summaries computed from their values.
type Column struct{ segments []segment }

// Result records matching values in storage order and actual scan work.
type Result struct {
	Values          []int
	SegmentsSkipped int
	SegmentsScanned int
	RowsExamined    int
}

// New creates an independent column with at most segmentSize values per segment.
func New(values []int, segmentSize int) (*Column, error) {
	if segmentSize <= 0 {
		return nil, fmt.Errorf("segment size must be positive")
	}
	c := &Column{}
	for start := 0; start < len(values); {
		end := start + min(segmentSize, len(values)-start)
		c.segments = append(c.segments, newSegment(values[start:end]))
		start = end
	}
	return c, nil
}

// Scan uses min/max to reject impossible segments; overlap only means MAY match.
// The exact inclusive predicate is always checked against surviving values.
func (c *Column) Scan(low, high int, trace io.Writer) Result {
	if trace == nil {
		trace = io.Discard
	}
	var result Result
	if low > high {
		return result
	}
	for i := range c.segments {
		s := &c.segments[i]
		fmt.Fprintf(trace, "segment %d min=%d max=%d: ", i, s.min, s.max)
		if !s.mayMatch(low, high) {
			result.SegmentsSkipped++
			fmt.Fprintln(trace, "skip")
			continue
		}
		fmt.Fprintln(trace, "scan (range overlaps; check actual values)")
		result.SegmentsScanned++
		result.RowsExamined += len(s.values)
		result.Values = append(result.Values, s.matchingValues(low, high)...)
	}
	return result
}

type segment struct {
	values   []int
	min, max int
}

func newSegment(values []int) segment {
	s := segment{values: append([]int(nil), values...), min: values[0], max: values[0]}
	for _, v := range values[1:] {
		s.min = min(s.min, v)
		s.max = max(s.max, v)
	}
	return s
}

func (s *segment) mayMatch(low, high int) bool { return low <= s.max && high >= s.min }

func (s *segment) matchingValues(low, high int) []int {
	var result []int
	for _, value := range s.values {
		if value >= low && value <= high {
			result = append(result, value)
		}
	}
	return result
}
