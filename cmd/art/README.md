# ART: a walkthrough of the implementation

An Adaptive Radix Tree routes a key by its bytes. It compresses shared byte
paths into prefixes and changes its child layout as a branch gets more children.
This walkthrough follows the actual files and the demo's insertion order:
`car`, `cart`, `carbon`, `cat`, then `dog`.

Run from the project root:

```sh
go run . art
```

Read the sections in order, opening the linked file beside each snippet.
All Go snippets below are excerpts from this implementation, not pseudocode.

## 1. Start with the caller

Source: [demo.go — Demo.Run](internal/demo/demo.go).

The command creates a demo, and the demo owns the sequence of operations.
The algorithm itself lives in the nested `tree` package. This loop is the
entry point for the changes we will follow:

```go
	var index tree.Tree
	fmt.Fprintln(d.out, "ART: compressed byte paths and adaptive child layouts.")
	for _, key := range []string{"car", "cart", "carbon", "cat", "dog"} {
		fmt.Fprintf(d.out, "\nInsert %q: %v\n", key, index.Insert(key))
		index.Print(d.out)
	}
```

`Insert` changes the tree; `Print` makes the result visible. Start by tracing
the first two iterations, then return to the loop after reading the split logic.

The CLI wiring is in [command.go](command.go); it delegates to this demo without
implementing tree behavior.

## 2. A tree points to a node; nodes can replace themselves

Source: [tree.go — Tree and node](internal/tree/tree.go).

The tree owns its root and the count of stored keys:

```go
type Tree struct {
	root node
	size int
}
```

Both leaves and branches implement this private interface:

```go
type node interface {
	insert(key string, depth int) node
	lookup(key string, depth int, trace io.Writer) bool
	print(out io.Writer, indent string)
}
```

The important return value is `insert(...) node`. Inserting can turn a leaf
into a branch, or put a new branch above an existing one. Returning a node lets
the caller replace its reference with the new subtree root.

`depth` is the number of bytes of the key already consumed by the path. It is
not the number of tree levels: one compressed prefix can consume several bytes.

Source: [tree.go — Tree.Insert](internal/tree/tree.go).

```go
func (t *Tree) Insert(key string) bool {
	if t.Search(key, nil) {
		return false
	}
	if t.root == nil {
		t.root = &leaf{key: key}
	} else {
		t.root = t.root.insert(key, 0)
	}
	t.size++
	return true
}
```

The first insertion creates a leaf containing the entire key, `car`. Later
inserts delegate to the current root and keep the returned root. The early
search rejects duplicates before any mutation and keeps `size` accurate.
That makes duplicate insertion return `false`; this example does not raise a
SQL constraint error.

## 3. Inserting cart turns the car leaf into a compressed branch

Source: [leaf.go — leaf.insert](internal/tree/leaf.go).

```go
func (l *leaf) insert(key string, depth int) node {
	if key == l.key {
		return l
	}
	shared := 0
	for depth+shared < len(key) && depth+shared < len(l.key) && key[depth+shared] == l.key[depth+shared] {
		shared++
	}
	b := &branch{prefix: key[depth : depth+shared], children: &node4{}}
	b.attach(l.key, depth+shared, l)
	b.attach(key, depth+shared, &leaf{key: key})
	return b
}
```

At this point `depth` is 0. Comparing `car` and `cart` finds three shared bytes:
`car`. The new branch stores those bytes once as its prefix. It then attaches
the old leaf and the new leaf under that branch.

The returned branch replaces the old root through the assignment in `Tree.Insert`.
The shape is now:

```text
branch prefix="car"
  <end> -> leaf "car"
  't'   -> leaf "cart"
```

This is path compression: the shared path does not need one branch for `c`,
another for `a`, and another for `r`. Leaves still store full keys so a lookup
can check exact equality at the end.

## 4. An end-of-key marker lets car and cart coexist

Source: [branch.go — branch and attach](internal/tree/branch.go).

```go
type branch struct {
	prefix   string
	terminal *leaf
	children childSet
}

func (b *branch) attach(key string, depth int, child node) {
	if depth == len(key) {
		b.terminal = child.(*leaf)
		return
	}
	b.children = b.children.put(key[depth], child)
}
```

After the prefix `car` is consumed, the old key has ended. `attach` stores its
leaf in `terminal`. The new key still has a `t` byte, so it uses an outgoing edge.

The `terminal` field means we do not need a special byte value for end-of-key.
All 256 byte values remain available for edges, including `0x00`. Empty strings
and keys containing zero bytes therefore work too.

When `carbon` is inserted next, its first three bytes match `car`. A second
outgoing edge, `b`, points to the new `carbon` leaf; `terminal` still holds `car`.

## 5. Matching a branch prefix means following one byte edge

Source: [branch.go — the matching-prefix path in branch.insert](internal/tree/branch.go).

After checking the prefix, the branch continues here:

```go
	depth += len(b.prefix)
	if depth == len(key) {
		b.terminal = &leaf{key: key}
		return b
	}
	edge := key[depth]
	child := b.children.get(edge)
	if child == nil {
		child = &leaf{key: key}
	} else {
		child = child.insert(key, depth+1)
	}
	b.children = b.children.put(edge, child)
	return b
}
```

The prefix consumes `len(b.prefix)` bytes. If the key ends there, it belongs
in `terminal`. Otherwise, its next byte selects a child.

For `carbon` under prefix `car`, the byte is `b`. A missing edge gets a new
leaf. An existing edge delegates insertion into its subtree. The recursive
call uses `depth+1` because the edge itself consumes one more byte.

The final `put` matters even when an edge already existed: its old child may
have become a different subtree root. `put` can also return a larger child layout.

## 6. Inserting cat splits an existing compressed prefix

Source: [branch.go — the mismatching-prefix path in branch.insert](internal/tree/branch.go).

The beginning of the same method handles a key that diverges inside the prefix:

```go
func (b *branch) insert(key string, depth int) node {
	shared := 0
	for shared < len(b.prefix) && depth+shared < len(key) && b.prefix[shared] == key[depth+shared] {
		shared++
	}
	if shared < len(b.prefix) {
		// Split a compressed path at its first mismatching byte.
		parent := &branch{prefix: b.prefix[:shared], children: &node4{}}
		oldEdge := b.prefix[shared]
		b.prefix = b.prefix[shared+1:]
		parent.children = parent.children.put(oldEdge, b)
		parent.attach(key, depth+shared, &leaf{key: key})
		return parent
	}
```

The branch currently has prefix `car`; the new key is `cat`. Only `ca` matches,
so `shared` is 2:

1. Create a new parent with prefix `ca`.
2. Use the old prefix's next byte, `r`, as the edge to the existing branch.
3. Remove `ca` **and the consumed `r` edge** from the old branch's prefix. Its
   remaining prefix becomes empty.
4. Attach `cat` under the new parent's `t` edge.

The existing `car`, `cart`, and `carbon` leaves remain in the old subtree:

```text
branch prefix="ca"
  'r' -> branch prefix=""
           <end> -> leaf "car"
           'b'   -> leaf "carbon"
           't'   -> leaf "cart"
  't' -> leaf "cat"
```

An empty prefix is valid: that branch distinguishes keys using the next edge
or the end-of-key marker. Inserting `dog` repeats the same operation with no
shared bytes, creating a root with an empty prefix and edges `c` and `d`.

## 7. Child layouts make the radix tree adaptive

Source: [children.go — childSet](internal/tree/children.go).

Branch navigation uses this interface rather than depending on a particular
array size:

```go
type childSet interface {
	get(edge byte) node
	put(edge byte, child node) childSet
	len() int
	kind() string
	each(visit func(byte, node))
}
```

`get` finds a child by byte label; `put` inserts or replaces it. `len`, `kind`,
and `each` support inspection and printing. As with `node.insert`, the return
value of `put` allows the implementation to replace itself.

Source: [node4.go — compact storage and growth](internal/tree/node4.go).

```go
type node4 struct {
	keys     [4]byte
	children [4]node
	count    int
}
```

Only the first `count` positions are occupied. `keys[i]` identifies the byte
edge whose child is in `children[i]`. Labels stay sorted, and `get` searches
these occupied positions.

Inserting a fifth distinct edge reaches this part of `node4.put`:

```go
	if n.count == len(n.keys) {
		grown := &node16{}
		n.each(func(edge byte, child node) { grown.put(edge, child) })
		return grown.put(edge, child)
	}
```

It allocates Node16, copies the four existing edges, and inserts the new edge.
Replacing the child of an existing edge happens earlier in the method and
does not increase `count` or trigger growth.

The same pattern in [node16.go](internal/tree/node16.go) grows Node16 into
Node48 at the seventeenth distinct edge. Node48 grows into Node256 at the
forty-ninth. These thresholds count outgoing byte edges on one branch, not
total keys in the tree; the branch's `terminal` does not consume an edge slot.

| Layout | Representation | Growth point |
| --- | --- | --- |
| Node4 | Four byte labels alongside four child slots | Fifth edge |
| Node16 | Sixteen byte labels alongside sixteen child slots | Seventeenth edge |
| Node48 | Byte-to-slot lookup table plus 48 child slots | Forty-ninth edge |
| Node256 | Directly addressable slot for every byte | Maximum 256 edges |

## 8. Node48 and Node256 trade space for direct addressing

Source: [children.go — node48 and get](internal/tree/children.go).

```go
type node48 struct {
	index    [256]uint8
	children [48]node
	count    int
}

func (n *node48) get(edge byte) node {
	if slot := n.index[edge]; slot != 0 {
		return n.children[int(slot)-1]
	}
	return nil
}
```

Node48 has an entry for every possible byte but only 48 child slots.
`index[edge]` stores **slot + 1**: zero means absent, one refers to child slot
zero, and so on. That avoids confusing a missing edge with the first slot.

For example, if `index[0x74]` is 3, the `t` edge points to `children[2]`.

Source: [children.go — node256](internal/tree/children.go).

```go
type node256 struct {
	children [256]node
	count    int
}

func (n *node256) get(edge byte) node { return n.children[edge] }
```

Node256 no longer needs a byte-to-slot table: the byte is the child index.
Sparse branches use smaller layouts; densely populated branches use this
direct representation. The demo's final loop grows a single branch through
all four layouts so you can see when each representation takes over.

## 9. Search checks prefixes, edges, and finally the full key

Source: [tree.go — Tree.Search](internal/tree/tree.go).

The public method handles an empty tree and optional trace output, then delegates
to the root's `lookup(key, 0, trace)` method.

Source: [branch.go — branch.lookup](internal/tree/branch.go).

```go
func (b *branch) lookup(key string, depth int, trace io.Writer) bool {
	fmt.Fprintf(trace, "%s prefix=%q: ", b.children.kind(), b.prefix)
	if !strings.HasPrefix(key[depth:], b.prefix) {
		fmt.Fprintln(trace, "prefix mismatch; not found")
		return false
	}
	depth += len(b.prefix)
	if depth == len(key) {
		fmt.Fprintln(trace, "check end-of-key")
		if b.terminal == nil {
			return false
		}
		return b.terminal.lookup(key, depth, trace)
	}
	edge := key[depth]
	fmt.Fprintf(trace, "follow byte 0x%02x\n", edge)
	child := b.children.get(edge)
	if child == nil {
		fmt.Fprintln(trace, "missing edge; not found")
		return false
	}
	return child.lookup(key, depth+1, trace)
}
```

A prefix mismatch fails immediately. Otherwise, lookup consumes the prefix
and either checks `terminal` for an ended key or follows its next byte edge.
If there is no matching edge, the key is absent.

Source: [leaf.go — leaf.lookup](internal/tree/leaf.go).

```go
func (l *leaf) lookup(key string, depth int, trace io.Writer) bool {
	found := l.key == key
	if found {
		fmt.Fprintf(trace, "leaf %q: found\n", l.key)
	} else {
		fmt.Fprintf(trace, "leaf %q: not found\n", l.key)
	}
	return found
}
```

Reaching a leaf is not enough to declare success: the stored full key must
equal the requested key. An unexpanded suffix may still differ.

After all five demo inserts, searching for `cart` takes this path:

| Current branch prefix | Depth on entry | Bytes consumed here | Next step |
| --- | --- | --- | --- |
| `""` | 0 | Edge `c` | Enter next branch at depth 1 |
| `"a"` | 1 | Prefix `a`, then edge `r` | Enter next branch at depth 3 |
| `""` | 3 | Edge `t` | Reach leaf at depth 4 |

The leaf confirms `cart == cart`. Searching for `car` follows the same first
two edges, but ends at depth 3 and checks `terminal` instead. Searching for
`cab` fails because the branch with prefix `a` has no `b` edge.

## 10. Read the tests as a list of invariants

Source: [tree_test.go](internal/tree/tree_test.go).

- `TestPrefixSplitsAndExactKeys` inserts prefix-related keys in different orders,
  including empty strings, zero bytes, and UTF-8 byte sequences.
- `TestAllAdaptiveLayouts` inserts every possible suffix byte and checks growth
  through Node4, Node16, Node48, and Node256 without losing keys.
- `TestRandomKeysAgainstSet` compares insertion and lookup with an independent
  Go map of expected keys.
- The private `validate` helper checks that leaf keys agree with their paths,
  terminal keys end at the right prefix, edges are ordered and reachable,
  fanout counts fit the layout, and every expected key remains searchable.

Run the algorithm tests directly:

```sh
go test ./cmd/art/internal/tree
```

Try pausing the demo loop after `cart` and after `cat`, then predict the next
printed tree before continuing. To inspect growth, follow `put` in
[node4.go](internal/tree/node4.go), then [node16.go](internal/tree/node16.go),
then [children.go](internal/tree/children.go).

## Scope and further reading

For k-byte keys, insertion and exact lookup take O(k) time with this fixed
256-byte alphabet. Duplicate checks add a preliminary search but keep the
same bound. Go strings are treated as bytes rather than locale-aware text.

This is an in-memory string set for learning ART's mechanics. It omits deletion,
row references, persistence, SQL constraints, and SIMD-specific optimizations.
It does not reproduce DuckDB's memory layout or query optimizer.

- [Original ART paper](https://db.in.tum.de/~leis/papers/ART.pdf).
- [DuckDB's ART explanation](https://duckdb.org/2022/07/27/art-storage) — historical implementation discussion.
- [DuckDB index documentation](https://duckdb.org/docs/current/sql/indexes) — current database behavior.
- [Project overview](../../README.md).
