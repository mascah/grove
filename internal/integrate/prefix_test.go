package integrate

import (
	"strings"
	"testing"
)

// The commit the checks read is what is merged, whatever the branch's name
// resolves to meanwhile; and a project under a prefix compares the record's
// path where git diff prints it, so code above the project counts and the
// record's own commits do not.
func TestIntegrateMergesTheInspectedCommitUnderAPrefix(t *testing.T) {
	t.Parallel()
	root, _, candidate := fixture(t, true, "sub")
	git(t, root, "tag", "feature", "main") // a tag of the branch's name, which git merge feature would take
	facts, err := run(t, root, root, false)
	if err != nil || len(facts) != 3 || !strings.HasPrefix(facts[1], "merge: fast-forward main from ") || !strings.HasPrefix(facts[2], "done: G-001 done at commit ") {
		t.Fatalf("%v; facts %q", err, facts)
	}
	if r := record(t, root); r.Status != "done" || r.Candidate != candidate || r.Approved != candidate {
		t.Fatalf("%+v", r)
	}
	git(t, root, "merge-base", "--is-ancestor", candidate, "HEAD") // fails the test when the candidate is not in main
}
