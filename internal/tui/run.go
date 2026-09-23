package tui

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/mascah/grove/internal/integrate"
	"github.com/mascah/grove/internal/update"
	"github.com/mascah/grove/internal/versions"
)

// Run shows the board on screen, reading keys from input, until the person
// selects a workspace or leaves. A nil workspace with a nil error is ordinary
// cancellation; context.Canceled is an interrupt. The terminal is restored and
// every read collected before Run returns, so the caller may then write its
// result. Nothing is written to a file except by the three review actions,
// each confirmed at a prompt.
func Run(ctx context.Context, root string, input, screen *os.File) (*versions.Workspace, error) {
	// The framework reads these from the process environment, not from the
	// environment a program is given, and each one makes it write a log file.
	for _, name := range []string{"TEA_DEBUG", "TEA_TRACE", "UV_DEBUG"} {
		os.Unsetenv(name)
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	// The framework handles SIGINT and SIGTERM. A hangup (the window closed,
	// the connection dropped) would otherwise end the process with a Git
	// child still running and the terminal in raw mode.
	hangup := make(chan os.Signal, 1)
	signal.Notify(hangup, syscall.SIGHUP)
	defer signal.Stop(hangup)
	go func() {
		select {
		case <-hangup:
			cancel()
		case <-ctx.Done():
		}
	}()
	m := New(ctx, root, Live())
	out := &watched{File: screen, stop: cancel}
	_, err := tea.NewProgram(m, tea.WithContext(ctx), tea.WithInput(input), tea.WithOutput(out)).Run()
	// Quitting does not stop a command that is still reading.
	cancel()
	m.reads.close()
	switch {
	case out.err != nil:
		return nil, fmt.Errorf("the terminal stopped accepting output, so nothing was selected: %w", out.err)
	case interrupted(err):
		return nil, context.Canceled
	case err != nil:
		return nil, err
	}
	return m.Workspace, nil
}

// Live is the backend over the real repository: the versions reads, and the
// three actions through the same functions the CLI runs, each in the
// checkout the review view named. An action ignores ctx for its Git
// commands: a key never cancels a write half-way.
func Live() Backend {
	return Backend{
		Inspect: versions.InspectContext, Resolve: versions.ResolveContext, History: versions.HistoryContext,
		Changes: versions.ChangesContext, Diff: versions.DiffContext,
		Approve: func(_ context.Context, root, id, verdict string) ([]string, error) {
			res, err := update.Approve(root, id, verdict, time.Now())
			if err != nil {
				return nil, err
			}
			return []string{fmt.Sprintf("approved: %s's candidate, in commit %s of %s", id, res.Commit[:min(len(res.Commit), 7)], root)}, nil
		},
		Feedback: func(_ context.Context, root, id, text string) ([]string, error) {
			res, err := update.Feedback(root, id, text, time.Now())
			if err != nil {
				return nil, err
			}
			branch, _ := update.Branch(root)
			return []string{
				fmt.Sprintf("feedback: %s is active again on branch %s in %s, commit %s", id, branch, root, res.Commit[:min(len(res.Commit), 7)]),
				fmt.Sprintf("next: continue there with /grove-work %s", id),
			}, nil
		},
		Integrate: func(_ context.Context, root, id string, cleanup bool) ([]string, error) {
			var facts []string
			cwd, _ := os.Getwd()
			err := integrate.Run(integrate.Request{Root: root, ID: id, Cwd: cwd, Cleanup: cleanup}, time.Now(), func(fact string) { facts = append(facts, fact) })
			return facts, err
		},
	}
}

// watched ends the session at the first failed write to the screen. The
// framework discards those errors, which would leave a person selecting a
// workspace on a screen they can no longer see.
type watched struct {
	*os.File
	stop func()
	once sync.Once
	err  error
}

func (w *watched) Write(p []byte) (int, error) {
	n, err := w.File.Write(p)
	if err != nil {
		w.once.Do(func() { w.err = err; w.stop() })
	}
	return n, err
}
