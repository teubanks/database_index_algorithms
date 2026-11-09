package tree

import (
	"math/rand"
	"testing"
)

// Validate structure after every insertion, including cascaded root splits.
func TestInsertInvariants(t *testing.T) {
	for _, order := range []string{"ascending", "descending", "shuffled"} {
		t.Run(order, func(t *testing.T) {
			keys := make([]int, 300)
			for i := range keys {
				keys[i] = i - 150
			}
			if order == "descending" {
				for i, j := 0, len(keys)-1; i < j; i, j = i+1, j-1 {
					keys[i], keys[j] = keys[j], keys[i]
				}
			} else if order == "shuffled" {
				rand.New(rand.NewSource(42)).Shuffle(len(keys), func(i, j int) {
					keys[i], keys[j] = keys[j], keys[i]
				})
			}
			var tree Tree
			for index, key := range keys {
				tree.Insert(key)
				leafDepth := -1
				count := validate(t, tree.root, true, nil, nil, 0, &leafDepth)
				if count != index+1 {
					t.Fatalf("got %d keys; want %d", count, index+1)
				}
				for _, inserted := range keys[:index+1] {
					if !tree.Search(inserted, nil) {
						t.Fatalf("lost key %d", inserted)
					}
				}
			}
			tree.Insert(keys[0])
			leafDepth := -1
			if validate(t, tree.root, true, nil, nil, 0, &leafDepth) != len(keys) {
				t.Fatal("duplicate insertion changed key count")
			}
			if tree.Search(1000, nil) {
				t.Fatal("found an absent key")
			}
		})
	}
}

func TestEmpty(t *testing.T) {
	var tree Tree
	if tree.Search(0, nil) {
		t.Fatal("empty tree found a key")
	}
}

func validate(t *testing.T, n node, root bool, low, high *int, depth int, leafDepth *int) int {
	t.Helper()
	if n == nil {
		t.Fatal("nil node")
	}
	keys := n.nodeKeys()
	if len(keys) < 1 || len(keys) > maxKeys || (!root && len(keys) < degree-1) {
		t.Fatal("invalid node occupancy")
	}
	for i, key := range keys {
		if (i > 0 && keys[i-1] >= key) || (low != nil && key <= *low) || (high != nil && key >= *high) {
			t.Fatalf("keys out of order or outside child bounds: %v", keys)
		}
	}
	if _, ok := n.(*leaf); ok {
		if *leafDepth == -1 {
			*leafDepth = depth
		}
		if depth != *leafDepth {
			t.Fatal("leaves have different depths")
		}
		return len(keys)
	}
	b, ok := n.(*branch)
	if !ok {
		t.Fatalf("unexpected node type %T", n)
	}
	if len(b.children) != len(keys)+1 {
		t.Fatal("wrong number of children")
	}
	count := len(keys)
	for i, child := range b.children {
		childLow, childHigh := low, high
		if i > 0 {
			childLow = &keys[i-1]
		}
		if i < len(keys) {
			childHigh = &keys[i]
		}
		count += validate(t, child, false, childLow, childHigh, depth+1, leafDepth)
	}
	return count
}

// Root growth must preserve concrete node types, including branch-to-branch links.
func TestRootNodeTypes(t *testing.T) {
	var tree Tree
	tree.Insert(10)
	if _, ok := tree.root.(*leaf); !ok {
		t.Fatalf("first root is %T; want *leaf", tree.root)
	}
	for _, key := range []int{20, 5, 6} {
		tree.Insert(key)
	}
	root, ok := tree.root.(*branch)
	if !ok {
		t.Fatalf("split root is %T; want *branch", tree.root)
	}
	for _, child := range root.children {
		if _, ok := child.(*leaf); !ok {
			t.Fatalf("child after first split is %T; want *leaf", child)
		}
	}
	for key := 0; key < 100; key++ {
		tree.Insert(key)
	}
	root = tree.root.(*branch)
	for _, child := range root.children {
		if _, ok := child.(*branch); !ok {
			t.Fatalf("child of taller tree is %T; want *branch", child)
		}
	}
	leafDepth := -1
	if count := validate(t, tree.root, true, nil, nil, 0, &leafDepth); count != 100 {
		t.Fatalf("got %d keys; want 100", count)
	}
}
