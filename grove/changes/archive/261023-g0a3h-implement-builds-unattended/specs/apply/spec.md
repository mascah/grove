# Changes to apply

The maintained spec is [grove/specs/apply/spec.md](../../../../specs/apply/spec.md).
This change retires the capability: every requirement is removed, its rules
move into `implement` ([../implement/spec.md](../implement/spec.md)), and
archiving this change deletes `grove/specs/apply/`. "Verify before
finishing a member" is added by change 261023-m2she, which this change
depends on.

## REMOVED Requirements

### Requirement: Apply a cycle

Apply is removed; `/grove-implement <cycle-id> --unattended` builds a cycle
asking no one, as implement's "Build unattended" says.

### Requirement: Apply one member of a cycle

Apply is removed; a member gets its cycle's plan from `grove instructions
implement <change-id>`, with or without `--unattended`, as implement's "A
member gets its cycle's plan" says.

### Requirement: Refuse a member that cannot be applied

Apply is removed. With no cycle argument there is no change outside the
cycle to refuse, and a change that waits is refused as implement's
"Implement an approved change" and "Build unattended" say.

### Requirement: Apply takes only a cycle

Apply is removed; `--unattended` takes a change ID as well as a cycle ID.

### Requirement: Verify before finishing a member

Apply is removed; implement's "Verify before finishing" covers a change
built with `--unattended`.
