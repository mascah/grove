package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mascah/grove-next2/internal/doc"
	"github.com/mascah/grove-next2/internal/gitx"
)

// check prints every problem in the Grove documents of this working tree.
func check(stdout io.Writer) error {
	p, err := openProject()
	if err != nil {
		return err
	}
	top, err := gitx.TopLevel(".")
	if err != nil {
		return err
	}
	c, err := doc.Load(os.DirFS(filepath.Join(top, p.Root)))
	if err != nil {
		return err
	}
	docs := append([]*doc.Document{}, c.Specs...)
	for _, ch := range c.Changes {
		docs = append(docs, ch.Proposal)
	}
	n := 0
	for _, d := range docs {
		for _, diag := range d.Diagnostics {
			fmt.Fprintf(stdout, "%s/%s\n", p.Root, diag)
			n++
		}
	}
	if n > 0 {
		return fmt.Errorf("%s found", plural(n, "problem"))
	}
	fmt.Fprintf(stdout, "No problems in %s/.\n", p.Root)
	return nil
}

func problems(stdout io.Writer, ds []doc.Diagnostic) {
	for _, d := range ds {
		fmt.Fprintf(stdout, "  problem: %s\n", d)
	}
}
