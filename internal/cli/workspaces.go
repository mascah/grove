package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mascah/grove-next2/internal/doc"
	"github.com/mascah/grove-next2/internal/gitx"
	"github.com/mascah/grove-next2/internal/instruct"
	"github.com/mascah/grove-next2/internal/session"
)

// workspace is one worktree Grove created, on branch grove/<id>-<slug>.
type workspace struct {
	ID, Path, Branch string
	Docs             *doc.Collection // the worktree's documents, uncommitted edits included
	Exploration      *doc.Document
}

func (p project) workspaces() ([]workspace, error) {
	worktrees, err := gitx.Worktrees(p.Main)
	if err != nil {
		return nil, err
	}
	var out []workspace
	for _, w := range worktrees {
		name, ok := strings.CutPrefix(w.Branch, "grove/")
		if !ok || len(name) < 12 {
			continue
		}
		ws := workspace{ID: name[:12], Path: w.Path, Branch: w.Branch}
		if ws.Docs, err = doc.Load(os.DirFS(filepath.Join(w.Path, p.Root))); err != nil {
			return nil, err
		}
		for _, e := range ws.Docs.Explorations {
			if e.Meta.ID == ws.ID {
				ws.Exploration = e
			}
		}
		out = append(out, ws)
	}
	return out, nil
}

func (p project) workspace(id string) (workspace, error) {
	all, err := p.workspaces()
	for _, ws := range all {
		if ws.ID == id {
			return ws, err
		}
	}
	if err == nil {
		err = fmt.Errorf("no workspace %s; workspaces lists them", id)
	}
	return workspace{}, err
}

func listWorkspaces(stdout io.Writer) error {
	p, err := openProject()
	if err != nil {
		return err
	}
	all, err := p.workspaces()
	if err != nil {
		return err
	}
	for _, ws := range all {
		n, err := gitx.Uncommitted(ws.Path)
		if err != nil {
			return err
		}
		title := "(no exploration)"
		if ws.Exploration != nil {
			title = ws.Exploration.Meta.Title
		}
		fmt.Fprintf(stdout, "%s  %s  %s, %s uncommitted\n", ws.ID, title, "no change yet", files(n))
	}
	return nil
}

// reopen starts a new conversation on a workspace's exploration, or resumes
// the last one.
func reopen(stdout io.Writer, id string, resume bool) error {
	p, err := openProject()
	if err != nil {
		return err
	}
	ws, err := p.workspace(id)
	if err != nil {
		return err
	}
	if resume {
		err = session.Resume(ws.Path)
	} else {
		ip, aerr := p.agreed()
		if aerr != nil {
			return aerr
		}
		ip.Workspace, ip.Branch = ws.Path, ws.Branch
		err = session.Start(ws.Path, instruct.Explore(ip, ws.Exploration))
	}
	if rerr := report(stdout, ws.Path, id); err == nil {
		err = rerr
	}
	return err
}
