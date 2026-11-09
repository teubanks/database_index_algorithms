# Project guidance

## SOLID and Go structure

- Follow SOLID using idiomatic Go. A Go package is the larger module boundary;
  structs and their methods serve the role of classes. Keep one Go module at
  the repository root unless a separate distribution boundary requires more.
- Use a struct with methods when there is identity, lifecycle, mutable state,
  an invariant, polymorphism, or several related behaviors that belong together.
- Group coherent calculations, transformations, normalization, parsing,
  mapping, and predicates on the type that owns that responsibility. A stateless
  struct is appropriate when it provides a useful, coherent abstraction.
- Use a standalone pure function only when the work does not belong to such a
  type and adding an object would provide no useful abstraction. Constructors,
  CLI factories, program entry points, and Go test functions/helpers may stand alone.
- Keep types and packages focused. Extend a feature through its owning package;
  avoid unrelated behavior, generic utility packages, and speculative frameworks.
- Use narrow interfaces where polymorphism or a real dependency boundary calls
  for them. Do not add an interface solely because a concrete type exists.
- Keep state and invariants private. Export the smallest useful API; avoid
  mutable package globals. Inject output or other external dependencies.

## Commands and nested packages

- Keep `main.go` as a thin entry point: execute the root command and report errors.
- `cmd` owns root command construction and subcommand registration only.
- Put each subcommand in `cmd/<name>` with an exported `New() *cobra.Command`.
  Nested CLI subcommands should follow the same package-per-command pattern.
- Put feature-specific logic in focused nested packages under that command's
  `internal` directory. Current examples use `internal/demo` for the walkthrough
  and `internal/tree` for tree algorithms; scans use `internal/column` or
  `internal/scan` for their focused logic.
- Keep Cobra argument/flag parsing in the command package, demonstration
  sequences in the demo package, and data structure invariants in the tree package.
- Dependencies point inward: main -> root command -> subcommand -> demo -> algorithm logic.
  Algorithm packages must not import Cobra or command/demo packages. One command must
  not import another command's private implementation.
- Move code shared by actual consumers to their nearest common owning package
  when needed. Do not share internals or create abstractions just for possible reuse.
- Use package context to simplify names: `tree.Tree`, `leaf`, and `branch` rather
  than repeating the command name on every type. Prefer zero-value-ready types.
- Keep tests beside their owning packages. Root CLI integration tests may test
  registration and dispatch through public APIs; tree tests may inspect internal
  invariants from within the tree package.

## Learning behavior and verification

- Preserve separate leaf and branch types, small visible node capacities, and
  explanatory comments around algorithm decisions and splits.
- Preserve existing command names, demo sequences, and output during structural
  refactors unless a requested behavior change requires otherwise.
- Keep examples small, in memory, and focused on the requested operation. The
  B-trees use integer sets, ART uses byte-string keys, and scans use small aligned
  columns. Do not add persistence, concurrency, or deletion without a task that
  calls for them. Document simplifications and distinguish educational models
  from actual DuckDB internals and performance.
- Run `gofmt` on changed Go files, `go test ./...`, and `go vet ./...`.
  For command or demo changes, also run the affected command and inspect output.
- Update the README when package layout, commands, or public behavior changes.

## Algorithm documentation

- Give each algorithm subcommand a `cmd/<name>/README.md` with a short overview,
  how its main operations work, a root-level run command, links to its owning
  code, and further reading (prefer the requested Grokipedia article when available).
- Link each algorithm README from the root README's algorithm table. Keep the
  root focused on navigation, comparison, project structure, and development.
- Verify external article links when possible. Explicitly mark links that could
  not be verified rather than claiming that an article exists or was reviewed.
