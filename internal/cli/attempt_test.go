package cli

import (
	"bytes"
	"context"
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
		{[]string{"run", "--budget", "1", "--permission-mode", "auto"}, 2, "run requires exactly one work ID"},
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
