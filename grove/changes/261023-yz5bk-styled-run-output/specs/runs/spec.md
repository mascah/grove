# Changes to runs

The changes that change 261023-yz5bk proposes to the requirements of
`runs` (`grove/specs/runs/spec.md` at 2f8f44d).

## MODIFIED Requirements

### Requirement: Show the session as it runs

While the session runs, the system shall write every event the session
emits to the attempt's output file as JSON Lines, and print to standard
output, as it happens, one row per shown event: the time Grove read it,
then a mark for its kind and its text, `▶` the session's start with its
model, `•` the first line of an assistant message, `›` a tool call with
the first line of its input, `✗` the first line of a tool result that is
an error, `·` thinking and other notices, and `■` the session's result
with its turn count and cost; a row identical to the one before it shall
be counted on that row as `(×N)` rather than printed again. Where standard
output is a terminal, rows shall be colored by kind and cut to its width,
and a counted row shall update in place; elsewhere, rows shall be plain
text, and a counted row printed once a different row follows or the
session ends.

#### Scenario: Watching a run

- GIVEN `grove run` launched a session on a cycle, in a terminal
- WHEN the agent runs a tool that succeeds
- THEN, before the session ends, standard output shows a `›` row naming
  the tool and the first line of its input, and no row for its result

#### Scenario: A tool fails

- GIVEN a running session
- WHEN a tool returns an error
- THEN standard output shows a `✗` row with the error's first line

#### Scenario: Thinking in a row

- GIVEN a running session in a terminal
- WHEN the agent thinks three times with nothing shown between
- THEN standard output shows one `·` thinking row ending `(×3)`

#### Scenario: A stopped session

- GIVEN a session that made several tool calls was stopped at the time
  limit
- WHEN the operator reads the attempt's output file
- THEN it holds each event up to the stop, one JSON object per line

#### Scenario: Started by a timer

- GIVEN launchd started `grove run` with standard output sent to a file
- WHEN the session runs
- THEN that file receives the same rows without color, each repeat counted
  on one line, and no row cut short

## Observed current behavior

- `relay` in `internal/cli/run.go` at 2f8f44d prints the first line of each
  assistant text, each tool call as `<name>: <first line of input>`, and
  each tool result's first line indented two spaces, with no time, mark or
  color, the same whether or not standard output is a terminal; thinking
  blocks, system events and the result event are written only to the
  `.log`.
