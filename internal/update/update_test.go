package update

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/mascah/grove/internal/create"
	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/repo"
)

const (
	work     = "---\nid: \"G-001\"\ntype: work\ntitle: First\nstatus: proposed\nrelates_to: [\"G-003\"]\ncreated: \"2026-09-19T12:00:00Z\"\nupdated: \"2026-09-19T12:00:00Z\"\n---\n\n## Outcome\n\nBody --- stays.\n"
	second   = "---\nid: \"G-002\"\ntype: work\ntitle: Second\nstatus: proposed\n---\nBody.\n"
	question = "---\nid: \"G-003\"\ntype: question\ntitle: Which?\nstatus: open\n---\nBody.\n"
	decision = "---\nid: \"G-004\"\ntype: decision\ntitle: Choose\nstatus: proposed\n---\nBody.\n"
)

var now = time.Date(2026, 9, 19, 18, 30, 0, 500, time.UTC)

func git(t *testing.T, dir string, args ...string) string {
	t.Helper()
	full := append([]string{"-c", "user.name=t", "-c", "user.email=t@t", "-c", "commit.gpgsign=false", "-c", "maintenance.auto=false", "-C", dir}, args...)
	out, err := exec.Command("git", full...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

func write(t *testing.T, root, path, source string) {
	t.Helper()
	path = filepath.Join(root, path)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
}

func read(t *testing.T, root, path string) string {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, path))
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

// gitProject returns a committed Git project with G-001, G-002, G-003, G-004.
func gitProject(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	git(t, root, "init", "-q", "-b", "main")
	write(t, root, "grove.yaml", "schema_version: 3\nrecords: grove\n")
	write(t, root, "grove/work/G-001-first.md", work)
	write(t, root, "grove/work/G-002-second.md", second)
	write(t, root, "grove/questions/G-003-which.md", question)
	write(t, root, "grove/decisions/G-004-choose.md", decision)
	git(t, root, "add", "-A")
	git(t, root, "commit", "-q", "-m", "init")
	return root
}

func revision(t *testing.T, root, path string) string {
	t.Helper()
	return project.Revision([]byte(read(t, root, path)))
}

func record(t *testing.T, root, id string) *project.Record {
	t.Helper()
	p, ds := project.Load(root, root)
	if len(ds) != 0 {
		t.Fatalf("project invalid: %v", ds)
	}
	i := slices.IndexFunc(p.Records, func(r *project.Record) bool { return r.ID == id })
	if i < 0 {
		t.Fatalf("%s missing", id)
	}
	return p.Records[i]
}

func apply(t *testing.T, root, id string, sets []Field, unsets ...string) Result {
	t.Helper()
	r := record(t, root, id)
	res, err := Apply(root, Request{ID: id, Expect: project.Revision(r.Source), Set: sets, Unset: unsets}, now, nil)
	if err != nil {
		t.Fatalf("update %s: %v", id, err)
	}
	if res.ID != id || res.Path != r.Path || res.Revision != revision(t, root, r.Path) {
		t.Fatalf("result %+v does not describe the file", res)
	}
	return res
}

func tempFiles(t *testing.T, root string) []string {
	t.Helper()
	var found []string
	filepath.WalkDir(filepath.Join(root, "grove"), func(path string, entry os.DirEntry, err error) error {
		if err == nil && !entry.IsDir() && filepath.Ext(path) != ".md" {
			found = append(found, path)
		}
		return nil
	})
	return found
}

func TestUpdateEveryFieldOnItsTypes(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	res := apply(t, root, "G-001", []Field{
		{"title", ` "Quoted": ünïcode — 版本 #1 `}, {"status", "active"}, {"kind", "tooling"}, {"priority", "1"},
		{"size", "large"}, {"members", `["G-002"]`}, {"depends_on", `["G-002"]`}, {"relates_to", `["G-003", "G-004"]`},
	})
	if !res.Changed {
		t.Fatal("expected a change")
	}
	r := record(t, root, "G-001")
	if r.Title != ` "Quoted": ünïcode — 版本 #1 ` || r.Status != "active" || r.Kind != "tooling" || *r.Priority != 1 || r.Size != "large" ||
		!slices.Equal(r.Members, []string{"G-002"}) || !slices.Equal(r.DependsOn, []string{"G-002"}) || !slices.Equal(r.RelatesTo, []string{"G-003", "G-004"}) {
		t.Fatalf("fields not applied: %+v", r)
	}
	if r.Created.Format(time.RFC3339) != "2026-09-19T12:00:00Z" || r.Updated.Format(time.RFC3339) != "2026-09-19T18:30:00Z" {
		t.Fatalf("dates: created %v updated %v", r.Created, r.Updated)
	}
	source := string(r.Source)
	if !strings.HasSuffix(source, "---\n\n## Outcome\n\nBody --- stays.\n") || !strings.HasPrefix(source, "---\nid: \"G-001\"\ntype: work\ntitle: \" \\\"Quoted\\\": ünïcode — 版本 #1 \"\nstatus: active\nrelates_to: [\"G-003\", \"G-004\"]\ncreated: \"2026-09-19T12:00:00Z\"\nupdated: \"2026-09-19T18:30:00Z\"\nkind: tooling\npriority: 1\nsize: large\nmembers: [\"G-002\"]\ndepends_on: [\"G-002\"]\n---") {
		t.Fatalf("unexpected source:\n%s", source)
	}
	apply(t, root, "G-003", []Field{{"blocks", `["G-001"]`}, {"status", "resolved"}, {"title", "Resolved?"}, {"relates_to", "[]"}})
	q := record(t, root, "G-003")
	if !slices.Equal(q.Blocks, []string{"G-001"}) || q.Status != "resolved" || q.RelatesTo == nil || len(q.RelatesTo) != 0 || q.Created != nil || q.Updated == nil {
		t.Fatalf("question: %+v", q)
	}
	apply(t, root, "G-004", []Field{{"status", "accepted"}, {"relates_to", `["G-001"]`}})
	if d := record(t, root, "G-004"); d.Status != "accepted" || d.Created != nil {
		t.Fatalf("decision: %+v", d)
	}
	// Optional removal and explicit empty lists.
	apply(t, root, "G-001", []Field{{"members", "[]"}}, "kind", "priority", "size", "depends_on")
	r = record(t, root, "G-001")
	if r.Kind != "" || r.Priority != nil || r.Size != "" || r.DependsOn != nil || r.Members == nil || len(r.Members) != 0 {
		t.Fatalf("removal: %+v", r)
	}
	if strings.Contains(string(r.Source), "kind") || !strings.Contains(string(r.Source), "members: []\n") {
		t.Fatalf("source after removal:\n%s", r.Source)
	}
	if files := tempFiles(t, root); len(files) != 0 {
		t.Fatalf("temporary files left behind: %v", files)
	}
}

func TestUpdateLifecycleAndReopening(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	// Untouched optional fields, including the pointer-valued priority, must
	// compare equal between the original and the candidate.
	apply(t, root, "G-001", []Field{{"priority", "2"}, {"size", "small"}, {"members", `["G-002"]`}})
	head := git(t, root, "rev-parse", "HEAD")
	for _, status := range []string{"active", "review", "done", "proposed", "abandoned", "done"} {
		sets := []Field{{"status", status}}
		if status == "review" { // review needs the candidate, and done keeps it
			sets = append(sets, Field{"candidate", head})
		}
		apply(t, root, "G-001", sets)
		r := record(t, root, "G-001")
		if r.Status != status || *r.Priority != 2 || r.Size != "small" || r.Created.Format(time.RFC3339) != "2026-09-19T12:00:00Z" || !strings.HasSuffix(string(r.Source), "\n---\n\n## Outcome\n\nBody --- stays.\n") {
			t.Fatalf("%s: %+v", status, r)
		}
		if status != "active" && (r.Candidate != head || !strings.Contains(string(r.Source), "candidate: \""+head+"\"")) {
			t.Fatalf("%s: candidate must stay and stay quoted: %+v", status, r)
		}
	}
	for _, status := range []string{"resolved", "open"} {
		apply(t, root, "G-003", []Field{{"status", status}})
	}
	for _, status := range []string{"accepted", "rejected", "proposed"} {
		apply(t, root, "G-004", []Field{{"status", status}})
	}
	if entries, _ := os.ReadDir(filepath.Join(root, "grove/work")); len(entries) != 2 {
		t.Fatalf("status changes must not rename or add files: %v", entries)
	}
}

func TestUpdateNoOpPreservesBytesAndRefusesStale(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	path := filepath.Join(root, "grove/work/G-001-first.md")
	if err := os.Chmod(path, 0o600); err != nil {
		t.Fatal(err)
	}
	past := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	os.Chtimes(path, past, past)
	before, _ := os.Stat(path)
	expect := revision(t, root, "grove/work/G-001-first.md")
	for _, req := range []Request{
		{Set: []Field{{"title", "First"}, {"status", "proposed"}, {"relates_to", `["G-003"]`}}},
		{Unset: []string{"kind", "priority", "members"}},
		{Set: []Field{{"relates_to", `["G-003"]`}}, Unset: []string{"size"}},
	} {
		req.ID, req.Expect = "G-001", expect
		// The clock is behind the record's dates; a no-op needs no clock.
		res, err := Apply(root, req, past, nil)
		if err != nil || res.Changed || res.Revision != expect {
			t.Fatalf("%+v: %+v %v", req, res, err)
		}
	}
	after, _ := os.Stat(path)
	if read(t, root, "grove/work/G-001-first.md") != work || !after.ModTime().Equal(before.ModTime()) || after.Mode() != before.Mode() {
		t.Fatal("no-op changed bytes, mtime, or permissions")
	}
	_, err := Apply(root, Request{ID: "G-001", Expect: "sha256:" + strings.Repeat("0", 64), Set: []Field{{"status", "proposed"}}}, now, nil)
	if err == nil || !strings.Contains(err.Error(), "changed since the expected revision") || !strings.Contains(err.Error(), expect) {
		t.Fatalf("a stale no-op must be refused and report the current revision: %v", err)
	}
	// Absent versus explicitly empty are different states.
	if res := apply(t, root, "G-001", []Field{{"members", "[]"}}); !res.Changed {
		t.Fatal("adding an empty list to an absent field is a change")
	}
	if res := apply(t, root, "G-001", []Field{{"members", "[]"}}); res.Changed {
		t.Fatal("an explicit empty list already present is a no-op")
	}
	if res := apply(t, root, "G-001", nil, "members"); !res.Changed {
		t.Fatal("removing a present empty list is a change")
	}
	if res := apply(t, root, "G-001", []Field{{"relates_to", `["G-003"]`}, {"members", "[]"}}, "kind"); !res.Changed {
		t.Fatal("one effective change among no-ops still applies")
	}
}

func TestUpdateRejectsInvalidRequestsWithoutWriting(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	expect := revision(t, root, "grove/work/G-001-first.md")
	for _, tc := range []struct {
		name string
		id   string
		set  []Field
		un   []string
		want string
	}{
		{"unknown field", "G-001", []Field{{"foo", "x"}}, nil, "foo"},
		{"id", "G-001", []Field{{"id", "G-009"}}, nil, "id cannot"},
		{"created", "G-001", []Field{{"created", "\"2026-01-01T00:00:00Z\""}}, nil, "created cannot"},
		{"updated", "G-001", nil, []string{"updated"}, "updated cannot"},
		{"unset id", "G-001", nil, []string{"id"}, "id cannot"},
		{"unset required", "G-001", nil, []string{"title"}, "required"},
		{"wrong type field", "G-001", []Field{{"blocks", "[]"}}, nil, "blocks"},
		{"empty title", "G-001", []Field{{"title", "  "}}, nil, "title"},
		{"status value", "G-001", []Field{{"status", "resolved"}}, nil, "status"},
		{"status injection", "G-001", []Field{{"status", "active\nkind: fix"}}, nil, "invalid"},
		{"kind value", "G-001", []Field{{"kind", "magic"}}, nil, "kind"},
		{"size value", "G-001", []Field{{"size", "huge"}}, nil, "size"},
		{"priority letters", "G-001", []Field{{"priority", "high"}}, nil, "priority"},
		{"priority negative", "G-001", []Field{{"priority", "-1"}}, nil, "priority"},
		{"priority range", "G-001", []Field{{"priority", "6"}}, nil, "priority"},
		{"list null", "G-001", []Field{{"members", "null"}}, nil, "members"},
		{"list not json", "G-001", []Field{{"members", "G-002"}}, nil, "members"},
		{"list numbers", "G-001", []Field{{"members", "[1]"}}, nil, "members"},
		{"list duplicate", "G-001", []Field{{"members", `["G-002", "G-002"]`}}, nil, "duplicate"},
		{"self link", "G-001", []Field{{"depends_on", `["G-001"]`}}, nil, "self"},
		{"unresolved", "G-001", []Field{{"depends_on", `["G-404"]`}}, nil, "unresolved"},
		{"non-work target", "G-001", []Field{{"members", `["G-003"]`}}, nil, "must be work"},
		{"missing record", "G-404", []Field{{"status", "active"}}, nil, "not found"},
		{"invalid utf8", "G-001", []Field{{"title", "bad\xff"}}, nil, "UTF-8"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := Apply(root, Request{ID: tc.id, Expect: expect, Set: tc.set, Unset: tc.un}, now, nil)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("wanted %q, got %v", tc.want, err)
			}
			if read(t, root, "grove/work/G-001-first.md") != work {
				t.Fatal("a rejected request changed the file")
			}
		})
	}
	if files := tempFiles(t, root); len(files) != 0 {
		t.Fatalf("temporary files left behind: %v", files)
	}
}

func TestUpdateRejectsCyclesAndInvalidProjects(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	apply(t, root, "G-001", []Field{{"depends_on", `["G-002"]`}, {"members", `["G-002"]`}})
	expect := revision(t, root, "grove/work/G-002-second.md")
	for _, field := range []string{"depends_on", "members"} {
		_, err := Apply(root, Request{ID: "G-002", Expect: expect, Set: []Field{{field, `["G-001"]`}}}, now, nil)
		if err == nil || !strings.Contains(err.Error(), "cycle") {
			t.Fatalf("%s: wanted a cycle refusal, got %v", field, err)
		}
	}
	// A group may depend on its own members; that is not a cycle.
	apply(t, root, "G-002", []Field{{"relates_to", `["G-001"]`}})
	write(t, root, "grove/work/broken.md", "---\nid: \"G-005\"\ntype: work\ntitle: Broken\nstatus: imaginary\n---\n")
	_, err := Apply(root, Request{ID: "G-001", Expect: revision(t, root, "grove/work/G-001-first.md"), Set: []Field{{"status", "active"}}}, now, nil)
	if err == nil || !strings.Contains(err.Error(), "not valid") || !strings.Contains(err.Error(), "broken.md") {
		t.Fatalf("an invalid project must be refused: %v", err)
	}
	if strings.Contains(read(t, root, "grove/work/G-001-first.md"), "status: active") {
		t.Fatal("refused update was written")
	}
}

func TestUpdateClockContract(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	expect := revision(t, root, "grove/work/G-001-first.md")
	early := time.Date(2026, 9, 19, 11, 59, 59, 0, time.UTC)
	_, err := Apply(root, Request{ID: "G-001", Expect: expect, Set: []Field{{"status", "active"}}}, early, nil)
	if err == nil || !strings.Contains(err.Error(), "clock") || read(t, root, "grove/work/G-001-first.md") != work {
		t.Fatalf("a clock behind the record must be refused without writing: %v", err)
	}
	// Equal to created is acceptable; two changes within one second share a stamp.
	same := time.Date(2026, 9, 19, 12, 0, 0, 999_999_999, time.UTC)
	first, err := Apply(root, Request{ID: "G-001", Expect: expect, Set: []Field{{"status", "active"}}}, same, nil)
	if err != nil {
		t.Fatal(err)
	}
	secondRes, err := Apply(root, Request{ID: "G-001", Expect: first.Revision, Set: []Field{{"status", "abandoned"}}}, same, nil)
	if err != nil || secondRes.Revision == first.Revision {
		t.Fatalf("same-second change: %+v %v", secondRes, err)
	}
	r := record(t, root, "G-001")
	if r.Updated.Format(time.RFC3339) != "2026-09-19T12:00:00Z" || r.Created.Format(time.RFC3339) != "2026-09-19T12:00:00Z" {
		t.Fatalf("stamp must truncate to seconds and preserve created: %+v", r)
	}
	// A record without created stays without one.
	apply(t, root, "G-002", []Field{{"status", "active"}})
	if r := record(t, root, "G-002"); r.Created != nil || r.Updated == nil || !strings.HasSuffix(string(r.Source), "status: active\nupdated: \"2026-09-19T18:30:00Z\"\n---\nBody.\n") {
		t.Fatalf("missing created must stay absent: %s", r.Source)
	}
}

func TestUpdateStaleAfterBodyOnlyEdit(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	expect := revision(t, root, "grove/work/G-001-first.md")
	write(t, root, "grove/work/G-001-first.md", work+"An appended paragraph with unchanged timestamps.\n")
	_, err := Apply(root, Request{ID: "G-001", Expect: expect, Set: []Field{{"status", "active"}}}, now, nil)
	if err == nil || !strings.Contains(err.Error(), "changed since") {
		t.Fatalf("body edits must invalidate the revision: %v", err)
	}
}

func TestUpdateDetectsChangesDuringPreparation(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name   string
		change func(t *testing.T, root string)
	}{
		{"neighbor body", func(t *testing.T, root string) { write(t, root, "grove/work/G-002-second.md", second+"more\n") }},
		{"inventory", func(t *testing.T, root string) {
			write(t, root, "grove/work/G-005-new.md", strings.Replace(second, "G-002", "G-005", 1))
		}},
		{"configuration", func(t *testing.T, root string) {
			os.Rename(filepath.Join(root, "grove"), filepath.Join(root, "records"))
			write(t, root, "grove.yaml", "schema_version: 3\nrecords: records\n")
		}},
		{"configuration comment only", func(t *testing.T, root string) {
			write(t, root, "grove.yaml", "# concurrent edit\nschema_version: 3\nrecords: grove\n")
		}},
		{"target permissions", func(t *testing.T, root string) { os.Chmod(filepath.Join(root, "grove/work/G-001-first.md"), 0o600) }},
		{"target replaced", func(t *testing.T, root string) {
			os.Remove(filepath.Join(root, "grove/work/G-001-first.md"))
			write(t, root, "grove/work/G-001-first.md", work)
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			root := gitProject(t)
			expect := revision(t, root, "grove/work/G-001-first.md")
			fault := func(step string) error {
				if step == "compare" {
					tc.change(t, root)
				}
				return nil
			}
			_, err := Apply(root, Request{ID: "G-001", Expect: expect, Set: []Field{{"status", "active"}}}, now, fault)
			if err == nil || errors.As(err, new(*Failure)) {
				t.Fatalf("expected a refusal before publication, got %v", err)
			}
			for _, dir := range []string{"grove", "records"} {
				if data, err := os.ReadFile(filepath.Join(root, dir, "work/G-001-first.md")); err == nil && string(data) != work {
					t.Fatal("refused update was written")
				}
			}
			if files := tempFiles(t, root); len(files) != 0 {
				t.Fatalf("temporary files left behind: %v", files)
			}
		})
	}
}

func TestUpdateInjectedFailures(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	path := filepath.Join(root, "grove/work/G-001-first.md")
	if err := os.Chmod(path, 0o640); err != nil {
		t.Fatal(err)
	}
	expect := revision(t, root, "grove/work/G-001-first.md")
	req := Request{ID: "G-001", Expect: expect, Set: []Field{{"status", "active"}}}
	for _, step := range []string{"write", "sync", "close", "compare", "rename"} {
		_, err := Apply(root, req, now, func(s string) error {
			if s == step {
				return errors.New("injected " + step)
			}
			return nil
		})
		if err == nil || !strings.Contains(err.Error(), "injected "+step) || !strings.Contains(err.Error(), "unchanged") || errors.As(err, new(*Failure)) {
			t.Fatalf("%s: %v", step, err)
		}
		if read(t, root, "grove/work/G-001-first.md") != work {
			t.Fatalf("%s: original bytes lost", step)
		}
		if files := tempFiles(t, root); len(files) != 0 {
			t.Fatalf("%s: temporary files left behind: %v", step, files)
		}
	}
	for _, step := range []string{"dirsync", "validate"} {
		root := gitProject(t)
		if err := os.Chmod(filepath.Join(root, "grove/work/G-001-first.md"), 0o640); err != nil {
			t.Fatal(err)
		}
		_, err := Apply(root, req, now, func(s string) error {
			if s == step {
				return errors.New("injected " + step)
			}
			return nil
		})
		var failure *Failure
		if !errors.As(err, &failure) || !strings.Contains(err.Error(), "applied") || failure.Path != "grove/work/G-001-first.md" {
			t.Fatalf("%s: expected an applied-state failure, got %v", step, err)
		}
		got := read(t, root, "grove/work/G-001-first.md")
		if !strings.Contains(got, "status: active") || failure.Revision != project.Revision([]byte(got)) {
			t.Fatalf("%s: applied state must describe the published bytes", step)
		}
		if info, _ := os.Stat(filepath.Join(root, "grove/work/G-001-first.md")); info.Mode().Perm() != 0o640 {
			t.Fatalf("%s: permissions not preserved: %v", step, info.Mode())
		}
		if step == "dirsync" && !strings.Contains(err.Error(), "durability") {
			t.Fatalf("directory sync failure must report uncertain durability: %v", err)
		}
	}
}

func TestUpdateSameRevisionRace(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	expect := revision(t, root, "grove/work/G-001-first.md")
	results := make(chan error, 2)
	var wg sync.WaitGroup
	for _, status := range []string{"active", "abandoned"} {
		wg.Add(1)
		go func() {
			defer wg.Done()
			res, err := Apply(root, Request{ID: "G-001", Expect: expect, Set: []Field{{"status", status}}}, now, nil)
			if err == nil && !res.Changed {
				err = errors.New("unexpected no-op")
			}
			results <- err
		}()
	}
	wg.Wait()
	close(results)
	var failures int
	for err := range results {
		if err != nil {
			failures++
			if !strings.Contains(err.Error(), "changed since") {
				t.Fatal(err)
			}
		}
	}
	if failures != 1 {
		t.Fatalf("exactly one writer must win, got %d refusals", failures)
	}
	if r := record(t, root, "G-001"); r.Status != "active" && r.Status != "abandoned" {
		t.Fatalf("winning change lost: %+v", r)
	}
}

func TestUpdateReciprocalDependenciesCannotFormACycle(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	pairs := [][2]string{{"G-001", "G-002"}, {"G-002", "G-001"}}
	errs := make([]error, 2)
	var wg sync.WaitGroup
	for i, pair := range pairs {
		wg.Add(1)
		expect := revision(t, root, "grove/work/"+pair[0]+"-"+map[string]string{"G-001": "first", "G-002": "second"}[pair[0]]+".md")
		go func() {
			defer wg.Done()
			_, errs[i] = Apply(root, Request{ID: pair[0], Expect: expect, Set: []Field{{"depends_on", `["` + pair[1] + `"]`}}}, now, nil)
		}()
	}
	wg.Wait()
	if (errs[0] == nil) == (errs[1] == nil) {
		t.Fatalf("exactly one must succeed: %v / %v", errs[0], errs[1])
	}
	if _, ds := project.Load(root, root); len(ds) != 0 {
		t.Fatalf("project invalid after concurrent updates: %v", ds)
	}
}

func TestNewAndUpdateShareTheWriteLockAcrossWorktrees(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	wt := filepath.Join(filepath.Dir(root), filepath.Base(root)+"-wt")
	git(t, root, "worktree", "add", "-q", "-b", "feature", wt)
	common, _, err := repo.CommonDir(root)
	if err != nil {
		t.Fatal(err)
	}
	unlock, err := repo.WriteLock(common)
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan string, 2)
	go func() {
		_, err := Apply(wt, Request{ID: "G-001", Expect: revision(t, wt, "grove/work/G-001-first.md"), Set: []Field{{"status", "active"}}}, now, nil)
		if err != nil {
			t.Error(err)
		}
		done <- "update"
	}()
	go func() {
		p, ds := project.Load(root, root)
		if len(ds) != 0 {
			t.Error(ds)
		}
		if _, err := create.New(p, "work", "Third", "third", now, &bytes.Buffer{}); err != nil {
			t.Error(err)
		}
		done <- "new"
	}()
	select {
	case who := <-done:
		t.Fatalf("%s completed while the write lock was held", who)
	case <-time.After(300 * time.Millisecond):
	}
	unlock()
	for range 2 {
		select {
		case <-done:
		case <-time.After(10 * time.Second):
			t.Fatal("commands did not proceed after the lock was released")
		}
	}
	if r := record(t, wt, "G-001"); r.Status != "active" {
		t.Fatal("worktree update lost")
	}
	if p, ds := project.Load(root, root); len(ds) != 0 || len(p.Records) != 5 {
		t.Fatalf("main checkout: %v", ds)
	}
	if _, err := os.Stat(filepath.Join(common, "grove", "write.lock")); err != nil {
		t.Fatal("write lock file must remain")
	}
}

func TestUpdateRequiresGit(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	write(t, root, "grove.yaml", "schema_version: 3\nrecords: grove\n")
	write(t, root, "grove/work/G-001-first.md", work)
	write(t, root, "grove/questions/G-003-which.md", question)
	_, err := Apply(root, Request{ID: "G-001", Expect: revision(t, root, "grove/work/G-001-first.md"), Set: []Field{{"status", "active"}}}, now, nil)
	if err == nil || !strings.Contains(err.Error(), "Git") || read(t, root, "grove/work/G-001-first.md") != work {
		t.Fatalf("expected a Git requirement without writes: %v", err)
	}
}

func TestUnchangedGuardCatchesEditorDrift(t *testing.T) {
	t.Parallel()
	before, _ := project.ParseRecord("w.md", []byte(work))
	after, _ := project.ParseRecord("w.md", []byte(strings.Replace(work, "title: First", "title: Other", 1)))
	if err := unchanged(before, after, []change{set("status", "active")}); err == nil || !strings.Contains(err.Error(), "title") {
		t.Fatalf("an untouched field that differs must be refused: %v", err)
	}
	if err := unchanged(before, after, []change{set("title", `"Other"`)}); err != nil {
		t.Fatal(err)
	}
	priority := 2
	before.Priority, after.Priority = &priority, new(int)
	*after.Priority = 2
	after.Title = before.Title
	if err := unchanged(before, after, nil); err != nil {
		t.Fatalf("equal priorities behind different pointers must compare equal: %v", err)
	}
}

// G-015: the review's update reproducers, at the level a user reaches them.
func TestUpdatePreservesAcceptedForms(t *testing.T) {
	t.Parallel()
	const stamped = "updated: \"2026-09-19T18:30:00Z\""
	t.Run("comment before a later-line value", func(t *testing.T) {
		t.Parallel()
		root := gitProject(t)
		write(t, root, "grove/work/G-001-first.md", strings.Replace(work, "title: First", "title: # retain\n  First", 1))
		apply(t, root, "G-001", []Field{{"title", "New"}})
		want := strings.NewReplacer("title: First", "title: # retain\n  \"New\"", "updated: \"2026-09-19T12:00:00Z\"", stamped).Replace(work)
		if got := read(t, root, "grove/work/G-001-first.md"); got != want {
			t.Fatalf("got:\n%s\nwant:\n%s", got, want)
		}
	})
	t.Run("final two flow entries", func(t *testing.T) {
		t.Parallel()
		for _, unsets := range [][]string{{"kind", "size"}, {"size", "kind"}} {
			root := gitProject(t)
			write(t, root, "grove/work/G-001-first.md", "---\n{id: G-001, type: work, title: T, status: proposed, kind: fix, size: small}\n---\nBody\n")
			apply(t, root, "G-001", nil, unsets...)
			want := "---\n{id: G-001, type: work, title: T, status: proposed, " + stamped + "}\n---\nBody\n"
			if got := read(t, root, "grove/work/G-001-first.md"); got != want {
				t.Fatalf("got:\n%s\nwant:\n%s", got, want)
			}
		}
	})
	t.Run("explicit keys", func(t *testing.T) {
		t.Parallel()
		root := gitProject(t)
		source := strings.NewReplacer("status: proposed", "? status\n: proposed", "relates_to: [\"G-003\"]", "? relates_to\n: [\"G-003\"] # why").Replace(work)
		write(t, root, "grove/work/G-001-first.md", source)
		apply(t, root, "G-001", []Field{{"status", "active"}}, "relates_to")
		want := strings.NewReplacer("status: proposed", "? status\n: active", "relates_to: [\"G-003\"]\n", "", "updated: \"2026-09-19T12:00:00Z\"", stamped).Replace(work)
		if got := read(t, root, "grove/work/G-001-first.md"); got != want {
			t.Fatalf("got:\n%s\nwant:\n%s", got, want)
		}
	})
}

// G-016: a main checkout named "new\nline" once put the write lock in a
// sibling ".../new/grove". Locks and counters belong under the real common
// directory, from the main and a linked checkout alike.
func TestCoordinationStateStaysUnderTheCommonDirectory(t *testing.T) {
	t.Parallel()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	parent, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	root, wt := filepath.Join(parent, "new\nline"), filepath.Join(parent, "linked\twt ")
	write(t, root, "grove.yaml", "schema_version: 3\nrecords: grove\n")
	write(t, root, "grove/work/G-001-first.md", work)
	write(t, root, "grove/questions/G-003-which.md", question)
	git(t, root, "init", "-q", "-b", "main")
	git(t, root, "add", "-A")
	git(t, root, "commit", "-q", "-m", "init")
	git(t, root, "worktree", "add", "-q", "-b", "feature", wt)

	apply(t, root, "G-001", []Field{{"status", "active"}})
	apply(t, wt, "G-001", []Field{{"status", "abandoned"}})
	for _, dir := range []string{root, wt} {
		p, ds := project.Load(dir, dir)
		if len(ds) != 0 {
			t.Fatal(ds)
		}
		if _, err := create.New(p, "work", "Odd", "odd", now, &bytes.Buffer{}); err != nil {
			t.Fatal(err)
		}
	}
	// The second creation came from the other checkout and the shared counter.
	if _, err := os.Stat(filepath.Join(wt, "grove/G-005-odd.md")); err != nil {
		t.Fatalf("linked checkouts must share one ID sequence: %v", err)
	}
	var state []string
	filepath.WalkDir(parent, func(path string, entry os.DirEntry, err error) error {
		if err == nil && slices.Contains([]string{"write.lock", "lock", "neutral-ids"}, entry.Name()) {
			state = append(state, path)
		}
		return nil
	})
	slices.Sort(state)
	common := filepath.Join(root, ".git", "grove")
	if want := []string{filepath.Join(common, "lock"), filepath.Join(common, "neutral-ids"), filepath.Join(common, "write.lock")}; !slices.Equal(state, want) {
		t.Fatalf("coordination state: %q\nwant %q", state, want)
	}
	if entries, _ := os.ReadDir(parent); len(entries) != 2 {
		t.Fatalf("nothing may appear beside the checkouts: %q", entries)
	}
}

// TestUpdateDoneMeansAnIntegratedCandidate covers the review lifecycle's one
// write-side rule: done needs a candidate that HEAD contains, so it is written
// on the target after the merge, while a done record from before the rule
// stays editable without one.
func TestUpdateDoneMeansAnIntegratedCandidate(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	base := git(t, root, "rev-parse", "HEAD")
	refuse := func(id, want string, sets ...Field) {
		t.Helper()
		r := record(t, root, id)
		_, err := Apply(root, Request{ID: id, Expect: project.Revision(r.Source), Set: sets}, now, nil)
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Fatalf("%s %v: got %v, want %q", id, sets, err, want)
		}
		if after := record(t, root, id); !bytes.Equal(after.Source, r.Source) {
			t.Fatalf("%s was written despite the refusal", id)
		}
	}
	// A candidate value of letters only must still be written as a quoted string.
	apply(t, root, "G-002", []Field{{"candidate", "abcdefa"}})
	if src := read(t, root, "grove/work/G-002-second.md"); !strings.Contains(src, "candidate: \"abcdefa\"\n") {
		t.Fatalf("candidate must be quoted:\n%s", src)
	}
	refuse("G-001", "candidate: required while status is review", Field{"status", "review"})
	refuse("G-001", "set candidate=COMMIT", Field{"status", "done"})
	refuse("G-001", "could not be checked against this checkout's HEAD", Field{"status", "done"}, Field{"candidate", strings.Repeat("a", 40)})
	// A candidate on an unmerged branch is refused on main until it is merged.
	git(t, root, "checkout", "-q", "-b", "feature")
	write(t, root, "grove/work/G-001-first.md", strings.Replace(work, "status: proposed", "status: active", 1))
	git(t, root, "commit", "-qam", "implement")
	candidate := git(t, root, "rev-parse", "HEAD")
	// On the work branch the candidate is an ancestor too: the CLI cannot tell
	// the target from the branch, so the guide, not this check, keeps done off
	// the branch. Undo it so the merge below is a fast-forward.
	apply(t, root, "G-001", []Field{{"status", "done"}, {"candidate", candidate}})
	git(t, root, "checkout", "-q", "--", ".")
	git(t, root, "checkout", "-q", "main")
	refuse("G-001", "is not an ancestor of this checkout's HEAD; mark done where it was merged", Field{"status", "done"}, Field{"candidate", candidate})
	git(t, root, "merge", "-q", "--ff-only", "feature")
	apply(t, root, "G-001", []Field{{"status", "done"}, {"candidate", candidate}})
	// Changing a done record's candidate is judged again; its other fields are
	// not, and the candidate cannot be dropped.
	refuse("G-001", "could not be checked against this checkout's HEAD", Field{"candidate", strings.Repeat("b", 40)})
	{
		r := record(t, root, "G-001")
		_, err := Apply(root, Request{ID: "G-001", Expect: project.Revision(r.Source), Unset: []string{"candidate"}}, now, nil)
		if err == nil || !strings.Contains(err.Error(), "stays done cannot lose its candidate") {
			t.Fatalf("unset candidate on a done record: %v", err)
		}
	}
	if r := record(t, root, "G-001"); r.Candidate != candidate {
		t.Fatalf("candidate changed: %+v", r)
	}
	apply(t, root, "G-001", []Field{{"candidate", base}, {"title", "Renamed"}})
	// Reopening may drop the candidate; closing from review while dropping it
	// gets the advice to set one, not the removal refusal.
	apply(t, root, "G-001", []Field{{"status", "review"}})
	{
		r := record(t, root, "G-001")
		_, err := Apply(root, Request{ID: "G-001", Expect: project.Revision(r.Source), Set: []Field{{"status", "done"}}, Unset: []string{"candidate"}}, now, nil)
		if err == nil || !strings.Contains(err.Error(), "set candidate=COMMIT") {
			t.Fatalf("done while dropping the candidate: %v", err)
		}
	}
	apply(t, root, "G-001", []Field{{"status", "active"}}, "candidate")
	if r := record(t, root, "G-001"); r.Candidate != "" || r.Status != "active" {
		t.Fatalf("reopening may drop the candidate: %+v", r)
	}
	// A historical done record has no candidate: editable, and never given one that HEAD lacks.
	write(t, root, "grove/work/G-002-second.md", strings.Replace(second, "status: proposed", "status: done", 1))
	apply(t, root, "G-002", []Field{{"title", "Still done"}})
	refuse("G-002", "could not be checked", Field{"candidate", strings.Repeat("c", 40)})
	if r := record(t, root, "G-002"); r.Candidate != "" || r.Title != "Still done" {
		t.Fatalf("historical done: %+v", r)
	}
}

// TestUpdateOptionalExpectAndCommit covers G-079: an omitted Expect applies to
// the file as it is while a stale one is still refused, and Commit commits the
// record's file alone, nothing for a no-op, and reports a failed commit as an
// applied update.
func TestUpdateOptionalExpectAndCommit(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	for _, kv := range [][2]string{{"user.name", "t"}, {"user.email", "t@t"}, {"commit.gpgsign", "false"}, {"maintenance.auto", "false"}} {
		git(t, root, "config", kv[0], kv[1]) // the product's commit uses the repository's own identity
	}
	head := func() string { return git(t, root, "rev-parse", "HEAD") }
	base := head()
	// A direct edit after the last commit: the omitted form applies to it, a stale revision is refused.
	write(t, root, "grove/work/G-001-first.md", strings.Replace(work, "Body --- stays.", "Edited body.", 1))
	if _, err := Apply(root, Request{ID: "G-001", Expect: project.Revision([]byte(work)), Set: []Field{{"status", "active"}}}, now, nil); err == nil || !strings.Contains(err.Error(), "changed since the expected revision") {
		t.Fatalf("stale expect: %v", err)
	}
	res, err := Apply(root, Request{ID: "G-001", Set: []Field{{"status", "active"}}}, now, nil)
	if err != nil || !res.Changed || res.Commit != "" || head() != base {
		t.Fatalf("omitted expect: %+v %v", res, err)
	}
	if src := read(t, root, "grove/work/G-001-first.md"); !strings.Contains(src, "status: active\n") || !strings.Contains(src, "Edited body.") {
		t.Fatalf("update not applied over the direct edit:\n%s", src)
	}
	// Other changes in the tree, staged and unstaged, are left alone by --commit.
	write(t, root, "grove/work/G-002-second.md", strings.Replace(second, "Second", "Second edited", 1))
	write(t, root, "notes.txt", "staged\n")
	git(t, root, "add", "notes.txt")
	res, err = Apply(root, Request{ID: "G-001", Set: []Field{{"status", "review"}, {"candidate", base}}, Commit: true}, now, nil)
	if err != nil || !res.Changed || res.Commit != head() || res.Commit == base {
		t.Fatalf("commit: %+v %v (HEAD %s)", res, err, head())
	}
	if files := git(t, root, "show", "--stat", "--format=", "--name-only", "HEAD"); files != "grove/work/G-001-first.md" {
		t.Fatalf("the commit must hold the record alone: %q", files)
	}
	if subject := git(t, root, "log", "-1", "--format=%s"); subject != "docs(G-001): set status=review candidate="+base {
		t.Fatalf("message: %q", subject)
	}
	if status := git(t, root, "status", "--porcelain"); status != "M grove/work/G-002-second.md\nA  notes.txt" { // git() trims the leading space
		t.Fatalf("other paths must stay as they were:\n%s", status)
	}
	// Acceptance 1: done from a clean tree in one step, with several fields and an unset in the message.
	git(t, root, "commit", "-qam", "the rest")
	clean := head()
	res, err = Apply(root, Request{ID: "G-001", Set: []Field{{"status", "done"}}, Unset: []string{"relates_to"}, Commit: true}, now, nil)
	if err != nil || res.Commit != head() || res.Commit == clean {
		t.Fatalf("done --commit: %+v %v", res, err)
	}
	if subject := git(t, root, "log", "-1", "--format=%s"); subject != "docs(G-001): set status=done unset relates_to" {
		t.Fatalf("message: %q", subject)
	}
	if status := git(t, root, "status", "--porcelain"); status != "" {
		t.Fatalf("tree must be clean after the commit:\n%s", status)
	}
	// A no-op commits nothing.
	done := head()
	res, err = Apply(root, Request{ID: "G-001", Set: []Field{{"status", "done"}}, Commit: true}, now, nil)
	if err != nil || res.Changed || res.Commit != "" || head() != done {
		t.Fatalf("no-op with commit: %+v %v", res, err)
	}
	// A commit that fails after publication: the file holds the update, the
	// error carries the applied revision, and nothing was committed.
	hooks := filepath.Join(root, "hooks")
	write(t, root, "hooks/pre-commit", "#!/bin/sh\necho refused by hook >&2\nexit 1\n")
	if err := os.Chmod(filepath.Join(hooks, "pre-commit"), 0o755); err != nil {
		t.Fatal(err)
	}
	git(t, root, "config", "core.hooksPath", hooks)
	res, err = Apply(root, Request{ID: "G-001", Set: []Field{{"title", "Renamed"}}, Commit: true}, now, nil)
	var failure *Failure
	if !errors.As(err, &failure) || failure.Revision != revision(t, root, "grove/work/G-001-first.md") || !strings.Contains(err.Error(), "refused by hook") || !strings.Contains(err.Error(), "the file is staged but nothing was committed") || !strings.Contains(err.Error(), "the update was applied to grove/work/G-001-first.md") {
		t.Fatalf("failed commit: %+v %v", res, err)
	}
	if head() != done || record(t, root, "G-001").Title != "Renamed" {
		t.Fatal("the file must hold the update and HEAD must not move")
	}
}
