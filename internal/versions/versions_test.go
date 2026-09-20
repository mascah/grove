package versions

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/mascah/grove/internal/project"
)

const config = "schema_version: 1\nrecords: grove\n"

func record(id, kind, status, body string) string {
	return "---\nid: \"" + id + "\"\ntype: " + kind + "\ntitle: T " + id + "\nstatus: " + status + "\n---\n" + body
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

func write(t *testing.T, root, rel, content string) {
	t.Helper()
	full := filepath.Join(root, rel)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func commit(t *testing.T, dir, message string) string {
	t.Helper()
	git(t, dir, "add", "-A")
	git(t, dir, "commit", "-q", "--allow-empty", "-m", message)
	return git(t, dir, "rev-parse", "HEAD")
}

// repoFixture returns a committed main worktree holding W-001 and Q-001.
func repoFixture(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	root = filepath.Join(root, "repo")
	if err := os.Mkdir(root, 0o755); err != nil {
		t.Fatal(err)
	}
	git(t, root, "init", "-q", "-b", "main")
	write(t, root, "grove.yaml", config)
	write(t, root, "grove/work/W-001-first.md", record("W-001", "work", "proposed", "Main body.\n"))
	write(t, root, "grove/questions/Q-001-q.md", record("Q-001", "question", "open", "Q.\n"))
	commit(t, root, "init")
	return root
}

// addWorktree registers a linked worktree beside root at commitish ("" for
// HEAD) with extra options, and returns its path.
func addWorktree(t *testing.T, root, name, commitish string, options ...string) string {
	t.Helper()
	path := filepath.Join(filepath.Dir(root), name)
	args := append([]string{"worktree", "add", "-q"}, options...)
	args = append(args, path)
	if commitish != "" {
		args = append(args, commitish)
	}
	git(t, root, args...)
	return path
}

func mustInspect(t *testing.T, root, id string) *Result {
	t.Helper()
	res, err := Inspect(root, id)
	if err != nil {
		t.Fatal(err)
	}
	return res
}

func group(t *testing.T, res *Result, id string) Group {
	t.Helper()
	for _, g := range res.Groups {
		if g.ID == id {
			return g
		}
	}
	t.Fatalf("no group %s in %+v", id, res.Groups)
	return Group{}
}

// find returns the version from a committed ref or a live locator.
func find(t *testing.T, g Group, kind, where string) Version {
	t.Helper()
	for _, v := range g.Versions {
		if v.Source.Kind == kind && (kind == "committed" && v.Source.Ref == where || kind == "live" && v.Source.Locator == where) {
			return v
		}
	}
	t.Fatalf("no %s version at %s for %s", kind, where, g.ID)
	return Version{}
}

func source(t *testing.T, res *Result, kind, where string) *Source {
	t.Helper()
	for _, s := range res.Sources {
		if s.Kind == kind && (kind == "committed" && s.Ref == where || kind == "live" && s.Locator == where) {
			return s
		}
	}
	t.Fatalf("no %s source %s", kind, where)
	return nil
}

func TestInspectMainAndFeature(t *testing.T) {
	root := repoFixture(t)
	wt := addWorktree(t, root, "feature", "", "-b", "feature")
	write(t, wt, "grove/work/W-001-first.md", record("W-001", "work", "active", "Feature body.\n"))
	write(t, wt, "grove/work/W-002-second.md", record("W-002", "work", "proposed", "Only here.\n"))
	tip := commit(t, wt, "feature progress")
	mainTip := git(t, root, "rev-parse", "HEAD")

	res := mustInspect(t, root, "")
	if !res.Complete || res.Prefix != "" || res.Project != root {
		t.Fatalf("unexpected result header: %+v", res)
	}
	if ids := groupIDs(res); !reflect.DeepEqual(ids, []string{"W-001", "W-002", "Q-001"}) {
		t.Fatalf("groups: %v", ids)
	}
	g := group(t, res, "W-001")
	if len(g.Versions) != 4 {
		t.Fatalf("expected four versions, got %+v", g.Versions)
	}
	order := []string{"committed refs/heads/feature", "committed refs/heads/main", "live .", "live feature"}
	for i, v := range g.Versions {
		got := v.Source.Kind + " " + v.Source.Ref
		if v.Source.Kind == "live" {
			got = "live " + v.Source.Locator
		}
		if got != order[i] {
			t.Fatalf("version %d is %s, expected %s", i, got, order[i])
		}
	}
	main, feature := find(t, g, "committed", "refs/heads/main"), find(t, g, "committed", "refs/heads/feature")
	if main.Record.Status != "proposed" || feature.Record.Status != "active" || main.Revision == feature.Revision {
		t.Fatalf("statuses must come from each branch: %s %s", main.Record.Status, feature.Record.Status)
	}
	if string(main.Record.Source) != record("W-001", "work", "proposed", "Main body.\n") || string(feature.Record.Source) != record("W-001", "work", "active", "Feature body.\n") {
		t.Fatal("committed sources must carry each branch's exact bytes")
	}
	if main.Source.Commit != mainTip || feature.Source.Commit != tip {
		t.Fatalf("commits: %s %s", main.Source.Commit, feature.Source.Commit)
	}
	liveMain, liveFeature := find(t, g, "live", "."), find(t, g, "live", "feature")
	if liveMain.Change != "unchanged" || liveFeature.Change != "unchanged" || liveMain.Source.Worktree != root || liveFeature.Source.Worktree != wt {
		t.Fatalf("live versions: %+v %+v", liveMain, liveFeature)
	}
	if liveMain.Source.Ref != "refs/heads/main" || liveFeature.Source.Ref != "refs/heads/feature" || liveFeature.Source.GitDir != filepath.Join(root, ".git", "worktrees", "feature") {
		t.Fatalf("live identity: %+v", liveFeature.Source)
	}
	// Identical bytes in different sources still get different selectors.
	if main.Revision != liveMain.Revision || main.Selector == liveMain.Selector || feature.Selector == liveFeature.Selector {
		t.Fatalf("selectors must bind source identity: %s %s", main.Selector, liveMain.Selector)
	}
	selectors := map[string]bool{}
	for _, g := range res.Groups {
		for _, v := range g.Versions {
			if selectors[v.Selector] {
				t.Fatalf("duplicate selector %s", v.Selector)
			}
			selectors[v.Selector] = true
		}
	}
	if !strings.HasPrefix(main.Selector, "committed:refs/heads/main@"+mainTip[:12]+":W-001@"+strings.TrimPrefix(main.Revision, "sha256:")[:12]+":") ||
		!strings.HasPrefix(liveFeature.Selector, "live:feature:refs/heads/feature@"+tip[:12]+":W-001@") {
		t.Fatalf("selector grammar: %s %s", main.Selector, liveFeature.Selector)
	}
	if two := group(t, res, "W-002"); len(two.Versions) != 2 || two.Versions[0].Source.Ref != "refs/heads/feature" || two.Versions[1].Source.Locator != "feature" {
		t.Fatalf("W-002 exists only on the feature: %+v", two.Versions)
	}
	if only := mustInspect(t, root, "W-002"); len(only.Groups) != 1 || only.Groups[0].ID != "W-002" || len(only.Sources) != 4 {
		t.Fatalf("an ID filter keeps every source: %+v", only.Groups)
	}
	if none := mustInspect(t, root, "W-009"); len(none.Groups) != 0 || !none.Complete {
		t.Fatalf("an absent ID is an empty, complete result: %+v", none.Groups)
	}
	// The same call from the feature worktree sees the same picture.
	again := mustInspect(t, wt, "")
	if again.Project != wt || !reflect.DeepEqual(selectorsOf(res), selectorsOf(again)) {
		t.Fatalf("selectors depend on the repository, not the invoking checkout")
	}
}

func dump(res *Result) string {
	var b strings.Builder
	for _, s := range res.Sources {
		b.WriteString("\n" + s.Kind + " " + s.Ref + " " + s.Locator + " valid=" + strconv.FormatBool(s.Valid) + " " + strings.Join(s.Diagnostics, "; "))
	}
	return b.String()
}

func groupIDs(res *Result) []string {
	ids := make([]string, len(res.Groups))
	for i, g := range res.Groups {
		ids[i] = g.ID
	}
	return ids
}

func selectorsOf(res *Result) []string {
	var out []string
	for _, g := range res.Groups {
		for _, v := range g.Versions {
			out = append(out, v.Selector)
		}
	}
	return out
}

func TestInspectBranchWithoutCheckoutAndDetached(t *testing.T) {
	root := repoFixture(t)
	wt := addWorktree(t, root, "temp", "", "-b", "orphan")
	write(t, wt, "grove/work/W-001-first.md", record("W-001", "work", "done", "Orphan.\n"))
	commit(t, wt, "orphan")
	git(t, root, "worktree", "remove", "--force", wt)
	detached := addWorktree(t, root, "det", "orphan", "--detach")

	res := mustInspect(t, root, "W-001")
	if !res.Complete {
		t.Fatalf("unexpected diagnostics: %+v", res.Sources)
	}
	g := group(t, res, "W-001")
	orphan := find(t, g, "committed", "refs/heads/orphan")
	if orphan.Record.Status != "done" {
		t.Fatalf("a branch without a checkout contributes its committed records: %+v", orphan)
	}
	for _, v := range g.Versions {
		if v.Source.Kind == "live" && v.Source.Worktree == wt {
			t.Fatal("a removed worktree must not appear")
		}
	}
	d := find(t, g, "live", "det")
	if !d.Source.Detached() || d.Source.Ref != "" || d.Source.Worktree != detached || d.Change != "unchanged" || d.Record.Status != "done" {
		t.Fatalf("detached checkout: %+v", d)
	}
	if !strings.HasPrefix(d.Selector, "live:det:detached@"+orphan.Source.Commit[:12]+":W-001@") {
		t.Fatalf("detached selector: %s", d.Selector)
	}
	if n := len(g.Versions); n != 4 { // main + orphan committed, main + det live
		t.Fatalf("expected four versions, got %d", n)
	}
}

func TestInspectLiveChanges(t *testing.T) {
	root := repoFixture(t)
	wt := addWorktree(t, root, "feature", "", "-b", "feature")
	write(t, wt, "grove/work/W-002-second.md", record("W-002", "work", "proposed", "Two.\n"))
	write(t, wt, "grove/work/W-004-fourth.md", record("W-004", "work", "proposed", "Four.\n"))
	write(t, wt, "grove/work/W-005-fifth.md", record("W-005", "work", "proposed", "Five.\n"))
	commit(t, wt, "records")
	write(t, wt, "grove/work/W-001-first.md", record("W-001", "work", "active", "Main body.\n"))
	if err := os.Remove(filepath.Join(wt, "grove/work/W-002-second.md")); err != nil {
		t.Fatal(err)
	}
	write(t, wt, "grove/work/W-003-third.md", record("W-003", "work", "proposed", "Three.\n"))
	if err := os.Rename(filepath.Join(wt, "grove/work/W-004-fourth.md"), filepath.Join(wt, "grove/work/W-004-renamed.md")); err != nil {
		t.Fatal(err)
	}
	write(t, wt, "grove/work/W-005-moved.md", record("W-005", "work", "done", "Five.\n"))
	if err := os.Remove(filepath.Join(wt, "grove/work/W-005-fifth.md")); err != nil {
		t.Fatal(err)
	}

	res := mustInspect(t, root, "")
	if !res.Complete {
		t.Fatalf("unexpected diagnostics: %+v", res.Sources)
	}
	want := map[string][2]string{
		"W-001": {"modified", ""}, "W-002": {"deleted", ""}, "W-003": {"added", ""},
		"W-004": {"renamed", "grove/work/W-004-fourth.md"}, "W-005": {"modified", "grove/work/W-005-fifth.md"},
	}
	for id, expect := range want {
		v := find(t, group(t, res, id), "live", "feature")
		if v.Change != expect[0] || v.HeadPath != expect[1] {
			t.Fatalf("%s: change=%s head_path=%s, expected %v", id, v.Change, v.HeadPath, expect)
		}
	}
	deleted := find(t, group(t, res, "W-002"), "live", "feature")
	if deleted.Record != nil || deleted.Selector != "" || deleted.Revision != "" || deleted.Path != "grove/work/W-002-second.md" {
		t.Fatalf("a deleted record is a row without content or selector: %+v", deleted)
	}
	if c := find(t, group(t, res, "W-002"), "committed", "refs/heads/feature"); c.Record.Status != "proposed" {
		t.Fatal("the committed W-002 stays visible on the branch")
	}
	if v := find(t, group(t, res, "W-004"), "live", "feature"); v.Path != "grove/work/W-004-renamed.md" {
		t.Fatalf("renamed path: %s", v.Path)
	}
	if v := find(t, group(t, res, "W-001"), "live", "."); v.Change != "unchanged" {
		t.Fatalf("main is clean: %s", v.Change)
	}
	// A checkout whose HEAD has no project reports every record as added.
	empty := addWorktree(t, root, "empty", git(t, root, "commit-tree", "-m", "empty", git(t, root, "hash-object", "-t", "tree", "--stdin", "-w")), "--detach")
	write(t, empty, "grove.yaml", config)
	write(t, empty, "grove/work/W-001-first.md", record("W-001", "work", "proposed", "x\n"))
	res = mustInspect(t, root, "W-001")
	if v := find(t, group(t, res, "W-001"), "live", "empty"); v.Change != "added" || !res.Complete {
		t.Fatalf("no project at HEAD means an empty baseline: %+v", v)
	}
}

func TestInspectPrefixAndConfig(t *testing.T) {
	root := repoFixture(t)
	// Move the project below the repository root with a different record folder.
	git(t, root, "rm", "-q", "-r", "grove.yaml", "grove")
	write(t, root, "sub/grove.yaml", "schema_version: 1\nrecords: docs/records\n")
	write(t, root, "sub/docs/records/work/W-001-first.md", record("W-001", "work", "proposed", "Main.\n"))
	write(t, root, "sub/docs/records/questions/Q-001-q.md", record("Q-001", "question", "open", "Q.\n"))
	commit(t, root, "nested")
	wt := addWorktree(t, root, "feature", "", "-b", "feature")
	git(t, wt, "rm", "-q", "-r", "sub/docs")
	write(t, wt, "sub/grove.yaml", config)
	write(t, wt, "sub/grove/work/W-001-first.md", record("W-001", "work", "active", "Feature.\n"))
	write(t, wt, "sub/grove/questions/Q-001-q.md", record("Q-001", "question", "open", "Q.\n"))
	commit(t, wt, "feature config")
	git(t, root, "branch", "-q", "bare-branch", git(t, root, "commit-tree", "-m", "no project", git(t, root, "hash-object", "-t", "tree", "--stdin", "-w")))

	res := mustInspect(t, filepath.Join(root, "sub"), "")
	if res.Prefix != "sub/" || !res.Complete {
		t.Fatalf("prefix=%q sources=%s", res.Prefix, dump(res))
	}
	g := group(t, res, "W-001")
	main, feature := find(t, g, "committed", "refs/heads/main"), find(t, g, "committed", "refs/heads/feature")
	if main.Path != "docs/records/work/W-001-first.md" || feature.Path != "grove/work/W-001-first.md" || main.Source.ConfigRevision == feature.Source.ConfigRevision {
		t.Fatalf("each source uses its own configuration: %+v %+v", main, feature)
	}
	if feature.Source.ConfigRevision != project.Revision([]byte(config)) {
		t.Fatalf("configuration revision hashes the exact grove.yaml bytes: %s", feature.Source.ConfigRevision)
	}
	if l := find(t, g, "live", "feature"); l.Path != "grove/work/W-001-first.md" || l.Change != "unchanged" || l.Source.Worktree != wt {
		t.Fatalf("live feature: %+v", l)
	}
	if b := source(t, res, "committed", "refs/heads/bare-branch"); b.Present || b.Valid || len(b.Diagnostics) != 0 {
		t.Fatalf("a branch without the project is absent, not an error: %+v", b)
	}
	if len(g.Versions) != 4 {
		t.Fatalf("expected four versions, got %d", len(g.Versions))
	}
}

func TestInspectSourceLocalValidation(t *testing.T) {
	root := repoFixture(t)
	wt := addWorktree(t, root, "feature", "", "-b", "feature")
	write(t, wt, "grove/work/W-002-second.md", "---\nid: \"W-002\"\ntype: work\ntitle: T\nstatus: proposed\ndepends_on: [\"W-009\"]\n---\n")
	commit(t, wt, "dangling dependency")
	write(t, root, "grove/work/W-009-ninth.md", record("W-009", "work", "done", "Only on main.\n"))
	commit(t, root, "ninth on main")

	res := mustInspect(t, root, "")
	if res.Complete {
		t.Fatal("a dependency missing in one source keeps that source invalid")
	}
	for _, where := range []string{"refs/heads/feature"} {
		s := source(t, res, "committed", where)
		if s.Valid || !s.Present || len(s.Diagnostics) != 1 || !strings.Contains(s.Diagnostics[0], "unresolved target W-009") {
			t.Fatalf("feature source: %+v", s)
		}
	}
	if s := source(t, res, "live", "feature"); s.Valid || !strings.Contains(strings.Join(s.Diagnostics, "\n"), "unresolved target W-009") {
		t.Fatalf("live feature source: %+v", s)
	}
	if s := source(t, res, "committed", "refs/heads/main"); !s.Valid {
		t.Fatalf("main stays valid: %+v", s)
	}
	if _, ok := lookup(res, "W-002"); ok {
		t.Fatal("records from invalid sources are not admitted as observations")
	}
	if g := group(t, res, "W-009"); len(g.Versions) != 2 {
		t.Fatalf("W-009 from the valid main sources: %+v", g.Versions)
	}
}

func lookup(res *Result, id string) (Group, bool) {
	for _, g := range res.Groups {
		if g.ID == id {
			return g, true
		}
	}
	return Group{}, false
}

func TestInspectIncomplete(t *testing.T) {
	root := repoFixture(t)
	bad := addWorktree(t, root, "bad", "", "-b", "bad-yaml")
	write(t, bad, "grove/work/W-002-second.md", "---\nid: \"W-002\"\ntype: work\ntitle: [unclosed\nstatus: proposed\n---\n")
	commit(t, bad, "bad yaml")
	dup := addWorktree(t, root, "dup", "", "-b", "dup-ids")
	write(t, dup, "grove/work/W-001-copy.md", record("W-001", "work", "done", "Copy.\n"))
	commit(t, dup, "duplicate id")
	gone := addWorktree(t, root, "gone", "", "-b", "gone")
	if err := os.RemoveAll(gone); err != nil {
		t.Fatal(err)
	}
	symlinked := addWorktree(t, root, "sym", "", "-b", "sym")
	if err := os.Remove(filepath.Join(symlinked, "grove.yaml")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("nowhere", filepath.Join(symlinked, "grove.yaml")); err != nil {
		t.Fatal(err)
	}
	commit(t, symlinked, "symlinked configuration")

	res := mustInspect(t, root, "")
	if res.Complete {
		t.Fatal("invalid sources make the result incomplete")
	}
	expect := map[string]string{
		"committed refs/heads/bad-yaml": "grove/work/W-002-second.md:", "live bad": "grove/work/W-002-second.md:",
		"committed refs/heads/dup-ids": "duplicate W-001", "live dup": "duplicate W-001",
		"committed refs/heads/sym": "grove.yaml: symlink", "live sym": "grove.yaml: symlink",
		"live gone": "prunable",
	}
	for _, s := range res.Sources {
		key := s.Kind + " " + s.Ref
		if s.Kind == "live" {
			key = "live " + s.Locator
			if s.Locator == "" {
				key = "live " + filepath.Base(s.Worktree)
			}
		}
		if want, ok := expect[key]; ok {
			if s.Valid || len(s.Diagnostics) == 0 || !strings.Contains(strings.Join(s.Diagnostics, "\n"), want) {
				t.Fatalf("%s: expected %q in %+v", key, want, s.Diagnostics)
			}
			delete(expect, key)
		} else if !s.Valid {
			t.Fatalf("unexpected invalid source %s: %+v", key, s.Diagnostics)
		}
	}
	if len(expect) != 0 {
		t.Fatalf("sources not reported: %v", expect)
	}
	if g := group(t, res, "W-001"); len(g.Versions) != 3 { // main and gone committed, main live
		t.Fatalf("valid sources still contribute: %s", dump(res))
	}
	if !res.Complete {
		for _, s := range res.Sources {
			if s.Kind == "committed" && s.Ref == "refs/heads/gone" && !s.Valid {
				t.Fatal("the branch of a prunable worktree is still readable")
			}
		}
	}
}

func TestInspectUnstable(t *testing.T) {
	root := repoFixture(t)
	moving := addWorktree(t, root, "moving", "", "-b", "moving")
	leaving := addWorktree(t, root, "leaving", "", "-b", "leaving")
	detaching := addWorktree(t, root, "detaching", "", "-b", "detaching")
	res, err := inspect(t.Context(), root, "", func() {
		write(t, moving, "grove/work/W-002-second.md", record("W-002", "work", "proposed", "x\n"))
		commit(t, moving, "moved")
		git(t, root, "worktree", "remove", "--force", leaving)
		git(t, detaching, "checkout", "-q", "--detach")
		addWorktree(t, root, "arriving", "", "-b", "arriving")
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.Complete {
		t.Fatal("identity changes during reading must make the result incomplete")
	}
	for locator, want := range map[string]string{"moving": "changed while being read: refs/heads/moving at", "leaving": "removed while being read", "detaching": "to detached at"} {
		s := source(t, res, "live", locator)
		if s.Valid || !strings.Contains(strings.Join(s.Diagnostics, "\n"), want) {
			t.Fatalf("%s: %+v", locator, s.Diagnostics)
		}
	}
	arrived := false
	for _, s := range res.Sources {
		if s.Kind == "live" && filepath.Base(s.Worktree) == "arriving" {
			arrived = !s.Valid && strings.Contains(s.Diagnostics[0], "appeared while being read")
		}
	}
	if !arrived {
		t.Fatal("a worktree registered during reading is reported")
	}
	if s := source(t, res, "live", "."); !s.Valid {
		t.Fatalf("the untouched checkout stays valid: %+v", s.Diagnostics)
	}
	if _, ok := lookup(res, "W-002"); ok {
		t.Fatal("records from an unstable source are not admitted")
	}
}

func TestInspectBytesAndSelectors(t *testing.T) {
	root := repoFixture(t)
	bom := "\ufeff---\r\nid: \"W-002\"\r\ntype: work\r\ntitle: \"Ünïcode\"\r\nstatus: proposed\r\n---\r\nBody without final newline"
	write(t, root, "grove/work/W-002-bom.md", bom)
	commit(t, root, "bom crlf")
	wt := addWorktree(t, root, "feature", "", "-b", "feature")

	res := mustInspect(t, root, "W-002")
	g := group(t, res, "W-002")
	sum := sha256.Sum256([]byte(bom))
	want := "sha256:" + hex.EncodeToString(sum[:])
	seen := map[string]bool{}
	for _, v := range g.Versions {
		if string(v.Record.Source) != bom || v.Revision != want || v.Record.Title != "Ünïcode" {
			t.Fatalf("bytes and revision must agree in every source: %+v", v)
		}
		if seen[v.Selector] {
			t.Fatalf("identical content in %s and another source share selector %s", v.Source.Kind, v.Selector)
		}
		seen[v.Selector] = true
	}
	if len(g.Versions) != 4 || find(t, g, "live", "feature").Source.Worktree != wt {
		t.Fatalf("expected four versions with identical bytes: %d", len(g.Versions))
	}
}

func TestInspectPathsAndRepeatedReads(t *testing.T) {
	root := repoFixture(t)
	wt := addWorktree(t, root, "odd\nname\twith space", "", "-b", "odd")
	write(t, wt, "grove/work/W-002-with space.md", record("W-002", "work", "proposed", "x\n"))
	commit(t, wt, "odd path")

	res := mustInspect(t, root, "")
	s := source(t, res, "live", "odd-name-with-space")
	if s.Worktree != wt || !s.Valid {
		t.Fatalf("worktree path must round-trip: %+v", s)
	}
	v := find(t, group(t, res, "W-002"), "live", "odd-name-with-space")
	if v.Path != "grove/work/W-002-with space.md" || !strings.HasPrefix(v.Selector, "live:odd-name-with-space:refs/heads/odd@") {
		t.Fatalf("record path and selector: %+v", v)
	}
	for i := 0; i < 3; i++ {
		again := mustInspect(t, root, "")
		if !reflect.DeepEqual(selectorsOf(again), selectorsOf(res)) || !slices.Equal(groupIDs(again), groupIDs(res)) {
			t.Fatal("repeated unchanged reads must order identically")
		}
	}
}

func TestInspectUnbornAndInvalidHEAD(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	root, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	git(t, root, "init", "-q", "-b", "main")
	write(t, root, "grove.yaml", config)
	write(t, root, "grove/work/W-001-first.md", record("W-001", "work", "proposed", "x\n"))
	res := mustInspect(t, root, "")
	if !res.Complete || len(res.Sources) != 1 || res.Sources[0].Locator != "." || !res.Sources[0].Valid {
		t.Fatalf("a project with no commits yet is one valid live source: %s", dump(res))
	}
	if v := find(t, group(t, res, "W-001"), "live", "."); v.Change != "added" {
		t.Fatalf("every record is added before the first commit: %+v", v)
	}
	// A HEAD whose project does not validate cannot be compared with, but the
	// valid live checkout still counts; its changes are unknown and noted.
	commit(t, root, "first")
	temp := addWorktree(t, root, "temp", "", "-b", "bad")
	write(t, temp, "grove/work/W-001-copy.md", record("W-001", "work", "done", "dup\n"))
	bad := commit(t, temp, "duplicate ids")
	git(t, root, "worktree", "remove", "--force", temp)
	wt := addWorktree(t, root, "det", bad, "--detach")
	git(t, root, "branch", "-q", "-D", "bad")
	if err := os.Remove(filepath.Join(wt, "grove/work/W-001-copy.md")); err != nil {
		t.Fatal(err)
	}
	res = mustInspect(t, root, "")
	if !res.Complete || len(res.Sources) != 3 {
		t.Fatalf("the fixed live checkout is valid; only its comparison is unknown: %s", dump(res))
	}
	det := source(t, res, "live", "det")
	if !det.Valid || !strings.Contains(det.Note, "HEAD "+bad+" does not validate") {
		t.Fatalf("expected a note on the detached source: %+v", det)
	}
	if v := find(t, group(t, res, "W-001"), "live", "det"); v.Change != "unknown" || v.Selector == "" {
		t.Fatalf("changes against an invalid HEAD are unknown but the version is selectable: %+v", v)
	}
}

func TestInspectRequiresGit(t *testing.T) {
	root := t.TempDir()
	write(t, root, "grove.yaml", config)
	write(t, root, "grove/work/W-001-first.md", record("W-001", "work", "proposed", "x\n"))
	if _, err := Inspect(root, ""); err == nil || !strings.Contains(err.Error(), "requires a Git repository") {
		t.Fatalf("plain directories need a Git diagnostic: %v", err)
	}
}

// A branch is read for what the project loader reads and nothing else:
// branches that differ only elsewhere share one loaded project, a differing
// record folder gets its own, and one git process serves them all.
func TestCommittedReadIsScopedAndShared(t *testing.T) {
	root := repoFixture(t)
	git(t, root, "checkout", "-q", "-b", "code")
	write(t, root, "src/notes.md", "not a record\n")
	write(t, root, "grove-extra/work/W-777.md", "not a record either\n")
	commit(t, root, "code only")
	git(t, root, "checkout", "-q", "-b", "records", "main")
	write(t, root, "grove/work/W-001-first.md", record("W-001", "work", "active", "Records body.\n"))
	commit(t, root, "records")
	git(t, root, "checkout", "-q", "main")

	trace := filepath.Join(t.TempDir(), "trace")
	t.Setenv("GIT_TRACE", trace)
	res := mustInspect(t, root, "")
	main, code, records := source(t, res, "committed", "refs/heads/main"), source(t, res, "committed", "refs/heads/code"), source(t, res, "committed", "refs/heads/records")
	if !res.Complete || !code.Valid || main.Commit == code.Commit {
		t.Fatalf("fixture: %s", dump(res))
	}
	if main.project != code.project {
		t.Error("branches that differ only outside the record folder should share one loaded project")
	}
	if main.project == records.project || find(t, group(t, res, "W-001"), "committed", "refs/heads/records").Record.Status != "active" {
		t.Error("a branch with a differing record folder needs its own project")
	}
	log, err := os.ReadFile(trace)
	if n := strings.Count(string(log), "git cat-file"); err != nil || n != 1 || strings.Contains(string(log), "ls-tree") {
		t.Errorf("three branches and a checkout should share one cat-file process and list no whole tree: %d, %v", n, err)
	}
}
