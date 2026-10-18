// Package instruct holds Grove's skills and its built-in change schema, and
// builds the instructions a session follows. Instructions are index first: the
// documents being worked on are included in full, everything else is listed
// with the commit to read it at.
package instruct

import (
	"bytes"
	"cmp"
	"embed"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"path"
	"regexp"
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

// Schema is a change schema as a project has it: its kinds of document
// and, for each skill, its guide and the kinds of document it writes.
type Schema struct {
	Name        string `yaml:"-"`
	Description string
	Documents   []Kind
	Skills      map[string]Step
}

// Step is what a change schema says for one skill.
type Step struct {
	Guide  string
	Writes []string
}

// builtin is Grove's change schema, default, as its schema.yaml has it.
var builtin Schema

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
	for name, v := range map[string]any{"skills.yaml": &skills, "schemas/default/schema.yaml": &builtin} {
		if err := yaml.Unmarshal([]byte(embedded(name)), v); err != nil {
			panic(err) // embedded at build time
		}
	}
}

// SchemaNames returns the project's change schemas: default, then each
// folder of <root>/schemas/ that holds a schema.yaml.
func (p Project) SchemaNames() []string {
	names := []string{"default"}
	found, _ := fs.Glob(p.Schemas, "*/schema.yaml") // only a bad pattern fails
	for _, f := range found {
		if dir := path.Dir(f); dir != "default" {
			names = append(names, dir)
		}
	}
	return names
}

// dirs are the folders of <root>/schemas/ that schema name's files come
// from, first first; Grove's default comes after them (D-55).
func dirs(name string) []string { return slices.Compact([]string{name, "default"}) }

// load returns change schema name: Grove's default, with the schema.yaml of
// the project's default folder over it, then the schema's own, field by
// field.
func (p Project) load(name string) (Schema, error) {
	if names := p.SchemaNames(); !slices.Contains(names, name) {
		return Schema{}, errors.New(notSchema(name, names))
	}
	s := builtin
	s.Name, s.Skills = name, maps.Clone(builtin.Skills)
	for _, dir := range slices.Backward(dirs(name)) {
		b, err := fs.ReadFile(p.Schemas, dir+"/schema.yaml")
		if errors.Is(err, fs.ErrNotExist) {
			continue
		}
		var l Schema
		if err == nil {
			err = yaml.Unmarshal(b, &l)
		}
		if err != nil {
			return s, fmt.Errorf("%s/schemas/%s/schema.yaml: %w", p.Root, dir, err)
		}
		s.Description = cmp.Or(l.Description, s.Description)
		if l.Documents != nil {
			s.Documents = l.Documents
		}
		for n, st := range l.Skills {
			cur := s.Skills[n]
			cur.Guide = cmp.Or(st.Guide, cur.Guide)
			if st.Writes != nil {
				cur.Writes = st.Writes
			}
			s.Skills[n] = cur
		}
	}
	return s, nil
}

func notSchema(name string, names []string) string {
	return fmt.Sprintf("schema %q is not a change schema; schemas are %s", name, and(names))
}

// and joins names as "a, b and c".
func and(names []string) string {
	if len(names) == 1 {
		return names[0]
	}
	return strings.Join(names[:len(names)-1], ", ") + " and " + names[len(names)-1]
}

// Kinds returns the names of the kinds of document of every change schema
// the project has, in order.
func (p Project) Kinds() []string {
	var names []string
	for _, n := range p.SchemaNames() {
		s, _ := p.load(n) // a broken schema.yaml is check's to report
		for _, k := range s.Documents {
			if !slices.Contains(names, k.Kind) {
				names = append(names, k.Kind)
			}
		}
	}
	return names
}

// MainDocument returns the file in a change folder that holds the main
// document of change schema name, or a problem when name is not a schema.
func (p Project) MainDocument(name string) (file, problem string) {
	if names := p.SchemaNames(); !slices.Contains(names, name) {
		return "", notSchema(name, names)
	}
	s, _ := p.load(name) // a broken schema.yaml is check's to report
	for _, k := range s.Documents {
		if k.Metadata == changeMetadata {
			return strings.TrimPrefix(k.Path, "changes/<change>/"), ""
		}
	}
	return "", ""
}

// Project is what a session is told about the repository.
type Project struct {
	Root    string          // document root, such as "grove"
	Target  string          // target branch, such as "main"
	Commit  string          // the target's commit
	Agreed  *doc.Collection // the documents at Commit
	Context string          // grove.yaml's context
	Rules   Rules           // grove.yaml's rules
	Schemas fs.FS           // <root>/schemas in the working tree, the project's own schema files
	Schema  string          // grove.yaml's schema, for exploring and new changes
}

// Rules are grove.yaml's extra lines: for a kind of document, given when it
// is written, and for a skill, given with its guide. They are apart because
// the two are named alike, such as proposal and propose.
type Rules struct {
	Documents map[string][]string
	Skills    map[string][]string
}

// Instructions is everything one skill's instructions say, built before
// it is printed (rule 10).
type Instructions struct {
	Project
	Skill       string
	Follows     Schema     // the change schema the instructions follow
	Schemas     []Schema   // the project's change schemas, listed when a new change may choose
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

// Check returns each problem in the project's change schemas: the
// structure every schema keeps, the files each names, and, which other
// commands ignore, unknown fields and unused files.
func (p Project) Check() []string {
	out, used := p.fields(), map[string]bool{}
	for _, n := range p.SchemaNames() {
		s, err := p.load(n)
		if err != nil {
			continue // invalid YAML, which fields reports
		}
		out = append(out, p.problems(s)...)
		for _, f := range s.files() {
			for _, dir := range dirs(n) {
				used[dir+"/"+f] = true
			}
		}
	}
	fs.WalkDir(p.Schemas, ".", func(name string, d fs.DirEntry, err error) error {
		dir, _, _ := strings.Cut(name, "/")
		if err == nil && !d.IsDir() && path.Base(name) != "schema.yaml" && !used[name] && slices.Contains(p.SchemaNames(), dir) {
			out = append(out, fmt.Sprintf("%s/schemas/%s: schema %s does not use this file", p.Root, name, dir))
		}
		return err
	})
	return out
}

var (
	yamlLine     = regexp.MustCompile(`^(?:yaml: )?line ([0-9]+): (.*)$`)
	unknownField = regexp.MustCompile(`^field (\S+) not found in type instruct\.(\w+)$`)
	takes        = map[string]string{"Schema": "schema.yaml takes description, documents and skills",
		"Kind": "a document takes kind, path, template, guide and metadata", "Step": "a skill takes guide and writes"}
)

// fields returns each unknown field and YAML error in the project's
// schema.yaml files.
func (p Project) fields() []string {
	var out []string
	for _, n := range p.SchemaNames() {
		file := n + "/schema.yaml"
		b, err := fs.ReadFile(p.Schemas, file)
		if err != nil {
			continue // default's is optional
		}
		d := yaml.NewDecoder(bytes.NewReader(b))
		d.KnownFields(true)
		msgs, te := []string{}, &yaml.TypeError{}
		if err := d.Decode(&Schema{}); errors.As(err, &te) {
			msgs = te.Errors
		} else if err != nil && !errors.Is(err, io.EOF) {
			msgs = []string{err.Error()}
		}
		for _, m := range msgs {
			line, msg := "", m
			if l := yamlLine.FindStringSubmatch(m); l != nil {
				line, msg = ":"+l[1], l[2]
			}
			if f := unknownField.FindStringSubmatch(msg); f != nil {
				msg = fmt.Sprintf("unknown field %q; %s", f[1], takes[f[2]])
			}
			out = append(out, fmt.Sprintf("%s/schemas/%s%s: %s", p.Root, file, line, msg))
		}
	}
	return out
}

// problems returns what makes schema s unusable: a structure Grove relies
// on that it lacks, a kind it writes but does not list, or a file it names
// that is in none of its places.
func (p Project) problems(s Schema) []string {
	var out []string
	report := func(format string, a ...any) {
		out = append(out, fmt.Sprintf("%s/schemas/%s/schema.yaml: ", p.Root, s.Name)+fmt.Sprintf(format, a...))
	}
	kinds, mains := map[string]Kind{}, 0
	for _, k := range s.Documents {
		kinds[k.Kind] = k
		if _, ok := metadataTemplates[k.Metadata]; k.Metadata != "" && !ok {
			report("%s: metadata must be %s or specs/<capability>/.grove.yaml", k.Kind, changeMetadata)
		}
		if k.Metadata == changeMetadata {
			mains++
		}
		if !strings.HasPrefix(k.Path, "changes/<change>/") && (k.Kind != "spec" || k.Path != "specs/<capability>/spec.md") {
			report("%s is at %s; kinds are in changes/<change>/, except spec at specs/<capability>/spec.md", k.Kind, k.Path)
		}
	}
	for _, k := range []string{"tasks", "implementation"} {
		if kinds[k].Path != "changes/<change>/"+k+".md" {
			report("%s must be at changes/<change>/%s.md; Grove reads it to archive", k, k)
		}
	}
	if mains != 1 {
		report("one kind, the change's main document, must have metadata %s", changeMetadata)
	}
	var names []string
	for _, sk := range skills {
		names = append(names, sk.Name)
		for _, w := range s.Skills[sk.Name].Writes {
			if _, ok := kinds[w]; !ok {
				report("%s writes %q, which this schema does not list", sk.Name, w)
			}
		}
	}
	for _, n := range slices.Sorted(maps.Keys(s.Skills)) {
		if !slices.Contains(names, n) {
			report("unknown skill %q; skills are %s", n, and(names))
		}
	}
	for _, f := range s.files() {
		if _, err := p.file(s.Name, f); err != nil {
			out = append(out, err.Error())
		}
	}
	return out
}

// files returns every file schema s names.
func (s Schema) files() []string {
	var out []string
	for _, k := range s.Documents {
		out = append(out, k.Template, k.Guide)
	}
	for _, sk := range skills {
		out = append(out, s.Skills[sk.Name].Guide)
	}
	return out
}

// file returns change schema s's file name, from the first of its folders
// that has it, else from Grove's default.
func (p Project) file(s, name string) (string, error) {
	for _, dir := range dirs(s) {
		b, err := fs.ReadFile(p.Schemas, dir+"/"+name)
		if !errors.Is(err, fs.ErrNotExist) {
			return string(b), err
		}
	}
	b, err := files.ReadFile("schemas/default/" + name)
	if err != nil {
		also := ""
		if s != "default" {
			also = ", in " + p.Root + "/schemas/default"
		}
		return "", fmt.Errorf("%s/schemas/%s: no %s here%s or in Grove's default", p.Root, s, name, also)
	}
	return string(b), nil
}

// build returns skill name's instructions in change schema s, with its guide
// and documents, fill's values put into their paths and templates, and
// change c's files when c is not nil.
func build(p Project, s, name, task string, fill *strings.Replacer, c *doc.Change, files map[string]string) (Instructions, error) {
	schema, err := p.load(s)
	if err != nil {
		return Instructions{}, err
	}
	if problems := p.problems(schema); len(problems) > 0 {
		return Instructions{}, errors.New(problems[0])
	}
	file := func(name string) string {
		f, e := p.file(s, name)
		err = cmp.Or(err, e)
		return f
	}
	st := schema.Skills[name]
	in := Instructions{Project: p, Skill: name, Follows: schema, Task: task, Requires: embedded("skills/" + name + ".md"), Guide: file(st.Guide)}
	for _, kn := range st.Writes {
		k := schema.Documents[slices.IndexFunc(schema.Documents, func(k Kind) bool { return k.Kind == kn })]
		d := Document{Kind: k.Kind, Path: p.Root + "/" + fill.Replace(k.Path),
			Guide: file(k.Guide), Template: fill.Replace(file(k.Template)), Rules: p.Rules.Documents[k.Kind]}
		if k.Metadata != "" {
			d.Metadata, d.MetadataTemplate = p.Root+"/"+fill.Replace(k.Metadata), fill.Replace(embedded(metadataTemplates[k.Metadata]))
		}
		in.Documents = append(in.Documents, d)
	}
	if c != nil {
		for _, name := range c.Files {
			in.Change = append(in.Change, File{in.Root + "/" + name, files[name]})
		}
	}
	return in, err
}

// choices returns the project's change schemas when it has more than
// default, for a new change to choose from.
func (p Project) choices() []Schema {
	names := p.SchemaNames()
	if len(names) == 1 {
		return nil
	}
	var out []Schema
	for _, n := range names {
		s, _ := p.load(n) // a broken schema.yaml is check's to report
		out = append(out, s)
	}
	return out
}

// with names change schema s in a task sentence, unless it is default.
func with(s string) string {
	if s == "default" {
		return ""
	}
	return ", with change schema " + s
}

// Explore returns the instructions for interviewing the owner about an idea.
func Explore(p Project) (Instructions, error) {
	in, err := build(p, p.Schema, "explore", "You are exploring an idea with the owner, in their own working tree.", strings.NewReplacer(), nil, nil)
	in.Schemas = p.choices()
	return in, err
}

// Propose returns the instructions for writing new change id in change
// schema s, created at created, from the conversation or from description
// when it is not empty.
func Propose(p Project, s, id, created, description string) (Instructions, error) {
	from := "this conversation"
	if description != "" {
		from = "the owner's description below"
	}
	in, err := build(p, s, "propose", fmt.Sprintf("You are proposing new change %s%s, from %s, in the owner's own working tree.", id, with(s), from),
		strings.NewReplacer("<change>", id+"-<slug>", "{{change}}", id, "{{created}}", created, "{{schema}}", s, "{{base}}", short(p.Commit)), nil, nil)
	in.Description, in.Schemas = description, p.choices()
	return in, err
}

// Revise returns the instructions for revising drafted change c, read from
// the working tree; files are c.Files as read from there.
func Revise(p Project, c *doc.Change, files map[string]string) (Instructions, error) {
	task := fmt.Sprintf("You are revising drafted change %s, %q%s, in the owner's own working tree.", c.ID, c.Meta.Title, with(c.Meta.Schema))
	return build(p, c.Meta.Schema, "propose", task, fill(p, c), c, files)
}

// Apply returns the instructions for applying change c, which the target
// holds; files are c.Files as read from there.
func Apply(p Project, c *doc.Change, files map[string]string) (Instructions, error) {
	task := fmt.Sprintf("You are applying approved change %s, %q%s, as %s@%s holds it, in your own working tree.",
		c.ID, c.Meta.Title, with(c.Meta.Schema), p.Target, short(p.Commit))
	return build(p, c.Meta.Schema, "apply", task, fill(p, c), c, files)
}

// Archive returns the instructions for archiving change c, read from the
// working tree; files are c.Files as read from there.
func Archive(p Project, c *doc.Change, files map[string]string) (Instructions, error) {
	task := fmt.Sprintf("You are archiving change %s, %q%s, in the owner's own working tree.", c.ID, c.Meta.Title, with(c.Meta.Schema))
	return build(p, c.Meta.Schema, "archive", task, fill(p, c), c, files)
}

// fill puts change c's values into paths and templates.
func fill(p Project, c *doc.Change) *strings.Replacer {
	return strings.NewReplacer("<change>", path.Base(c.Dir), "{{title}}", c.Meta.Title,
		"{{change}}", c.ID, "{{created}}", c.Meta.Created, "{{schema}}", c.Meta.Schema, "{{base}}", short(p.Commit))
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
	rules(&b, "skill", in.Rules.Skills[in.Skill])
	if len(in.Schemas) > 0 {
		fmt.Fprintf(&b, "## Change schemas\n\nThe project's change schemas; these instructions follow `%s`.\n\n", in.Follows.Name)
		for _, s := range in.Schemas {
			fmt.Fprintf(&b, "- `%s`: %s\n", s.Name, s.Description)
		}
		b.WriteString("\n")
	}
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
