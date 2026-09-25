package cli

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/mascah/grove/internal/deps"
)

func TestDepsUsage(t *testing.T) {
	t.Parallel()
	root := projectFixture(t)
	for _, args := range [][]string{{"deps", "--include", "x"}, {"deps", "G-001", "--interaction", "headless"}, {"deps", "--json", "--json"}} {
		if code, out, errOut := run(t, root, args...); code != 2 || out != "" || !strings.Contains(errOut, "Usage:") {
			t.Errorf("%v: %d %q %q", args, code, out, errOut)
		}
	}
	for _, args := range [][]string{{"deps", "G-002"}, {"deps", "G-404"}, {"deps", "G-001", "G-001"}} {
		if code, out, errOut := run(t, root, args...); code != 1 || out != "" || strings.Contains(errOut, "context") {
			t.Errorf("%v: %d %q %q", args, code, out, errOut)
		}
	}
}

// A done prerequisite on the target, a candidate in review off it, and an
// uncommitted edit, read from a real repository.
func TestDepsCLI(t *testing.T) {
	t.Parallel()
	root := gitFixture(t)
	base := gitIn(t, root, "rev-parse", "HEAD")
	gitIn(t, root, "switch", "-q", "-c", "feature")
	gitIn(t, root, "commit", "-q", "--allow-empty", "-m", "candidate")
	candidate := gitIn(t, root, "rev-parse", "HEAD")
	gitIn(t, root, "switch", "-q", "main")
	write(t, root, "grove.yaml", "schema_version: 3\nrecords: docs/records\ntarget: main\n")
	write(t, root, "docs/records/G-003.md", "---\nid: G-003\ntype: work\ntitle: Done before\nstatus: done\ncandidate: "+base+"\n---\n")
	write(t, root, "docs/records/G-005.md", "---\nid: G-005\ntype: work\ntitle: In review\nstatus: review\ncandidate: "+candidate+"\n---\n")
	next := "---\nid: G-004\ntype: work\ntitle: Next\nstatus: proposed\ndepends_on: [G-003, G-005]\n---\n"
	write(t, root, "docs/records/G-004.md", next)
	gitIn(t, root, "add", "-A")
	gitIn(t, root, "commit", "-q", "-m", "records")

	code, out, errOut := run(t, root, "deps")
	if code != 0 {
		t.Fatalf("deps: %d %s", code, errOut)
	}
	for _, want := range []string{
		"; target main\n", "GROUP  LAYER  ID     STATUS    NEEDS        UNLOCKS  DELIVERY",
		"1      0      G-001  proposed  -            -        awaiting implementation",
		"2      0      G-005  review    -            G-004    awaiting review; candidate " + candidate[:7] + " not in HEAD, not on main  In review\n",
		"2      1      G-004  proposed  G-003 G-005  -        awaiting implementation",
		"G-003  done    G-004      candidate " + base[:7] + " in HEAD, on main  Done before\n",
		"Equal layers have no declared order",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("overview lacks %q:\n%s", want, out)
		}
	}

	write(t, root, "docs/records/G-004.md", next+"Edited.\n")
	code, out, errOut = run(t, root, "deps", "G-004", "G-005", "--json")
	if code != 0 {
		t.Fatalf("deps --json: %d %s", code, errOut)
	}
	var got struct {
		Checkout map[string]any `json:"checkout"`
		Target   string         `json:"target"`
		deps.View
	}
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatal(err)
	}
	if got.Target != "main" || got.Checkout["ref"] != "refs/heads/main" || got.Checkout["root"] != root ||
		!reflect.DeepEqual(got.Selected, []string{"G-004", "G-005"}) || !reflect.DeepEqual(got.Order, []string{"G-005", "G-004"}) {
		t.Fatalf("%s", out)
	}
	delivery := map[string]string{}
	for _, it := range got.Items {
		delivery[it.ID] = it.Delivery
	}
	if want := map[string]string{
		"G-005": "awaiting review; candidate " + candidate[:7] + " not in HEAD, not on main",
		"G-004": "awaiting implementation",
		"G-003": "candidate " + base[:7] + " in HEAD, on main",
	}; !reflect.DeepEqual(delivery, want) {
		t.Errorf("delivery %v", delivery)
	}
	if !reflect.DeepEqual(got.Notes, []string{"G-004 has uncommitted changes (modified) in this checkout, which is what is read here."}) {
		t.Errorf("notes %q", got.Notes)
	}
	// context orders the same selection the same way.
	_, contextOut, _ := run(t, root, "context", "G-004", "G-005")
	if !strings.Contains(contextOut, "Order: G-005 G-004\n") {
		t.Errorf("context disagrees:\n%s", contextOut)
	}
}
