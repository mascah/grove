package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/mascah/grove/internal/handoff"
	"github.com/mascah/grove/internal/project"
)

const plannedWork = "---\nid: G-003\ntype: work\ntitle: Planned\nstatus: proposed\ndepends_on: [G-001]\n---\nFollow the [plan](../../plans/G-003.md).\n"

func contextFixture(t *testing.T, root string) {
	t.Helper()
	write(t, root, "docs/records/work/planned.md", plannedWork)
	write(t, root, "docs/plans/G-003.md", "# Plan\n")
	write(t, root, "AGENTS.md", "Follow the guide.\n")
}

func TestContextUsage(t *testing.T) {
	t.Parallel()
	root := projectFixture(t)
	for _, args := range [][]string{
		{"context"}, {"context", "--json"},
		{"context", "G-001", "--interaction", "auto"},
		{"context", "G-001", "--interaction", "headless", "--interaction=headless"},
		{"context", "G-001", "--max-bytes", "0"}, {"context", "G-001", "--max-bytes", "8388609"},
		{"context", "G-001", "--max-bytes", "99999999999999999999"}, {"context", "G-001", "--max-bytes", "1k"},
		{"context", "G-001", "--max-bytes", "10", "--max-bytes", "10"},
		{"context", "G-001", "--include"}, {"context", "G-001", "--include", " "},
		{"context", "--work", "G-001"}, {"context", "G-001", "--phase", "implement"},
		{"context", "G-001", "--slug", "x"}, {"context", "G-001", "--source", "x"},
		// Context options never reach another command, nor open the board.
		{"--include", "AGENTS.md"}, {"--interaction", "headless"}, {"--max-bytes", "10"},
		{"list", "--include", "AGENTS.md"}, {"show", "G-001", "--interaction", "headless"},
	} {
		var out, errOut bytes.Buffer
		if code := Run(args, root, &out, &errOut); code != 2 || out.Len() != 0 || !strings.Contains(errOut.String(), "Usage:") {
			t.Errorf("%v: code %d, stdout %q, stderr %q", args, code, out.String(), errOut.String())
		}
	}
	var out, errOut bytes.Buffer
	if code := Run([]string{"context", "--help"}, t.TempDir(), &out, &errOut); code != 0 || !strings.Contains(out.String(), "context WORK_ID...") {
		t.Fatalf("help outside a project: %d %s", code, errOut.String())
	}
}

func TestContextCLI(t *testing.T) {
	t.Parallel()
	root := projectFixture(t)
	contextFixture(t, root)
	var out, errOut bytes.Buffer
	args := []string{"context", "G-003", "G-001", "--json", "--interaction=headless", "--include", "AGENTS.md", "--include", "grove.yaml", "--project", root}
	if code := Run(args, t.TempDir(), &out, &errOut); code != 0 {
		t.Fatalf("context: %d: %s", code, errOut.String())
	}
	var got handoff.Bundle
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got.Selected, []string{"G-003", "G-001"}) || !reflect.DeepEqual(got.Order, []string{"G-001", "G-003"}) ||
		got.Interaction != "headless" || got.Root != root || got.FormatVersion != 2 {
		t.Fatalf("%+v", got)
	}
	// The selected work, the configuration, and the includes are read in full. The
	// related question and the linked plan are listed for a later, explicit read.
	want := map[string]string{
		"AGENTS.md": "Follow the guide.\n", "docs/records/work/planned.md": plannedWork,
		"docs/records/work/renamed.md": work, "grove.yaml": "schema_version: 3\nrecords: docs/records\n",
	}
	for _, s := range got.Sources {
		if want[s.Path] != s.Content || s.Revision != project.Revision([]byte(s.Content)) {
			t.Errorf("%s: %q %s", s.Path, s.Content, s.Revision)
		}
		delete(want, s.Path)
	}
	if len(want) != 0 {
		t.Fatalf("missing sources: %v", want)
	}
	// Records sort by ID, so the selected G-001 and G-003 now come before the
	// listed related question G-002.
	var listed *handoff.Record
	for i := range got.Records {
		if got.Records[i].ID == "G-002" {
			listed = &got.Records[i]
		}
	}
	if listed == nil {
		t.Fatalf("G-002 missing from records: %+v", got.Records)
	}
	if listed.Included || listed.Revision != project.Revision([]byte(question)) || listed.Title == "" ||
		len(got.References) != 1 || got.References[0].Path != "docs/plans/G-003.md" || !strings.Contains(got.References[0].Reason, "not opened") {
		t.Fatalf("%+v %+v", *listed, got.References)
	}

	// Naming the listed plan reads it, with its exact revision.
	out.Reset()
	if code := Run([]string{"context", "G-003", "--json", "--include", "docs/plans/G-003.md"}, root, &out, &errOut); code != 0 {
		t.Fatal(errOut.String())
	}
	var staged handoff.Bundle
	if err := json.Unmarshal(out.Bytes(), &staged); err != nil || staged.Sources[0].Path != "docs/plans/G-003.md" ||
		staged.Sources[0].Content != "# Plan\n" || staged.Sources[0].Revision != project.Revision([]byte("# Plan\n")) ||
		!strings.Contains(staged.References[0].Reason, "included in full") {
		t.Fatalf("%v %+v", err, staged)
	}

	// A changed record changes its revision in the next context.
	revision := func() string {
		for _, s := range got.Sources {
			if s.Path == "docs/records/work/planned.md" {
				return s.Revision
			}
		}
		return ""
	}
	before := revision()
	write(t, root, "docs/records/work/planned.md", strings.Replace(plannedWork, "Follow", "Still follow", 1))
	out.Reset()
	if code := Run([]string{"context", "G-003", "--json"}, root, &out, &errOut); code != 0 {
		t.Fatal(errOut.String())
	}
	if err := json.Unmarshal(out.Bytes(), &got); err != nil || revision() == "" || revision() == before {
		t.Fatalf("%v %s", err, revision())
	}

	out.Reset()
	if code := Run([]string{"context", "G-003"}, root, &out, &errOut); code != 0 || !strings.Contains(out.String(), "G-003 depends on G-001: proposed, not selected") ||
		!strings.Contains(out.String(), "G-001  work  proposed  listed  prerequisite of G-003") || strings.Contains(out.String(), "Source: docs/records/work/renamed.md") {
		t.Fatalf("text: %d %s", code, out.String())
	}
	if code := Run([]string{"context", "G-003"}, root, brokenWriter{}, &errOut); code != 1 {
		t.Fatalf("broken output: %d", code)
	}
}

func TestContextRefusalsWriteNoResult(t *testing.T) {
	t.Parallel()
	root := projectFixture(t)
	contextFixture(t, root)
	os.Remove(filepath.Join(root, "docs/plans/G-003.md"))
	for want, args := range map[string][]string{
		"--include names docs/plans/G-003.md, which does not exist": {"context", "G-003", "--include", "docs/plans/G-003.md"},
		"does not fit":         {"context", "G-001", "--json", "--max-bytes", "64"},
		"--include names nope": {"context", "G-001", "--include", "nope"},
		"not in this checkout": {"context", "G-404"},
		"only work":            {"context", "G-002", "--json"},
		"more than once":       {"context", "G-001", "G-001"},
	} {
		var out, errOut bytes.Buffer
		if code := Run(args, root, &out, &errOut); code != 1 || out.Len() != 0 || !strings.Contains(errOut.String(), want) {
			t.Errorf("%v: code %d, stdout %q, stderr %q", args, code, out.String(), errOut.String())
		}
	}
}

func TestContextLeavesEverythingUnchanged(t *testing.T) {
	t.Parallel()
	root, wt := featureFixture(t)
	contextFixture(t, wt) // uncommitted files in the linked checkout
	before := map[string]map[string][32]byte{root: hashes(t, root), wt: hashes(t, wt)}
	for args, code := range map[*[]string]int{
		{"context", "G-001", "--json"}:                            0,
		{"--project", wt, "context", "G-003", "G-001"}:            0,
		{"--project", wt, "context", "G-003", "--max-bytes", "8"}: 1,
		{"context", "G-003"}:                                      1, // only in the other checkout
	} {
		var out, errOut bytes.Buffer
		if got := Run(*args, root, &out, &errOut); got != code {
			t.Fatalf("%v: %d: %s", *args, got, errOut.String())
		}
	}
	if !reflect.DeepEqual(before[root], hashes(t, root)) || !reflect.DeepEqual(before[wt], hashes(t, wt)) {
		t.Fatal("context must leave Git metadata, records, and dirty files unchanged")
	}
	if _, err := os.Stat(filepath.Join(root, ".git", "grove")); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("context must not create coordination state: %v", err)
	}
}
