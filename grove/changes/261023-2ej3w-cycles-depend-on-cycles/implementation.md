# Implementation: A cycle can depend on cycles

<!-- Progress on change 261023-2ej3w, written by whoever implements it. It is not
part of what was approved and does not amend it. If the owner drops the change,
write "Status: dropped" here with their reason. -->

## Tasks

### T-01 — Read a cycle's depends_on and check it

Status: done

Built on main@9308040, as the owner agreed, though every named file moved
since base ab885af (#63–#67). `doc.Load` reads a cycle's optional
`depends_on` into `Meta.DependsOn`; `Cycle.WaitsOn` holds what it names that
is not a cycle done; `Collection.follows`, beside `dependencies`, reports a
`depends_on` naming no cycle and a loop. `grove changes` adds `, waits on
cycle <id>`; the project index adds `follows <id>` under a cycle.
`TestCycles` covers each acceptance line: `go test ./...` passed.

### T-02 — Run waits on the cycles a cycle depends on

Status: not started

### T-03 — Plan asks, and the docs say it

Status: not started

## Scenarios

## Reconciliation
