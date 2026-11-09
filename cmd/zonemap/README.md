# Zone maps: min/max segment pruning

A zone map summarizes a segment of column values with its minimum and maximum.
It answers whether a segment might match a filter, allowing a scan to skip
segments that cannot match. It does not locate individual matching rows.

## How it works

1. Split the input column into segments and compute each segment's min/max.
2. For an inclusive query `[low, high]`, skip a segment when `max < low` or
   `min > high`. Disjoint ranges cannot contain a match.
3. If ranges overlap, scan the actual values and apply the exact predicate.
   For example, `[1, 100]` overlaps a query for 50 but contains no matching value.

Ordered data tends to make segment ranges narrow, so more segments can be
skipped. Mixed values can make every segment's min/max overlap the filter.
Building summaries takes O(n); scanning takes O(s + r), for s segment summaries
and r values in surviving segments. The worst case still examines every value.

## Run

```sh
go run . zonemap
```

The demo scans the same 12 values in ordered and mixed layouts, using four-value
segments. For `[10,12]`, the ordered layout skips two segments and checks four
values; the mixed layout checks all 12. A final false-positive example shows
why summary overlap never guarantees a matching row. Counts measure logical
work, not elapsed time or actual disk reads.

The `column.New(values, segmentSize)` constructor copies inputs and rejects
nonpositive sizes. `Column.Scan(low, high, trace)` returns matching values in
storage order and counts of skipped/scanned segments and examined values.
Empty data and reversed ranges produce no matches. Pass nil to suppress tracing.

## Connection to DuckDB

DuckDB automatically maintains min/max indexes for general-purpose column data
types. This in-memory example illustrates the pruning rule with immutable integer
segments; it does not reproduce DuckDB row-group sizes, storage, NULL semantics,
updates, compression, or predicate-pushdown implementation.

## Code and further reading

- [Command](command.go), [walkthrough](internal/demo/demo.go), [column and segment logic](internal/column/column.go).
- [DuckDB indexing guide](https://duckdb.org/docs/current/guides/performance/indexing).
- [DuckDB index types](https://duckdb.org/docs/current/sql/indexes).
- [Compare with column scans](../columnscan/README.md).
- [Project overview](../../README.md).
