// Package instruct holds Grove's embedded templates and guides and builds
// the instructions a session follows. Instructions are index first: the
// document being worked on and the brief are included in full, everything
// else is listed with the commit to read it at.
package instruct

import (
	"embed"
	"fmt"
	"maps"
	"slices"
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

// Spec returns a new maintained spec, frontmatter included.
func Spec(id, title string) string {
	return doc.Format(doc.Meta{ID: id, Title: title}, "\n"+strings.ReplaceAll(Template("spec.md"), "{{title}}", title))
}

// Project is what a session is told about the repository.
type Project struct {
	Root      string          // document root, such as "grove"
	Target    string          // target branch, such as "main"
	Commit    string          // the target's commit
	Agreed    *doc.Collection // the documents at Commit
	Workspace string          // absolute path of the session's worktree
	Branch    string          // branch checked out there
	Exe       string          // the Grove binary, to run its commands
	Drafts    []*doc.Document // specs in the workspace that differ from the target's
}

// Explore returns the instructions for exploring exploration, a document of
// the workspace's own working tree.
func Explore(p Project, exploration *doc.Document) string {
	var b strings.Builder
	b.WriteString("# Grove instructions: explore\n\n")
	p.preamble(&b, fmt.Sprintf("exploration %s, %q", exploration.Meta.ID, exploration.Meta.Title))
	fmt.Fprintf(&b, "## How to explore\n\n%s\n", embedded("guides/explore.md"))
	b.WriteString("## Baseline spec template\n\nOnly when the owner asks for a baseline. Choose the capability's ID, lowercase kebab-case such as `record-editing`.\n\n")
	include(&b, p.Root+"/specs/<capability>/spec.md", Spec("<capability>", "<title>"))
	fmt.Fprintf(&b, "## Switching to propose\n\nWhen the owner asks for the proposal, run `%s instructions propose %s` and follow what it prints, in this conversation.\n\n", p.Exe, exploration.Meta.ID)
	p.index(&b)
	p.brief(&b)
	b.WriteString("## This exploration\n\n")
	include(&b, p.Root+"/"+exploration.Path, exploration.Source)
	return b.String()
}

// Change is the change a propose session works on.
type Change struct {
	ID, Dir string            // Dir is relative to the document root
	Created string            // RFC 3339, for a change not yet written
	Files   map[string]string // path relative to the root to content; empty for a new change
}

// Propose returns the instructions for proposing the change that follows
// from exploration.
func Propose(p Project, exploration *doc.Document, c Change) string {
	var b strings.Builder
	b.WriteString("# Grove instructions: propose\n\n")
	p.preamble(&b, fmt.Sprintf("change %s, proposed from exploration %s", c.ID, exploration.Meta.ID))
	fmt.Fprintf(&b, "## How to propose\n\n%s\n", embedded("guides/propose.md"))
	dir := p.Root + "/" + c.Dir
	if len(c.Files) == 0 {
		fmt.Fprintf(&b, "## Templates\n\nCreate the change in `%s/`. Choose each affected capability's ID, lowercase kebab-case such as `record-editing`, and replace the bracketed placeholders.\n\n", dir)
		r := strings.NewReplacer("{{change}}", c.ID, "{{base}}", short(p.Commit))
		meta := doc.Meta{ID: c.ID, Title: "<title>", Created: c.Created, Exploration: exploration.Meta.ID}
		include(&b, dir+"/proposal.md", doc.Format(meta, "\n"+r.Replace(strings.ReplaceAll(Template("proposal.md"), "{{title}}", "<title>"))))
		include(&b, dir+"/specs/<capability>/spec.md", r.Replace(strings.ReplaceAll(Template("change-spec.md"), "{{title}}", "Changes to <capability>")))
		include(&b, dir+"/tasks.md", r.Replace(Template("tasks.md")))
	}
	p.index(&b)
	p.brief(&b)
	b.WriteString("## The exploration this follows\n\n")
	include(&b, p.Root+"/"+exploration.Path, exploration.Source)
	if len(c.Files) > 0 {
		b.WriteString("## This change, as drafted so far\n\n")
		for _, path := range slices.Sorted(maps.Keys(c.Files)) {
			include(&b, p.Root+"/"+path, c.Files[path])
		}
	}
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
	fmt.Fprintf(b, "You are working on %s in the workspace `%s`, on branch `%s`, started from %s@%s. Write only in this workspace. Run Grove as `%s`.\n\n%s\n",
		subject, p.Workspace, p.Branch, p.Target, short(p.Commit), p.Exe, meanings)
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
	if len(p.Drafts) > 0 {
		b.WriteString("Drafts in this workspace (not agreed until merged):\n")
		for _, d := range p.Drafts {
			entry(b, p.Root+"/"+d.Path, d)
		}
		b.WriteString("\n")
	}
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
