package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mascah/grove-next2/internal/gitx"
	"github.com/mascah/grove-next2/internal/instruct"
)

const config = `# Grove: where documents live and which branch holds what is agreed.
root: grove
target: main
`

// initProject writes the project's files at the top of the worktree. It
// creates each file exclusively, so it never overwrites anything.
func initProject(stdout io.Writer) error {
	top, err := gitx.TopLevel(".")
	if err != nil {
		return err
	}
	if err := create(filepath.Join(top, "grove.yaml"), config); err != nil {
		return err
	}
	if err := create(filepath.Join(top, "grove", "brief.md"), instruct.Template("brief.md")); err != nil {
		return err
	}
	f, err := os.OpenFile(filepath.Join(top, ".gitignore"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	if _, err := f.WriteString(".grove/\n"); err != nil {
		f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	fmt.Fprintln(stdout, "Wrote grove.yaml, grove/brief.md and a .gitignore entry for .grove/; nothing is committed.")
	return nil
}

func create(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return err
	}
	if _, err := f.WriteString(content); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}
