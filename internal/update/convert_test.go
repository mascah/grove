package update

import (
	"bytes"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/mascah/grove/internal/project"
)

// legacy is a schema-3 project still laid out as schema 2 left it: typed IDs
// in type folders, a plan shared by two work items, a review with an examined
// commit, and a plan that predates plan records.
func legacy(t *testing.T) string {
	t.Helper()
	root := gitProject(t)
	write(t, root, "grove.yaml", "schema_version: 3\nrecords: grove\n")
	write(t, root, "grove/work/W-002-second.md", "---\nid: \"W-002\"\ntype: work\ntitle: Second\nstatus: proposed\ndepends_on: [\"W-001\"]  # keep\nmembers:\n  - \"W-001\"\n---\nSee [first](W-001-first.md) and W-001.\n")
	write(t, root, "grove/questions/Q-001-which.md", "---\nid: \"Q-001\"\ntype: question\ntitle: Which?\nstatus: open\nblocks: [\"W-001\", \"W-002\"]\n---\nBody.\n")
	write(t, root, "grove/plans/P-001-shared.md", "---\nid: \"P-001\"\ntype: plan\ntitle: Shared\nstatus: superseded\nwork: [\"W-002\", \"W-001\"]\n---\nPlan.\n")
	write(t, root, "grove/reviews/R-001-first.md", "---\nid: \"R-001\"\ntype: review\ntitle: Review\nstatus: current\nwork: [\"W-001\"]\nexamined: \"fc9bef1\"\ncreated: \"2026-09-20T10:00:00Z\"\nupdated: \"2026-09-21T10:00:00Z\"\n---\nFindings.\n")
	write(t, root, "docs/plans/W-001-old.md", "# Old plan\n\nWritten before plan records.\n")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-q", "-m", "legacy layout under schema 3")
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

func TestConvertRecordsKeepsMetadataAndRebuildsReferences(t *testing.T) {
	t.Parallel()
	root := legacy(t)
	if c := convert(t, root, ConvertRequest{Source: "W-001"}); c != (Conversion{"W-001", "grove/work/W-001-first.md", "G-001", "grove/G-001-first.md"}) {
		t.Fatalf("mapping = %+v", c)
	}
	// Only id changes and formerly is appended; dates and body are untouched.
	want := strings.Replace(work, "id: \"W-001\"", "id: \"G-001\"", 1)
	want = strings.Replace(want, "updated: \"2026-09-19T12:00:00Z\"\n", "updated: \"2026-09-19T12:00:00Z\"\nformerly: \"W-001\"\n", 1)
	if got := read(t, root, "grove/G-001-first.md"); got != want {
		t.Fatalf("converted record:\n%s", got)
	}
	if _, err := os.Stat(filepath.Join(root, "grove/work/W-001-first.md")); !os.IsNotExist(err) {
		t.Fatal("the old file remains")
	}
	// Referrers change only in the lists; comments, block style order and prose stay.
	second := read(t, root, "grove/work/W-002-second.md")
	for _, part := range []string{"depends_on: [\"G-001\"]  # keep\n", "members: [\"G-001\"]\n", "See [first](W-001-first.md) and W-001.\n"} {
		if !strings.Contains(second, part) {
			t.Fatalf("W-002 lacks %q:\n%s", part, second)
		}
	}
	if got := read(t, root, "grove/plans/P-001-shared.md"); !strings.Contains(got, "work: [\"W-002\", \"G-001\"]\n") {
		t.Fatalf("shared plan:\n%s", got)
	}
	convert(t, root, ConvertRequest{Source: "R-001", Slug: "w-019-review"})
	review := read(t, root, "grove/G-002-w-019-review.md")
	for _, part := range []string{"work: [\"G-001\"]\n", "examined: \"fc9bef1\"\n", "updated: \"2026-09-21T10:00:00Z\"\n", "formerly: \"R-001\"\n"} {
		if !strings.Contains(review, part) {
			t.Fatalf("review lacks %q:\n%s", part, review)
		}
	}
	convert(t, root, ConvertRequest{Source: "P-001"})
	p, ds := project.Load(root, root)
	if len(ds) != 0 {
		t.Fatalf("project invalid after conversions:\n%s", diagnostics(ds))
	}
	for _, r := range p.Records {
		if r.ID == "G-003" && (r.Status != "superseded" || strings.Join(r.Work, " ") != "W-002 G-001") {
			t.Fatalf("plan lost metadata: %+v", r)
		}
	}
}

func TestConvertDocument(t *testing.T) {
	t.Parallel()
	root := legacy(t)
	c := convert(t, root, ConvertRequest{Source: "docs/plans/W-001-old.md", Type: "plan", Title: "Old plan"})
	if c != (Conversion{"docs/plans/W-001-old.md", "docs/plans/W-001-old.md", "G-001", "grove/G-001-old-plan.md"}) {
		t.Fatalf("mapping = %+v", c)
	}
	want := "---\nid: \"G-001\"\ntype: plan\ntitle: \"Old plan\"\nstatus: current\nformerly: \"docs/plans/W-001-old.md\"\n---\n\n# Old plan\n\nWritten before plan records.\n"
	if got := read(t, root, c.Path); got != want {
		t.Fatalf("got:\n%s", got)
	}
	write(t, root, "docs/bom.md", "\ufeff# BOM\r\n")
	if c := convert(t, root, ConvertRequest{Source: "docs/bom.md", Type: "page", Title: "B", Slug: "a--b"}); c.Path != "grove/G-002-a--b.md" || !strings.HasSuffix(read(t, root, c.Path), "---\n\n# BOM\r\n") {
		t.Fatalf("new and convert share one slug rule, and a BOM does not move mid-file: %+v\n%q", c, read(t, root, c.Path))
	}
	if read(t, root, "docs/plans/W-001-old.md") == "" {
		t.Fatal("the original is the caller's to remove")
	}
}

func TestConvertRefusals(t *testing.T) {
	t.Parallel()
	root := legacy(t)
	convert(t, root, ConvertRequest{Source: "W-001"})
	convert(t, root, ConvertRequest{Source: "docs/plans/W-001-old.md", Type: "page", Title: "Old"})
	before := read(t, stateDir(t, root), "neutral-ids")
	for _, tc := range []struct {
		req  ConvertRequest
		want string
	}{
		{ConvertRequest{Source: "W-001"}, "W-001 was already converted to G-001 in grove/G-001-first.md"},
		{ConvertRequest{Source: "docs/plans/W-001-old.md", Type: "plan", Title: "Again"}, "was already converted to G-002"},
		{ConvertRequest{Source: "docs/Plans/w-001-OLD.md", Type: "plan", Title: "Again"}, "was already converted to G-002"}, // one file on macOS
		{ConvertRequest{Source: "W-002", Slug: "a_b"}, "slug must contain only"},
		{ConvertRequest{Source: "G-001"}, "G-001 already has a neutral ID"},
		{ConvertRequest{Source: "W-404"}, "record W-404 not found"},
		{ConvertRequest{Source: "W-002", Type: "page"}, "--type and --title apply only to converting a document"},
		{ConvertRequest{Source: "W-002", Slug: "Bad Slug"}, "slug must contain only"},
		{ConvertRequest{Source: "docs/missing.md", Type: "page", Title: "T"}, "no such file"},
		{ConvertRequest{Source: "../outside.md", Type: "page", Title: "T"}, "clean path"},
		{ConvertRequest{Source: "grove/G-001-first.md", Type: "page", Title: "T"}, "inside the record root"},
		{ConvertRequest{Source: "docs/plans/W-001-old.md"}, "already converted"},
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
	write(t, root, "docs/G-003-second.md", "")
	if err := os.Rename(filepath.Join(root, "docs/G-003-second.md"), filepath.Join(root, "grove/G-003-second.md")); err != nil {
		t.Fatal(err)
	}
	if _, err := Convert(root, ConvertRequest{Source: "W-002"}, &bytes.Buffer{}); err == nil {
		t.Fatal("converted into an invalid project")
	}
}

func TestConvertNeedsSchema3(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	// Earlier schemas have no formerly, and say so as they do of any unknown field.
	p, _ := project.Load(root, root)
	if _, err := Apply(root, Request{ID: "W-002", Expect: project.Revision(p.Records[slices.IndexFunc(p.Records, func(r *project.Record) bool { return r.ID == "W-002" })].Source), Set: []Field{{"formerly", "x"}}}, now, nil); err == nil || !strings.Contains(err.Error(), "formerly is not a field that update accepts on work records") {
		t.Fatalf("err = %v", err)
	}
	if _, err := Convert(root, ConvertRequest{Source: "W-001"}, &bytes.Buffer{}); err == nil || !strings.Contains(err.Error(), "convert needs schema_version 3") {
		t.Fatalf("err = %v", err)
	}
}

// Classification changes in place: same ID, same path, and the new type's
// whole contract or nothing.
func TestReclassifyKeepsIdentityAndPath(t *testing.T) {
	t.Parallel()
	root := legacy(t)
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
	if _, err := Apply(root, Request{ID: "W-002", Expect: rev("W-002"), Set: []Field{{"type", "decision"}}}, now, nil); err == nil {
		t.Fatal("work that a question blocks and a plan names became a decision")
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
