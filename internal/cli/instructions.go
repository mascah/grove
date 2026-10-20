package cli

import (
	"cmp"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mascah/grove/internal/doc"
	"github.com/mascah/grove/internal/gitx"
	"github.com/mascah/grove/internal/instruct"
)

func instructionsExplore(stdout io.Writer) error {
	p, err := openProject()
	if err != nil {
		return err
	}
	ip, err := p.agreed()
	if err != nil {
		return err
	}
	return show(stdout)(instruct.Explore(ip))
}

// instructionsPropose prints the propose instructions: for a new change in
// change schema schema, or grove.yaml's when it is empty, when arg is empty
// or a description, or for revising the drafted change whose ID arg is.
func instructionsPropose(stdout io.Writer, schema, arg string) error {
	p, err := openProject()
	if err != nil {
		return err
	}
	ip, err := p.agreed()
	if err != nil {
		return err
	}
	if !doc.ValidID(arg) {
		now := time.Now()
		return show(stdout)(instruct.Propose(ip, cmp.Or(schema, p.Schema), doc.NewID(now), now.UTC().Format(time.RFC3339), arg))
	}
	if schema != "" {
		return errors.New("a change keeps its schema; to change it, edit schema in the change's .grove.yaml")
	}
	c, files, err := change(os.DirFS(filepath.Join(p.Top, p.Root)), ip.Layout, arg)
	if err != nil {
		return err
	}
	return show(stdout)(instruct.Revise(ip, c, files))
}

// instructionsArchive prints the archive instructions for change id as the
// working tree holds it.
func instructionsArchive(stdout io.Writer, id string) error {
	p, err := openProject()
	if err != nil {
		return err
	}
	ip, err := p.agreed()
	if err != nil {
		return err
	}
	c, files, err := change(os.DirFS(filepath.Join(p.Top, p.Root)), ip.Layout, id)
	if err != nil {
		return err
	}
	return show(stdout)(instruct.Archive(ip, c, files))
}

// archive moves change id's folder in the working tree to changes/archive/.
func archive(stdout io.Writer, id string) error {
	p, err := openProject()
	if err != nil {
		return err
	}
	root := filepath.Join(p.Top, p.Root)
	c, _, err := change(os.DirFS(root), p.instruct().Layout, id)
	if err != nil {
		return err
	}
	for _, t := range c.Tasks {
		if !c.Dropped && t.Status != "done" {
			return fmt.Errorf("%s is not done in implementation.md; archive a change when every task is done or the file says Status: dropped", t.ID)
		}
	}
	to := "changes/archive/" + filepath.Base(c.Dir)
	if err := os.MkdirAll(filepath.Join(root, "changes", "archive"), 0o755); err != nil {
		return err
	}
	if err := os.Rename(filepath.Join(root, c.Dir), filepath.Join(root, to)); err != nil {
		return err
	}
	fmt.Fprintf(stdout, "Archived %s to %s/%s/; nothing is committed.\n", id, p.Root, to)
	return nil
}

// instructionsApply prints the apply instructions for change id as the
// target holds it, unless it waits on a change (D-56).
func instructionsApply(stdout io.Writer, id string) error {
	p, err := openProject()
	if err != nil {
		return err
	}
	ip, err := p.agreed()
	if err != nil {
		return err
	}
	switch c := ip.Agreed.Change(id); {
	case c == nil:
		return fmt.Errorf("%s holds no approved change %s", p.Target, id)
	case len(c.WaitsOn) > 0:
		return fmt.Errorf("%s waits on %s, which %s has not archived as finished", id, strings.Join(c.WaitsOn, " and "), p.Target)
	}
	fsys, err := gitx.TreeFS(p.Top, ip.Commit, p.Root)
	if err != nil {
		return err
	}
	c, files, err := change(fsys, ip.Layout, id)
	if err != nil {
		return err
	}
	return show(stdout)(instruct.Apply(ip, c, files))
}

// show returns a function printing instructions unless building them failed.
func show(stdout io.Writer) func(instruct.Instructions, error) error {
	return func(in instruct.Instructions, err error) error {
		if err == nil {
			fmt.Fprint(stdout, in.Markdown())
		}
		return err
	}
}

// change loads change id from the document root fsys, with every file it
// holds.
func change(fsys fs.FS, layout doc.Layout, id string) (*doc.Change, map[string]string, error) {
	all, err := doc.Load(fsys, layout)
	if err != nil {
		return nil, nil, err
	}
	c := all.Change(id)
	if c == nil {
		return nil, nil, fmt.Errorf("no change %s in this working tree", id)
	}
	files := map[string]string{}
	for _, name := range c.Files {
		b, err := fs.ReadFile(fsys, name)
		if err != nil {
			return nil, nil, err
		}
		files[name] = string(b)
	}
	return c, files, nil
}
