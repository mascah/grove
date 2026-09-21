package update

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mascah/grove/internal/project"
)

// sampleProject builds on gitProject's G-001..G-004 (two work, a question, a
// decision) with a plan shared by two work items, a review with an examined
// commit, and legacy documents outside the record root for convert to bring
// in: one predating plan records, one plain.
func sampleProject(t *testing.T) string {
	t.Helper()
	root := gitProject(t)
	write(t, root, "grove/G-005-shared.md", "---\nid: \"G-005\"\ntype: plan\ntitle: Shared\nstatus: superseded\nwork: [\"G-002\", \"G-001\"]\n---\nPlan.\n")
	write(t, root, "grove/G-006-first-review.md", "---\nid: \"G-006\"\ntype: review\ntitle: Review\nstatus: current\nwork: [\"G-001\"]\nexamined: \"fc9bef1\"\ncreated: \"2026-09-20T10:00:00Z\"\nupdated: \"2026-09-21T10:00:00Z\"\n---\nFindings.\n")
	write(t, root, "docs/plans/old-plan.md", "# Old plan\n\nWritten before plan records.\n")
	write(t, root, "docs/other.md", "# Other\n\nA legacy note.\n")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-q", "-m", "sample project")
	return root
}

func stateDir(t *testing.T, root string) string {
	t.Helper()
	return filepath.Join(git(t, root, "rev-parse", "--path-format=absolute", "--git-common-dir"), "grove")
}

func convert(t *testing.T, root string, req ConvertRequest) Conversion {
	t.Helper()
	c, err := Convert(root, req, &bytes.Buffer{})
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestConvertDocument(t *testing.T) {
	t.Parallel()
	root := sampleProject(t)
	c := convert(t, root, ConvertRequest{Source: "docs/plans/old-plan.md", Type: "plan", Title: "Old plan"})
	if c != (Conversion{"docs/plans/old-plan.md", "docs/plans/old-plan.md", "G-007", "grove/G-007-old-plan.md"}) {
		t.Fatalf("mapping = %+v", c)
	}
	want := "---\nid: \"G-007\"\ntype: plan\ntitle: \"Old plan\"\nstatus: current\nformerly: \"docs/plans/old-plan.md\"\n---\n\n# Old plan\n\nWritten before plan records.\n"
	if got := read(t, root, c.Path); got != want {
		t.Fatalf("got:\n%s", got)
	}
	write(t, root, "docs/bom.md", string(rune(0xFEFF))+"# BOM\r\n")
	if c := convert(t, root, ConvertRequest{Source: "docs/bom.md", Type: "page", Title: "B", Slug: "a--b"}); c.Path != "grove/G-008-a--b.md" || !strings.HasSuffix(read(t, root, c.Path), "---\n\n# BOM\r\n") {
		t.Fatalf("new and convert share one slug rule, and a BOM does not move mid-file: %+v\n%q", c, read(t, root, c.Path))
	}
	if read(t, root, "docs/plans/old-plan.md") == "" {
		t.Fatal("the original is the caller's to remove")
	}
}

func TestConvertRefusals(t *testing.T) {
	t.Parallel()
	root := sampleProject(t)
	convert(t, root, ConvertRequest{Source: "docs/plans/old-plan.md", Type: "page", Title: "Old"})
	before := read(t, stateDir(t, root), "neutral-ids")
	for _, tc := range []struct {
		req  ConvertRequest
		want string
	}{
		{ConvertRequest{Source: "docs/plans/old-plan.md", Type: "plan", Title: "Again"}, "was already converted to G-007"},
		{ConvertRequest{Source: "docs/Plans/Old-Plan.md", Type: "plan", Title: "Again"}, "was already converted to G-007"}, // one file on macOS
		{ConvertRequest{Source: "docs/other.md", Slug: "a_b", Type: "page", Title: "T"}, "slug must contain only"},
		{ConvertRequest{Source: "docs/other.md", Slug: "Bad Slug", Type: "page", Title: "T"}, "slug must contain only"},
		{ConvertRequest{Source: "docs/missing.md", Type: "page", Title: "T"}, "no such file"},
		{ConvertRequest{Source: "../outside.md", Type: "page", Title: "T"}, "clean path"},
		{ConvertRequest{Source: "grove/whatever.md", Type: "page", Title: "T"}, "inside the record root"},
		{ConvertRequest{Source: "docs/plans/old-plan.md"}, "already converted"},
		{ConvertRequest{Source: "grove.yaml", Type: "page", Title: "T"}, "clean path"},
		{ConvertRequest{Source: "docs/other.md", Type: "note", Title: "T"}, "requires --type"},
	} {
		if _, err := Convert(root, tc.req, &bytes.Buffer{}); err == nil || !strings.Contains(err.Error(), tc.want) {
			t.Fatalf("%+v: err = %v, want %q", tc.req, err, tc.want)
		}
	}
	if after := read(t, stateDir(t, root), "neutral-ids"); after != before {
		t.Fatalf("a refused conversion reserved an ID: %q -> %q", before, after)
	}
	// An existing target file is refused before anything else is touched.
	write(t, root, "docs/blocker.md", "")
	if err := os.Rename(filepath.Join(root, "docs/blocker.md"), filepath.Join(root, "grove/G-008-second.md")); err != nil {
		t.Fatal(err)
	}
	if _, err := Convert(root, ConvertRequest{Source: "docs/other.md", Type: "page", Title: "Blocked", Slug: "second"}, &bytes.Buffer{}); err == nil {
		t.Fatal("converted into an invalid project")
	}
}

// Classification changes in place: same ID, same path, and the new type's
// whole contract or nothing.
func TestReclassifyKeepsIdentityAndPath(t *testing.T) {
	t.Parallel()
	root := sampleProject(t)
	write(t, root, "grove/notes/G-050-idea.md", "---\nid: \"G-050\"\ntype: page\ntitle: Idea\n---\nProse that says status: done and approved.\n")
	rev := func(id string) string {
		p, _ := project.Load(root, root)
		for _, r := range p.Records {
			if r.ID == id {
				return project.Revision(r.Source)
			}
		}
		t.Fatalf("%s missing", id)
		return ""
	}
	if _, err := Apply(root, Request{ID: "G-050", Expect: rev("G-050"), Set: []Field{{"type", "work"}}}, now, nil); err == nil || !strings.Contains(err.Error(), "status: required field is missing") {
		t.Fatalf("a page became work without a status: %v", err)
	}
	res, err := Apply(root, Request{ID: "G-050", Expect: rev("G-050"), Set: []Field{{"type", "work"}, {"status", "proposed"}, {"priority", "2"}}}, now, nil)
	if err != nil || res.Path != "grove/notes/G-050-idea.md" || res.ID != "G-050" {
		t.Fatalf("res=%+v err=%v", res, err)
	}
	if _, err := Apply(root, Request{ID: "G-050", Expect: rev("G-050"), Set: []Field{{"type", "page"}}, Unset: []string{"status"}}, now, nil); err == nil || !strings.Contains(err.Error(), "priority: unknown field") {
		t.Fatalf("work fields survived on a page: %v", err)
	}
	if _, err := Apply(root, Request{ID: "G-050", Expect: rev("G-050"), Set: []Field{{"type", "page"}}, Unset: []string{"status", "priority"}}, now, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := Apply(root, Request{ID: "G-002", Expect: rev("G-002"), Set: []Field{{"type", "decision"}}}, now, nil); err == nil {
		t.Fatal("work that a plan names became a decision")
	}
	if _, err := Apply(root, Request{ID: "G-050", Expect: rev("G-050"), Unset: []string{"type"}}, now, nil); err == nil {
		t.Fatal("type was unset")
	}
	for _, name := range []string{"id", "formerly"} {
		if _, err := Apply(root, Request{ID: "G-050", Expect: rev("G-050"), Set: []Field{{name, "G-051"}}}, now, nil); err == nil || !strings.Contains(err.Error(), "cannot be changed by update") {
			t.Fatalf("%s: %v", name, err)
		}
	}
}
