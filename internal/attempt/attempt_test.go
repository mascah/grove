package attempt

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"syscall"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/mascah/grove/internal/repo"
)

// TestMain lets the test binary be the owner the launcher starts, exactly as
// cmd/grove does, so the lifecycle tests exercise the real owner process.
func TestMain(m *testing.M) {
	if dir := os.Getenv(OwnerEnv); dir != "" {
		if os.Getenv("GROVE_TEST_OWNER_DIES") != "" { // an owner that exits before it does anything
			os.Exit(3)
		}
		os.Exit(Own(dir))
	}
	os.Exit(m.Run())
}

const config = "schema_version: 3\nrecords: grove\ntarget: main\n"

const work = "---\nid: \"G-001\"\ntype: work\ntitle: First\nstatus: %s\ncreated: \"2026-09-22T10:00:00Z\"\nupdated: \"2026-09-22T10:00:00Z\"\n---\n\n## Outcome\n\nA thing.\n\n## Acceptance\n\n1. It is.\n"

const question = "---\nid: \"G-002\"\ntype: question\ntitle: Which colour?\nstatus: open\ncreated: \"2026-09-22T11:00:00Z\"\nupdated: \"2026-09-22T11:00:00Z\"\nblocks: [\"G-001\"]\n---\n\nRed or blue.\n"

var now = time.Date(2026, 9, 22, 18, 30, 0, 0, time.UTC)

func git(t *testing.T, dir string, args ...string) string {
	t.Helper()
	full := append([]string{"-c", "user.name=t", "-c", "user.email=t@t", "-c", "commit.gpgsign=false", "-c", "maintenance.auto=false"}, args...)
	out, err := repo.Command(context.Background(), dir, full...).CombinedOutput()
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

// fixture is a checkout on main with G-001 proposed, committed.
func fixture(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	dir, err := filepath.EvalSymlinks(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	root := filepath.Join(dir, "repo")
	if err := os.Mkdir(root, 0o755); err != nil {
		t.Fatal(err)
	}
	git(t, root, "init", "-q", "-b", "main") // identity comes from git()'s -c flags; the fakes that commit pass their own
	write(t, root, "grove.yaml", config)
	write(t, root, "grove/G-001-first.md", fmt.Sprintf(work, "proposed"))
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "init")
	return root
}

// fake installs a shell script as the provider. It answers --version, finds
// its --session-id, then runs body with $SID, $MARK (a file to append starts
// to) and $RELEASE (a file whose appearance ends a wait) set.
func fake(t *testing.T, body string) (mark, release string) {
	t.Helper()
	dir := t.TempDir()
	mark, release = filepath.Join(dir, "starts"), filepath.Join(dir, "release")
	script := "#!/bin/sh\n" +
		"[ \"$1\" = --version ] && { echo 'fake 0.1'; exit 0; }\n" +
		"SID=; while [ $# -gt 0 ]; do [ \"$1\" = --session-id ] && SID=$2; shift; done\n" +
		"MARK=" + mark + "; RELEASE=" + release + "\n" +
		"echo start >> \"$MARK\"\n" +
		"env | grep -E '^(CLAUDE|GROVE_ATTEMPT|GIT_DIR|GIT_WORK_TREE)' > env.txt\n" +
		body + "\n"
	path := filepath.Join(dir, "claude")
	if err := os.WriteFile(path, []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv(ClaudeEnv, path)
	return mark, release
}

const initLine = `echo '{"type":"system","subtype":"init","model":"fake-model","permissionMode":"acceptEdits","claude_code_version":"2.1.280","tools":["Bash","Edit"],"capabilities":["cap_v1"]}'`

func resultLine(subtype string, isError bool) string {
	return `echo '{"type":"result","subtype":"` + subtype + `","is_error":` + fmt.Sprint(isError) + `,"session_id":"'"$SID"'","total_cost_usd":0.5,"num_turns":3,"duration_ms":10,"permission_denials":[{"a":1}],"result":"the text"}'`
}

// waiting emits init, then waits for $RELEASE, then the result; SIGINT ends it
// with an interrupted result, as the provider ends its turn.
const waiting = `trap 'echo "{\"type\":\"result\",\"subtype\":\"interrupted\",\"is_error\":true,\"session_id\":\"$SID\"}"; exit 130' INT
` + initLine + `
echo partial > partial.txt
while [ ! -e "$RELEASE" ]; do sleep 0.05; done
` + `echo '{"type":"assistant"}'
`

// start launches the default request at a distinct time per call.
func start(t *testing.T, root string, at time.Time) (*Launch, []string) {
	t.Helper()
	var facts []string
	l, err := Start(Request{Root: root, ID: "G-001", BudgetUSD: "1", PermissionMode: "acceptEdits"}, at, func(f string) { facts = append(facts, f) })
	if err != nil {
		t.Fatal(err)
	}
	return l, facts
}

// await polls through List, which reads files only (Show adds a Git process
// for the inputs check), and returns the Show of the attempt once it is want.
func await(t *testing.T, root, attempt string, want Status) *View {
	t.Helper()
	deadline := time.Now().Add(20 * time.Second)
	for {
		views, err := List(root, "")
		if err != nil {
			t.Fatal(err)
		}
		i := slices.IndexFunc(views, func(v View) bool { return v.Launch.Attempt == attempt })
		if i < 0 {
			t.Fatalf("attempt %s is not listed", attempt)
		}
		if v := views[i]; v.Status == want {
			shown, err := Show(root, attempt)
			if err != nil {
				t.Fatal(err)
			}
			return shown
		}
		if time.Now().After(deadline) {
			t.Fatalf("attempt %s is %s, not %s", attempt, views[i].Status, want)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func release(t *testing.T, path string) {
	t.Helper()
	if err := os.WriteFile(path, nil, 0o644); err != nil {
		t.Fatal(err)
	}
}

func starts(t *testing.T, mark string) int {
	t.Helper()
	data, _ := os.ReadFile(mark)
	return strings.Count(string(data), "start")
}

func skipShort(t *testing.T) {
	if testing.Short() {
		t.Skip("starts an owner process and waits for a fake provider") // AGENTS.md: process-spawning tests skip under -short
	}
}

// awaitChild waits until the owner has recorded the provider's process group.
func awaitChild(t *testing.T, root, attempt string) *View {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for {
		v, err := Show(root, attempt)
		if err != nil {
			t.Fatal(err)
		}
		if v.ChildPGID != 0 {
			return v
		}
		if time.Now().After(deadline) {
			t.Fatalf("attempt %s never recorded its provider", attempt)
		}
		time.Sleep(20 * time.Millisecond)
	}
}

func TestRunToResult(t *testing.T) {
	skipShort(t)
	root := fixture(t)
	for _, kv := range [][2]string{{"CLAUDECODE", "1"}, {"CLAUDE_CODE_ENTRYPOINT", "cli"}, {"CLAUDE_PID", "1"}, {"GIT_DIR", "/nowhere"}, {"GIT_WORK_TREE", "/nowhere"}, {"CLAUDE_CONFIG_DIR", "/kept"}} {
		t.Setenv(kv[0], kv[1]) // repo.Command drops the Git ones for the test's own Git processes
	}
	fake(t, initLine+"\necho '{\"type\":\"assistant\"}'\necho '{\"type\":\"weird\"}'\n"+resultLine("success", false))
	l, facts := start(t, root, now)
	if len(facts) != 1 || !strings.Contains(facts[0], "worktree-G-001 created at") {
		t.Fatalf("facts %q", facts)
	}
	want := filepath.Join(root, ".claude", "worktrees", "worktree-G-001")
	if l.Worktree != want || l.Branch != "worktree-G-001" || l.WorktreeReused || l.Base != git(t, root, "rev-parse", "HEAD") {
		t.Fatalf("launch %+v", l)
	}
	if got := git(t, l.Worktree, "symbolic-ref", "--short", "HEAD"); got != "worktree-G-001" {
		t.Fatalf("worktree branch %s", got)
	}
	cmd := strings.Join(l.Command, " ")
	for _, part := range []string{"-p /grove-work G-001 --interaction headless", "--output-format stream-json --verbose", "--session-id " + l.SessionID, "--max-budget-usd 1", "--permission-mode acceptEdits --permission-prompts none"} {
		if !strings.Contains(cmd, part) {
			t.Fatalf("command %q lacks %q", cmd, part)
		}
	}
	if l.ClaudeVersion != "fake 0.1" || l.Owner <= 0 || l.Attempt != "G-001.20260922T183000Z" {
		t.Fatalf("launch %+v", l)
	}
	v := await(t, root, l.Attempt, Finished)
	r := v.Result
	if r.ExitCode != 0 || r.Stopped || !r.Dirty || r.Head != l.Base { // dirty: the fake left env.txt untracked
		t.Fatalf("result %+v", r)
	}
	if r.Record == nil || r.Record.Status != "proposed" || r.RecordError != "" || r.Record.Path != "grove/G-001-first.md" || r.RecordUncommitted {
		t.Fatalf("record %+v %q", r.Record, r.RecordError)
	}
	ev := r.Events
	if ev.Init == nil || ev.Init.Model != "fake-model" || ev.Init.Tools != 2 || ev.Init.Version != "2.1.280" || len(ev.Init.Capabilities) != 1 {
		t.Fatalf("init %+v", ev.Init)
	}
	if ev.Result == nil || ev.Result.SessionID != l.SessionID || ev.Result.CostUSD != 0.5 || ev.Result.PermissionDenials != 1 || ev.Result.Turns != 3 {
		t.Fatalf("result event %+v", ev.Result)
	}
	if ev.Lines != 4 || ev.Unknown != 1 || ev.Types["assistant"] != 1 || ev.ResultText != len("the text") || ev.Partial {
		t.Fatalf("events %+v", ev)
	}
	env, err := os.ReadFile(filepath.Join(l.Worktree, "env.txt"))
	if err != nil || strings.TrimSpace(string(env)) != "CLAUDE_CONFIG_DIR=/kept" {
		t.Fatalf("the provider saw %q (%v); only CLAUDE_CONFIG_DIR may reach it of the CLAUDE*, GROVE_ATTEMPT_OWNER and GIT_* variables", env, err)
	}
	log, _ := os.ReadFile(filepath.Join(v.Dir, "owner.log"))
	sid, _ := syscall.Getsid(0)
	if !strings.Contains(string(log), fmt.Sprintf("owner pid %d sid %d;", l.Owner, l.Owner)) || strings.Contains(string(log), fmt.Sprintf(" sid %d;", sid)) {
		t.Fatalf("the owner must lead its own session, apart from the launcher's %d:\n%s", sid, log)
	}
	exclude, _ := os.ReadFile(filepath.Join(root, ".git", "info", "exclude"))
	if !strings.Contains(string(exclude), "/.claude/worktrees/worktree-G-001/\n") {
		t.Fatalf("exclude %q", exclude)
	}
	if out := git(t, root, "status", "--porcelain"); out != "" {
		t.Fatalf("the launching checkout shows %q", out)
	}
	views, err := List(root, "")
	if err != nil || len(views) != 1 || views[0].Status != Finished {
		t.Fatalf("list %v %v", views, err)
	}

	// Provider failures, in the same checkout: an error result with exit 1
	// after an uncommitted record edit, no result event with exit 2, no
	// executable at all.
	fake(t, initLine+"\necho edited >> grove/G-001-first.md\n"+resultLine("error_max_budget_usd", true)+"\nexit 1")
	l, _ = start(t, root, now.Add(time.Minute))
	v = await(t, root, l.Attempt, Finished)
	if v.Result.ExitCode != 1 || v.Result.Events.Result == nil || !v.Result.Events.Result.IsError || v.Result.Events.Result.Subtype != "error_max_budget_usd" || !v.Result.RecordUncommitted {
		t.Fatalf("%+v", v.Result)
	}
	// No result event at all, and a nonzero exit.
	fake(t, initLine+"\necho '{\"type\":\"assistant\"}'\nexit 2")
	l2, _ := start(t, root, now.Add(2*time.Minute))
	v = await(t, root, l2.Attempt, Finished)
	if v.Result.ExitCode != 2 || v.Result.Events.Result != nil || v.Result.Events.Lines != 2 {
		t.Fatalf("%+v", v.Result)
	}
	// The provider cannot start at all.
	t.Setenv(ClaudeEnv, filepath.Join(t.TempDir(), "missing"))
	if _, err := Start(Request{Root: root, ID: "G-001", BudgetUSD: "1", PermissionMode: "acceptEdits"}, now.Add(3*time.Minute), func(string) {}); err == nil || !strings.Contains(err.Error(), "provider executable is not available") {
		t.Fatalf("%v", err)
	}
}

func TestCompetingStartAndReconnect(t *testing.T) {
	skipShort(t)
	root := fixture(t)
	mark, rel := fake(t, waiting+resultLine("success", false))
	l, _ := start(t, root, now)
	v := await(t, root, l.Attempt, Running)
	if v.Events == nil || v.Events.Init == nil {
		// init is written before the wait; give the file a moment
		deadline := time.Now().Add(5 * time.Second)
		for v.Events == nil || v.Events.Init == nil {
			if time.Now().After(deadline) {
				t.Fatalf("no init yet: %+v", v.Events)
			}
			time.Sleep(50 * time.Millisecond)
			v, _ = Show(root, l.Attempt)
		}
	}
	_, err := Start(Request{Root: root, ID: "G-001", BudgetUSD: "1", PermissionMode: "acceptEdits"}, now.Add(time.Second), func(string) {})
	if err == nil || !strings.Contains(err.Error(), "attempt "+l.Attempt+" of G-001 is running since") {
		t.Fatalf("competing start: %v", err)
	}
	for range 3 { // reconnecting is a read
		if v, err := Show(root, l.Attempt); err != nil || v.Status != Running {
			t.Fatalf("%v %v", v, err)
		}
	}
	if n := starts(t, mark); n != 1 {
		t.Fatalf("the provider started %d times", n)
	}
	release(t, rel)
	v = await(t, root, l.Attempt, Finished)
	if v.Result.ExitCode != 0 || v.Result.Events.Result == nil || v.Result.Events.Types["assistant"] != 1 {
		t.Fatalf("result %+v", v.Result)
	}
	if n := starts(t, mark); n != 1 {
		t.Fatalf("the provider started %d times", n)
	}
	// A next attempt reuses the branch's worktree and keeps what was left there.
	l2, facts := start(t, root, now.Add(time.Minute))
	if !l2.WorktreeReused || l2.Attempt == l.Attempt || len(facts) != 1 || !strings.Contains(facts[0], "reusing") {
		t.Fatalf("%+v %q", l2, facts)
	}
	release(t, rel)
	await(t, root, l2.Attempt, Finished)
	if _, err := os.Stat(filepath.Join(l.Worktree, "partial.txt")); err != nil {
		t.Fatal("partial work was not preserved:", err)
	}
	views, _ := List(root, "G-001")
	if len(views) != 2 || views[0].Launch.Attempt != l2.Attempt {
		t.Fatalf("list %+v", views)
	}
}

func TestStopAfterReconnect(t *testing.T) {
	skipShort(t)
	root := fixture(t)
	mark, _ := fake(t, waiting+resultLine("success", false))
	l, _ := start(t, root, now)
	await(t, root, l.Attempt, Running)
	var facts []string
	if err := Stop(root, l.Attempt, func(f string) { facts = append(facts, f) }); err != nil {
		t.Fatal(err)
	}
	if len(facts) != 1 || !strings.Contains(facts[0], fmt.Sprintf("SIGTERM sent to owner %d", l.Owner)) {
		t.Fatalf("%q", facts)
	}
	v := await(t, root, l.Attempt, Finished)
	r := v.Result
	if !r.Stopped || r.ExitCode != 130 || r.Events.Result == nil || r.Events.Result.Subtype != "interrupted" || !r.Events.Result.IsError {
		t.Fatalf("result %+v %+v", r, r.Events.Result)
	}
	if _, err := os.Stat(filepath.Join(l.Worktree, "partial.txt")); err != nil {
		t.Fatal("partial work was not preserved:", err)
	}
	if n := starts(t, mark); n != 1 {
		t.Fatalf("the provider started %d times", n)
	}
	if err := Stop(root, l.Attempt, func(string) {}); err == nil || !strings.Contains(err.Error(), "is finished; nothing to stop") {
		t.Fatalf("second stop: %v", err)
	}

	// A Stop right after Start, with no read in between, races the owner's start-up and must be seen either way.
	l, _ = start(t, root, now.Add(time.Minute))
	// No read in between: the signal races the owner's start-up and must be
	// seen either way.
	if err := syscall.Kill(l.Owner, syscall.SIGTERM); err != nil {
		t.Fatal(err)
	}
	v = await(t, root, l.Attempt, Finished)
	if !v.Result.Stopped {
		t.Fatalf("%+v", v.Result)
	}
	if n := starts(t, mark); n > 2 {
		t.Fatalf("the provider started %d times", n)
	}
	if n := starts(t, mark); n == 2 && v.Result.ExitCode != 130 {
		t.Fatalf("started then stopped, but exit %d", v.Result.ExitCode)
	}
}

func TestOwnerLost(t *testing.T) {
	skipShort(t)
	root := fixture(t)
	_, rel := fake(t, waiting+resultLine("success", false))
	l, _ := start(t, root, now)
	v := awaitChild(t, root, l.Attempt)
	if err := syscall.Kill(l.Owner, syscall.SIGKILL); err != nil {
		t.Fatal(err)
	}
	v = await(t, root, l.Attempt, Orphaned)
	if _, err := Start(Request{Root: root, ID: "G-001", BudgetUSD: "1", PermissionMode: "acceptEdits"}, now.Add(time.Second), func(string) {}); err == nil || !strings.Contains(err.Error(), "is orphaned since") {
		t.Fatalf("start over an orphan: %v", err)
	}
	var facts []string
	if err := Stop(root, l.Attempt, func(f string) { facts = append(facts, f) }); err != nil {
		t.Fatal(err)
	}
	v = await(t, root, l.Attempt, Finished)
	if v.Result.ReconciledBy != "stop" || !v.Result.Stopped || v.Result.Events.Init == nil || len(facts) == 0 || !strings.Contains(facts[len(facts)-1], "was orphaned") {
		t.Fatalf("%+v %q", v.Result, facts)
	}
	if syscall.Kill(-v.ChildPGID, 0) == nil {
		t.Fatal("the provider's process group is still alive")
	}

	// Owner lost and nothing left alive: interrupted, and a new attempt may start.
	l2, _ := start(t, root, now.Add(time.Minute))
	v = awaitChild(t, root, l2.Attempt)
	syscall.Kill(l2.Owner, syscall.SIGKILL)
	syscall.Kill(-v.ChildPGID, syscall.SIGKILL)
	v = await(t, root, l2.Attempt, Interrupted)
	if v.Result != nil || v.Events == nil || v.Events.Init == nil {
		t.Fatalf("%+v", v)
	}
	if err := Stop(root, l2.Attempt, func(string) {}); err == nil || !strings.Contains(err.Error(), "is interrupted; nothing to stop") {
		t.Fatalf("%v", err)
	}
	l3, _ := start(t, root, now.Add(2*time.Minute))
	release(t, rel)
	await(t, root, l3.Attempt, Finished)
}

func TestBlockingQuestionStopsTheNextRun(t *testing.T) {
	skipShort(t)
	root := fixture(t)
	// The attempt persists a question that blocks its work, on its branch.
	fake(t, initLine+"\nprintf '%s' '"+strings.ReplaceAll(question, "'", "'\\''")+"' > grove/G-002-q.md\ngit add -A && git -c user.name=t -c user.email=t@t -c commit.gpgsign=false -c maintenance.auto=false commit -qm 'question' || exit 3\n"+resultLine("success", false))
	l, _ := start(t, root, now)
	v := await(t, root, l.Attempt, Finished)
	if v.Result.ExitCode != 0 || v.Result.Head == l.Base {
		t.Fatalf("%+v", v.Result)
	}
	_, err := Start(Request{Root: root, ID: "G-001", BudgetUSD: "1", PermissionMode: "acceptEdits"}, now.Add(time.Second), func(string) {})
	if err == nil || !strings.Contains(err.Error(), "G-001 is blocked by open question G-002 (Which colour?)") {
		t.Fatalf("%v", err)
	}
	if _, err := Start(Request{Root: root, ID: "G-001", BudgetUSD: "1", PermissionMode: "acceptEdits"}, now.Add(2*time.Second), func(string) {}); err == nil || !strings.Contains(err.Error(), "blocked by open question G-002") {
		t.Fatalf("an unchanged wait must refuse the same way: %v", err)
	}
	views, _ := List(root, "G-001")
	if len(views) != 1 {
		t.Fatalf("a refused run left %d attempts", len(views))
	}
}

// TestInputsChanged also covers a project below the checkout's top: the
// guards, the record state and the inputs check must use the prefix.
func TestInputsChanged(t *testing.T) {
	skipShort(t)
	top := fixture(t)
	root := filepath.Join(top, "sub")
	if err := os.Mkdir(root, 0o755); err != nil {
		t.Fatal(err)
	}
	git(t, top, "mv", "grove.yaml", "grove", "sub/")
	git(t, top, "commit", "-qm", "move the project below the top")
	fake(t, initLine+"\n"+resultLine("success", false))
	l, err := Start(Request{Root: root, ID: "G-001", BudgetUSD: "1", PermissionMode: "acceptEdits"}, now, func(string) {})
	if err != nil {
		t.Fatal(err)
	}
	if l.Prefix != "sub" || l.Worktree != filepath.Join(root, ".claude", "worktrees", "worktree-G-001") {
		t.Fatalf("%+v", l)
	}
	v := await(t, root, l.Attempt, Finished)
	if v.InputsChanged != "" || v.Result.Record == nil || v.Result.Record.Status != "proposed" {
		t.Fatalf("%q %+v %q", v.InputsChanged, v.Result.Record, v.Result.RecordError)
	}
	if _, err := os.Stat(filepath.Join(l.Worktree, "sub", "env.txt")); err != nil {
		t.Fatal("the provider did not run in the project directory:", err)
	}
	// The branch's record enters review: a second run is refused from the project below the top.
	write(t, l.Worktree, "sub/grove/G-001-first.md", strings.Replace(fmt.Sprintf(work, "review"), "---\n\n## Outcome", "candidate: \""+l.Base+"\"\n---\n\n## Outcome", 1))
	git(t, l.Worktree, "commit", "-qam", "review")
	// Reached through a symlink: Git still names the prefix, and the branch's review refuses.
	link := filepath.Join(t.TempDir(), "link")
	if err := os.Symlink(root, link); err != nil {
		t.Fatal(err)
	}
	if _, err := Start(Request{Root: link, ID: "G-001", BudgetUSD: "1", PermissionMode: "acceptEdits"}, now.Add(time.Minute), func(string) {}); err == nil || !strings.Contains(err.Error(), "G-001 is review on worktree-G-001") {
		t.Fatal(err)
	}
	write(t, root, "grove/G-001-first.md", fmt.Sprintf(work, "active"))
	git(t, top, "commit", "-qam", "activate")
	v, _ = Show(root, l.Attempt)
	if !strings.Contains(v.InputsChanged, "grove/G-001-first.md on main is sha256:") || !strings.Contains(v.InputsChanged, "launched from "+l.RecordRevision) {
		t.Fatal(v.InputsChanged)
	}
}

func TestRefusals(t *testing.T) {
	root := fixture(t)
	fake(t, resultLine("success", false))
	try := func(req Request, want string) {
		t.Helper()
		if req.Root == "" {
			req.Root = root
		}
		if req.ID == "" {
			req.ID = "G-001"
		}
		_, err := Start(req, now, func(string) {})
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Fatalf("want %q, got %v", want, err)
		}
	}
	try(Request{BudgetUSD: "1"}, "run requires --budget USD and --permission-mode MODE")
	try(Request{PermissionMode: "auto"}, "run requires --budget USD and --permission-mode MODE")
	try(Request{ID: "G-009", BudgetUSD: "1", PermissionMode: "auto"}, "G-009 is not in this checkout")
	try(Request{ID: "nope", BudgetUSD: "1", PermissionMode: "auto"}, "nope is not a record ID")
	try(Request{BudgetUSD: "1", PermissionMode: "auto", Expect: "sha256:old"}, "G-001 changed since it was read: grove/G-001-first.md is sha256:")
	write(t, root, "grove/G-002-q.md", question)
	git(t, root, "add", "-A")
	git(t, root, "commit", "-qm", "question")
	try(Request{BudgetUSD: "1", PermissionMode: "auto"}, "G-001 is blocked by open question G-002 (Which colour?)")
	try(Request{ID: "G-002", BudgetUSD: "1", PermissionMode: "auto"}, "G-002 is a question, not work")
	git(t, root, "rm", "-q", "grove/G-002-q.md")
	git(t, root, "commit", "-qm", "resolved")
	write(t, root, "grove/G-001-first.md", fmt.Sprintf(work, "active"))
	try(Request{BudgetUSD: "1", PermissionMode: "auto"}, "grove/G-001-first.md has uncommitted changes in this checkout")
	git(t, root, "checkout", "-q", "--", "grove/G-001-first.md")
	write(t, root, "grove/G-001-first.md", strings.Replace(fmt.Sprintf(work, "review"), "---\n\n## Outcome", "candidate: \""+git(t, root, "rev-parse", "HEAD")+"\"\n---\n\n## Outcome", 1))
	git(t, root, "commit", "-qam", "review")
	try(Request{BudgetUSD: "1", PermissionMode: "auto"}, "G-001 is review; only proposed or active work runs")
	git(t, root, "revert", "--no-edit", "HEAD")
	// The branch's own record already in review: judgment, not another attempt.
	wt := filepath.Join(root, ".claude", "worktrees", "worktree-G-001")
	git(t, root, "worktree", "add", "-q", "-b", "worktree-G-001", wt)
	write(t, wt, "grove/G-001-first.md", strings.Replace(fmt.Sprintf(work, "review"), "---\n\n## Outcome", "candidate: \""+git(t, root, "rev-parse", "HEAD")+"\"\n---\n\n## Outcome", 1))
	git(t, wt, "commit", "-qam", "review")
	if _, err := Start(Request{Root: root, ID: "G-001", BudgetUSD: "1", PermissionMode: "auto"}, now, func(string) {}); err == nil || !strings.Contains(err.Error(), "G-001 is review on worktree-G-001 at "+wt+"; judge that candidate") {
		t.Fatal(err)
	}
	git(t, root, "worktree", "remove", "--force", wt)
	git(t, root, "branch", "-qD", "worktree-G-001")
	// A directory in the way that is not the branch's worktree.
	write(t, wt, "stray.txt", "x")
	try(Request{BudgetUSD: "1", PermissionMode: "auto"}, wt+" exists but is not a registered worktree of worktree-G-001")
	if views, err := List(root, ""); err != nil || len(views) != 0 {
		t.Fatalf("refusals wrote attempts: %v %v", views, err)
	}
	// An owner that exits before doing anything is reported, not announced as started.
	if testing.Short() {
		return // this case starts an owner process and waits on it (AGENTS.md)
	}
	if err := os.RemoveAll(wt); err != nil {
		t.Fatal(err)
	}
	t.Setenv("GROVE_TEST_OWNER_DIES", "1")
	if l, err := Start(Request{BudgetUSD: "1", PermissionMode: "auto", Root: root, ID: "G-001"}, now.Add(time.Hour), func(string) {}); err == nil || !strings.Contains(err.Error(), "exited while starting") {
		t.Fatalf("%+v %v", l, err)
	}
	if views, _ := List(root, ""); len(views) != 1 || views[0].Status != Interrupted {
		t.Fatalf("%+v", views)
	}
	if _, err := Show(root, "G-001.20260922T183000Z"); err == nil || !strings.Contains(err.Error(), "does not exist in this repository") {
		t.Fatal(err)
	}
	if _, err := Show(root, "bogus"); err == nil || !strings.Contains(err.Error(), "is not an attempt id") {
		t.Fatal(err)
	}
}

func TestReadEventsBounded(t *testing.T) {
	path := filepath.Join(t.TempDir(), "events.jsonl")
	huge := strings.Repeat("x", MaxLine+10)
	content := `{"type":"system","subtype":"init","model":"m","tools":[]}` + "\n" +
		"not json\n" +
		"\n" +
		`{"type":"assistant","message":{}}` + "\n" +
		`{"type":"mystery"}` + "\n" +
		`{"type":"assistant","big":"` + huge + `"}` + "\n" +
		`{"type":"result","subtype":"success","is_error":true,"session_id":"s","result":"abc"}` + "\n" +
		`{"type":"assistant","cut`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	ev, err := ReadEvents(path)
	if err != nil {
		t.Fatal(err)
	}
	if ev.Lines != 8 || ev.Malformed != 3 || ev.Unknown != 1 || ev.Oversized != 1 || !ev.Partial || ev.Bytes != int64(len(content)) {
		t.Fatalf("%+v", ev)
	}
	if ev.Types["assistant"] != 1 || ev.Types["system"] != 1 || ev.Types["result"] != 1 {
		t.Fatalf("%+v", ev.Types)
	}
	if ev.Init == nil || ev.Init.Model != "m" || ev.Result == nil || !ev.Result.IsError || ev.Result.SessionID != "s" || ev.ResultText != 3 {
		t.Fatalf("%+v %+v", ev.Init, ev.Result)
	}
	if ev, err := ReadEvents(filepath.Join(t.TempDir(), "none")); err != nil || ev.Lines != 0 {
		t.Fatalf("%+v %v", ev, err)
	}
}

func TestReadActivityBounded(t *testing.T) {
	path := filepath.Join(t.TempDir(), "events.jsonl")
	if a, err := ReadActivity(path, ActivityWindow); err != nil || len(a.Lines) != 0 || a.Cut {
		t.Fatalf("%+v %v", a, err)
	}
	var b strings.Builder
	b.WriteString(`{"type":"system","subtype":"init","model":"m"}` + "\n")
	for i := range 300 {
		fmt.Fprintf(&b, `{"type":"assistant","message":{"content":[{"type":"text","text":"step %d\nmore"},{"type":"tool_use","name":"Bash","input":{"command":"go test\n-v"}}]}}`+"\n", i)
	}
	b.WriteString(`{"type":"user","message":{"content":[{"type":"tool_result","is_error":true}]}}` + "\n")
	b.WriteString(`{"type":"assistant","message":{"content":[{"type":"text","text":"` + strings.Repeat("é", 400) + `"}]}}` + "\n")
	b.WriteString("not json\n")
	b.WriteString(`{"type":"result","subtype":"success","is_error":false,"result":"## Done\nAll of it."}` + "\n")
	b.WriteString(`{"type":"assistant","message":{"content":[{"type":"text","text":"still being writ`)
	if err := os.WriteFile(path, []byte(b.String()), 0o644); err != nil {
		t.Fatal(err)
	}
	a, err := ReadActivity(path, ActivityWindow)
	if err != nil {
		t.Fatal(err)
	}
	n := len(a.Lines)
	if n != maxActivity || !a.Cut || a.Report != "## Done\nAll of it." {
		t.Fatalf("%d lines, cut %v, report %q", n, a.Cut, a.Report)
	}
	if a.Lines[n-1] != "result: success" || a.Lines[n-3] != "tool error" || a.Lines[n-4] != "tool: Bash go test" || a.Lines[n-5] != "step 299" {
		t.Fatalf("%q", a.Lines[n-5:])
	}
	if long := a.Lines[n-2]; len(long) > maxActivityLine+len("…") || !strings.HasSuffix(long, "…") || !utf8.ValidString(long) {
		t.Fatalf("%q", long)
	}
	// A window that starts mid-file drops the line it cut into.
	a, err = ReadActivity(path, 200)
	if err != nil || !a.Cut || len(a.Lines) != 1 || a.Lines[0] != "result: success" {
		t.Fatalf("%+v %v", a, err)
	}
}
