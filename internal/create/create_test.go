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

const config = "schema_version: 1\nrecords: grove\n"

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

// repo returns a committed Git project containing W-001.
func repo(t *testing.T) string {
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
	write(t, root, "grove/work/W-001-first.md", record("W-001", "work", "done"))
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
	root := repo(t)
	wt := filepath.Join(filepath.Dir(root), filepath.Base(root)+"-wt")
	git(t, root, "worktree", "add", "-q", "-b", "feature", wt)
	write(t, wt, "grove/work/W-003-committed-on-branch.md", record("W-003", "work", "proposed"))
	git(t, wt, "add", "-A")
	git(t, wt, "commit", "-q", "-m", "branch record")
	write(t, wt, "grove/work/W-005-live-only.md", record("W-005", "work", "proposed"))

	var report bytes.Buffer
	n, err := Allocate(root, "grove", "W", &report)
	if err != nil || n != 6 {
		t.Fatalf("got %d, %v; want 6 from committed W-003 and live W-005", n, err)
	}
	if !strings.Contains(strings.ToLower(report.String()), "initialized") {
		t.Fatalf("first allocation should report initialization: %q", report.String())
	}
	report.Reset()
	n, err = Allocate(wt, "grove", "W", &report)
	if err != nil || n != 7 {
		t.Fatalf("got %d, %v; want 7 from the shared counter", n, err)
	}
	if report.Len() != 0 {
		t.Fatalf("steady-state allocation should be silent: %q", report.String())
	}
	if n, err := Allocate(root, "grove", "Q", &report); err != nil || n != 1 {
		t.Fatalf("questions start at 1: got %d, %v", n, err)
	}
}

func TestAllocateCorrectsCounterBelowFloor(t *testing.T) {
	root := repo(t)
	write(t, root, "grove/work/W-004-later.md", record("W-004", "work", "proposed"))
	dir := stateDir(t, root)
	write(t, dir, "next-ids", "W 2\n")
	var report bytes.Buffer
	n, err := Allocate(root, "grove", "W", &report)
	if err != nil || n != 5 {
		t.Fatalf("got %d, %v; want 5", n, err)
	}
	if !strings.Contains(report.String(), "below") {
		t.Fatalf("should report the corrected counter: %q", report.String())
	}
	data, _ := os.ReadFile(filepath.Join(dir, "next-ids"))
	if !strings.Contains(string(data), "W 6\n") {
		t.Fatalf("counter not persisted: %q", data)
	}
}

func TestAllocateConcurrentAcrossWorktrees(t *testing.T) {
	root := repo(t)
	wt := filepath.Join(filepath.Dir(root), filepath.Base(root)+"-wt")
	git(t, root, "worktree", "add", "-q", "-b", "feature", wt)
	roots := []string{root, wt}
	var mu sync.Mutex
	var wg sync.WaitGroup
	var got []int
	for i := range 20 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			n, err := Allocate(roots[i%2], "grove", "W", &bytes.Buffer{})
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

func TestAllocateRequiresGit(t *testing.T) {
	root := t.TempDir()
	write(t, root, "grove.yaml", config)
	if _, err := Allocate(root, "grove", "W", &bytes.Buffer{}); err == nil || !strings.Contains(err.Error(), "Git") {
		t.Fatalf("expected a Git requirement error, got %v", err)
	}
	entries, _ := os.ReadDir(root)
	if len(entries) != 1 {
		t.Fatalf("no state may be created outside Git: %v", entries)
	}
}

func TestNewCreatesValidRecords(t *testing.T) {
	root := repo(t)
	now := time.Date(2026, 9, 19, 16, 0, 0, 0, time.UTC)
	var report bytes.Buffer
	path, err := New(load(t, root), "work", `Title: with "quotes" & more`, "", now, &report)
	if err != nil || path != "grove/work/W-002-title-with-quotes-more.md" {
		t.Fatalf("got %q, %v", path, err)
	}
	p := load(t, root)
	i := slices.IndexFunc(p.Records, func(r *project.Record) bool { return r.ID == "W-002" })
	if i < 0 {
		t.Fatal("W-002 not readable")
	}
	r := p.Records[i]
	if r.Title != `Title: with "quotes" & more` || r.Status != "proposed" || r.Created == nil || !r.Created.Equal(now) || !r.Updated.Equal(now) {
		t.Fatalf("unexpected record %+v", r)
	}
	if !strings.Contains(string(r.Source), "## Outcome") {
		t.Fatalf("missing body skeleton: %s", r.Source)
	}
	if path, err := New(p, "question", "Which?", "which", now, &report); err != nil || path != "grove/questions/Q-001-which.md" {
		t.Fatalf("got %q, %v", path, err)
	}
	if path, err := New(p, "decision", "Choose", "", now, &report); err != nil || path != "grove/decisions/D-001-choose.md" {
		t.Fatalf("got %q, %v", path, err)
	}
	load(t, root)
}

func TestNewNeverOverwritesAndConsumesReservation(t *testing.T) {
	root := repo(t)
	// A directory at the target name makes O_EXCL creation fail while the
	// reader (which skips directories) still sees a valid project.
	if err := os.MkdirAll(filepath.Join(root, "grove/work/W-002-x.md"), 0o755); err != nil {
		t.Fatal(err)
	}
	p := load(t, root)
	now := time.Now().UTC()
	if _, err := New(p, "work", "X", "x", now, &bytes.Buffer{}); err == nil {
		t.Fatal("expected creation to fail on an existing path")
	}
	path, err := New(p, "work", "Y", "y", now, &bytes.Buffer{})
	if err != nil || path != "grove/work/W-003-y.md" {
		t.Fatalf("failed creation must consume W-002: got %q, %v", path, err)
	}
}

func TestNewRejectsBadSlugAndKind(t *testing.T) {
	root := repo(t)
	p := load(t, root)
	for _, c := range [][2]string{{"work", "Bad Slug"}, {"work", "UPPER"}, {"release", "ok"}} {
		if _, err := New(p, c[0], "T", c[1], time.Now(), &bytes.Buffer{}); err == nil {
			t.Fatalf("expected %v to be rejected", c)
		}
	}
	if entries, _ := os.ReadDir(filepath.Join(root, "grove/work")); len(entries) != 1 {
		t.Fatalf("rejected input must not create files: %v", entries)
	}
}

func TestSlug(t *testing.T) {
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
