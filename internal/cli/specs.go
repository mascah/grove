package cli

import (
	"fmt"
	"io"
	"io/fs"
	"path"
	"slices"
	"strings"

	"github.com/mascah/grove-next2/internal/doc"
	"github.com/mascah/grove-next2/internal/gitx"
	"github.com/mascah/grove-next2/internal/tui"
)

func listSpecs(stdout io.Writer) error {
	p, err := openProject()
	if err != nil {
		return err
	}
	ip, err := p.agreed()
	if err != nil {
		return err
	}
	for _, s := range ip.Agreed.Specs {
		fmt.Fprintf(stdout, "%s  %s  %s/%s/spec.md\n", s.ID, s.Meta.Title, p.Root, s.Dir)
		problems(stdout, s.Diagnostics)
	}
	return nil
}

// reader opens the reader on the target's maintained specs and approved
// changes, each change with all of its files, implementation.md as written.
func reader() error {
	p, err := openProject()
	if err != nil {
		return err
	}
	ip, err := p.agreed()
	if err != nil {
		return err
	}
	fsys, err := gitx.TreeFS(p.Top, ip.Commit, p.Root)
	if err != nil {
		return err
	}
	var items []tui.Item
	for _, s := range ip.Agreed.Specs {
		b, err := fs.ReadFile(fsys, s.Dir+"/spec.md")
		if err != nil {
			return err
		}
		items = append(items, tui.Item{Label: "maintained spec", ID: s.ID, Title: s.Meta.Title, Sections: []string{doc.Body(string(b))}})
	}
	for _, c := range ip.Agreed.Changes {
		it := tui.Item{Label: "approved change", ID: c.ID, Title: c.Meta.Title}
		for _, name := range readingOrder(c) {
			b, err := fs.ReadFile(fsys, name)
			if err != nil {
				return err
			}
			it.Sections = append(it.Sections, "`"+p.Root+"/"+name+"`\n\n"+doc.Body(string(b)))
		}
		items = append(items, it)
	}
	return tui.Run(fmt.Sprintf("%s@%s: maintained specs and approved changes", p.Target, ip.Commit[:7]), items)
}

// readingOrder is the proposal, the change specs, the tasks, then whatever
// else the folder holds, implementation.md among it. It leaves out
// .grove.yaml, whose title the reader already shows.
func readingOrder(c *doc.Change) []string {
	rank := func(name string) int {
		switch rest := strings.TrimPrefix(name, c.Dir+"/"); {
		case rest == "proposal.md":
			return 0
		case strings.HasPrefix(rest, "specs/"):
			return 1
		case rest == "tasks.md":
			return 2
		}
		return 3
	}
	files := slices.DeleteFunc(slices.Clone(c.Files), func(name string) bool { return path.Base(name) == ".grove.yaml" })
	slices.SortStableFunc(files, func(a, b string) int { return rank(a) - rank(b) })
	return files
}
