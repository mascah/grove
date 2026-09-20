package versions

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

// blockingGit puts a fake git first on PATH. It runs the real git unless the
// returned block file exists and names one of the arguments; then it reports
// its PID on the returned FIFO and sleeps under that same PID until killed.
// Build fixtures before calling it.
func blockingGit(t *testing.T) (block, fifo string) {
	t.Helper()
	real, err := exec.LookPath("git")
	if err != nil {
		t.Skip("git not installed")
	}
	dir := t.TempDir()
	block, fifo = filepath.Join(dir, "block"), filepath.Join(dir, "started")
	must(t, syscall.Mkfifo(fifo, 0o600))
	script := "#!/bin/sh\nif [ -e \"$GROVE_TEST_BLOCK\" ]; then\n\tcase \" $* \" in *\" $(cat \"$GROVE_TEST_BLOCK\") \"*)\n\t\techo $$ > \"$GROVE_TEST_FIFO\"\n\t\texec sleep 600\n\tesac\nfi\nexec \"$GROVE_TEST_GIT\" \"$@\"\n"
	must(t, os.MkdirAll(filepath.Join(dir, "bin"), 0o755))
	must(t, os.WriteFile(filepath.Join(dir, "bin", "git"), []byte(script), 0o755))
	t.Setenv("GROVE_TEST_GIT", real)
	t.Setenv("GROVE_TEST_BLOCK", block)
	t.Setenv("GROVE_TEST_FIFO", fifo)
	t.Setenv("PATH", filepath.Join(dir, "bin")+string(os.PathListSeparator)+os.Getenv("PATH"))
	return block, fifo
}

// cancelDuring runs call, cancels its context once the fake git has reported
// that it started, and requires a prompt context.Canceled with no value and a
// collected git process. The timeouts only detect failure.
func cancelDuring(t *testing.T, fifo string, call func(context.Context) (returned bool, err error)) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	type outcome struct {
		returned bool
		err      error
	}
	done := make(chan outcome, 1)
	go func() {
		returned, err := call(ctx)
		done <- outcome{returned, err}
	}()
	started := make(chan int, 1)
	go func() {
		data, _ := os.ReadFile(fifo) // blocks until the fake git opens the FIFO
		pid, _ := strconv.Atoi(strings.TrimSpace(string(data)))
		started <- pid
	}()
	var pid int
	select {
	case pid = <-started:
	case o := <-done:
		t.Fatalf("returned before any git started: %+v", o)
	case <-time.After(10 * time.Second):
		t.Fatal("no git process started")
	}
	if pid <= 0 {
		t.Fatal("the fake git reported no PID")
	}
	cancel()
	select {
	case o := <-done:
		if !errors.Is(o.err, context.Canceled) || o.returned {
			t.Errorf("expected context.Canceled and no value, got value=%v err=%v", o.returned, o.err)
		}
	case <-time.After(10 * time.Second):
		t.Error("cancellation did not end the call")
	}
	if err := syscall.Kill(pid, 0); !errors.Is(err, syscall.ESRCH) {
		syscall.Kill(pid, syscall.SIGKILL)
		t.Fatalf("git process %d was not killed and collected: %v", pid, err)
	}
}

func TestCancellationKillsGitAndWritesNothing(t *testing.T) {
	root, _ := nestedFixture(t)
	project := filepath.Join(root, "sub")
	live := selectorFor(t, project, "W-001", "live", "feature")
	committed := selectorFor(t, project, "W-001", "committed", "refs/heads/feature")
	block, fifo := blockingGit(t)
	blockAt := func(arg string) func() {
		return func() { must(t, os.WriteFile(block, []byte(arg), 0o644)) }
	}
	before := treeHashes(t, filepath.Dir(root))

	cases := []struct {
		name string
		call func(ctx context.Context) (bool, error)
	}{
		{"inspect locating the repository", func(ctx context.Context) (bool, error) {
			blockAt("--git-common-dir")()
			res, err := InspectContext(ctx, project, "")
			return res != nil, err
		}},
		// The next two fail a source, not the inspection, when Git fails.
		{"inspect reading a committed tree", func(ctx context.Context) (bool, error) {
			blockAt("cat-file")()
			res, err := InspectContext(ctx, project, "")
			return res != nil, err
		}},
		{"inspect entering a worktree", func(ctx context.Context) (bool, error) {
			blockAt("--git-dir")()
			res, err := InspectContext(ctx, project, "")
			return res != nil, err
		}},
		{"inspect at the second inventory", func(ctx context.Context) (bool, error) {
			res, err := inspect(ctx, project, "", blockAt("worktree"))
			return res != nil, err
		}},
		// Nothing after this re-entry would fail the inspection by itself.
		{"inspect re-entering a worktree after the second inventory", func(ctx context.Context) (bool, error) {
			res, err := inspect(ctx, project, "", blockAt("--git-dir"))
			return res != nil, err
		}},
		{"resolve at the start", func(ctx context.Context) (bool, error) {
			blockAt("--git-common-dir")()
			w, err := ResolveContext(ctx, project, live)
			return w != nil, err
		}},
		{"resolve live at the final check", func(ctx context.Context) (bool, error) {
			w, err := resolveWith(ctx, project, live, blockAt("worktree"))
			return w != nil, err
		}},
		// A failed tip check reads as a moved branch unless cancellation wins.
		{"resolve committed at the final tip check", func(ctx context.Context) (bool, error) {
			w, err := resolveWith(ctx, project, committed, blockAt("--verify"))
			return w != nil, err
		}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			defer os.Remove(block)
			cancelDuring(t, fifo, c.call)
		})
	}
	if after := treeHashes(t, filepath.Dir(root)); !reflect.DeepEqual(before, after) {
		t.Fatal("cancelled reads changed files or Git state")
	}
	// The same calls still succeed once nothing blocks or cancels them.
	if w, err := ResolveContext(context.Background(), project, committed); err != nil || w == nil {
		t.Fatalf("resolution after cancelled attempts: %+v %v", w, err)
	}
}

func TestCancelledBeforeStart(t *testing.T) {
	root, _ := nestedFixture(t)
	project := filepath.Join(root, "sub")
	live := selectorFor(t, project, "W-001", "live", "feature")
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if res, err := InspectContext(ctx, project, ""); !errors.Is(err, context.Canceled) || res != nil {
		t.Fatalf("inspect: %+v %v", res, err)
	}
	if w, err := ResolveContext(ctx, project, live); !errors.Is(err, context.Canceled) || w != nil {
		t.Fatalf("resolve: %+v %v", w, err)
	}
}
