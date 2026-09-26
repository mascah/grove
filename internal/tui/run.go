package tui

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	tea "charm.land/bubbletea/v2"

	"github.com/mascah/grove/internal/attempt"
	"github.com/mascah/grove/internal/deps"
	"github.com/mascah/grove/internal/integrate"
	"github.com/mascah/grove/internal/update"
	"github.com/mascah/grove/internal/versions"
)

// Run shows the board on screen, reading keys from input, until the person
// selects a workspace or leaves. A nil workspace with a nil error is ordinary
// cancellation; context.Canceled is an interrupt. The terminal is restored and
// every read collected before Run returns, so the caller may then write its
// result. Nothing is written to a file except by the three review actions,
// the feedback before a resolution attempt, and the resolve of an answered
// question, each confirmed at a prompt, and by the owner's editor with the
// Answer heading it is handed.
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
	b := Live()
	liveAnswer(&b, input, screen)
	m := New(ctx, root, b)
	out := &watched{File: screen, stop: cancel}
	_, err := tea.NewProgram(m, tea.WithContext(ctx), tea.WithInput(input), tea.WithOutput(out)).Run()
	// Quitting does not stop a command that is still reading.
	cancel()
	m.reads.close()
	// A session ended while the editor had a question, by a hangup, never
	// saw it exit: an unused Answer heading is taken back as it would be.
	if m.editing != nil {
		m.editing.takeBack()
	}
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
	b := Backend{
		Inspect: versions.InspectContext, Resolve: versions.ResolveContext, History: versions.HistoryContext,
		Changes: versions.ChangesContext, Diff: versions.DiffContext, Ancestry: deps.Ancestry, Predict: versions.PredictContext,
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
			facts := []string{fmt.Sprintf("feedback: %s is active again on branch %s in %s, commit %s", id, branch, root, res.Commit[:min(len(res.Commit), 7)])}
			ids := []string{id}
			for _, o := range res.Reopened {
				facts = append(facts, fmt.Sprintf("reopened: %s shared the candidate and is active again, commit %s", o.ID, o.Commit[:min(len(o.Commit), 7)]))
				ids = append(ids, o.ID)
			}
			return append(facts, "next: continue there with /grove-work "+strings.Join(ids, " ")), nil
		},
		Integrate: func(_ context.Context, root, id string, cleanup bool) ([]string, error) {
			var facts []string
			cwd, _ := os.Getwd()
			err := integrate.Run(integrate.Request{Root: root, ID: id, Cwd: cwd, Cleanup: cleanup}, time.Now(), func(fact string) { facts = append(facts, fact) })
			return facts, err
		},
	}
	liveAttempts(&b)
	b.Conflict = func(_ context.Context, req attempt.Request, shown *versions.Merge) ([]string, error) {
		var facts []string
		l, err := attempt.Resolve(req, shown, time.Now(), func(f string) { facts = append(facts, f) })
		if err != nil {
			return facts, err
		}
		return append(facts, launched(l)...), nil
	}
	return b
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
