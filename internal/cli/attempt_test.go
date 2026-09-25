package cli

import (
	"bytes"
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mascah/grove/internal/repo"
)

func TestAttemptCommandsUsage(t *testing.T) {
	root := projectFixture(t)
	if out, err := repo.Command(context.Background(), root, "init", "-q").CombinedOutput(); err != nil {
		t.Fatalf("%v %s", err, out)
	}
	for _, c := range []struct {
		args []string
		code int
		want string
	}{
		{[]string{"run", "G-001"}, 2, "run requires --budget USD and --permission-mode MODE, or their defaults under run: in grove.yaml"},
		{[]string{"run", "G-001", "--budget", "1"}, 2, "run requires --budget USD and --permission-mode MODE"},
		{[]string{"run", "G-001", "--budget", "NaN", "--permission-mode", "auto"}, 2, "--budget must be a positive decimal dollar amount"},
		{[]string{"run", "G-001", "--budget", "0.0", "--permission-mode", "auto"}, 2, "--budget must be a positive decimal dollar amount"},
		{[]string{"run", "--budget", "1", "--permission-mode", "auto"}, 2, "run requires at least one work ID"},
		{[]string{"run", "G-001", "--dry-run", "--expect", "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}, 2, "give one or the other"},
		{[]string{"run", "G-001", "--expect", "sha256:abc"}, 2, "--expect must be the digest --dry-run printed"},
		{[]string{"show", "G-001", "--dry-run"}, 2, "--dry-run applies only to run"},
		{[]string{"show", "G-001", "--expect", "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}, 2, "--expect applies only to update and run"},
		{[]string{"show", "G-001", "--model", "x"}, 2, "--budget, --permission-mode, --until, --model, --effort, --branch, and --worktree apply only to run"},
		{[]string{"show", "G-001", "--effort", "high"}, 2, "apply only to run"},
		{[]string{"run", "G-001", "--budget", "1", "--permission-mode", "auto", "--until", "review"}, 2, "--until must be plan"},
		{[]string{"run", "G-001", "--budget", "1", "--permission-mode", "auto", "--effort", "high", "--effort", "low"}, 2, "--effort may only be supplied once"},
		{[]string{"attempts", "a", "b"}, 2, "attempts takes at most one work ID"},
		{[]string{"attempt"}, 2, "attempt requires exactly one attempt id"},
		{[]string{"stop", "--json", "x"}, 2, "--json applies only to"},
		{[]string{"attempt", "bogus"}, 1, "grove: bogus is not an attempt id"},
		{[]string{"stop", "G-001.20260922T183000Z"}, 1, "does not exist in this repository"},
	} {
		var out, errOut bytes.Buffer
		if code := Run(append([]string{"--project", root}, c.args...), t.TempDir(), &out, &errOut); code != c.code || !strings.Contains(errOut.String(), c.want) {
			t.Fatalf("%v: exit %d\n%s", c.args, code, errOut.String())
		}
	}
	// With run: defaults the flags are optional: the refusal comes from Start.
	write(t, root, "grove.yaml", "schema_version: 3\nrecords: docs/records\nrun: {budget: 50, permission_mode: auto}\n")
	var out, errOut bytes.Buffer
	if code := Run([]string{"--project", root, "run", "G-009"}, t.TempDir(), &out, &errOut); code != 1 || !strings.Contains(errOut.String(), "G-009 is not in this checkout") {
		t.Fatalf("exit %d\n%s", code, errOut.String())
	}
	out.Reset()
	errOut.Reset()
	if code := Run([]string{"--project", root, "attempts"}, t.TempDir(), &out, &errOut); code != 0 || out.String() != "ATTEMPT  WORK  STATUS  STARTED  BRANCH  EXIT  COST\n" {
		t.Fatalf("exit %d\n%s%s", code, out.String(), errOut.String())
	}
}

// TestRunDryRun previews a selection, writing nothing, and prints the
// digest a launch can be held to.
func TestRunDryRun(t *testing.T) {
	t.Parallel()
	root := gitFixture(t)
	write(t, root, "docs/records/work/second.md", "---\nid: G-003\ntype: work\ntitle: Build on it\nstatus: proposed\ndepends_on: [G-001]\n---\nMore.\n")
	write(t, root, ".claude/skills/grove-work/SKILL.md", "---\nname: grove-work\n---\n")
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-qm", "selection")
	var out, errOut bytes.Buffer
	if code := Run([]string{"--project", root, "run", "G-003", "G-001", "--dry-run", "--budget", "2", "--permission-mode", "auto"}, root, &out, &errOut); code != 0 {
		t.Fatalf("exit %d\n%s", code, errOut.String())
	}
	for _, want := range []string{
		"Preview: nothing was written or started",
		"Worktree: " + filepath.Join(root, ".claude", "worktrees", "worktree-G-003-G-001") + " on worktree-G-003-G-001",
		"Bounds: budget 2 USD over the whole selection",
		"Selected: G-003 G-001; order G-001, G-003",
		"Member G-003: proposed at docs/records/work/second.md",
		"Digest: sha256:",
		"Launch exactly this: grove run G-003 G-001 --expect sha256:",
	} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("missing %q in\n%s", want, out.String())
		}
	}
	if branches := gitIn(t, root, "branch", "--list", "worktree-*"); branches != "" {
		t.Fatalf("a preview made %q", branches)
	}
}
