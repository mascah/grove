package create

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/mascah/grove/internal/project"
)

const config = "schema_version: 3\nrecords: grove\n"

func record(id, kind, status string) string {
	return "---\nid: \"" + id + "\"\ntype: " + kind + "\ntitle: T\nstatus: " + status + "\n---\nBody.\n"
}

func git(t *testing.T, dir string, args ...string) string {
	t.Helper()
	full := append([]string{"-c", "user.name=t", "-c", "user.email=t@t", "-c", "commit.gpgsign=false", "-C", dir}, args...)
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

// gitProject returns a committed Git project containing G-001.
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
	write(t, root, "grove.yaml", config)
	write(t, root, "grove/G-001-first.md", record("G-001", "work", "done"))
	git(t, root, "add", "-A")
	git(t, root, "commit", "-q", "-m", "init")
	return root
}

func load(t *testing.T, root string) *project.Project {
	t.Helper()
	p, ds := project.Load(root, root)
	if len(ds) != 0 {
		t.Fatalf("unexpected diagnostics: %v", ds)
	}
	return p
}

func stateDir(t *testing.T, root string) string {
	t.Helper()
	return filepath.Join(git(t, root, "rev-parse", "--path-format=absolute", "--git-common-dir"), "grove")
}

func TestAllocateFloorsFromRefsAndWorktrees(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	wt := filepath.Join(filepath.Dir(root), filepath.Base(root)+"-wt")
	git(t, root, "worktree", "add", "-q", "-b", "feature", wt)
	// G-007 exists only in the feature branch's history, so the ref scan alone
	// can find it; G-005 exists only as a live file in the worktree. Both are
	// nested: the scan follows discovery, which is recursive.
	write(t, wt, "grove/deep/G-007-committed-on-branch.md", record("G-007", "work", "proposed"))
	git(t, wt, "add", "-A")
	git(t, wt, "commit", "-q", "-m", "branch record")
	if err := os.Remove(filepath.Join(wt, "grove/deep/G-007-committed-on-branch.md")); err != nil {
		t.Fatal(err)
	}
	write(t, wt, "grove/any/where/G-005-live-only.md", record("G-005", "work", "proposed"))

	var report bytes.Buffer
	n, err := Allocate(root, "grove", &report)
	if err != nil || n != 8 {
		t.Fatalf("got %d, %v; want 8 from committed G-007 and live G-005", n, err)
	}
	if !strings.Contains(strings.ToLower(report.String()), "initialized") {
		t.Fatalf("first allocation should report initialization: %q", report.String())
	}
	report.Reset()
	n, err = Allocate(wt, "grove", &report)
	if err != nil || n != 9 {
		t.Fatalf("got %d, %v; want 9 from the shared counter", n, err)
	}
	if report.Len() != 0 {
		t.Fatalf("steady-state allocation should be silent: %q", report.String())
	}
}

func TestAllocateCorrectsCounterBelowFloor(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	write(t, root, "grove/G-004-later.md", record("G-004", "work", "proposed"))
	dir := stateDir(t, root)
	write(t, dir, "neutral-ids", "G 2\n")
	var report bytes.Buffer
	n, err := Allocate(root, "grove", &report)
	if err != nil || n != 5 {
		t.Fatalf("got %d, %v; want 5", n, err)
	}
	if !strings.Contains(report.String(), "below") {
		t.Fatalf("should report the corrected counter: %q", report.String())
	}
	data, _ := os.ReadFile(filepath.Join(dir, "neutral-ids"))
	if !strings.Contains(string(data), "G 6\n") {
		t.Fatalf("counter not persisted: %q", data)
	}
}

func TestAllocateConcurrentAcrossWorktrees(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	wt := filepath.Join(filepath.Dir(root), filepath.Base(root)+"-odd\n\twt ") // W-008: a path Git quotes for display
	git(t, root, "worktree", "add", "-q", "-b", "feature", wt)
	roots := []string{root, wt}
	var mu sync.Mutex
	var wg sync.WaitGroup
	var got []int
	for i := range 20 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			n, err := Allocate(roots[i%2], "grove", &bytes.Buffer{})
			if err != nil {
				t.Error(err)
				return
			}
			mu.Lock()
			got = append(got, n)
			mu.Unlock()
		}()
	}
	wg.Wait()
	slices.Sort(got)
	for i, n := range got {
		if n != i+2 {
			t.Fatalf("expected 2..21 without duplicates, got %v", got)
		}
	}
}

func TestAllocateRefusesWhenAWorktreeCannotBeScanned(t *testing.T) {
	t.Parallel()
	if os.Geteuid() == 0 {
		t.Skip("permission checks do not apply to root")
	}
	root := gitProject(t)
	wt := filepath.Join(filepath.Dir(root), filepath.Base(root)+"-wt")
	git(t, root, "worktree", "add", "-q", "-b", "feature", wt)
	hidden := filepath.Join(wt, "grove/G-030-unreadable.md")
	write(t, wt, "grove/G-030-unreadable.md", record("G-030", "work", "proposed"))
	if err := os.Chmod(hidden, 0); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(hidden, 0o644) })
	if n, err := Allocate(root, "grove", &bytes.Buffer{}); err == nil {
		t.Fatalf("issued %d although a worktree record could not be read", n)
	}
	if _, err := os.Stat(filepath.Join(stateDir(t, root), "neutral-ids")); !os.IsNotExist(err) {
		t.Fatalf("no counter may be written when the floor is unknown: %v", err)
	}
}

func TestPersistenceFailureIssuesNoIDAndCreatesNoFile(t *testing.T) {
	t.Parallel()
	if os.Geteuid() == 0 {
		t.Skip("permission checks do not apply to root")
	}
	root := gitProject(t)
	dir := stateDir(t, root)
	write(t, dir, "lock", "")
	if err := os.Chmod(dir, 0o500); err != nil { // lock opens, counter temp file cannot be created
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(dir, 0o755) })
	p := load(t, root)
	if _, err := New(p, "work", "Blocked", "blocked", time.Now(), &bytes.Buffer{}); err == nil || !strings.Contains(err.Error(), "no ID issued") {
		t.Fatalf("expected a persistence diagnostic, got %v", err)
	}
	if entries, _ := os.ReadDir(filepath.Join(root, "grove")); len(entries) != 1 {
		t.Fatalf("failed allocation must create nothing: %v", entries)
	}
}

func TestAllocateRequiresGit(t *testing.T) {
	t.Parallel()
	root := t.TempDir()
	write(t, root, "grove.yaml", config)
	if _, err := Allocate(root, "grove", &bytes.Buffer{}); err == nil || !strings.Contains(err.Error(), "Git") {
		t.Fatalf("expected a Git requirement error, got %v", err)
	}
	entries, _ := os.ReadDir(root)
	if len(entries) != 1 {
		t.Fatalf("no state may be created outside Git: %v", entries)
	}
}

func TestNewCreatesValidRecords(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	now := time.Date(2026, 9, 19, 16, 0, 0, 0, time.UTC)
	var report bytes.Buffer
	path, err := New(load(t, root), "work", `Title: with "quotes" & more`, "", now, &report)
	if err != nil || path != "grove/G-002-title-with-quotes-more.md" {
		t.Fatalf("got %q, %v", path, err)
	}
	p := load(t, root)
	i := slices.IndexFunc(p.Records, func(r *project.Record) bool { return r.ID == "G-002" })
	if i < 0 {
		t.Fatal("G-002 not readable")
	}
	r := p.Records[i]
	if r.Title != `Title: with "quotes" & more` || r.Status != "proposed" || r.Created == nil || !r.Created.Equal(now) || !r.Updated.Equal(now) {
		t.Fatalf("unexpected record %+v", r)
	}
	if !strings.Contains(string(r.Source), "## Outcome") {
		t.Fatalf("missing body skeleton: %s", r.Source)
	}
	if path, err := New(p, "question", "Which?", "which", now, &report); err != nil || path != "grove/G-003-which.md" {
		t.Fatalf("got %q, %v", path, err)
	}
	if path, err := New(p, "decision", "Choose", "", now, &report); err != nil || path != "grove/G-004-choose.md" {
		t.Fatalf("got %q, %v", path, err)
	}
	load(t, root)
}

func TestNewNeverOverwritesAndConsumesReservation(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	// A directory at the target name makes O_EXCL creation fail while the
	// reader (which skips directories) still sees a valid project.
	if err := os.MkdirAll(filepath.Join(root, "grove/G-002-x.md"), 0o755); err != nil {
		t.Fatal(err)
	}
	p := load(t, root)
	now := time.Now().UTC()
	if _, err := New(p, "work", "X", "x", now, &bytes.Buffer{}); err == nil {
		t.Fatal("expected creation to fail on an existing path")
	}
	path, err := New(p, "work", "Y", "y", now, &bytes.Buffer{})
	if err != nil || path != "grove/G-003-y.md" {
		t.Fatalf("failed creation must consume G-002: got %q, %v", path, err)
	}
}

func TestNewRejectsBadSlugAndKind(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	p := load(t, root)
	for _, c := range [][2]string{{"work", "Bad Slug"}, {"work", "UPPER"}, {"release", "ok"}} {
		if _, err := New(p, c[0], "T", c[1], time.Now(), &bytes.Buffer{}); err == nil {
			t.Fatalf("expected %v to be rejected", c)
		}
	}
	if entries, _ := os.ReadDir(filepath.Join(root, "grove")); len(entries) != 1 {
		t.Fatalf("rejected input must not create files: %v", entries)
	}
}

func TestSlug(t *testing.T) {
	t.Parallel()
	cases := map[string]string{
		"  Hello,   World!! 1234567890123456789012345678 ": "hello-world-12345678901234567890",
		"Ünïcode ünd Emoji 🎉":                              "n-code-nd-emoji",
		"!!!":                                              "record",
		"trailing-hyphen-at-limit-exactly-32-x":            "trailing-hyphen-at-limit-exactly",
	}
	for in, want := range cases {
		if got := Slug(in); got != want {
			t.Errorf("Slug(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestNewRefusesWhenRecordRootChangedAfterLoad(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	p := load(t, root)
	write(t, root, "other/G-001-first.md", record("G-001", "work", "done"))
	write(t, root, "grove.yaml", "schema_version: 3\nrecords: other\n")
	_, err := New(p, "work", "Moved", "moved", time.Now(), &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "record root changed") {
		t.Fatalf("expected a configuration-change refusal, got %v", err)
	}
	for _, dir := range []string{"grove", "other"} {
		if entries, _ := os.ReadDir(filepath.Join(root, dir)); len(entries) != 1 {
			t.Fatalf("%s: nothing may be created: %v", dir, entries)
		}
	}
	write(t, root, "grove.yaml", config)
	if path, err := New(load(t, root), "work", "Next", "next", time.Now(), &bytes.Buffer{}); err != nil || path != "grove/G-003-next.md" {
		t.Fatalf("the refused reservation must stay consumed: got %q, %v", path, err)
	}
}

// A configuration edit that keeps the parsed meaning is still an observed
// change between the loaded allocation input and publication.
func TestNewRefusesWhenConfigurationBytesChangedAfterLoad(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	p := load(t, root)
	write(t, root, "grove.yaml", "# concurrent edit\n"+config)
	_, err := New(p, "work", "Late", "late", time.Now(), &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "reserved but not created: grove.yaml changed") {
		t.Fatalf("expected a configuration-change refusal, got %v", err)
	}
	if entries, _ := os.ReadDir(filepath.Join(root, "grove")); len(entries) != 1 {
		t.Fatalf("nothing may be created: %v", entries)
	}
	if path, err := New(load(t, root), "work", "Next", "next", time.Now(), &bytes.Buffer{}); err != nil || path != "grove/G-003-next.md" {
		t.Fatalf("the refused reservation must stay consumed: got %q, %v", path, err)
	}
}

// W-008: a live ID in a checkout whose path Git would quote for display still
// raises the floor, with no counter file, at the root and at a nested prefix.
func TestAllocateScansOddlyNamedWorktrees(t *testing.T) {
	t.Parallel()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	for _, prefix := range []string{"", "sub\nproject"} {
		parent, err := filepath.EvalSymlinks(t.TempDir())
		if err != nil {
			t.Fatal(err)
		}
		root := filepath.Join(parent, "new\nline")
		write(t, filepath.Join(root, prefix), "grove.yaml", config)
		write(t, filepath.Join(root, prefix), "grove/G-001-first.md", record("G-001", "work", "done"))
		git(t, root, "init", "-q", "-b", "main")
		git(t, root, "add", "-A")
		git(t, root, "commit", "-q", "-m", "init")
		odd := filepath.Join(parent, "odd \"quoted\"\n\twt ")
		git(t, root, "worktree", "add", "-q", "-b", "odd", odd)
		write(t, filepath.Join(odd, prefix), "grove/G-090-live-only.md", record("G-090", "work", "proposed"))
		write(t, odd, "elsewhere/grove/G-500-unrelated.md", record("G-500", "work", "proposed"))
		// A checkout without the record folder is normal, not a scan error.
		absent := filepath.Join(parent, "absent\nwt")
		git(t, root, "worktree", "add", "-q", "--detach", absent)
		if err := os.RemoveAll(filepath.Join(absent, prefix, "grove")); err != nil {
			t.Fatal(err)
		}

		if n, err := Allocate(filepath.Join(root, prefix), "grove", &bytes.Buffer{}); err != nil || n != 91 {
			t.Fatalf("prefix %q: got %d, %v; want 91 from the live G-090", prefix, n, err)
		}
		if n, err := Allocate(filepath.Join(odd, prefix), "grove", &bytes.Buffer{}); err != nil || n != 92 {
			t.Fatalf("prefix %q: got %d, %v; want 92 from the shared counter", prefix, n, err)
		}
		if _, err := os.Stat(filepath.Join(root, ".git", "grove", "neutral-ids")); err != nil {
			t.Fatalf("the counter belongs under the real common directory: %v", err)
		}
		if entries, _ := os.ReadDir(parent); len(entries) != 3 {
			t.Fatalf("nothing may appear beside the checkouts: %q", entries)
		}
	}
}

// Another process can define the term between this caller's load and its
// turn at the write lock; the second record must never reach the disk.
func TestNewRefusesATermDefinedAfterLoad(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	stale := load(t, root)
	if _, err := New(load(t, root), "term", "Attempt", "", time.Now(), &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	_, err := New(stale, "term", "attempt", "", time.Now(), &bytes.Buffer{})
	if err == nil || !strings.Contains(err.Error(), "G-003 reserved but not created: the term attempt is already defined by G-002") {
		t.Fatalf("expected a refusal under the write lock, got %v", err)
	}
	if entries, _ := os.ReadDir(filepath.Join(root, "grove")); len(entries) != 2 {
		t.Fatalf("nothing may be written for the refused attempt: %v", entries)
	}
	if len(load(t, root).Records) == 0 {
		t.Fatal("the project must stay loadable")
	}
}

func TestAllocationConcurrentAcrossWorktrees(t *testing.T) {
	t.Parallel()
	root := gitProject(t)
	wt := filepath.Join(filepath.Dir(root), filepath.Base(root)+"-wt")
	git(t, root, "worktree", "add", "-q", "-b", "feature", wt)
	projects := []*project.Project{load(t, root), load(t, wt)}
	var wg sync.WaitGroup
	for i := range 12 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := New(projects[i%2], "page", "P", "", time.Unix(0, 0), &bytes.Buffer{}); err != nil {
				t.Error(err)
			}
		}()
	}
	wg.Wait()
	seen := map[string]bool{}
	for _, dir := range []string{root, wt} {
		for _, r := range load(t, dir).Records {
			if r.ID == "G-001" {
				continue // committed baseline, held by both worktrees
			}
			if seen[r.ID] {
				t.Fatalf("%s issued twice", r.ID)
			}
			seen[r.ID] = true
		}
	}
	if len(seen) != 12 || !seen["G-013"] || seen["G-014"] {
		t.Fatalf("expected G-002..G-013, got %v", seen)
	}
}
