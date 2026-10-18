// Package instruct holds Grove's skills and its built-in change schema, and
// builds the instructions a session follows. Instructions are index first: the
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

//go:embed skills.yaml skills metadata schemas
var files embed.FS

func embedded(name string) string {
	b, err := files.ReadFile(name)
	if err != nil {
		panic(err) // embedded at build time: a missing name is a programming error
	}
	return string(b)
}

// Skill is one of Grove's skills: the description of the skill init writes
// for it, and what that skill passes to `instructions <skill>`. What Grove
// requires of it is skills/<name>.md, the same in every change schema.
type Skill struct{ Name, Description, Args string }

// Skills returns Grove's skills, in order.
func Skills() []Skill { return skills }

var skills []Skill

// Kind is one kind of document: where it goes, the template it starts from,
// the guide for writing it, and the .grove.yaml written beside it, if any.
// Paths are relative to the document root.
type Kind struct{ Kind, Path, Template, Guide, Metadata string }

// Schema is a change schema: its kinds of document and, for each skill, its
// guide and the kinds of document it writes.
type Schema struct {
	Description string
	Documents   []Kind
	Skills      map[string]Step
}

// Step is what a change schema says for one skill.
type Step struct {
	Guide  string
	Writes []string
}

// schema is Grove's built-in change schema, default.
var schema Schema

// changeMetadata is where a change's metadata goes; the kind that names it
// is the change's main document.
const changeMetadata = "changes/<change>/.grove.yaml"

// metadataTemplates are Grove's own templates for the .grove.yaml a kind's
// metadata names. They are not a schema's, so no schema changes the fields
// Grove reads.
var metadataTemplates = map[string]string{
	changeMetadata:                   "metadata/change.grove.yaml",
	"specs/<capability>/.grove.yaml": "metadata/spec.grove.yaml",
}

func init() {
	for name, v := range map[string]any{"skills.yaml": &skills, "schemas/default/schema.yaml": &schema} {
		if err := yaml.Unmarshal([]byte(embedded(name)), v); err != nil {
			panic(err) // embedded at build time
		}
	}
}

// Kinds returns the names of the schema's kinds of document, in order.
func Kinds() []string {
	var names []string
	for _, k := range schema.Documents {
		names = append(names, k.Kind)
	}
	return names
}

// MainDocument returns the file of a change folder that holds its main
// document.
func MainDocument() string {
	for _, k := range schema.Documents {
		if k.Metadata == changeMetadata {
			return strings.TrimPrefix(k.Path, "changes/<change>/")
		}
	}
	return ""
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
	Context string              // grove.yaml's context
	Rules   map[string][]string // grove.yaml's rules, by kind of document or skill
}

// Instructions is everything one skill's instructions say, built before
// it is printed (rule 10).
type Instructions struct {
	Project
	Skill       string
	Task        string     // what the session is doing, one sentence
	Requires    string     // what Grove requires of the skill, in every schema
	Guide       string     // how the skill goes, from the change schema
	Description string     // the owner's description, when given
	Documents   []Document // the kinds of document the skill writes
	Change      []File     // the change being worked on, in full
}

// Document is one kind of document to write, with its paths and templates
// filled in.
type Document struct {
	Kind, Path, Guide, Template string
	Metadata, MetadataTemplate  string   // empty when the kind has no metadata
	Rules                       []string // the project's own, from grove.yaml
}

// File is one document included in full; Path is from the repository root.
type File struct{ Path, Content string }

// build returns skill name's instructions with its guide and documents,
// fill's values put into their paths and templates.
func build(p Project, name, task string, fill *strings.Replacer) Instructions {
	st, file := schema.Skills[name], func(name string) string { return embedded("schemas/default/" + name) }
	in := Instructions{Project: p, Skill: name, Task: task, Requires: embedded("skills/" + name + ".md"), Guide: file(st.Guide)}
	for _, kn := range st.Writes {
		k := kind(kn)
		d := Document{Kind: k.Kind, Path: p.Root + "/" + fill.Replace(k.Path),
			Guide: file(k.Guide), Template: fill.Replace(file(k.Template)), Rules: p.Rules[k.Kind]}
		if k.Metadata != "" {
			d.Metadata, d.MetadataTemplate = p.Root+"/"+fill.Replace(k.Metadata), fill.Replace(embedded(metadataTemplates[k.Metadata]))
		}
		in.Documents = append(in.Documents, d)
	}
	return in
}

// Explore returns the instructions for interviewing the owner about an idea.
func Explore(p Project) Instructions {
	return build(p, "explore", "You are exploring an idea with the owner, in their own working tree.", strings.NewReplacer())
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
	task := fmt.Sprintf("You are revising drafted change %s, %q, in the owner's own working tree.", c.ID, c.Meta.Title)
	return withChange(build(p, "propose", task, fill(p, c)), c, files)
}

// Apply returns the instructions for applying change c, which the target
// holds; files are c.Files as read from there.
func Apply(p Project, c *doc.Change, files map[string]string) Instructions {
	task := fmt.Sprintf("You are applying approved change %s, %q, as %s@%s holds it, in your own working tree.",
		c.ID, c.Meta.Title, p.Target, short(p.Commit))
	return withChange(build(p, "apply", task, fill(p, c)), c, files)
}

// Archive returns the instructions for archiving change c, read from the
// working tree; files are c.Files as read from there.
func Archive(p Project, c *doc.Change, files map[string]string) Instructions {
	task := fmt.Sprintf("You are archiving change %s, %q, in the owner's own working tree.", c.ID, c.Meta.Title)
	return withChange(build(p, "archive", task, fill(p, c)), c, files)
}

// fill puts change c's values into paths and templates.
func fill(p Project, c *doc.Change) *strings.Replacer {
	return strings.NewReplacer("<change>", path.Base(c.Dir), "{{title}}", c.Meta.Title,
		"{{change}}", c.ID, "{{created}}", c.Meta.Created, "{{base}}", short(p.Commit))
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
	fmt.Fprintf(&b, "# Grove instructions: %s\n\n%s Run Grove as `grove`.\n\n%s\n", in.Skill, in.Task, meanings)
	if in.Context != "" {
		fmt.Fprintf(&b, "## Project context\n\n%s\n\n", strings.TrimRight(in.Context, "\n"))
	}
	fmt.Fprintf(&b, "## What Grove requires of %s\n\n%s\n", in.Skill, in.Requires)
	fmt.Fprintf(&b, "## How to %s\n\n%s\n", in.Skill, in.Guide)
	rules(&b, "skill", in.Rules[in.Skill])
	if in.Description != "" {
		fmt.Fprintf(&b, "## Owner's description\n\n%s\n\n", in.Description)
	}
	if len(in.Documents) > 0 {
		b.WriteString("## Documents\n\n")
	}
	for _, d := range in.Documents {
		fmt.Fprintf(&b, "### %s: `%s`\n\n%s\n", d.Kind, d.Path, d.Guide)
		rules(&b, "document", d.Rules)
		if d.Metadata != "" {
			include(&b, d.Metadata, d.MetadataTemplate)
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
  implemented. Whoever applies it records progress in its
  ` + "`implementation.md`" + `.
- What the code does is observed behavior: evidence, not a rule.

A Grove document in the working tree that the target does not hold as it is,
is a draft: not agreed until it is merged.

Grove itself commits nothing. Commit only when the owner or the project's
rules ask; merging into the target branch is the only approval.
`

func (p Project) index(b *strings.Builder) {
	fmt.Fprintf(b, "## Project index\n\nWhat %s says at %s. Read a listed document with `git show %s:<path>`.\n\n",
		p.Target, short(p.Commit), short(p.Commit))
	b.WriteString("Maintained specs (agreed behavior):\n")
	for _, d := range p.Agreed.Specs {
		entry(b, p.Root+"/"+d.Dir+"/spec.md", d)
	}
	if len(p.Agreed.Specs) == 0 {
		b.WriteString("- None. Where nothing is agreed, what the code does is observed behavior.\n")
	}
	b.WriteString("\nApproved changes (future work; implementation is not established):\n")
	for _, c := range p.Agreed.Changes {
		entry(b, p.Root+"/"+c.Dir+"/", &c.Document)
	}
	if len(p.Agreed.Changes) == 0 {
		b.WriteString("- None.\n")
	}
	b.WriteString("\n")
}

func entry(b *strings.Builder, path string, d *doc.Document) {
	fmt.Fprintf(b, "- %s %q: `%s`", d.ID, d.Meta.Title, path)
	if d.Meta.Description != "" {
		fmt.Fprintf(b, " — %s", d.Meta.Description)
	}
	b.WriteString("\n")
	for _, diag := range d.Diagnostics {
		fmt.Fprintf(b, "  - problem: %s\n", diag)
	}
}

func rules(b *strings.Builder, of string, rs []string) {
	if len(rs) == 0 {
		return
	}
	fmt.Fprintf(b, "The project's own rules for this %s:\n\n", of)
	for _, r := range rs {
		fmt.Fprintf(b, "- %s\n", r)
	}
	b.WriteString("\n")
}

func include(b *strings.Builder, path, content string) {
	fmt.Fprintf(b, "<file path=%q>\n%s\n</file>\n\n", path, strings.TrimRight(content, "\n"))
}

func short(commit string) string { return commit[:min(7, len(commit))] }
