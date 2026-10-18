package cli

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
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
	return show(stdout)(instruct.Explore(ip))
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
		return show(stdout)(instruct.Propose(ip, doc.NewID(now), now.UTC().Format(time.RFC3339), arg))
	}
	fsys := os.DirFS(filepath.Join(p.Top, p.Root))
	c, files, err := change(fsys, arg)
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
	c, files, err := change(os.DirFS(filepath.Join(p.Top, p.Root)), id)
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
	c, _, err := change(os.DirFS(root), id)
	if err != nil {
		return err
	}
	if t := unfinished(os.DirFS(root), c); t != "" {
		return fmt.Errorf("%s is not done in implementation.md; archive a change when every task is done or the file says Status: dropped", t)
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

var taskHeading = regexp.MustCompile(`^#+ (T-\d+)\b`)

// unfinished returns the first task in change c's tasks.md that its
// implementation.md does not mark done, or "" when every task is done or
// the change is dropped (rule 3).
func unfinished(fsys fs.FS, c *doc.Change) string {
	tasks, _ := fs.ReadFile(fsys, path.Join(c.Dir, "tasks.md"))
	impl, _ := fs.ReadFile(fsys, path.Join(c.Dir, "implementation.md"))
	status, task := map[string]string{}, ""
	for _, l := range strings.Split(string(impl), "\n") {
		if m := taskHeading.FindStringSubmatch(l); m != nil {
			task = m[1]
		}
		if s, ok := strings.CutPrefix(l, "Status: "); ok {
			if s = strings.TrimSpace(s); s == "dropped" {
				return ""
			}
			status[task] = strings.TrimSpace(s)
		}
	}
	for _, l := range strings.Split(string(tasks), "\n") {
		if m := taskHeading.FindStringSubmatch(l); m != nil && status[m[1]] != "done" {
			return m[1]
		}
	}
	return ""
}

// instructionsApply prints the apply instructions for change id as the
// target holds it.
func instructionsApply(stdout io.Writer, id string) error {
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
func change(fsys fs.FS, id string) (*doc.Change, map[string]string, error) {
	all, err := doc.Load(fsys, instruct.MainDocument())
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
