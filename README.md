# Learn B-trees with Go

Run `go run . btree` to see the tree after every insertion, followed by search traces.
Run `go test ./...` to check ordering, node occupancy, child ranges, equal leaf
depth, duplicate handling, and successful searches after every insertion.
The CLI uses Cobra for subcommands. Run `go run . --help` to list commands,
or `go run . btree --help` for help with the B-tree demonstration.

`main.go` starts the CLI; `root.go` registers subcommands;
`btree_command.go` runs the existing demo; `btree.go` contains the tree itself.
Add future demonstrations as separate commands through `root.AddCommand` in
`newRootCommand`.

## Start with Search

A `node` is an interface implemented by `*leaf` and `*branch`. Both contain
sorted `keys`, but only a branch contains `children []node`. A branch has one
more child than keys. Children can themselves be leaves or branches, allowing
the tree to grow to any height. The root starts as a leaf and becomes a branch
when it splits. Splitting preserves the child type and promotes the middle key
into the parent; this remains a classic B-tree.
For example, `[10 20]` has three children: values below 10, values between 10
and 20, and values above 20. Search either finds a key in the current node or
chooses exactly one child. A missing key stops at a leaf.

## Then follow Insert

We use minimum degree 2. Each node holds at most 3 keys; each non-root node
holds at least 1. All leaves have the same depth.

1. Insert 10, 20, 5: the root becomes `[5 10 20]`.
2. Before inserting 6, split that full root. Promote 10 into a new root, leaving
   `[5]` and `[20]` as children. Insert 6 into the left leaf: `[5 6]`.
3. Continue inserting 12, 30, 7, 17. Before descending into a full child,
   split it and choose the appropriate resulting child.

The final tree is:

```text
[10 20]
  [5 6 7]
  [12 17]
  [30]
```

Search for 17 visits `[10 20]`, then `[12 17]`. Search for 15 follows the same
path but fails at the leaf. Search for 10 succeeds at the root.

Try adding 8 to the insertion sequence. Predict which leaf splits and which
key moves into the root before running it.

## Connection to database indexes

This is an in-memory classic B-tree storing unique integers. In a classic
B-tree, keys can be found in internal nodes, as the search for 10 demonstrates.
Many database indexes use a B+ tree: internal keys route searches, while leaf
entries carry row references or records, and linked leaves support range scans.
Real nodes correspond to pages holding many more keys, keeping the tree shallow.

Search and insertion take O(log n) time with this fixed node capacity. The
duplicate check adds a second search to insertion but preserves that bound.
Deletion, persistence, row references, and concurrency are left out so the
search and split algorithms remain easy to read.

## B+ tree

Run `go run . bplus` to compare the same insertion sequence with a B+ tree.
Use `go run . bplus --help` for command help. The implementation is in
`bplus_tree.go`, and its demonstration is in `bplus_command.go`.

This example stores unique integer keys, with up to 3 keys per node. It has
separate `bplusLeaf` and `bplusBranch` types behind a `bplusNode` interface:

- Stored keys live in leaves. Branch separators are copies used for navigation.
- A separator equals the smallest key in its right child's subtree. Equality
  follows the right child; search always finishes at a leaf.
- Insertion splits a node after it overflows. A leaf split copies the first key
  of its new right sibling into the parent, keeping that key in the leaf.
- A branch split moves its middle separator into the parent. Splits propagate
  upward and can create a new root. All leaves remain at the same depth.
- Leaves link to their next neighbor. An inclusive range scan descends once,
  then follows these links, stopping when it passes the upper bound.

The demo finishes with:

```text
branch [10 20]
  leaf [5 6 7]
  leaf [10 12 17]
  leaf [20 30]
Leaf chain: [5 6 7] -> [10 12 17] -> [20 30] -> nil
```

Notice that 10 and 20 remain in leaves as well as appearing in the branch.
Searching for 10 descends to `[10 12 17]`. Scanning `[6, 20]` returns
`[6 7 10 12 17 20]` by following three leaves.

The zero-value `BPlusTree` supports `Insert(key)`, `Search(key, trace)`,
`Range(low, high, trace)`, and `Print(writer)`. Pass `nil` to suppress tracing.
Duplicate inserts are ignored; empty or reversed ranges return no keys.
Search and insertion take O(log n) time; a range scan takes O(log n + k) for
k returned keys with this fixed capacity. Deletion, payloads, persistence, and
concurrency are left out, matching the scope of the classic B-tree example.

Tests check ordering, exact separator values, occupancy, equal leaf depth,
linked-leaf order, duplicate handling, range results, and extreme integer keys.
Ascending, descending, and shuffled inserts exercise splits through multiple levels.
