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

// Workflow is one workflow: its skill, its guide and the kinds of document
// it writes.
type Workflow struct {
	Name, Skill, Args, Guide string
	Writes                   []string
}

// Workflows returns the schema's workflows, in order.
func Workflows() []Workflow { return schema.Workflows }

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
	Root    string              // document root, such as "grove"
	Target  string              // target branch, such as "main"
	Commit  string              // the target's commit
	Agreed  *doc.Collection     // the documents at Commit
	Exe     string              // the Grove binary, to run its commands
	Context string              // grove.yaml's context
	Rules   map[string][]string // grove.yaml's rules, by kind of document
}

// Instructions is everything one workflow's instructions say, built before
// it is printed (rule 10).
type Instructions struct {
	Project
	Workflow    string
	Task        string // what the session is doing, one sentence
	Guide       string
	Description string     // the owner's description, when given
	Documents   []Document // the kinds of document the workflow writes
	Change      []File     // the change being worked on, in full
}

// Document is one kind of document to write, with its path and template
// filled in.
type Document struct {
	Kind, Path, Guide, Template string
	Rules                       []string // the project's own, from grove.yaml
}

// File is one document included in full; Path is from the repository root.
type File struct{ Path, Content string }

// build returns workflow name's instructions with its guide and documents,
// fill's values put into their paths and templates.
func build(p Project, name, task string, fill *strings.Replacer) Instructions {
	w := workflow(name)
	in := Instructions{Project: p, Workflow: name, Task: task, Guide: embedded(w.Guide)}
	for _, kn := range w.Writes {
		k := kind(kn)
		in.Documents = append(in.Documents, Document{Kind: k.Kind, Path: p.Root + "/" + fill.Replace(k.Path),
			Guide: embedded(k.Guide), Template: fill.Replace(embedded(k.Template)), Rules: p.Rules[k.Kind]})
	}
	return in
}

// Explore returns the instructions for interviewing the owner about an idea.
func Explore(p Project) Instructions {
	return build(p, "explore", "You are exploring an idea with the owner, in their own working tree.",
		strings.NewReplacer("{{base}}", short(p.Commit)))
}

// Propose returns the instructions for writing new change id, created at
// created, from the conversation or from description when it is not empty.
func Propose(p Project, id, created, description string) Instructions {
	from := "this conversation"
	if description != "" {
		from = "the owner's description below"
	}
	in := build(p, "propose", fmt.Sprintf("You are proposing new change %s from %s, in the owner's own working tree.", id, from),
		strings.NewReplacer("<change>", id+"-<slug>", "{{change}}", id, "{{created}}", created, "{{base}}", short(p.Commit)))
	in.Description = description
	return in
}

// Revise returns the instructions for revising drafted change c, read from
// the working tree; files are c.Files as read from there.
func Revise(p Project, c *doc.Change, files map[string]string) Instructions {
	task := fmt.Sprintf("You are revising drafted change %s, %q, in the owner's own working tree.", c.Proposal.Meta.ID, c.Proposal.Meta.Title)
	return withChange(build(p, "propose", task, fill(p, c)), c, files)
}

// Implement returns the instructions for implementing change c, which the
// target holds; files are c.Files as read from there.
func Implement(p Project, c *doc.Change, files map[string]string) Instructions {
	task := fmt.Sprintf("You are implementing approved change %s, %q, as %s@%s holds it, in your own working tree.",
		c.Proposal.Meta.ID, c.Proposal.Meta.Title, p.Target, short(p.Commit))
	return withChange(build(p, "implement", task, fill(p, c)), c, files)
}

// fill puts change c's values into paths and templates.
func fill(p Project, c *doc.Change) *strings.Replacer {
	return strings.NewReplacer("<change>", path.Base(c.Dir), "{{title}}", c.Proposal.Meta.Title,
		"{{change}}", c.Proposal.Meta.ID, "{{created}}", c.Proposal.Meta.Created, "{{base}}", short(p.Commit))
}

func withChange(in Instructions, c *doc.Change, files map[string]string) Instructions {
	for _, name := range c.Files {
		in.Change = append(in.Change, File{in.Root + "/" + name, files[name]})
	}
	return in
}

// Markdown renders the instructions as the text a session reads.
func (in Instructions) Markdown() string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Grove instructions: %s\n\n%s Run Grove as `%s`.\n\n%s\n", in.Workflow, in.Task, in.Exe, meanings)
	if in.Context != "" {
		fmt.Fprintf(&b, "## Project context\n\n%s\n\n", strings.TrimRight(in.Context, "\n"))
	}
	fmt.Fprintf(&b, "## How to %s\n\n%s\n", in.Workflow, in.Guide)
	if in.Description != "" {
		fmt.Fprintf(&b, "## Owner's description\n\n%s\n\n", in.Description)
	}
	if len(in.Documents) > 0 {
		b.WriteString("## Documents\n\n")
	}
	for _, d := range in.Documents {
		fmt.Fprintf(&b, "### %s: `%s`\n\n%s\n", d.Kind, d.Path, d.Guide)
		if len(d.Rules) > 0 {
			b.WriteString("The project's own rules for this document:\n\n")
			for _, r := range d.Rules {
				fmt.Fprintf(&b, "- %s\n", r)
			}
			b.WriteString("\n")
		}
		include(&b, d.Path, d.Template)
	}
	in.index(&b)
	if len(in.Change) > 0 {
		b.WriteString("## The change\n\n")
	}
	for _, f := range in.Change {
		include(&b, f.Path, f.Content)
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

A Grove document in the working tree that the target does not hold as it is,
is a draft: not agreed until it is merged.

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
