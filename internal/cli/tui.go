package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/charmbracelet/x/term"

	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/tui"
)

// runBoard opens the terminal board, which is what grove does without a
// command. The board draws on stderr and reads stdin; stdout carries only the
// result, written like workspace's after the terminal has been restored.
// Leaving without a selection prints nothing and succeeds.
func runBoard(a invocation, cwd string, out, errOut io.Writer) int {
	// Refuse before anything touches terminal modes, and never fall back to
	// another output format that a script might mistake for a result.
	screen, ok := errOut.(*os.File)
	if !ok || !term.IsTerminal(screen.Fd()) || !term.IsTerminal(os.Stdin.Fd()) {
		fmt.Fprint(errOut, "grove: the board needs a terminal on stdin and stderr, so nothing was opened\n"+
			"Without one, use grove list, grove versions, or grove workspace; grove --help describes them.\n")
		return 1
	}
	// Only discovery happens here. The board reads the project itself, so a
	// large or invalid project shows up inside it instead of delaying it.
	root, d := project.Discover(cwd, a.project)
	if d != nil {
		fmt.Fprintln(errOut, visible(d.String()))
		return 1
	}
	if _, err := fmt.Fprintf(errOut, "Project: %s\n", visible(root)); err != nil {
		return 1
	}
	w, err := tui.Run(context.Background(), root, os.Stdin, screen)
	switch {
	case errors.Is(err, context.Canceled):
		fmt.Fprintln(errOut, "grove: interrupted; no workspace was selected")
		return 1
	case err != nil:
		report(errOut, err)
		return 1
	case w == nil:
		return 0
	}
	return writeWorkspace(w, a.json, out, errOut)
}
