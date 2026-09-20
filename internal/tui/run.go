package tui

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	tea "charm.land/bubbletea/v2"

	"github.com/mascah/grove/internal/versions"
)

// Run shows the board on screen, reading keys from input, until the person
// selects a workspace or leaves. A nil workspace with a nil error is ordinary
// cancellation; context.Canceled is an interrupt. The terminal is restored and
// every read collected before Run returns, so the caller may then write its
// result. Nothing is ever written to a file.
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
	m := New(ctx, root, Backend{Inspect: versions.InspectContext, Resolve: versions.ResolveContext})
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
