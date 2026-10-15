// Package instruct holds Grove's embedded templates and guides and builds
// the instructions a session follows. Instructions are index first: the
// document being worked on and the brief are included in full, everything
// else is listed with the commit to read it at.
package instruct

import (
	"embed"
	"fmt"
	"strings"

	"github.com/mascah/grove-next2/internal/doc"
)

//go:embed templates guides
var files embed.FS

// Template returns the embedded template called name, such as "brief.md".
func Template(name string) string { return embedded("templates/" + name) }

func embedded(name string) string {
	b, err := files.ReadFile(name)
	if err != nil {
		panic(err) // embedded at build time: a missing name is a programming error
	}
	return string(b)
}

// Exploration returns the body of a new exploration.
func Exploration(title, description, base string) string {
	return strings.NewReplacer("{{title}}", title, "{{description}}", description, "{{base}}", base).
		Replace(Template("exploration.md"))
}

// Project is what a session is told about the repository.
type Project struct {
	Root      string          // document root, such as "grove"
	Target    string          // target branch, such as "main"
	Commit    string          // the target's commit
	Agreed    *doc.Collection // the documents at Commit
	Workspace string          // absolute path of the session's worktree
	Branch    string          // branch checked out there
}

// Explore returns the instructions for exploring exploration, a document of
// the workspace's own working tree.
func Explore(p Project, exploration *doc.Document) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Grove instructions: explore\n\n")
	p.preamble(&b, fmt.Sprintf("exploration %s, %q", exploration.Meta.ID, exploration.Meta.Title))
	fmt.Fprintf(&b, "## How to explore\n\n%s\n", embedded("guides/explore.md"))
	p.index(&b)
	p.brief(&b)
	fmt.Fprintf(&b, "## This exploration\n\n")
	include(&b, p.Root+"/"+exploration.Path, exploration.Source)
	return b.String()
}

const meanings = `Grove keeps four meanings apart; keep them apart when you speak and write:

- A maintained spec (` + "`specs/<id>/spec.md`" + ` in the target branch) is agreed behavior.
- A change (` + "`changes/<id>-<slug>/`" + `) in the target branch is approved future work.
  Approval is the proposal being in the target branch; it never means
  implemented. Whoever implements it records progress in its
  ` + "`implementation.md`" + `.
- An exploration is saved findings, not agreement.
- What the code does is observed behavior: evidence, not a rule.

Grove commits nothing. The owner commits, and merging into the target branch
is the only approval.
`

func (p Project) preamble(b *strings.Builder, subject string) {
	fmt.Fprintf(b, "You are working on %s in the workspace `%s`, on branch `%s`, started from %s@%s. Write only in this workspace.\n\n%s\n",
		subject, p.Workspace, p.Branch, p.Target, short(p.Commit), meanings)
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
	b.WriteString("\nExplorations (findings, not agreement):\n")
	for _, d := range p.Agreed.Explorations {
		entry(b, p.Root+"/"+d.Path, d)
	}
	if len(p.Agreed.Explorations) == 0 {
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

func (p Project) brief(b *strings.Builder) {
	b.WriteString("## Brief\n\n")
	switch p.Agreed.Brief {
	case "", Template("brief.md"):
		b.WriteString("The project has no brief yet.\n\n")
	default:
		include(b, p.Root+"/brief.md", p.Agreed.Brief)
	}
}

func include(b *strings.Builder, path, content string) {
	fmt.Fprintf(b, "<file path=%q>\n%s\n</file>\n\n", path, strings.TrimRight(content, "\n"))
}

func short(commit string) string { return commit[:min(7, len(commit))] }
