# B+ tree

A B+ tree keeps every stored key in a leaf. Branches contain separator keys
that guide searches to the appropriate leaf. A separator can appear in a branch
and a leaf, but it counts as stored data only in the leaf. Neighboring leaves
are linked in sorted order.

Insertion descends to a leaf and splits nodes that overflow. A leaf split copies
its right sibling's first key into the parent while keeping that key in the
leaf. Branch splits move a separator upward; splits can reach the root and add
a level. All leaves remain at the same depth.

For a range query, the tree finds the starting leaf once, then follows leaf
links until it passes the upper bound. With this example's fixed node capacity,
search and insertion take O(log n), and a scan returning k keys takes
O(log n + k).

## Run the example

From the project root:

```sh
go run . bplus
```

Run `go run . bplus` to compare the same insertion sequence with a B+ tree.
Use `go run . bplus --help` for command help. The implementation is in
[internal/tree](internal/tree), and its demonstration is in [internal/demo](internal/demo).

This example stores unique integer keys, with up to 3 keys per node. It has
separate `leaf` and `branch` types behind a `node` interface:

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

The zero-value B+ `tree.Tree` supports `Insert(key)`, `Search(key, trace)`,
`Range(low, high, trace)`, and `Print(writer)`. Pass `nil` to suppress tracing.
Duplicate inserts are ignored; empty or reversed ranges return no keys.
Search and insertion take O(log n) time; a range scan takes O(log n + k) for
k returned keys with this fixed capacity. Deletion, payloads, persistence, and
concurrency are left out, matching the scope of the classic B-tree example.

Tests check ordering, exact separator values, occupancy, equal leaf depth,
linked-leaf order, duplicate handling, range results, and extreme integer keys.
Ascending, descending, and shuffled inserts exercise splits through multiple levels.

## Code and further reading

- [Command](command.go): Cobra entry point.
- [Walkthrough](internal/demo/demo.go): insertion sequence, searches, and range scan.
- [Tree implementation](internal/tree/tree.go): search, insertion, range scanning, and printing.
- [Leaf](internal/tree/leaf.go) and [branch](internal/tree/branch.go): node-specific insertion and split behavior.
- [Grokipedia: B+ tree](https://grokipedia.com/page/B%2B_tree) — unverified article link; site access was blocked during documentation.
- [Compare with the classic B-tree](../btree/README.md).
- [Back to the project overview](../../README.md).
