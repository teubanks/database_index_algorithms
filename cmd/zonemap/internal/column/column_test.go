package column

import (
	"math/rand"
	"slices"
	"testing"
)

func TestPrunedScanMatchesFullScan(t *testing.T) {
	rng := rand.New(rand.NewSource(42))
	for _, size := range []int{1, 4, 17, 300} {
		values := make([]int, 211)
		for i := range values {
			values[i] = rng.Intn(200) - 100
		}
		column, err := New(values, size)
		if err != nil {
			t.Fatal(err)
		}
		for _, bounds := range [][2]int{{-20, 30}, {-100, 100}, {0, 0}, {300, 400}, {20, -20}} {
			var want []int
			for _, v := range values {
				if v >= bounds[0] && v <= bounds[1] {
					want = append(want, v)
				}
			}
			result := column.Scan(bounds[0], bounds[1], nil)
			if !slices.Equal(result.Values, want) {
				t.Fatalf("size=%d bounds=%v: result differs from full scan", size, bounds)
			}
			if bounds[0] <= bounds[1] && result.SegmentsScanned+result.SegmentsSkipped != len(column.segments) {
				t.Fatal("wrong segment accounting")
			}
			if result.RowsExamined > len(values) {
				t.Fatal("overcounted scan work")
			}
		}
	}
}

func TestOrderedPruningAndFalsePositive(t *testing.T) {
	c, _ := New([]int{1, 2, 3, 4, 10, 11, 12, 13, 20, 21, 22, 23}, 4)
	r := c.Scan(10, 12, nil)
	if !slices.Equal(r.Values, []int{10, 11, 12}) || r.SegmentsSkipped != 2 || r.RowsExamined != 4 {
		t.Fatalf("unexpected ordered scan: %+v", r)
	}
	c, _ = New([]int{1, 100}, 2)
	r = c.Scan(50, 50, nil)
	if len(r.Values) != 0 || r.SegmentsScanned != 1 || r.RowsExamined != 2 {
		t.Fatalf("summary overlap was treated as an exact match: %+v", r)
	}
}

func TestEmptyBoundariesAndInputOwnership(t *testing.T) {
	if _, err := New(nil, 0); err == nil {
		t.Fatal("accepted zero segment size")
	}
	if _, err := New(nil, -1); err == nil {
		t.Fatal("accepted negative segment size")
	}
	c, _ := New(nil, 4)
	if r := c.Scan(0, 10, nil); r.RowsExamined != 0 || len(r.Values) != 0 {
		t.Fatal("incorrect empty scan")
	}
	maxInt := int(^uint(0) >> 1)
	minInt := -maxInt - 1
	values := []int{minInt, 0, maxInt}
	c, _ = New(values, 2)
	values[0] = 99
	if r := c.Scan(minInt, maxInt, nil); !slices.Equal(r.Values, []int{minInt, 0, maxInt}) {
		t.Fatal("changed input modified summaries/data")
	}
}
