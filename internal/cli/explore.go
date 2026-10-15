package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/mascah/grove-next2/internal/doc"
	"github.com/mascah/grove-next2/internal/gitx"
	"github.com/mascah/grove-next2/internal/instruct"
	"github.com/mascah/grove-next2/internal/session"
)

// explore creates a workspace from the target, saves a new exploration of
// description in it and opens a session there.
func explore(stdout io.Writer, description string) error {
	p, err := openProject()
	if err != nil {
		return err
	}
	ip, err := p.agreed()
	if err != nil {
		return err
	}
	now := time.Now()
	id := doc.NewID(now)
	name := id + "-" + slug(description)
	ip.Workspace, ip.Branch = filepath.Join(p.worktrees(), name), "grove/"+name
	if err := gitx.AddWorktree(p.Main, ip.Workspace, ip.Branch, ip.Commit); err != nil {
		return err
	}
	exploration := &doc.Document{
		Path: "explorations/" + name + ".md",
		Meta: doc.Meta{ID: id, Title: description, Created: now.UTC().Format(time.RFC3339)},
	}
	exploration.Source = doc.Format(exploration.Meta, "\n"+instruct.Exploration(description, description, ip.Commit[:7]))
	if err := create(filepath.Join(ip.Workspace, p.Root, exploration.Path), exploration.Source); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "Workspace %s on branch %s, from %s@%s.\n", shown(ip.Workspace), ip.Branch, p.Target, ip.Commit[:7])
	err = session.Start(ip.Workspace, instruct.Explore(ip, exploration))
	if rerr := report(stdout, ip.Workspace, id); err == nil {
		err = rerr
	}
	return err
}

// report is the two lines printed when a session ends: where the draft is
// and how to continue.
func report(stdout io.Writer, workspace, id string) error {
	n, err := gitx.Uncommitted(workspace)
	if err != nil {
		return err
	}
	exe, _ := os.Executable()
	fmt.Fprintf(stdout, "Draft: exploration %s in %s, %s uncommitted.\n", id, shown(workspace), files(n))
	fmt.Fprintf(stdout, "Continue: %s explore %s --resume\n", exe, id)
	return nil
}
