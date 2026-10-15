package cli

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/mascah/grove-next2/internal/doc"
	"github.com/mascah/grove-next2/internal/gitx"
	"github.com/mascah/grove-next2/internal/instruct"
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
	fmt.Fprint(stdout, instruct.Explore(ip).Markdown())
	return nil
}

// instructionsPropose prints the propose instructions: for a new change
// when arg is empty or a description, or for revising the drafted change
// whose ID arg is.
func instructionsPropose(stdout io.Writer, arg string) error {
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
		fmt.Fprint(stdout, instruct.Propose(ip, doc.NewID(now), now.UTC().Format(time.RFC3339), arg).Markdown())
		return nil
	}
	fsys := os.DirFS(filepath.Join(p.Top, p.Root))
	c, files, err := change(fsys, arg)
	if err != nil {
		return err
	}
	fmt.Fprint(stdout, instruct.Revise(ip, c, files).Markdown())
	return nil
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
	c, files, err := change(os.DirFS(filepath.Join(p.Top, p.Root)), id)
	if err != nil {
		return err
	}
	fmt.Fprint(stdout, instruct.Archive(ip, c, files).Markdown())
	return nil
}

// archive moves change id's folder in the working tree to changes/archive/.
func archive(stdout io.Writer, id string) error {
	p, err := openProject()
	if err != nil {
		return err
	}
	root := filepath.Join(p.Top, p.Root)
	c, _, err := change(os.DirFS(root), id)
	if err != nil {
		return err
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

// instructionsImplement prints the implement instructions for change id as
// the target holds it.
func instructionsImplement(stdout io.Writer, id string) error {
	p, err := openProject()
	if err != nil {
		return err
	}
	ip, err := p.agreed()
	if err != nil {
		return err
	}
	if ip.Agreed.Change(id) == nil {
		return fmt.Errorf("%s holds no approved change %s", p.Target, id)
	}
	fsys, err := gitx.TreeFS(p.Top, ip.Commit, p.Root)
	if err != nil {
		return err
	}
	c, files, err := change(fsys, id)
	if err != nil {
		return err
	}
	fmt.Fprint(stdout, instruct.Implement(ip, c, files).Markdown())
	return nil
}

// change loads change id from the document root fsys, with every file it
// holds.
func change(fsys fs.FS, id string) (*doc.Change, map[string]string, error) {
	all, err := doc.Load(fsys)
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
