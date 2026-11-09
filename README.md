# Learn B-trees with Go

Small Go implementations and terminal walkthroughs for learning how balanced
search trees work. Each example stores unique integers and prints the changes
after insertion, followed by search traces.

## Algorithms

| Algorithm | High-level idea | Run | Overview and further reading |
| --- | --- | --- | --- |
| B-tree | Store keys in leaves and branches; split full nodes to keep the tree balanced. | `go run . btree` | [B-tree README](cmd/btree/README.md) |
| B+ tree | Store keys in linked leaves; use branch separators for navigation and leaf links for range scans. | `go run . bplus` | [B+ tree README](cmd/bplus/README.md) |

Each algorithm README explains search and insertion, shows the example's final
tree, and includes a Grokipedia link for further reading. Those external links
could not be verified from this environment.

Run `go run . --help` to list commands, or add `--help` after a subcommand for
its help. These are learning examples: deletion, persistence, payloads, and
concurrency are outside their current scope.

## Project structure

```text
main.go                         # Thin executable entry point
cmd/
  root.go                       # Root command and registration
  root_test.go                  # CLI dispatch and behavior tests
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
writer and a fresh tree for each run. Nested `internal` packages keep the example's
implementation private to its command. Tree packages have no Cobra dependency;
their exported `Tree` types are ready to use at the zero value. Leaves and
branches own their insertion and split behavior through private methods.

To add an example, create `cmd/<name>/command.go`, place its walkthrough and
isolated logic in focused nested packages, then register `<name>.New()` in
`cmd/root.go`. Put tests beside the package they exercise. Follow `AGENTS.md`
for the project's Go adaptation of SOLID and type/package ownership rules.

Run `go test ./...` and `go vet ./...` from the repository root.
