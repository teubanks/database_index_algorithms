# B-tree

A B-tree is a balanced search tree whose nodes hold several sorted keys. Each
key separates child subtrees into ordered ranges, so a search chooses one child
at a time. Keys are stored in both branches and leaves, so a search can succeed
before reaching a leaf.

When inserting, this implementation splits full nodes before descending into
them. The middle key moves into the parent, while smaller and larger keys stay
in separate children. Splitting the root creates a new level; every leaf stays
at the same depth. With the fixed node capacity used here, search and insertion
take O(log n) time.

## Run the example

From the project root:

```sh
go run . btree
```

The walkthrough uses at most three keys per node so each split is easy to see.
Duplicate inserts are ignored.

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

## Code and further reading

- [Command](command.go): Cobra entry point.
- [Walkthrough](internal/demo/demo.go): insertion sequence and search examples.
- [Tree implementation](internal/tree/tree.go): public tree operations.
- [Leaf](internal/tree/leaf.go) and [branch](internal/tree/branch.go): insertion and split behavior.
- [Grokipedia: B-tree](https://grokipedia.com/page/B-tree) — unverified article link; site access was blocked during documentation.
- [Compare with the B+ tree](../bplus/README.md).
- [Back to the project overview](../../README.md).
