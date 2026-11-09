# Learn B-trees with Go

Small Go implementations and terminal walkthroughs for learning how balanced
search trees, indexes, and column-oriented access methods work. Tree examples
show insertion and lookup; scan examples expose filtering and projection work.

## Algorithms

| Algorithm | High-level idea | Run | Overview and further reading |
| --- | --- | --- | --- |
| B-tree | Store keys in leaves and branches; split full nodes to keep the tree balanced. | `go run . btree` | [B-tree README](cmd/btree/README.md) |
| B+ tree | Store keys in linked leaves; use branch separators for navigation and leaf links for range scans. | `go run . bplus` | [B+ tree README](cmd/bplus/README.md) |
| ART | Compress byte-key paths and adapt child layouts to fanout. | `go run . art` | [ART README](cmd/art/README.md) |
| Zone maps | Use segment min/max summaries to skip impossible matches. | `go run . zonemap` | [Zone-map README](cmd/zonemap/README.md) |
| Columnar scans | Filter batches of aligned column values and project matching positions. | `go run . columnscan` | [Column-scan README](cmd/columnscan/README.md) |

Each README explains its operations, shows how to run the example, and links to
further reading. The B-tree Grokipedia links remain unverified; the DuckDB-related
examples use primary documentation and the original ART paper.

DuckDB distinguishes ART and zone maps as index types. Columnar scanning is an
execution method. Explore these examples in the order `art`, `zonemap`, then
`columnscan`: find individual keys, skip irrelevant regions, then process columns
in batches. These are standalone Go models, not a DuckDB implementation.
See [DuckDB's indexes](https://duckdb.org/docs/current/sql/indexes) and
[execution format](https://duckdb.org/docs/current/internals/vector).

Run `go run . --help` to list commands, or add `--help` after a subcommand for
its help. These are learning examples: deletion, persistence, SQL processing, and
concurrency are outside their current scope.

## Project structure

```text
main.go                         # Thin executable entry point
cmd/
  root.go                       # Root command and registration
  root_test.go                  # CLI dispatch and behavior tests
  art/                          # ART: internal/demo and internal/tree
  zonemap/                      # Zone maps: internal/demo and internal/column
  columnscan/                   # Batch scans: internal/demo and internal/scan
  btree/
    command.go                  # Cobra command factory
    internal/
      demo/                     # Stateful B-tree walkthrough
      tree/                     # B-tree, nodes, algorithms, invariant tests
  bplus/
    command.go                  # Cobra command factory
    internal/
      demo/                     # Stateful B+ tree walkthrough
      tree/                     # B+ tree, nodes, algorithms, invariant tests
```

Each command package exposes `New() *cobra.Command`. Its demo owns the output
writer and fresh example data for each run. Nested `internal` packages keep the example's
implementation private to its command. Tree packages have no Cobra dependency;
tree types are ready to use at the zero value, while scan constructors validate
alignment and positive batch/segment sizes. Types own their state and algorithms.

To add an example, create `cmd/<name>/command.go`, place its walkthrough and
isolated logic in focused nested packages, then register `<name>.New()` in
`cmd/root.go`. Put tests beside the package they exercise. Follow `AGENTS.md`
for the project's Go adaptation of SOLID and type/package ownership rules.

Run `go test ./...` and `go vet ./...` from the repository root.
