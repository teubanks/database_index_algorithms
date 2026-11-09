# Columnar scans: batches, filtering, and projection

A columnar scan processes aligned arrays of column values in batches. It is an
access/execution method, not an additional index. Position i represents the same
row across all columns, so selected positions can be reused for projection.

## How it works

1. Store `id`, `amount`, and `notes` in separate, equal-length columns.
2. Read a batch of IDs and evaluate `id BETWEEN low AND high`.
3. Record matching positions in a selection vector, local to that batch.
4. Read amounts only at selected positions, collect matching IDs, and accumulate
   `SUM(amount)`. Do not read the unrelated notes column during the scan.
5. Repeat for all batches, including a partial final batch.

The filter still examines all n IDs: this example has no index or segment
pruning. Scanning takes O(n) time plus output space. Batching organizes execution
and column projection avoids unused values; it is not itself SIMD or a promise
of a speedup. Zone maps could prune segments before this stage in a larger engine.

## Run

```sh
go run . columnscan
```

The demo uses batches of four over ten rows. Filtering IDs 3 through 7 produces
`[3 4 5 6 7]` and a sum of 250, after three batches, ten ID predicate checks,
and five amount reads. Printed selections show how batch-local positions map
back to aligned rows. Work counts describe this algorithm, not physical I/O.

`scan.NewTable(ids, amounts, notes)` checks alignment and copies inputs.
`scan.NewScanner(table, batchSize)` rejects nil tables and nonpositive sizes.
`Scanner.Scan(low, high, trace)` returns selected IDs, an int64 sum, and work
counts. Pass nil for quiet execution. Empty data and reversed ranges yield no
rows. Sums assume the result fits int64.

## Connection to DuckDB

DuckDB's execution engine operates on vectors; its documented default vector
size is 2048 tuples. Its data chunks collect vectors for multiple columns.
This example uses four rows so selections are readable. It models only flat
arrays, filtering, projection, and aggregation, without SQL planning, SIMD,
compressed vectors, NULLs, parallelism, or DuckDB integration.

## Code and further reading

- [Command](command.go), [walkthrough](internal/demo/demo.go).
- [Aligned column storage](internal/scan/table.go), [scanner](internal/scan/scanner.go).
- [DuckDB execution format](https://duckdb.org/docs/current/internals/vector).
- [Compare with zone-map pruning](../zonemap/README.md).
- [Project overview](../../README.md).
