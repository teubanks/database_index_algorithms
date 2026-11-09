package scan

import (
	"math/rand"
	"slices"
	"testing"
)

func TestBatchesMatchScalarReference(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	ids, amounts, notes := make([]int, 103), make([]int64, 103), make([]string, 103)
	for i := range ids {
		ids[i] = rng.Intn(50) - 25
		amounts[i] = int64(rng.Intn(100) - 50)
		notes[i] = "unused"
	}
	table, err := NewTable(ids, amounts, notes)
	if err != nil {
		t.Fatal(err)
	}
	for _, size := range []int{1, 4, 16, 1000} {
		scanner, err := NewScanner(table, size)
		if err != nil {
			t.Fatal(err)
		}
		for _, bounds := range [][2]int{{-10, 10}, {-25, 25}, {0, 0}, {100, 200}, {5, -5}} {
			var wantIDs []int
			var wantSum int64
			for i, id := range ids {
				if id >= bounds[0] && id <= bounds[1] {
					wantIDs = append(wantIDs, id)
					wantSum += amounts[i]
				}
			}
			r := scanner.Scan(bounds[0], bounds[1], nil)
			if !slices.Equal(r.IDs, wantIDs) || r.Sum != wantSum {
				t.Fatalf("batch size=%d bounds=%v: differs from scalar scan", size, bounds)
			}
			wantBatches, wantRead := 0, 0
			if bounds[0] <= bounds[1] {
				wantBatches = (len(ids) + size - 1) / size
				wantRead = len(ids)
			}
			if r.Batches != wantBatches || r.FilterValuesRead != wantRead || r.AmountValuesRead != len(wantIDs) {
				t.Fatalf("wrong work accounting: %+v", r)
			}
		}
	}
}

func TestValidationAndOwnership(t *testing.T) {
	if _, err := NewTable([]int{1}, nil, nil); err == nil {
		t.Fatal("accepted unaligned columns")
	}
	if _, err := NewScanner(nil, 4); err == nil {
		t.Fatal("accepted nil table")
	}
	table, _ := NewTable(nil, nil, nil)
	for _, size := range []int{0, -1} {
		if _, err := NewScanner(table, size); err == nil {
			t.Fatal("accepted invalid batch size")
		}
	}
	scanner, _ := NewScanner(table, 4)
	if r := scanner.Scan(0, 10, nil); r.Batches != 0 || r.Sum != 0 {
		t.Fatal("incorrect empty scan")
	}
	ids, amounts, notes := []int{1, 2}, []int64{10, 20}, []string{"a", "b"}
	table, _ = NewTable(ids, amounts, notes)
	ids[0], amounts[0], notes[0] = 99, 999, "changed"
	scanner, _ = NewScanner(table, 4)
	r := scanner.Scan(1, 2, nil)
	if !slices.Equal(r.IDs, []int{1, 2}) || r.Sum != 30 {
		t.Fatal("external mutation changed table")
	}
}
