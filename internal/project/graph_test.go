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
		{"missing", record("G-001", "work", "depends_on: [G-003]\n"), record("G-002", "work", ""), "depends_on"},
		{"self", record("G-001", "work", "relates_to: [G-001]\n"), record("G-002", "work", ""), "self"},
		{"dependency cycle", record("G-001", "work", "depends_on: [G-002]\n"), record("G-002", "work", "depends_on: [G-001]\n"), "depends_on: cycle"},
		{"member cycle", record("G-001", "work", "members: [G-002]\n"), record("G-002", "work", "members: [G-001]\n"), "members: cycle"},
		{"duplicate identity", record("G-001", "work", ""), record("G-001", "work", ""), "duplicate"},
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
	put(t, root, "grove/work/a.md", record("G-001", "work", "depends_on: [G-003]\nmembers: [G-002]\n"))
	put(t, root, "grove/questions/q.md", record("G-002", "question", "blocks: [G-003]\n"))
	put(t, root, "grove/decisions/d.md", record("G-003", "decision", "relates_to: [G-999]\n"))
	_, ds := Load(root, "")
	msg := diagnostics(ds)
	for _, field := range []string{"depends_on", "members", "blocks", "relates_to"} {
		if !strings.Contains(msg, field) {
			t.Fatalf("missing %s diagnostic: %s", field, msg)
		}
	}
	put(t, root, "grove/work/b.md", record("G-004", "work", ""))
	put(t, root, "grove/work/c.md", record("G-004", "work", ""))
	put(t, root, "grove/questions/q.md", record("G-002", "question", "blocks: [G-004]\n"))
	_, ds = Load(root, "")
	if !strings.Contains(diagnostics(ds), "ambiguous") {
		t.Fatal(diagnostics(ds))
	}
}

func TestMembershipIsSeparateFromDependencyOrdering(t *testing.T) {
	t.Parallel()
	root := fixture(t)
	put(t, root, "grove/work/a.md", record("G-001", "work", "members: [G-002, G-003]\ndepends_on: [G-002]\n"))
	put(t, root, "grove/work/b.md", record("G-002", "work", "members: [G-003]\n"))
	put(t, root, "grove/work/c.md", record("G-003", "work", "depends_on: [G-001]\n"))
	put(t, root, "grove/questions/q.md", record("G-004", "question", "blocks: [G-001]\nrelates_to: [G-005]\n"))
	put(t, root, "grove/decisions/d.md", record("G-005", "decision", "relates_to: [G-004]\n"))
	_, ds := Load(root, "")
	if len(ds) != 0 {
		t.Fatal(diagnostics(ds))
	}
}

func TestDuplicateIDsAcrossDifferentTypes(t *testing.T) {
	t.Parallel()
	root := fixture(t)
	put(t, root, "grove/work/work.md", record("G-001", "work", ""))
	put(t, root, "grove/questions/question.md", record("G-001", "question", ""))
	_, ds := Load(root, "")
	message := diagnostics(ds)
	if !strings.Contains(message, "duplicate G-001 in grove/questions/question.md, grove/work/work.md") {
		t.Fatalf("duplicate identity must be reported across differing types: %s", message)
	}
}
