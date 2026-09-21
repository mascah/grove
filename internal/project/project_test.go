package project

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func put(t *testing.T, root, path, content string) {
	t.Helper()
	path = filepath.Join(root, path)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}

func fixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	root, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatal(err)
	}
	put(t, root, "grove.yaml", "schema_version: 3\nrecords: grove\n")
	if err := os.Mkdir(filepath.Join(root, "grove"), 0755); err != nil {
		t.Fatal(err)
	}
	return root
}

func record(id, kind, extra string) string {
	status := "proposed"
	if kind == "question" {
		status = "open"
	}
	return "---\nid: " + id + "\ntype: " + kind + "\ntitle: Example\nstatus: " + status + "\n" + extra + "---\nBody with --- inside.\n"
}

func diagnostics(ds []Diagnostic) string {
	var out []string
	for _, d := range ds {
		out = append(out, d.String())
	}
	return strings.Join(out, "\n")
}

func TestLoadRenamedNestedRecordsAndPreserveSource(t *testing.T) {
	t.Parallel()
	root := fixture(t)
	source := record("G-001", "work", "kind: investigation\npriority: 2\nsize: small\ncreated: \"2026-09-19T14:08:40Z\"\n")
	source = string(rune(0xFEFF)) + strings.ReplaceAll(source, "\n", "\r\n")
	put(t, root, "grove/work/nested/custom.md", source)
	put(t, root, "grove/questions/anything.md", record("G-002", "question", ""))
	put(t, root, "grove/notes.txt", "not a record")
	p, ds := Load(filepath.Join(root, "grove", "work", "nested"), "")
	if len(ds) != 0 {
		t.Fatal(diagnostics(ds))
	}
	if len(p.Records) != 2 || p.Records[0].ID != "G-001" {
		t.Fatalf("records: %+v", p.Records)
	}
	if string(p.Records[0].Source) != source || p.Records[0].Path != "grove/work/nested/custom.md" {
		t.Fatal("source bytes or path changed")
	}
	if p.Records[1].Priority != nil || p.Records[1].Size != "" {
		t.Fatal("missing planning values must remain unspecified")
	}
}

func TestDiscoveryAndExplicitProject(t *testing.T) {
	t.Parallel()
	root := fixture(t)
	put(t, root, "nested/.git", "gitdir: unused-by-reader\n")
	_, ds := Load(filepath.Join(root, "nested"), "")
	if !strings.Contains(diagnostics(ds), "grove.yaml") {
		t.Fatal("nested Git checkout inherited parent configuration")
	}
	p, ds := Load(filepath.Join(root, "nested"), "..")
	if len(ds) != 0 || p.Root != root {
		t.Fatalf("explicit parent project: %v", diagnostics(ds))
	}
	put(t, root, "nested/grove.yaml", "schema_version: 3\nrecords: records\n")
	if err := os.Mkdir(filepath.Join(root, "nested", "records"), 0755); err != nil {
		t.Fatal(err)
	}
	p, ds = Load(filepath.Join(root, "nested"), "")
	if len(ds) != 0 || p.Root != filepath.Join(root, "nested") {
		t.Fatalf("nearest project: %v", diagnostics(ds))
	}
	_, ds = Load(root, "grove")
	if len(ds) == 0 {
		t.Fatal("explicit project must not search upward")
	}
}

func TestInvalidConfiguration(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, source, field string }{
		{"version missing", "records: grove\n", "schema_version"},
		{"unsupported", "schema_version: 4\nrecords: grove\n", "schema_version"},
		{"schema 1 no longer supported", "schema_version: 1\nrecords: grove\n", "schema_version"},
		{"schema 2 no longer supported", "schema_version: 2\nrecords: grove\n", "schema_version"},
		{"quoted version", "schema_version: \"3\"\nrecords: grove\n", "schema_version"},
		{"fractional version", "schema_version: 3.0\nrecords: grove\n", "schema_version"},
		{"null root", "schema_version: 3\nrecords: null\n", "records"},
		{"escape", "schema_version: 3\nrecords: ../outside\n", "records"},
		{"hidden escape", "schema_version: 3\nrecords: a/../grove\n", "records"},
		{"project root", "schema_version: 3\nrecords: .\n", "records"},
		{"absolute root", "schema_version: 3\nrecords: /tmp\n", "records"},
		{"unknown", "schema_version: 3\nrecords: grove\nextra: true\n", "extra"},
		{"duplicate", "schema_version: 3\nrecords: grove\nrecords: other\n", "records"},
		{"extra document", "schema_version: 3\nrecords: grove\n---\n{}\n", "document"},
		{"missing root", "schema_version: 3\nrecords: absent\n", "records"},
		{"non mapping", "- records\n", "mapping"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := fixture(t)
			put(t, root, "grove.yaml", tc.source)
			_, ds := Load(root, "")
			msg := diagnostics(ds)
			if !strings.Contains(msg, "grove.yaml") || !strings.Contains(msg, tc.field) {
				t.Fatalf("wanted file and %q diagnostic; got %s", tc.field, msg)
			}
		})
	}
}

func TestStrictRecordMetadata(t *testing.T) {
	t.Parallel()
	base := record("G-001", "work", "")
	for _, tc := range []struct{ name, source, field string }{
		{"missing title", strings.Replace(base, "title: Example\n", "", 1), "title"},
		{"numeric title", strings.Replace(base, "title: Example", "title: 123", 1), "title"},
		{"empty title", strings.Replace(base, "title: Example", "title: \"  \"", 1), "title"},
		{"null", record("G-001", "work", "size: null\n"), "size"},
		{"priority float", record("G-001", "work", "priority: 2.0\n"), "priority"},
		{"priority range", record("G-001", "work", "priority: 6\n"), "priority"},
		{"size", record("G-001", "work", "size: spike\n"), "size"},
		{"kind", record("G-001", "work", "kind: magic\n"), "kind"},
		{"unknown", record("G-001", "work", "priorty: 2\n"), "priorty"},
		{"duplicate key", record("G-001", "work", "title: Again\n"), "title"},
		{"duplicate relationship", record("G-001", "work", "relates_to: [G-002, G-002]\n"), "relates_to"},
		{"sequence required", record("G-001", "work", "depends_on: G-002\n"), "depends_on"},
		{"nonstring target", record("G-001", "work", "depends_on: [42]\n"), "depends_on"},
		{"wrong type field", record("G-001", "work", "blocks: []\n"), "blocks"},
		{"bare timestamp", record("G-001", "work", "created: 2026-09-19T12:00:00Z\n"), "created"},
		{"offset timestamp", record("G-001", "work", "created: \"2026-09-19T12:00:00+00:00\"\n"), "created"},
		{"fraction timestamp", record("G-001", "work", "created: \"2026-09-19T12:00:00.123Z\"\n"), "created"},
		{"invalid date", record("G-001", "work", "created: \"2026-02-30T12:00:00Z\"\n"), "created"},
		{"backward date", record("G-001", "work", "created: \"2026-09-19T12:00:00Z\"\nupdated: \"2026-09-18T12:00:00Z\"\n"), "updated"},
		{"backward zero date", record("G-001", "work", "created: \"0001-01-01T00:00:00Z\"\nupdated: \"0000-01-01T00:00:00Z\"\n"), "updated"},
		{"zero id", record("G-000", "work", ""), "id"},
		{"short id", record("G-1", "work", ""), "id"},
		{"extra padding", record("G-0001", "work", ""), "id"},
		{"unknown prefix", record("X-001", "work", ""), "id"},
		{"lifecycle", strings.Replace(base, "status: proposed", "status: resolved", 1), "status"},
		{"unclosed header", strings.TrimSuffix(base, "---\nBody with --- inside.\n"), "frontmatter"},
		{"no header", "Hello\n", "frontmatter"},
		{"alias", "---\nid: G-001\ntype: work\ntitle: &label Example\nstatus: *label\n---\n", "alias"},
		{"custom tag", strings.Replace(base, "title: Example", "title: !special Example", 1), "title"},
		{"invalid utf8", base + "\xff", "UTF-8"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := fixture(t)
			put(t, root, "grove/work/custom.md", tc.source)
			_, ds := Load(root, "")
			msg := diagnostics(ds)
			if !strings.Contains(msg, "custom.md") || !strings.Contains(msg, tc.field) {
				t.Fatalf("wanted file and %q diagnostic; got %s", tc.field, msg)
			}
		})
	}
}

func TestRejectSymlinks(t *testing.T) {
	t.Parallel()
	for _, target := range []string{"grove.yaml", "grove", "grove/work", "grove/work/link.md", "grove/link.txt"} {
		t.Run(target, func(t *testing.T) {
			root := fixture(t)
			outside := t.TempDir()
			put(t, outside, "target", record("G-001", "work", ""))
			path := filepath.Join(root, target)
			os.Remove(path)
			os.MkdirAll(filepath.Dir(path), 0755)
			if err := os.Symlink(filepath.Join(outside, "target"), path); err != nil {
				t.Fatal(err)
			}
			_, ds := Load(root, "")
			if !strings.Contains(diagnostics(ds), "symlink") {
				t.Fatalf("followed %s: %s", target, diagnostics(ds))
			}
		})
	}
}

func TestOrderUsesCreationThenNumericID(t *testing.T) {
	t.Parallel()
	root := fixture(t)
	for _, id := range []string{"G-1000", "G-999", "G-002", "G-001"} {
		extra := ""
		if id == "G-002" {
			extra = "created: \"2026-09-19T12:00:00Z\"\n"
		}
		put(t, root, "grove/work/"+id+".md", record(id, "work", extra))
	}
	p, ds := Load(root, "")
	if len(ds) != 0 {
		t.Fatal(diagnostics(ds))
	}
	var ids []string
	for _, r := range p.Records {
		ids = append(ids, r.ID)
	}
	if strings.Join(ids, ",") != "G-002,G-001,G-999,G-1000" {
		t.Fatal(ids)
	}
}

func TestExplicitZeroTimeIsNotAnAbsentDate(t *testing.T) {
	t.Parallel()
	root := fixture(t)
	put(t, root, "grove/work/undated.md", record("G-001", "work", ""))
	put(t, root, "grove/work/dated.md", record("G-999", "work", "created: \"0001-01-01T00:00:00Z\"\n"))
	p, ds := Load(root, "")
	if len(ds) != 0 || p.Records[0].ID != "G-999" {
		t.Fatalf("a present timestamp must sort before undated records: %s", diagnostics(ds))
	}
}

func schema3(t *testing.T, config string) string {
	t.Helper()
	root := fixture(t)
	put(t, root, "grove.yaml", "schema_version: 3\nrecords: grove\n"+config)
	return root
}

func typed(id, kind, status, extra string) string {
	return "---\nid: " + id + "\ntype: " + kind + "\ntitle: " + id + " title\nstatus: " + status + "\n" + extra + "---\nBody.\n"
}

func TestKnowledgeRecords(t *testing.T) {
	t.Parallel()
	root := schema3(t, "")
	put(t, root, "grove/G-001.md", typed("G-001", "work", "proposed", ""))
	put(t, root, "grove/terms/G-002-attempt.md", typed("G-002", "term", "settled", "relates_to: [\"G-001\"]\n"))
	put(t, root, "grove/plans/G-003-shared.md", typed("G-003", "plan", "current", "work: [\"G-001\"]\n"))
	put(t, root, "grove/reviews/G-004-first.md", typed("G-004", "review", "current", "work: [\"G-001\"]\nexamined: \"42c077d\"\n"))
	p, ds := Load(root, root)
	if len(ds) != 0 {
		t.Fatalf("unexpected diagnostics:\n%s", diagnostics(ds))
	}
	if len(p.Records) != 4 {
		t.Fatalf("got %d records", len(p.Records))
	}
	for _, r := range p.Records {
		if r.ID == "G-004" && (r.Examined != "42c077d" || len(r.Work) != 1) {
			t.Fatalf("review fields not read: %+v", r)
		}
	}
}

func TestKnowledgeRecordProblemsNameFileAndField(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, path, source, want string }{
		{"missing work target", "grove/plans/G-002.md", typed("G-002", "plan", "current", "work: [\"G-009\"]\n"), "grove/plans/G-002.md: work: unresolved target G-009"},
		{"work target not work", "grove/plans/G-002.md", typed("G-002", "plan", "current", "work: [\"G-002\"]\n"), "work: self-reference"},
		{"work names a term", "grove/reviews/G-002.md", typed("G-002", "review", "current", "work: [\"G-001\"]\n"), "grove/reviews/G-002.md: work: target G-001 must be work"},
		{"unquoted numeric commit", "grove/reviews/G-002.md", typed("G-002", "review", "current", "examined: 1234567\n"), "examined: expected a nonempty string"},
		{"not a commit", "grove/reviews/G-002.md", typed("G-002", "review", "current", "examined: \"main\"\n"), "examined: expected a quoted Git commit"},
		{"examined on a plan", "grove/plans/G-002.md", typed("G-002", "plan", "current", "examined: \"42c077d\"\n"), "examined: unknown field"},
		{"work on a term", "grove/terms/G-003.md", typed("G-003", "term", "proposed", "work: [\"G-001\"]\n"), "work: unknown field"},
		{"bad status", "grove/terms/G-003.md", typed("G-003", "term", "done", ""), "status: unsupported lifecycle value for term"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := schema3(t, "")
			put(t, root, "grove/G-001.md", typed("G-001", "term", "settled", ""))
			put(t, root, tc.path, tc.source)
			_, ds := Load(root, root)
			if got := diagnostics(ds); !strings.Contains(got, tc.want) {
				t.Fatalf("wanted %q; got %s", tc.want, got)
			}
		})
	}
}

func TestDuplicateTermTitle(t *testing.T) {
	t.Parallel()
	root := schema3(t, "")
	put(t, root, "grove/terms/G-001.md", "---\nid: G-001\ntype: term\ntitle: Attempt\nstatus: settled\n---\n")
	put(t, root, "grove/terms/G-002.md", "---\nid: G-002\ntype: term\ntitle: \" attempt\"\nstatus: proposed\n---\n")
	_, ds := Load(root, root)
	if got := diagnostics(ds); !strings.Contains(got, "grove/terms/G-002.md: title: term already defined by G-001") {
		t.Fatalf("got %s", got)
	}
}

// briefFixture is a project with one baseline work record, for tests
// that also need the brief: configuration key.
func briefFixture(t *testing.T, config string) string {
	t.Helper()
	root := schema3(t, config)
	put(t, root, "grove/G-001.md", record("G-001", "work", ""))
	return root
}

func TestBrief(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, config, file, want string }{
		{"under the record root", "brief: grove/brief.md\n", "grove/brief.md", ""},
		{"elsewhere in the project", "brief: docs/restart-brief.md\n", "docs/restart-brief.md", ""},
		{"missing", "brief: grove/brief.md\n", "", "grove.yaml: brief: "},
		{"escaping", "brief: ../brief.md\n", "", "brief: must name a project-relative .md file"},
		{"unclean", "brief: ./grove/brief.md\n", "grove/brief.md", "brief: must name a project-relative .md file"},
		{"not Markdown", "brief: grove/brief.txt\n", "grove/brief.txt", "brief: must name a project-relative .md file"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := briefFixture(t, tc.config)
			if tc.file != "" {
				put(t, root, tc.file, "# Brief\n")
			}
			p, ds := Load(root, root)
			got := diagnostics(ds)
			if tc.want == "" {
				if got != "" {
					t.Fatalf("unexpected diagnostics: %s", got)
				}
				if source, err := p.ReadBrief(); err != nil || string(source) != "# Brief\n" {
					t.Fatalf("ReadBrief = %q, %v", source, err)
				}
				if len(p.Records) != 1 {
					t.Fatalf("the brief must not be a record; got %d records", len(p.Records))
				}
			} else if !strings.Contains(got, tc.want) {
				t.Fatalf("wanted %q; got %s", tc.want, got)
			}
		})
	}
	t.Run("spelled with another case than the disk", func(t *testing.T) {
		t.Parallel()
		root := briefFixture(t, "brief: grove/Notes/brief.md\n")
		put(t, root, "grove/notes/brief.md", "# Brief\n")
		if _, err := os.Stat(filepath.Join(root, "grove/Notes/brief.md")); err != nil {
			t.Skip("case-sensitive filesystem")
		}
		if _, ds := Load(root, root); len(ds) != 0 {
			t.Fatalf("the brief must be exempt however it is spelled: %s", diagnostics(ds))
		}
	})
	t.Run("symlinked parent directory", func(t *testing.T) {
		t.Parallel()
		root := briefFixture(t, "brief: docs/brief.md\n")
		put(t, root, "real/brief.md", "# Brief\n")
		if err := os.Symlink("real", filepath.Join(root, "docs")); err != nil {
			t.Fatal(err)
		}
		_, ds := Load(root, root)
		if got := diagnostics(ds); !strings.Contains(got, "grove.yaml: brief: docs is a symlink") {
			t.Fatalf("got %s", got)
		}
	})
	t.Run("symlinked brief", func(t *testing.T) {
		t.Parallel()
		root := briefFixture(t, "brief: docs/brief.md\n")
		put(t, root, "real.md", "# Brief\n")
		if err := os.Mkdir(filepath.Join(root, "docs"), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink("../real.md", filepath.Join(root, "docs", "brief.md")); err != nil {
			t.Fatal(err)
		}
		_, ds := Load(root, root)
		if got := diagnostics(ds); !strings.Contains(got, "brief: symlink files are not supported") {
			t.Fatalf("got %s", got)
		}
	})
}

const page = "---\nid: G-002\ntype: page\ntitle: Synthesis\nrelates_to: [\"G-001\"]\n---\nNotes.\n"

// Location carries no meaning: the root, a legacy-named folder, and an
// arbitrary nested folder all hold any type, under recursive discovery.
func TestDiscoversByIdentityNotLocation(t *testing.T) {
	t.Parallel()
	root := schema3(t, "brief: grove/brief.md\n")
	put(t, root, "grove/brief.md", "# Brief, not a record\n")
	put(t, root, "grove/G-001-work.md", typed("G-001", "work", "proposed", ""))
	put(t, root, "grove/G-002-synthesis.md", page)
	put(t, root, "grove/decisions/deep/er/G-003.md", typed("G-003", "plan", "current", "work: [\"G-001\"]\n"))
	put(t, root, "grove/work/G-004-legacy-folder.md", typed("G-004", "work", "done", "depends_on: [\"G-001\"]\n"))
	put(t, root, "grove/G-005-reclassified.md", typed("G-005", "term", "settled", "formerly: \"docs/old.md\"\n"))
	put(t, root, "grove/notes.txt", "ignored")
	p, ds := Load(root, root)
	if len(ds) != 0 {
		t.Fatalf("unexpected diagnostics:\n%s", diagnostics(ds))
	}
	if len(p.Records) != 5 {
		t.Fatalf("got %d records", len(p.Records))
	}
	for _, r := range p.Records {
		if r.ID == "G-002" && (r.Type != "page" || r.Status != "" || r.Title != "Synthesis") {
			t.Fatalf("page not read: %+v", r)
		}
	}
}

func TestRecordProblems(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct{ name, path, source, want string }{
		{"page with a lifecycle", "grove/a.md", "---\nid: G-002\ntype: page\ntitle: T\nstatus: proposed\n---\n", "status: unknown field"},
		{"page with a work field", "grove/a.md", "---\nid: G-002\ntype: page\ntitle: T\npriority: 1\n---\n", "priority: unknown field"},
		{"no type is not a page", "grove/a.md", "---\nid: G-002\ntitle: T\n---\n", "type: required field is missing"},
		{"unknown type", "grove/a.md", typed("G-002", "note", "open", ""), "type: unknown record type"},
		{"work record missing status", "grove/a.md", "---\nid: G-002\ntype: work\ntitle: T\n---\n", "status: required field is missing"},
		{"work with a bad status", "grove/a.md", typed("G-002", "work", "settled", ""), "status: unsupported lifecycle value for work"},
		{"plain Markdown", "grove/a.md", "# Just prose\n", "frontmatter: expected an opening --- line"},
		{"uncanonical ID", "grove/a.md", typed("G-02", "work", "proposed", ""), "id: expected a canonical positive ID, e.g. G-001"},
		{"unknown prefix", "grove/a.md", typed("X-002", "work", "proposed", ""), "id: expected a canonical positive ID"},
		{"duplicate ID across folders", "grove/x/y/a.md", typed("G-001", "term", "settled", ""), "id: duplicate G-001 in grove/G-001.md, grove/x/y/a.md"},
		{"non-work target", "grove/a.md", typed("G-002", "plan", "current", "work: [\"G-009\"]\n"), "work: target G-009 must be work"},
		{"formerly twice", "grove/a.md", typed("G-002", "work", "done", "formerly: \"W-007\"\n"), "formerly: W-007 was already converted to G-001"},
		{"formerly twice by another case", "grove/a.md", typed("G-002", "work", "done", "formerly: \"w-007\"\n"), "formerly: w-007 was already converted to G-001"},
		{"formerly still present", "grove/a.md", typed("G-002", "work", "done", "formerly: \"G-009\"\n"), "formerly: G-009 still exists in grove/G-009.md"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := schema3(t, "")
			put(t, root, "grove/G-001.md", typed("G-001", "work", "proposed", "formerly: \"W-007\"\n"))
			put(t, root, "grove/G-009.md", "---\nid: G-009\ntype: page\ntitle: P\n---\n")
			put(t, root, tc.path, tc.source)
			_, ds := Load(root, root)
			if got := diagnostics(ds); !strings.Contains(got, tc.path+":") && !strings.Contains(got, tc.path+": ") || !strings.Contains(got, tc.want) {
				t.Fatalf("wanted %q naming %s; got:\n%s", tc.want, tc.path, got)
			}
		})
	}
}
