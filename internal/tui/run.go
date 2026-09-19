package tui

import (
	"context"
	"os"

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
	m := New(ctx, root, Backend{Inspect: versions.InspectContext, Resolve: versions.ResolveContext})
	_, err := tea.NewProgram(m, tea.WithContext(ctx), tea.WithInput(input), tea.WithOutput(screen)).Run()
	// Quitting does not stop a command that is still reading.
	cancel()
	m.reads.close()
	switch {
	case interrupted(err):
		return nil, context.Canceled
	case err != nil:
		return nil, err
	}
	return m.Workspace, nil
}
