package main

import (
	"bytes"
	"math/rand"
	"slices"
	"strings"
	"testing"
)

func TestBPlusInsertInvariants(t *testing.T) {
	for _, order := range []string{"ascending", "descending", "shuffled"} {
		t.Run(order, func(t *testing.T) {
			keys := make([]int, 200)
			for i := range keys {
				keys[i] = 2*i - 200 // Leave gaps to test absent keys and range bounds.
			}
			if order == "descending" {
				slices.Reverse(keys)
			} else if order == "shuffled" {
				rand.New(rand.NewSource(42)).Shuffle(len(keys), func(i, j int) {
					keys[i], keys[j] = keys[j], keys[i]
				})
			}
			var tree BPlusTree
			var want []int
			for _, key := range keys {
				tree.Insert(key)
				want = append(want, key)
				slices.Sort(want)
				validateBPlus(t, &tree, want)
				for _, inserted := range want {
					if !tree.Search(inserted, nil) {
						t.Fatalf("lost key %d", inserted)
					}
				}
				if tree.Search(key+1, nil) {
					t.Fatalf("found absent key %d", key+1)
				}
				for _, bounds := range [][2]int{{-201, 201}, {-51, 51}, {key, key}, {key + 1, key + 1}, {5, -5}} {
					var expected []int
					for _, candidate := range want {
						if candidate >= bounds[0] && candidate <= bounds[1] {
							expected = append(expected, candidate)
						}
					}
					if got := tree.Range(bounds[0], bounds[1], nil); !slices.Equal(got, expected) {
						t.Fatalf("range %v = %v; want %v", bounds, got, expected)
					}
				}
				// A duplicate must not change leaf data, separators, or links.
				tree.Insert(key)
				validateBPlus(t, &tree, want)
			}
		})
	}
}

func validateBPlus(t *testing.T, tree *BPlusTree, want []int) {
	t.Helper()
	var leaves []*bplusLeaf
	var stored []int
	leafDepth := -1
	var visit func(bplusNode, bool, int) (int, int)
	visit = func(n bplusNode, root bool, depth int) (int, int) {
		t.Helper()
		if n == nil {
			t.Fatal("nil child")
		}
		keys := n.bplusKeys()
		if len(keys) == 0 || len(keys) > bplusMaxKeys {
			t.Fatalf("invalid occupancy: %v", keys)
		}
		for i := 1; i < len(keys); i++ {
			if keys[i-1] >= keys[i] {
				t.Fatalf("unsorted or duplicate keys: %v", keys)
			}
		}
		switch n := n.(type) {
		case *bplusLeaf:
			if !root && len(n.keys) < (bplusMaxKeys+1)/2 {
				t.Fatal("underfilled leaf")
			}
			if leafDepth == -1 {
				leafDepth = depth
			}
			if leafDepth != depth {
				t.Fatal("leaves have different depths")
			}
			leaves = append(leaves, n)
			stored = append(stored, n.keys...)
			return n.keys[0], n.keys[len(n.keys)-1]
		case *bplusBranch:
			if len(n.children) != len(n.keys)+1 {
				t.Fatal("wrong child count")
			}
			min, previousMax := visit(n.children[0], false, depth+1)
			for i, child := range n.children[1:] {
				childMin, childMax := visit(child, false, depth+1)
				if n.keys[i] != childMin || previousMax >= n.keys[i] {
					t.Fatalf("invalid separator %d: left max %d, right min %d", n.keys[i], previousMax, childMin)
				}
				previousMax = childMax
			}
			return min, previousMax
		default:
			t.Fatalf("unexpected node type %T", n)
			return 0, 0
		}
	}
	visit(tree.root, true, 0)
	if !slices.Equal(stored, want) {
		t.Fatalf("stored leaf keys = %v; want %v", stored, want)
	}
	// Compare pointer identity to tree traversal: detects skipped leaves and cycles.
	for i, l := range leaves {
		var wantNext *bplusLeaf
		if i+1 < len(leaves) {
			wantNext = leaves[i+1]
		}
		if l.next != wantNext {
			t.Fatalf("incorrect next link for leaf %v", l.keys)
		}
	}
}

func TestBPlusEmptyAndExtremeKeys(t *testing.T) {
	var tree BPlusTree
	var out bytes.Buffer
	tree.Print(&out)
	if out.String() != "(empty)\n" || tree.Search(0, nil) || len(tree.Range(-1, 1, nil)) != 0 {
		t.Fatal("incorrect empty tree behavior")
	}
	maxInt := int(^uint(0) >> 1)
	minInt := -maxInt - 1
	want := []int{minInt, -1, 0, 1, maxInt}
	for _, key := range want {
		tree.Insert(key)
	}
	validateBPlus(t, &tree, want)
	if got := tree.Range(minInt, maxInt, nil); !slices.Equal(got, want) {
		t.Fatalf("extreme range = %v; want %v", got, want)
	}
	if !tree.Search(minInt, nil) || !tree.Search(maxInt, nil) {
		t.Fatal("extreme keys not found")
	}
}

func TestBPlusSeparatorSearchAndRangeTrace(t *testing.T) {
	var tree BPlusTree
	for _, key := range []int{10, 20, 5, 6, 12, 30, 7, 17} {
		tree.Insert(key)
	}
	var out bytes.Buffer
	if !tree.Search(10, &out) || !strings.Contains(out.String(), "branch [10 20]: follow child 1") || !strings.Contains(out.String(), "leaf [10 12 17]: found") {
		t.Fatalf("separator search did not reach its leaf: %s", out.String())
	}
	out.Reset()
	if got := tree.Range(6, 20, &out); !slices.Equal(got, []int{6, 7, 10, 12, 17, 20}) {
		t.Fatalf("unexpected range result: %v", got)
	}
	if strings.Count(out.String(), "follow next leaf") != 2 {
		t.Fatalf("range did not follow the leaf chain: %s", out.String())
	}
}
