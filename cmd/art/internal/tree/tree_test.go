package tree

import (
	"bytes"
	"fmt"
	"math/rand"
	"strings"
	"testing"
)

func TestPrefixSplitsAndExactKeys(t *testing.T) {
	keys := []string{"abcdefgh1", "abcdefgh2", "abc", "abx", "a", "", "\x00", "\x00foo", "abc\x00def", "abc\x00", "世界", "世"}
	for seed := int64(0); seed < 10; seed++ {
		order := append([]string(nil), keys...)
		rand.New(rand.NewSource(seed)).Shuffle(len(order), func(i, j int) { order[i], order[j] = order[j], order[i] })
		var tree Tree
		want := map[string]bool{}
		if tree.Search("", nil) {
			t.Fatal("empty tree found a key")
		}
		for _, key := range order {
			if !tree.Insert(key) {
				t.Fatalf("new key rejected: %q", key)
			}
			want[key] = true
			validate(t, &tree, want)
			if tree.Insert(key) {
				t.Fatalf("duplicate accepted: %q", key)
			}
			validate(t, &tree, want)
		}
		for _, key := range []string{"abcdef", "abcdefgh", "abc\x00d", "世界x", "missing"} {
			if tree.Search(key, nil) {
				t.Fatalf("absent key found: %q", key)
			}
		}
	}
}

func TestAllAdaptiveLayouts(t *testing.T) {
	var tree Tree
	want := map[string]bool{}
	// Reverse byte order tests sorted insertion in Node4/16, not just append.
	for i := 255; i >= 0; i-- {
		key := "prefix/" + string([]byte{byte(i)})
		tree.Insert(key)
		want[key] = true
		validate(t, &tree, want)
		count := 256 - i
		if count < 2 {
			continue
		}
		kind := "Node4"
		if count > 4 {
			kind = "Node16"
		}
		if count > 16 {
			kind = "Node48"
		}
		if count > 48 {
			kind = "Node256"
		}
		if tree.root.(*branch).children.kind() != kind {
			t.Fatalf("count=%d: want %s", count, kind)
		}
	}
	// Exact prefix and longer child paths must work even with all 256 edges occupied.
	for _, key := range []string{"prefix/", "prefix/\x00/longer", "prefix/\xff/longer"} {
		tree.Insert(key)
		want[key] = true
		validate(t, &tree, want)
	}
}

func TestRandomKeysAgainstSet(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	var tree Tree
	want := map[string]bool{}
	for i := 0; i < 500; i++ {
		data := make([]byte, rng.Intn(15))
		rng.Read(data)
		key := string(data)
		inserted := tree.Insert(key)
		if inserted == want[key] {
			t.Fatalf("incorrect insertion result for %q", key)
		}
		want[key] = true
	}
	validate(t, &tree, want)
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("absent-%d", i)
		if tree.Search(key, nil) != want[key] {
			t.Fatalf("lookup mismatch: %q", key)
		}
	}
}

func validate(t *testing.T, tree *Tree, want map[string]bool) {
	t.Helper()
	got := map[string]bool{}
	var visit func(node, string)
	visit = func(n node, path string) {
		switch n := n.(type) {
		case *leaf:
			if !strings.HasPrefix(n.key, path) || got[n.key] {
				t.Fatalf("invalid/duplicate leaf %q under %q", n.key, path)
			}
			got[n.key] = true
		case *branch:
			path += n.prefix
			if n.terminal != nil {
				if n.terminal.key != path || got[path] {
					t.Fatalf("invalid terminal at %q", path)
				}
				got[path] = true
			}
			count, previous := 0, -1
			n.children.each(func(edge byte, child node) {
				if int(edge) <= previous || child == nil || n.children.get(edge) != child {
					t.Fatal("invalid child layout")
				}
				previous = int(edge)
				count++
				visit(child, path+string([]byte{edge}))
			})
			if count != n.children.len() || count == 0 {
				t.Fatal("invalid fanout")
			}
			capacity := map[string]int{"Node4": 4, "Node16": 16, "Node48": 48, "Node256": 256}[n.children.kind()]
			if count > capacity || capacity == 0 {
				t.Fatal("invalid adaptive capacity")
			}
		default:
			t.Fatalf("unexpected node %T", n)
		}
	}
	if tree.root != nil {
		visit(tree.root, "")
	}
	if len(got) != len(want) || tree.Len() != len(want) {
		t.Fatal("key count mismatch")
	}
	for key := range want {
		if !got[key] || !tree.Search(key, nil) {
			t.Fatalf("lost key %q", key)
		}
	}
}

func TestPrintAndTrace(t *testing.T) {
	var tree Tree
	var out bytes.Buffer
	tree.Print(&out)
	if out.String() != "(empty)\n" {
		t.Fatal("bad empty print")
	}
	for _, key := range []string{"car", "cart", "carbon"} {
		tree.Insert(key)
	}
	out.Reset()
	tree.Print(&out)
	if !strings.Contains(out.String(), `prefix="car"`) || !strings.Contains(out.String(), "<end>") {
		t.Fatal(out.String())
	}
	out.Reset()
	if !tree.Search("cart", &out) || !strings.Contains(out.String(), "found") {
		t.Fatal(out.String())
	}
}
