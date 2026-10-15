package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

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

// change is the change proposed from the workspace's exploration, or nil.
func (ws workspace) change() *doc.Change {
	for _, c := range ws.Docs.Changes {
		if c.Proposal.Meta.Exploration == ws.ID {
			return c
		}
	}
	return nil
}

// drafts are the workspace's specs that the target does not hold as they are.
func (ws workspace) drafts(agreed *doc.Collection) []*doc.Document {
	var out []*doc.Document
	for _, d := range ws.Docs.Specs {
		if !slices.ContainsFunc(agreed.Specs, func(a *doc.Document) bool { return a.Path == d.Path && a.Source == d.Source }) {
			out = append(out, d)
		}
	}
	return out
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
	ip, err := p.agreed()
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
		what := "no change yet"
		if c := ws.change(); c != nil {
			what = "change " + c.Proposal.Meta.ID + " draft"
			if ip.Agreed.Change(c.Proposal.Meta.ID) != nil {
				what = "change " + c.Proposal.Meta.ID + " approved"
			}
		}
		fmt.Fprintf(stdout, "%s  %s  %s, %s uncommitted\n", ws.ID, title, what, files(n))
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
		ip.Workspace, ip.Branch, ip.Drafts = ws.Path, ws.Branch, ws.drafts(ip.Agreed)
		err = session.Start(ws.Path, instruct.Explore(ip, ws.Exploration))
	}
	if rerr := report(stdout, p, id); err == nil {
		err = rerr
	}
	return err
}

// instructionsPropose prints the propose instructions for the change that
// follows from exploration id, a new one unless its workspace has one.
func instructionsPropose(stdout io.Writer, id string) error {
	p, err := openProject()
	if err != nil {
		return err
	}
	ip, err := p.agreed()
	if err != nil {
		return err
	}
	ws, err := p.workspace(id)
	if err != nil {
		return err
	}
	ip.Workspace, ip.Branch, ip.Drafts = ws.Path, ws.Branch, ws.drafts(ip.Agreed)
	var c instruct.Change
	if ch := ws.change(); ch != nil {
		c = instruct.Change{ID: ch.Proposal.Meta.ID, Dir: ch.Dir, Files: map[string]string{}}
		for _, name := range ch.Files {
			b, err := os.ReadFile(filepath.Join(ws.Path, p.Root, name))
			if err != nil {
				return err
			}
			c.Files[name] = string(b)
		}
	} else {
		now := time.Now()
		c.ID = doc.NewID(now)
		c.Dir, c.Created = "changes/"+c.ID+"-"+slug(ws.Exploration.Meta.Title), now.UTC().Format(time.RFC3339)
	}
	fmt.Fprint(stdout, instruct.Propose(ip, ws.Exploration, c))
	return nil
}
