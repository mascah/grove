// Package attempt runs one bounded implementation of a work record as a
// Grove-owned `claude -p` process that outlives the launching terminal, and
// reads what such a run left behind (G-045, plan G-100).
//
// An attempt is a directory under the repository's Git common directory,
// <common>/grove/attempts/<WORK.TIMESTAMP>/, shared by every worktree and
// never committed:
//
//	attempt.json  what was launched: inputs, worktree, command, budget, ids
//	owner.lock    an flock the owner process holds for its whole lifetime
//	owner.log     the owner's own notes: pids, signals, exit
//	events.jsonl  the provider's stdout, raw, written by the kernel
//	stderr.log    the provider's stderr, raw
//	result.json   written once when the process is gone: exit, fields, counts
//
// Liveness is the lock, never a pid alone: running means the owner holds
// owner.lock; owner lost with the child's process group still alive is
// orphaned; owner lost with nothing alive and no result is interrupted.
// Reading an attempt starts no process and resumes no conversation.
package attempt

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/repo"
)

// OwnerEnv names the attempt directory to a Grove binary that should run as
// that attempt's owner instead of the CLI. The child never sees it.
const OwnerEnv = "GROVE_ATTEMPT_OWNER"

// ClaudeEnv overrides the provider executable, so tests run a fake process
// through the real owner; the value is recorded in attempt.json.
const ClaudeEnv = "GROVE_CLAUDE"

// MaxLine bounds one event line a reader parses; longer lines stay in the raw
// file and are counted as oversized.
const MaxLine = 1 << 20

// stopGrace is how long the owner waits after SIGINT, which ends the turn,
// before SIGKILL to the child's process group.
const stopGrace = 15 * time.Second

// Launch is what the launcher records before the owner starts.
type Launch struct {
	Attempt        string    `json:"attempt"`
	Work           string    `json:"work"`
	Project        string    `json:"project"`         // the launching checkout's project root
	Target         string    `json:"target"`          // grove.yaml's target branch, "" when none
	RecordPath     string    `json:"record_path"`     // project-relative
	RecordRevision string    `json:"record_revision"` // as HEAD held it at launch
	Base           string    `json:"base"`            // the commit the worktree started from, or continues on
	Branch         string    `json:"branch"`
	Worktree       string    `json:"worktree"`        // absolute checkout the process runs in
	WorktreeReused bool      `json:"worktree_reused"` // it existed before this attempt
	Command        []string  `json:"command"`         // the exact argv, command[0] the executable as resolved
	Executable     string    `json:"executable"`      // command[0] as given (claude or GROVE_CLAUDE)
	ClaudeVersion  string    `json:"claude_version"`  // `--version` at launch
	GroveVersion   string    `json:"grove_version"`
	Model          string    `json:"model,omitempty"` // requested; the actual one is in the result's init
	BudgetUSD      string    `json:"budget_usd"`
	PermissionMode string    `json:"permission_mode"`
	SessionID      string    `json:"session_id"` // generated here, passed as --session-id
	Started        time.Time `json:"started"`
	Owner          int       `json:"owner_pid"` // the owner, set by the launcher after the spawn
}

// Init is what the reader keeps of the provider's system/init event.
type Init struct {
	Model          string   `json:"model,omitempty"`
	PermissionMode string   `json:"permission_mode,omitempty"`
	Version        string   `json:"claude_code_version,omitempty"`
	Tools          int      `json:"tools"`
	Capabilities   []string `json:"capabilities,omitempty"`
	Agents         []string `json:"agents,omitempty"`
}

// Final is what the reader keeps of the provider's result event.
type Final struct {
	Subtype           string  `json:"subtype"`
	IsError           bool    `json:"is_error"`
	SessionID         string  `json:"session_id"`
	CostUSD           float64 `json:"total_cost_usd"`
	Turns             int     `json:"num_turns"`
	DurationMS        int64   `json:"duration_ms"`
	PermissionDenials int     `json:"permission_denials"`
}

// Events counts what a bounded read of events.jsonl found.
type Events struct {
	Lines      int            `json:"lines"`
	Bytes      int64          `json:"bytes"`
	Types      map[string]int `json:"types"`     // known event types
	Unknown    int            `json:"unknown"`   // JSON objects whose type is not one the reader knows
	Malformed  int            `json:"malformed"` // lines that are not a JSON object
	Oversized  int            `json:"oversized"` // lines over MaxLine, skipped
	Partial    bool           `json:"partial"`   // the file ends without a newline
	Init       *Init          `json:"init"`
	Result     *Final         `json:"result"`
	ResultText int            `json:"result_text_bytes"` // length of the result's text, never printed
}

// Result is what the owner writes once the process is gone.
type Result struct {
	Finished     time.Time `json:"finished"`
	ExitCode     int       `json:"exit_code"` // -1 when the process died of a signal
	Signal       string    `json:"signal,omitempty"`
	Stopped      bool      `json:"stopped"`                 // a Stop was requested
	ReconciledBy string    `json:"reconciled_by,omitempty"` // "stop" when written for a lost owner
	Events       Events    `json:"events"`
	Head         string    `json:"head"`  // the worktree's HEAD after the exit
	Dirty        bool      `json:"dirty"` // uncommitted changes to tracked files
	Record       *State    `json:"record"`
	RecordError  string    `json:"record_error,omitempty"`
}

// State is the work record as the worktree holds it.
type State struct {
	Status    string `json:"status"`
	Candidate string `json:"candidate,omitempty"`
	Revision  string `json:"revision"`
}

// Status classifies an attempt from its files and the kernel alone.
type Status string

const (
	Running     Status = "running"     // the owner holds the lock
	Finished    Status = "finished"    // result.json exists
	Orphaned    Status = "orphaned"    // owner lost, child's process group alive
	Interrupted Status = "interrupted" // owner lost, nothing alive, no result
)

// View is one attempt as read.
type View struct {
	Dir           string  `json:"dir"`
	Launch        Launch  `json:"launch"`
	Status        Status  `json:"status"`
	ChildPGID     int     `json:"child_pgid,omitempty"`
	Result        *Result `json:"result,omitempty"`
	Events        *Events `json:"events,omitempty"`         // while there is no result: a bounded read so far
	InputsChanged string  `json:"inputs_changed,omitempty"` // the record on the target no longer hashes to the launch revision
	EventsPath    string  `json:"events_path"`
	StderrPath    string  `json:"stderr_path"`
}

// Request is one launch.
type Request struct {
	Root, ID       string // the launching project root and the work
	BudgetUSD      string
	PermissionMode string
	Model          string
	Branch         string // default worktree-ID
	Worktree       string // default <root>/.claude/worktrees/<branch>
}

var idPattern = regexp.MustCompile(`^[A-Z]+-[0-9]+$`)
var attemptPattern = regexp.MustCompile(`^[A-Z]+-[0-9]+\.[0-9]{8}T[0-9]{6}Z$`)

// Dir is where root's repository keeps attempts.
func Dir(root string) (string, error) {
	common, _, err := repo.CommonDir(root)
	if err != nil {
		return "", err
	}
	return filepath.Join(common, "grove", "attempts"), nil
}

// Start launches one attempt of req.ID and returns what was launched. Every
// refusal comes before anything is written; after the worktree exists,
// a failure to start the owner is reported with the worktree kept.
func Start(req Request, now time.Time, report func(string)) (*Launch, error) {
	if !idPattern.MatchString(req.ID) {
		return nil, fmt.Errorf("%s is not a record ID", req.ID)
	}
	if req.BudgetUSD == "" || req.PermissionMode == "" {
		return nil, errors.New("run requires --budget USD and --permission-mode MODE: Grove sets no default spend or permission profile")
	}
	p, ds := project.Load(req.Root, req.Root)
	if len(ds) != 0 {
		var lines []string
		for _, d := range ds {
			lines = append(lines, d.String())
		}
		return nil, fmt.Errorf("the project is not valid; fix it before running:\n%s", strings.Join(lines, "\n"))
	}
	root := p.Root
	var r *project.Record
	for _, c := range p.Records {
		if c.ID == req.ID {
			r = c
		}
	}
	switch {
	case r == nil:
		return nil, fmt.Errorf("%s is not in this checkout", req.ID)
	case r.Type != "work":
		return nil, fmt.Errorf("%s is a %s, not work", req.ID, r.Type)
	case r.Status != "proposed" && r.Status != "active":
		return nil, fmt.Errorf("%s is %s; only proposed or active work runs", req.ID, r.Status)
	}
	if err := blocked(p, req.ID); err != nil {
		return nil, err
	}
	if dirty, err := repo.Git(root, "status", "--porcelain", "--", r.Path); err != nil {
		return nil, err
	} else if strings.TrimSpace(dirty) != "" {
		return nil, fmt.Errorf("%s has uncommitted changes in this checkout; commit them so the attempt sees them", r.Path)
	}
	head, err := repo.Git(root, "rev-parse", "HEAD")
	if err != nil {
		return nil, fmt.Errorf("this checkout's HEAD could not be read: %v", err)
	}
	head = strings.TrimSpace(head)
	dir, err := Dir(root)
	if err != nil {
		return nil, err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	// One launch per repository at a time: the running check and the
	// attempt id are decided under this lock.
	unlock, err := repo.Lock(filepath.Join(dir, "launch.lock"))
	if err != nil {
		return nil, err
	}
	defer unlock()
	views, err := List(root, req.ID)
	if err != nil {
		return nil, err
	}
	for _, v := range views {
		if v.Status == Running || v.Status == Orphaned {
			return nil, fmt.Errorf("attempt %s of %s is %s since %s; stop it or wait for its result before another attempt", v.Launch.Attempt, req.ID, v.Status, v.Launch.Started.UTC().Format(time.RFC3339))
		}
	}
	branch := req.Branch
	if branch == "" {
		branch = "worktree-" + req.ID
	}
	worktree := req.Worktree
	if worktree == "" {
		worktree = filepath.Join(root, ".claude", "worktrees", branch)
	} else if !filepath.IsAbs(worktree) {
		worktree = filepath.Join(root, worktree)
	}
	worktree = filepath.Clean(worktree)
	exe := os.Getenv(ClaudeEnv)
	if exe == "" {
		exe = "claude"
	}
	resolved, err := exec.LookPath(exe)
	if err != nil {
		return nil, fmt.Errorf("the provider executable is not available: %v", err)
	}
	version, err := exec.Command(resolved, "--version").Output()
	if err != nil {
		return nil, fmt.Errorf("%s --version failed: %v", exe, err)
	}
	session, err := uuid()
	if err != nil {
		return nil, err
	}
	attempt := req.ID + "." + now.UTC().Format("20060102T150405Z")
	adir := filepath.Join(dir, attempt)
	if _, err := os.Stat(adir); err == nil {
		return nil, fmt.Errorf("attempt %s already exists; try again in a second", attempt)
	}
	base, reused, err := prepareWorktree(root, branch, worktree, head, report)
	if err != nil {
		return nil, err
	}
	// The branch may hold what an earlier attempt persisted: a candidate in
	// review awaiting the owner, or the question the headless guide writes
	// for a missing decision. Either is a wait, not a reason to spend again.
	if reused || base != head {
		if wp, _ := project.Load(worktree, worktree); wp != nil {
			for _, c := range wp.Records {
				if c.ID == req.ID && c.Status != "proposed" && c.Status != "active" {
					return nil, fmt.Errorf("%s is %s on %s at %s; judge that candidate (approve, feedback) before another attempt", req.ID, c.Status, branch, worktree)
				}
			}
			if err := blocked(wp, req.ID); err != nil {
				return nil, fmt.Errorf("%v (on %s at %s)", err, branch, worktree)
			}
		}
	}
	command := []string{resolved, "-p", "/grove-work " + req.ID + " --interaction headless",
		"--output-format", "stream-json", "--verbose",
		"--session-id", session,
		"--max-budget-usd", req.BudgetUSD,
		"--permission-mode", req.PermissionMode, "--permission-prompts", "none"}
	if req.Model != "" {
		command = append(command, "--model", req.Model)
	}
	l := &Launch{
		Attempt: attempt, Work: req.ID, Project: root, Target: p.Target,
		RecordPath: r.Path, RecordRevision: project.Revision(r.Source),
		Base: base, Branch: branch, Worktree: worktree, WorktreeReused: reused,
		Command: command, Executable: exe, ClaudeVersion: strings.TrimSpace(string(version)),
		GroveVersion: groveVersion(), Model: req.Model, BudgetUSD: req.BudgetUSD,
		PermissionMode: req.PermissionMode, SessionID: session, Started: now.UTC(),
	}
	if err := os.Mkdir(adir, 0o755); err != nil {
		return nil, err
	}
	if err := writeJSON(filepath.Join(adir, "attempt.json"), l); err != nil {
		return nil, err
	}
	// The lock is taken here and inherited, so it is held from before the
	// owner exists until it exits; the launcher's own descriptor closes on
	// return. A launcher that dies between Start and the owner's first
	// instruction leaves the lock with the owner, which already has it.
	lock, err := os.OpenFile(filepath.Join(adir, "owner.lock"), os.O_CREATE|os.O_RDWR, 0o644)
	if err != nil {
		return nil, err
	}
	defer lock.Close()
	if err := syscall.Flock(int(lock.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		return nil, fmt.Errorf("attempt %s is already owned", attempt)
	}
	self, err := os.Executable()
	if err != nil {
		return nil, err
	}
	log, err := os.OpenFile(filepath.Join(adir, "owner.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, err
	}
	defer log.Close()
	owner := exec.Command(self)
	owner.Dir = worktree
	owner.Env = append(environ(os.Environ(), OwnerEnv), OwnerEnv+"="+adir)
	owner.Stdin = nil // /dev/null
	owner.Stdout, owner.Stderr = log, log
	owner.ExtraFiles = []*os.File{lock} // fd 3
	owner.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	if err := owner.Start(); err != nil {
		return nil, fmt.Errorf("the owner could not start: %v (worktree %s is kept)", err, worktree)
	}
	l.Owner = owner.Process.Pid
	owner.Process.Release()
	if err := writeJSON(filepath.Join(adir, "attempt.json"), l); err != nil {
		return l, fmt.Errorf("the owner started as pid %d but attempt.json could not be rewritten: %v", l.Owner, err)
	}
	return l, nil
}

// blocked is the open question that stops id, if any, as p holds it.
func blocked(p *project.Project, id string) error {
	for _, q := range p.Records {
		if q.Type == "question" && q.Status == "open" && slices.Contains(q.Blocks, id) {
			return fmt.Errorf("%s is blocked by open question %s (%s); resolve it before another attempt", id, q.ID, q.Title)
		}
	}
	return nil
}

// prepareWorktree makes branch's checkout at worktree from head, or reuses
// the registered one, reporting what it did. It returns the commit the
// attempt starts from.
func prepareWorktree(root, branch, worktree, head string, report func(string)) (base string, reused bool, err error) {
	worktrees, err := repo.Worktrees(root)
	if err != nil {
		return "", false, err
	}
	ref := "refs/heads/" + branch
	for _, w := range worktrees {
		same := samePath(w.Path, worktree)
		switch {
		case w.Branch == ref && same:
			report(fmt.Sprintf("worktree: reusing %s on %s at %s", worktree, branch, short(w.Head)))
			return w.Head, true, nil
		case w.Branch == ref:
			return "", false, fmt.Errorf("branch %s is checked out at %s, not %s; pass --worktree %s to continue there", branch, w.Path, worktree, w.Path)
		case same:
			return "", false, fmt.Errorf("%s is the worktree of %s, not %s", worktree, orDetached(w.Branch), branch)
		}
	}
	if _, err := os.Lstat(worktree); err == nil {
		return "", false, fmt.Errorf("%s exists but is not a registered worktree of %s; remove it or pass --worktree elsewhere", worktree, branch)
	}
	if err := os.MkdirAll(filepath.Dir(worktree), 0o755); err != nil {
		return "", false, err
	}
	if _, err := repo.Git(root, "rev-parse", "-q", "--verify", ref); err == nil {
		tip, err := repo.Git(root, "rev-parse", ref)
		if err != nil {
			return "", false, err
		}
		if _, err := repo.Git(root, "worktree", "add", worktree, branch); err != nil {
			return "", false, err
		}
		report(fmt.Sprintf("worktree: %s checked out at %s from existing branch %s at %s", branch, worktree, branch, short(tip)))
		base = strings.TrimSpace(tip)
	} else {
		if _, err := repo.Git(root, "worktree", "add", "-b", branch, worktree, head); err != nil {
			return "", false, err
		}
		report(fmt.Sprintf("worktree: %s created at %s from %s", branch, worktree, short(head)))
		base = head
	}
	// A worktree inside the checkout would otherwise show as untracked there.
	if top, err := repo.GitPath(root, "--show-toplevel"); err == nil {
		if rel, err := filepath.Rel(top, worktree); err == nil && !strings.HasPrefix(rel, "..") {
			if repo.Command(context.Background(), root, "check-ignore", "-q", "--", worktree).Run() != nil {
				if common, _, err := repo.CommonDir(root); err == nil {
					if f, err := os.OpenFile(filepath.Join(common, "info", "exclude"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644); err == nil {
						fmt.Fprintf(f, "/%s/\n", filepath.ToSlash(rel))
						f.Close()
					}
				}
			}
		}
	}
	return base, false, nil
}

// Own runs as the attempt's owner: it holds the inherited lock, runs the
// provider in the worktree with its output in files, handles Stop, and
// writes the result. It returns the process exit code.
func Own(dir string) int {
	lock := os.NewFile(3, "owner.lock")
	if lock == nil {
		fmt.Fprintln(os.Stderr, "owner: no inherited lock on fd 3")
		return 1
	}
	syscall.CloseOnExec(3) // the child must not inherit the lock: it would read as running after the owner is gone
	defer lock.Close()
	logf := func(format string, args ...any) {
		fmt.Fprintf(os.Stderr, time.Now().UTC().Format(time.RFC3339)+" "+format+"\n", args...)
	}
	var l Launch
	if err := readJSON(filepath.Join(dir, "attempt.json"), &l); err != nil {
		logf("owner: %v", err)
		return 1
	}
	events, err := os.OpenFile(filepath.Join(dir, "events.jsonl"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		logf("owner: %v", err)
		return 1
	}
	stderr, err := os.OpenFile(filepath.Join(dir, "stderr.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		logf("owner: %v", err)
		return 1
	}
	child := exec.Command(l.Command[0], l.Command[1:]...)
	child.Dir = l.Worktree
	// A session's own variables would make the provider a child of the
	// launching session (CLAUDECODE guards nesting); the config directory is
	// the one CLAUDE variable that describes the machine, not a session.
	child.Env = environ(os.Environ(), OwnerEnv, "CLAUDE*", "!CLAUDE_CONFIG_DIR",
		"GIT_DIR", "GIT_WORK_TREE", "GIT_INDEX_FILE", "GIT_COMMON_DIR", "GIT_OBJECT_DIRECTORY", "GIT_ALTERNATE_OBJECT_DIRECTORIES", "GIT_NAMESPACE")
	child.Stdin = nil
	child.Stdout, child.Stderr = events, stderr
	child.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := child.Start(); err != nil {
		logf("owner: the provider could not start: %v", err)
		res := Result{Finished: time.Now().UTC(), ExitCode: -1, Signal: "not started: " + err.Error()}
		reconcile(dir, &l, &res, logf)
		return 1
	}
	events.Close()
	stderr.Close()
	pgid := child.Process.Pid
	if err := writeJSON(filepath.Join(dir, "child.json"), map[string]int{"pid": child.Process.Pid, "pgid": pgid, "owner_pid": os.Getpid()}); err != nil {
		logf("owner: %v", err)
	}
	logf("owner pid %d sid %d; provider pid %d pgid %d; %s", os.Getpid(), sid(), child.Process.Pid, pgid, strings.Join(l.Command, " "))

	waited := make(chan error, 1)
	go func() { waited <- child.Wait() }()
	signals := make(chan os.Signal, 2)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
	stopped := false
	var werr error
	select {
	case werr = <-waited:
	case s := <-signals:
		stopped = true
		logf("owner: %v received; SIGINT to the provider's group %d, SIGKILL after %s", s, pgid, stopGrace)
		syscall.Kill(-pgid, syscall.SIGINT)
		select {
		case werr = <-waited:
		case <-time.After(stopGrace):
			syscall.Kill(-pgid, syscall.SIGKILL)
			werr = <-waited
		}
	}
	res := Result{Finished: time.Now().UTC(), Stopped: stopped}
	res.ExitCode, res.Signal = exitOf(werr)
	logf("owner: the provider exited: code %d signal %q stopped %v", res.ExitCode, res.Signal, stopped)
	reconcile(dir, &l, &res, logf)
	return 0
}

// reconcile completes res from the files and the worktree and writes it.
func reconcile(dir string, l *Launch, res *Result, logf func(string, ...any)) {
	ev, err := ReadEvents(filepath.Join(dir, "events.jsonl"))
	if err != nil {
		logf("owner: events: %v", err)
	}
	res.Events = ev
	if head, err := repo.Git(l.Worktree, "rev-parse", "HEAD"); err == nil {
		res.Head = strings.TrimSpace(head)
	} else {
		logf("owner: HEAD: %v", err)
	}
	if dirty, err := repo.Git(l.Worktree, "status", "--porcelain", "--untracked-files=no"); err == nil {
		res.Dirty = strings.TrimSpace(dirty) != ""
	}
	res.Record, res.RecordError = recordState(l.Worktree, l.Work)
	if err := writeJSON(filepath.Join(dir, "result.json"), res); err != nil {
		logf("owner: result: %v", err)
	}
}

// recordState reads the work record as a checkout holds it, whether or not
// the rest of the project validates there.
func recordState(root, id string) (*State, string) {
	p, ds := project.Load(root, root)
	if p == nil {
		return nil, ds[0].String()
	}
	for _, r := range p.Records {
		if r.ID == id {
			return &State{Status: r.Status, Candidate: r.Candidate, Revision: project.Revision(r.Source)}, ""
		}
	}
	msg := id + " is not in the worktree"
	if len(ds) != 0 {
		msg += "; " + ds[0].String()
	}
	return nil, msg
}

// ReadEvents parses events.jsonl bounded: one line at a time, a line over
// MaxLine skipped and counted, keeping only the init and result fields.
func ReadEvents(path string) (Events, error) {
	ev := Events{Types: map[string]int{}}
	f, err := os.Open(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ev, nil
		}
		return ev, err
	}
	defer f.Close()
	r := bufio.NewReaderSize(f, 64<<10)
	for {
		var line []byte
		size, over := 0, false
		var err error
		for {
			var chunk []byte
			chunk, err = r.ReadSlice('\n')
			size += len(chunk)
			if over = over || size > MaxLine; !over {
				line = append(line, chunk...)
			}
			if !errors.Is(err, bufio.ErrBufferFull) {
				break
			}
		}
		ev.Bytes += int64(size)
		if size == 0 {
			break
		}
		ev.Lines++
		if err != nil { // the file ended before this line's newline
			ev.Partial = true
		}
		if over {
			ev.Oversized++
		} else {
			ev.note(bytes.TrimSpace(line))
		}
		if err != nil {
			break
		}
	}
	return ev, nil
}

func (ev *Events) note(line []byte) {
	if len(line) == 0 {
		ev.Malformed++
		return
	}
	var head struct {
		Type    string `json:"type"`
		Subtype string `json:"subtype"`
	}
	if err := json.Unmarshal(line, &head); err != nil || head.Type == "" {
		ev.Malformed++
		return
	}
	switch head.Type {
	case "system", "assistant", "user", "result", "stream_event", "rate_limit_event":
		ev.Types[head.Type]++
	default:
		ev.Unknown++
		return
	}
	switch {
	case head.Type == "system" && head.Subtype == "init":
		var init struct {
			Model          string   `json:"model"`
			PermissionMode string   `json:"permissionMode"`
			Version        string   `json:"claude_code_version"`
			Tools          []any    `json:"tools"`
			Capabilities   []string `json:"capabilities"`
			Agents         []string `json:"agents"`
		}
		if json.Unmarshal(line, &init) == nil {
			ev.Init = &Init{Model: init.Model, PermissionMode: init.PermissionMode, Version: init.Version, Tools: len(init.Tools), Capabilities: init.Capabilities, Agents: init.Agents}
		}
	case head.Type == "result":
		var final struct {
			Subtype    string  `json:"subtype"`
			IsError    bool    `json:"is_error"`
			SessionID  string  `json:"session_id"`
			CostUSD    float64 `json:"total_cost_usd"`
			Turns      int     `json:"num_turns"`
			DurationMS int64   `json:"duration_ms"`
			Denials    []any   `json:"permission_denials"`
			Result     string  `json:"result"`
		}
		if json.Unmarshal(line, &final) == nil {
			ev.Result = &Final{Subtype: final.Subtype, IsError: final.IsError, SessionID: final.SessionID, CostUSD: final.CostUSD, Turns: final.Turns, DurationMS: final.DurationMS, PermissionDenials: len(final.Denials)}
			ev.ResultText = len(final.Result)
		}
	}
}

// List reads every attempt of root's repository, newest first, or those of
// one work ID.
func List(root, id string) ([]View, error) {
	dir, err := Dir(root)
	if err != nil {
		return nil, err
	}
	entries, err := os.ReadDir(dir)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var views []View
	for _, e := range entries {
		name := e.Name()
		if !e.IsDir() || !attemptPattern.MatchString(name) || (id != "" && !strings.HasPrefix(name, id+".")) {
			continue
		}
		v, err := read(filepath.Join(dir, name))
		if err != nil {
			return nil, err
		}
		views = append(views, *v)
	}
	slices.SortFunc(views, func(a, b View) int { return b.Launch.Started.Compare(a.Launch.Started) })
	return views, nil
}

// Show reads one attempt by its id and adds whether its inputs changed.
func Show(root, attempt string) (*View, error) {
	if !attemptPattern.MatchString(attempt) {
		return nil, fmt.Errorf("%s is not an attempt id (WORK.YYYYMMDDTHHMMSSZ, from attempts)", attempt)
	}
	dir, err := Dir(root)
	if err != nil {
		return nil, err
	}
	v, err := read(filepath.Join(dir, attempt))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("attempt %s does not exist in this repository", attempt)
		}
		return nil, err
	}
	if v.Launch.Target != "" {
		if source, err := repo.Git(root, "show", "refs/heads/"+v.Launch.Target+":"+v.Launch.RecordPath); err != nil {
			v.InputsChanged = fmt.Sprintf("%s is not readable on %s: %v", v.Launch.RecordPath, v.Launch.Target, err)
		} else if rev := project.Revision([]byte(source)); rev != v.Launch.RecordRevision {
			v.InputsChanged = fmt.Sprintf("%s on %s is %s, launched from %s", v.Launch.RecordPath, v.Launch.Target, rev, v.Launch.RecordRevision)
		}
	}
	return v, nil
}

func read(dir string) (*View, error) {
	v := &View{Dir: dir, EventsPath: filepath.Join(dir, "events.jsonl"), StderrPath: filepath.Join(dir, "stderr.log")}
	if err := readJSON(filepath.Join(dir, "attempt.json"), &v.Launch); err != nil {
		return nil, err
	}
	var child struct {
		PGID int `json:"pgid"`
	}
	if readJSON(filepath.Join(dir, "child.json"), &child) == nil {
		v.ChildPGID = child.PGID
	}
	var res Result
	switch err := readJSON(filepath.Join(dir, "result.json"), &res); {
	case err == nil:
		v.Status, v.Result = Finished, &res
		return v, nil
	case !errors.Is(err, os.ErrNotExist):
		return nil, err
	}
	switch {
	case locked(filepath.Join(dir, "owner.lock")):
		v.Status = Running
	case v.ChildPGID != 0 && syscall.Kill(-v.ChildPGID, 0) == nil:
		v.Status = Orphaned // ponytail: a reused pgid would read as alive; pids are not trusted for anything but this
	default:
		v.Status = Interrupted
	}
	ev, err := ReadEvents(v.EventsPath)
	if err != nil {
		return nil, err
	}
	v.Events = &ev
	return v, nil
}

// locked reports whether another process holds the flock on path.
func locked(path string) bool {
	f, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		return false
	}
	defer f.Close()
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		return true
	}
	syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	return false
}

// Stop ends a running attempt through its owner, or an orphaned one
// directly, and reports each fact. A finished or interrupted attempt is
// refused: there is nothing to stop.
func Stop(root, attempt string, report func(string)) error {
	v, err := Show(root, attempt)
	if err != nil {
		return err
	}
	switch v.Status {
	case Running:
		if v.Launch.Owner <= 0 {
			return fmt.Errorf("attempt %s is starting and its owner is not recorded yet; try again", attempt)
		}
		if err := syscall.Kill(v.Launch.Owner, syscall.SIGTERM); err != nil {
			return fmt.Errorf("the owner (pid %d) could not be signalled: %v", v.Launch.Owner, err)
		}
		report(fmt.Sprintf("stop: SIGTERM sent to owner %d of %s; it ends the turn with SIGINT, kills the provider after %s, and writes the result", v.Launch.Owner, attempt, stopGrace))
		return nil
	case Orphaned:
		syscall.Kill(-v.ChildPGID, syscall.SIGINT)
		deadline := time.Now().Add(stopGrace)
		for syscall.Kill(-v.ChildPGID, 0) == nil && time.Now().Before(deadline) {
			time.Sleep(200 * time.Millisecond)
		}
		if syscall.Kill(-v.ChildPGID, 0) == nil {
			syscall.Kill(-v.ChildPGID, syscall.SIGKILL)
			for syscall.Kill(-v.ChildPGID, 0) == nil {
				time.Sleep(50 * time.Millisecond)
			}
		}
		res := Result{Finished: time.Now().UTC(), ExitCode: -1, Signal: "owner lost; stopped by grove stop", Stopped: true, ReconciledBy: "stop"}
		reconcile(v.Dir, &v.Launch, &res, func(format string, args ...any) { report("stop: " + fmt.Sprintf(format, args...)) })
		report(fmt.Sprintf("stop: %s was orphaned (owner %d gone); its process group %d is ended and the result written", attempt, v.Launch.Owner, v.ChildPGID))
		return nil
	default:
		return fmt.Errorf("attempt %s is %s; nothing to stop", attempt, v.Status)
	}
}

func exitOf(err error) (int, string) {
	if err == nil {
		return 0, ""
	}
	var exit *exec.ExitError
	if errors.As(err, &exit) {
		if ws, ok := exit.Sys().(syscall.WaitStatus); ok && ws.Signaled() {
			return -1, ws.Signal().String()
		}
		return exit.ExitCode(), ""
	}
	return -1, err.Error()
}

// environ drops the named variables; a name ending in * drops a prefix, and
// a name starting with ! keeps that variable regardless.
func environ(env []string, drop ...string) []string {
	kept := make([]string, 0, len(env))
	for _, entry := range env {
		name, _, _ := strings.Cut(entry, "=")
		if !slices.Contains(drop, "!"+name) && slices.ContainsFunc(drop, func(d string) bool {
			return name == d || strings.HasSuffix(d, "*") && strings.HasPrefix(name, strings.TrimSuffix(d, "*"))
		}) {
			continue
		}
		kept = append(kept, entry)
	}
	return kept
}

func uuid() (string, error) {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "", err
	}
	b[6] = b[6]&0x0f | 0x40
	b[8] = b[8]&0x3f | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16]), nil
}

func groveVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "grove (no build information)"
	}
	line := "grove " + info.Main.Version
	for _, s := range info.Settings {
		if s.Key == "vcs.revision" {
			line += " " + s.Value
		}
	}
	return line
}

func writeJSON(path string, v any) error {
	data, err := json.MarshalIndent(v, "", " ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, append(data, '\n'), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func readJSON(path string, v any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func samePath(a, b string) bool {
	if a == b {
		return true
	}
	ra, err1 := filepath.EvalSymlinks(a)
	rb, err2 := filepath.EvalSymlinks(b)
	return err1 == nil && err2 == nil && ra == rb
}

func orDetached(branch string) string {
	if branch == "" {
		return "a detached HEAD"
	}
	return strings.TrimPrefix(branch, "refs/heads/")
}

func short(commit string) string { return commit[:min(len(commit), 7)] }

func sid() int {
	id, _ := syscall.Getsid(0)
	return id
}
