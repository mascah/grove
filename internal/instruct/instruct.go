// Package instruct holds Grove's embedded templates and guides and builds
// the instructions a session follows. Instructions are index first: the
// documents being worked on are included in full, everything else is listed
// with the commit to read it at.
package instruct

import (
	"embed"
	"fmt"
	"path"
	"slices"
	"strings"

	"github.com/mascah/grove-next2/internal/doc"
)

//go:embed templates guides
var files embed.FS

// Template returns the embedded template called name, such as "spec.md".
func Template(name string) string { return embedded("templates/" + name) }

func embedded(name string) string {
	b, err := files.ReadFile(name)
	if err != nil {
		panic(err) // embedded at build time: a missing name is a programming error
	}
	return string(b)
}

// Spec returns a new maintained spec, frontmatter included.
func Spec(id, title string) string {
	return doc.Format(doc.Meta{ID: id, Title: title}, "\n"+strings.ReplaceAll(Template("spec.md"), "{{title}}", title))
}

// Project is what a session is told about the repository.
type Project struct {
	Root   string          // document root, such as "grove"
	Target string          // target branch, such as "main"
	Commit string          // the target's commit
	Agreed *doc.Collection // the documents at Commit
	Exe    string          // the Grove binary, to run its commands
}

// Implement returns the instructions for implementing change c, which the
// target holds; c.Files are read from there.
func Implement(p Project, c *doc.Change, files map[string]string) string {
	var b strings.Builder
	b.WriteString("# Grove instructions: implement\n\n")
	fmt.Fprintf(&b, "You are implementing approved change %s, %q, as %s@%s holds it, in your own working tree. Run Grove as `%s`.\n\n%s\n",
		c.Proposal.Meta.ID, c.Proposal.Meta.Title, p.Target, short(p.Commit), p.Exe, meanings)
	fmt.Fprintf(&b, "## How to implement\n\n%s\n", embedded("guides/implement.md"))
	dir := p.Root + "/" + c.Dir
	b.WriteString("## Templates\n\n")
	include(&b, dir+"/implementation.md", strings.NewReplacer("{{title}}", c.Proposal.Meta.Title, "{{change}}", c.Proposal.Meta.ID).Replace(Template("implementation.md")))
	for _, name := range c.Files {
		capability, ok := strings.CutPrefix(path.Dir(name), c.Dir+"/specs/")
		if ok && !slices.ContainsFunc(p.Agreed.Specs, func(d *doc.Document) bool { return d.Meta.ID == capability }) {
			include(&b, p.Root+"/specs/"+capability+"/spec.md", Spec(capability, "<title>"))
		}
	}
	p.index(&b)
	b.WriteString("## The change\n\n")
	for _, name := range c.Files {
		include(&b, p.Root+"/"+name, files[name])
	}
	return b.String()
}

const meanings = `Grove keeps three meanings apart; keep them apart when you speak and write:

- A maintained spec (` + "`specs/<id>/spec.md`" + ` in the target branch) is agreed behavior.
- A change (` + "`changes/<id>-<slug>/`" + `) in the target branch is approved future work.
  Approval is the proposal being in the target branch; it never means
  implemented. Whoever implements it records progress in its
  ` + "`implementation.md`" + `.
- What the code does is observed behavior: evidence, not a rule.

Grove commits nothing. The owner commits, and merging into the target branch
is the only approval.
`

func (p Project) index(b *strings.Builder) {
	fmt.Fprintf(b, "## Project index\n\nWhat %s says at %s. Read a listed document with `git show %s:<path>`.\n\n",
		p.Target, short(p.Commit), short(p.Commit))
	b.WriteString("Maintained specs (agreed behavior):\n")
	for _, d := range p.Agreed.Specs {
		entry(b, p.Root+"/"+d.Path, d)
	}
	if len(p.Agreed.Specs) == 0 {
		b.WriteString("- None. Where nothing is agreed, what the code does is observed behavior.\n")
	}
	b.WriteString("\nApproved changes (future work; implementation is not established):\n")
	for _, c := range p.Agreed.Changes {
		entry(b, p.Root+"/"+c.Dir+"/", c.Proposal)
	}
	if len(p.Agreed.Changes) == 0 {
		b.WriteString("- None.\n")
	}
	b.WriteString("\n")
}

func entry(b *strings.Builder, path string, d *doc.Document) {
	fmt.Fprintf(b, "- %s %q: `%s`", d.Meta.ID, d.Meta.Title, path)
	if d.Meta.Description != "" {
		fmt.Fprintf(b, " — %s", d.Meta.Description)
	}
	b.WriteString("\n")
	for _, diag := range d.Diagnostics {
		fmt.Fprintf(b, "  - problem: %s\n", diag)
	}
}

func include(b *strings.Builder, path, content string) {
	fmt.Fprintf(b, "<file path=%q>\n%s\n</file>\n\n", path, strings.TrimRight(content, "\n"))
}

func short(commit string) string { return commit[:min(7, len(commit))] }
