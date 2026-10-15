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

	"gopkg.in/yaml.v3"

	"github.com/mascah/grove-next2/internal/doc"
)

//go:embed schema.yaml templates guides
var files embed.FS

func embedded(name string) string {
	b, err := files.ReadFile(name)
	if err != nil {
		panic(err) // embedded at build time: a missing name is a programming error
	}
	return string(b)
}

// Kind is one kind of document: where it goes, the template it starts from
// and the guide for writing it. Path is relative to the document root.
type Kind struct{ Kind, Path, Template, Guide string }

// Workflow is one workflow: its guide and the kinds of document it writes.
type Workflow struct {
	Name, Guide string
	Writes      []string
}

// schema is the embedded document set and workflows, schema.yaml.
var schema struct {
	Documents []Kind
	Workflows []Workflow
}

func init() {
	if err := yaml.Unmarshal([]byte(embedded("schema.yaml")), &schema); err != nil {
		panic(err) // embedded at build time
	}
}

func workflow(name string) Workflow {
	i := slices.IndexFunc(schema.Workflows, func(w Workflow) bool { return w.Name == name })
	return schema.Workflows[i]
}

func kind(name string) Kind {
	i := slices.IndexFunc(schema.Documents, func(k Kind) bool { return k.Kind == name })
	return schema.Documents[i]
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
	w := workflow("implement")
	fmt.Fprintf(&b, "## How to implement\n\n%s\n", embedded(w.Guide))
	p.documents(&b, w, strings.NewReplacer("<change>", path.Base(c.Dir), "{{title}}", c.Proposal.Meta.Title, "{{change}}", c.Proposal.Meta.ID, "{{base}}", short(p.Commit)))
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

// documents gives each kind w writes: its path, guide and template, with
// fill's values put in.
func (p Project) documents(b *strings.Builder, w Workflow, fill *strings.Replacer) {
	b.WriteString("## Documents\n\n")
	for _, name := range w.Writes {
		k := kind(name)
		path := p.Root + "/" + fill.Replace(k.Path)
		fmt.Fprintf(b, "### %s: `%s`\n\n%s\n", k.Kind, path, embedded(k.Guide))
		include(b, path, fill.Replace(embedded(k.Template)))
	}
}

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
