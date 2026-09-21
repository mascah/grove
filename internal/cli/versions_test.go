package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func gitIn(t *testing.T, dir string, args ...string) string {
	t.Helper()
	full := append([]string{"-c", "user.name=t", "-c", "user.email=t@t", "-c", "commit.gpgsign=false", "-C", dir}, args...)
	out, err := exec.Command("git", full...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
	return strings.TrimSpace(string(out))
}

// featureFixture adds a linked worktree whose G-001 is active with a changed
// body, committed on branch feature, beside the main checkout from gitFixture.
func featureFixture(t *testing.T) (root, wt string) {
	t.Helper()
	root = gitFixture(t)
	wt = filepath.Join(filepath.Dir(root), "feature wt")
	gitIn(t, root, "worktree", "add", "-q", "-b", "feature", wt)
	write(t, wt, "docs/records/work/renamed.md", strings.Replace(work, "status: proposed", "status: active", 1))
	gitIn(t, wt, "add", "-A")
	gitIn(t, wt, "commit", "-q", "-m", "feature")
	return root, wt
}

func TestVersionsUsage(t *testing.T) {
	t.Parallel()
	for _, args := range [][]string{
		{"versions", "G-001", "G-003"}, {"versions", "--slug", "x"}, {"list", "--json"}, {"versions", "--json", "--json"},
		{"versions", "--expect", rev},
	} {
		var out, errOut bytes.Buffer
		if code := Run(args, t.TempDir(), &out, &errOut); code != 2 || out.Len() != 0 || !strings.Contains(errOut.String(), "Usage:") {
			t.Fatalf("%v: code=%d stderr=%s", args, code, errOut.String())
		}
	}
}

// Every file under .git (refs, index, worktree metadata, allocator state) and
// every file in both checkouts hashes identically after text and JSON reads,
// and no coordination folder appears.
func TestVersionsLeavesGitUnchanged(t *testing.T) {
	t.Parallel()
	root, wt := featureFixture(t)
	write(t, wt, "docs/records/questions/dirty.md", strings.Replace(question, "G-002", "G-003", 1))
	before := map[string]map[string][32]byte{root: hashes(t, root), wt: hashes(t, wt)}
	for _, args := range [][]string{{"versions"}, {"versions", "G-001", "--json"}, {"--project", wt, "versions", "G-003"}} {
		var out, errOut bytes.Buffer
		if code := Run(args, root, &out, &errOut); code != 0 {
			t.Fatalf("%v: %s", args, errOut.String())
		}
	}
	if !reflect.DeepEqual(before[root], hashes(t, root)) || !reflect.DeepEqual(before[wt], hashes(t, wt)) {
		t.Fatal("versions must leave Git metadata, records, and dirty files unchanged")
	}
	if _, err := os.Stat(filepath.Join(root, ".git", "grove")); !errors.Is(err, fs.ErrNotExist) {
		t.Fatalf("versions must not create coordination state: %v", err)
	}
}

func TestVersionsCLI(t *testing.T) {
	t.Parallel()
	root, wt := featureFixture(t)
	var out, errOut bytes.Buffer
	if code := Run([]string{"versions", "G-001"}, root, &out, &errOut); code != 0 {
		t.Fatal(errOut.String())
	}
	rows := rowsOf(out.String())
	want := [][]string{
		{"ID", "STATUS", "SOURCE", "CHANGE", "SELECTOR"},
		{"G-001", "active", "committed refs/heads/feature", "-", "committed:refs/heads/feature@"},
		{"G-001", "proposed", "committed refs/heads/main", "-", "committed:refs/heads/main@"},
		{"G-001", "proposed", "live . refs/heads/main", "unchanged", "live:.:refs/heads/main@"},
		{"G-001", "active", "live feature-wt refs/heads/feature", "unchanged", "live:feature-wt:refs/heads/feature@"},
	}
	if len(rows) != len(want) {
		t.Fatalf("stdout:\n%s", out.String())
	}
	for i, row := range rows {
		if len(row) != 5 || !reflect.DeepEqual(row[:4], want[i][:4]) || !strings.HasPrefix(row[4], want[i][4]) {
			t.Fatalf("row %d: %q, expected %q", i, row, want[i])
		}
	}
	for _, want := range []string{"Project: " + root, "Repository: " + filepath.Join(root, ".git"), "Source: committed refs/heads/main ", "Source: live feature-wt refs/heads/feature ", wt} {
		if !strings.Contains(errOut.String(), want) {
			t.Fatalf("missing %q in stderr:\n%s", want, errOut.String())
		}
	}
	out.Reset()
	errOut.Reset()
	if code := Run([]string{"versions", "G-001", "--json"}, wt, &out, &errOut); code != 0 {
		t.Fatal(errOut.String())
	}
	var got struct {
		Project, Repository, Prefix string
		Complete                    bool
		Sources                     []map[string]any
		Records                     []struct {
			ID       string
			Versions []map[string]any
		}
	}
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatal(err)
	}
	if got.Project != wt || !got.Complete || len(got.Sources) != 4 || len(got.Records) != 1 || len(got.Records[0].Versions) != 4 {
		t.Fatalf("json: %s", out.String())
	}
	live := got.Records[0].Versions[3]
	if live["kind"] != "live" || live["worktree"] != wt || live["locator"] != "feature-wt" || live["ref"] != "refs/heads/feature" || live["change"] != "unchanged" ||
		live["path"] != "docs/records/work/renamed.md" || live["status"] != "active" || live["type"] != "work" || live["title"] != "Inspect records" ||
		!strings.HasPrefix(live["selector"].(string), "live:feature-wt:refs/heads/feature@") || live["source"] != strings.Replace(work, "status: proposed", "status: active", 1) ||
		!strings.HasPrefix(live["revision"].(string), "sha256:") || !strings.HasPrefix(live["config_revision"].(string), "sha256:") || live["detached"] != false {
		t.Fatalf("live version: %v", live)
	}
	if committed := got.Records[0].Versions[1]; committed["ref"] != "refs/heads/main" || committed["status"] != "proposed" || committed["change"] != nil || committed["worktree"] != nil {
		t.Fatalf("committed version: %v", committed)
	}
	for _, s := range got.Sources {
		if s["valid"] != true || s["present"] != true || len(s["diagnostics"].([]any)) != 0 {
			t.Fatalf("source: %v", s)
		}
	}
	// A control character in a worktree path survives JSON and is escaped in text.
	odd := filepath.Join(filepath.Dir(root), "odd\nwt")
	gitIn(t, root, "worktree", "add", "-q", "--detach", odd)
	out.Reset()
	errOut.Reset()
	if code := Run([]string{"versions", "G-002", "--json"}, root, &out, &errOut); code != 0 {
		t.Fatal(errOut.String())
	}
	if !strings.Contains(out.String(), `"worktree":"`+strings.Replace(odd, "\n", `\n`, 1)+`"`) || strings.Contains(out.String(), "\n"+"wt") {
		t.Fatalf("json path: %s", out.String())
	}
	out.Reset()
	errOut.Reset()
	if code := Run([]string{"versions", "G-002"}, root, &out, &errOut); code != 0 || !strings.Contains(errOut.String(), `odd\nwt`) || !strings.Contains(out.String(), "live odd-wt detached") {
		t.Fatalf("text output escapes control characters: %s\n%s", out.String(), errOut.String())
	}
	rows = rowsOf(out.String())
	last := rows[len(rows)-1]
	revision := strings.TrimPrefix(showJSON(t, root, "G-002")["revision"].(string), "sha256:")[:12]
	if !reflect.DeepEqual(last[:4], []string{"G-002", "open", "live odd-wt detached", "unchanged"}) || !strings.HasPrefix(last[4], "live:odd-wt:detached@"+gitIn(t, odd, "rev-parse", "HEAD")[:12]+":G-002@"+revision+":") {
		t.Fatalf("detached row: %q", last)
	}
}

// rowsOf splits table output into rows of cells at runs of two or more spaces.
func rowsOf(text string) [][]string {
	var rows [][]string
	for _, line := range strings.Split(strings.TrimSpace(text), "\n") {
		var cells []string
		for _, cell := range strings.Split(line, "  ") {
			if cell = strings.TrimSpace(cell); cell != "" {
				cells = append(cells, cell)
			}
		}
		rows = append(rows, cells)
	}
	return rows
}

func TestVersionsIncompleteAndNotFound(t *testing.T) {
	t.Parallel()
	root, wt := featureFixture(t)
	write(t, wt, "docs/records/work/broken.md", "---\nid: G-001\ntype: work\ntitle: dup\nstatus: proposed\n---\n")
	var out, errOut bytes.Buffer
	code := Run([]string{"versions"}, root, &out, &errOut)
	if code != 1 || !strings.Contains(errOut.String(), "duplicate G-001") || !strings.Contains(errOut.String(), "the result is incomplete") {
		t.Fatalf("code=%d stderr=%s", code, errOut.String())
	}
	if !strings.Contains(out.String(), "committed refs/heads/feature") || strings.Contains(out.String(), "live feature-wt") {
		t.Fatalf("valid sources still print; the invalid live source does not:\n%s", out.String())
	}
	out.Reset()
	errOut.Reset()
	if code := Run([]string{"versions", "--json"}, root, &out, &errOut); code != 1 || !strings.Contains(out.String(), `"complete":false`) {
		t.Fatalf("json marks incompleteness: %d %s", code, out.String())
	}
	if err := os.Remove(filepath.Join(wt, "docs/records/work/broken.md")); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	errOut.Reset()
	if code := Run([]string{"versions", "G-404"}, root, &out, &errOut); code != 1 || !strings.Contains(errOut.String(), "record G-404 not found in any valid source") || !strings.HasPrefix(out.String(), "ID  ") {
		t.Fatalf("code=%d stdout=%q stderr=%s", code, out.String(), errOut.String())
	}
	// An invalid current checkout is one invalid live source, not a hard stop.
	write(t, root, "grove.yaml", "schema_version: 4\nrecords: docs/records\n")
	out.Reset()
	errOut.Reset()
	if code := Run([]string{"versions", "G-001"}, root, &out, &errOut); code != 1 || !strings.Contains(errOut.String(), "unsupported version 4") || !strings.Contains(out.String(), "live feature-wt refs/heads/feature") {
		t.Fatalf("code=%d stdout=%s stderr=%s", code, out.String(), errOut.String())
	}
	if strings.Count(errOut.String(), "unsupported version 4") != 1 {
		t.Fatalf("the diagnostic is attributed once, to the live source:\n%s", errOut.String())
	}
	plain := projectFixture(t)
	out.Reset()
	errOut.Reset()
	if code := Run([]string{"versions"}, plain, &out, &errOut); code != 1 || out.Len() != 0 || !strings.Contains(errOut.String(), "requires a Git repository") {
		t.Fatalf("plain directories: %d %s", code, errOut.String())
	}
}
