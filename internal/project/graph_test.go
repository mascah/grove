package project

import (
	"strings"
	"testing"
)

func TestGraphErrors(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, first, second, want string
	}{
		{"missing", record("W-001", "work", "depends_on: [W-003]\n"), record("W-002", "work", ""), "depends_on"},
		{"self", record("W-001", "work", "relates_to: [W-001]\n"), record("W-002", "work", ""), "self"},
		{"dependency cycle", record("W-001", "work", "depends_on: [W-002]\n"), record("W-002", "work", "depends_on: [W-001]\n"), "depends_on: cycle"},
		{"member cycle", record("W-001", "work", "members: [W-002]\n"), record("W-002", "work", "members: [W-001]\n"), "members: cycle"},
		{"duplicate identity", record("W-001", "work", ""), record("W-001", "work", ""), "duplicate"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := fixture(t)
			put(t, root, "grove/work/first.md", tc.first)
			put(t, root, "grove/work/second.md", tc.second)
			_, ds := Load(root, "")
			if !strings.Contains(diagnostics(ds), tc.want) {
				t.Fatalf("wanted %s, got %s", tc.want, diagnostics(ds))
			}
			if tc.name == "duplicate identity" && (!strings.Contains(diagnostics(ds), "first.md") || !strings.Contains(diagnostics(ds), "second.md")) {
				t.Fatal("duplicate diagnostic must identify both files")
			}
		})
	}
}

func TestTypedTargetsAndAmbiguity(t *testing.T) {
	t.Parallel()
	root := fixture(t)
	put(t, root, "grove/work/a.md", record("W-001", "work", "depends_on: [D-001]\nmembers: [Q-001]\n"))
	put(t, root, "grove/questions/q.md", record("Q-001", "question", "blocks: [D-001]\n"))
	put(t, root, "grove/decisions/d.md", record("D-001", "decision", "relates_to: [W-999]\n"))
	_, ds := Load(root, "")
	msg := diagnostics(ds)
	for _, field := range []string{"depends_on", "members", "blocks", "relates_to"} {
		if !strings.Contains(msg, field) {
			t.Fatalf("missing %s diagnostic: %s", field, msg)
		}
	}
	put(t, root, "grove/work/b.md", record("W-002", "work", ""))
	put(t, root, "grove/work/c.md", record("W-002", "work", ""))
	put(t, root, "grove/questions/q.md", record("Q-001", "question", "blocks: [W-002]\n"))
	_, ds = Load(root, "")
	if !strings.Contains(diagnostics(ds), "ambiguous") {
		t.Fatal(diagnostics(ds))
	}
}

func TestMembershipIsSeparateFromDependencyOrdering(t *testing.T) {
	t.Parallel()
	root := fixture(t)
	put(t, root, "grove/work/a.md", record("W-001", "work", "members: [W-002, W-003]\ndepends_on: [W-002]\n"))
	put(t, root, "grove/work/b.md", record("W-002", "work", "members: [W-003]\n"))
	put(t, root, "grove/work/c.md", record("W-003", "work", "depends_on: [W-001]\n"))
	put(t, root, "grove/questions/q.md", record("Q-001", "question", "blocks: [W-001]\nrelates_to: [D-001]\n"))
	put(t, root, "grove/decisions/d.md", record("D-001", "decision", "relates_to: [Q-001]\n"))
	_, ds := Load(root, "")
	if len(ds) != 0 {
		t.Fatal(diagnostics(ds))
	}
}

func TestDuplicateIDsAreReportedEvenWhenARecordHasWrongType(t *testing.T) {
	t.Parallel()
	root := fixture(t)
	put(t, root, "grove/work/work.md", record("W-001", "work", ""))
	put(t, root, "grove/questions/question.md", record("W-001", "question", ""))
	_, ds := Load(root, "")
	message := diagnostics(ds)
	if !strings.Contains(message, "duplicate W-001") || !strings.Contains(message, "matching type prefix") {
		t.Fatalf("both identity and type errors must be reported: %s", message)
	}
}
