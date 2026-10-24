package cli

import (
	"fmt"
	"io"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/mascah/grove/internal/doc"
	"github.com/mascah/grove/internal/gitx"
	"github.com/mascah/grove/internal/instruct"
)

// check prints every problem in grove.yaml and in the Grove documents of
// this working tree, and each requirement a change is stale on.
func check(stdout io.Writer) error {
	p, err := openProject()
	if err != nil {
		return err
	}
	ip := p.instruct()
	unknown, err := unknownConfig(p.Top, ip)
	if err != nil {
		return err
	}
	unknown = append(unknown, ip.Check()...)
	for _, u := range unknown {
		fmt.Fprintln(stdout, u)
	}
	root := os.DirFS(filepath.Join(p.Top, p.Root))
	c, err := doc.Load(root, ip.Layout)
	if err != nil {
		return err
	}
	if err := stale(p, root, c); err != nil {
		return err
	}
	if err := stalePlans(p, c); err != nil {
		return err
	}
	if err := movedSinceBase(p, c); err != nil {
		return err
	}
	shapes(root, c)
	docs := append([]*doc.Document{}, c.Specs...)
	for _, ch := range c.Changes {
		docs = append(docs, &ch.Document)
	}
	for _, cy := range c.Cycles {
		docs = append(docs, &cy.Document)
	}
	n := len(unknown)
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

// stale reports each requirement a change in c modifies or removes that the
// maintained spec in root gives a text that is neither what it was at the
// change's base, nor what the change itself or a change it depends on gives
// it. Reading root rather than the target, it sees a dependency archived on
// this branch; allowing the change's own text, it stays quiet after archive
// has reconciled the spec.
func stale(p project, root fs.FS, c *doc.Collection) error {
	for _, ch := range c.Changes {
		if ch.Meta.Base == "" {
			continue
		}
		base, err := gitx.TreeFS(p.Top, ch.Meta.Base, p.Root)
		if err != nil {
			ch.Diagnostics = append(ch.Diagnostics, doc.Diagnostic{Path: ch.Dir + "/.grove.yaml", Message: "base: " + ch.Meta.Base + " is not a commit in this repository"})
			continue
		}
		for _, file := range slices.Sorted(maps.Keys(ch.Touches)) {
			spec := "specs/" + ch.Touches[file] + "/spec.md"
			b, _ := fs.ReadFile(root, file)
			for _, r := range doc.Requirements(string(b)) {
				was, ok := requirement(base, spec, r.Name)
				if (r.Section != "MODIFIED" && r.Section != "REMOVED") || !ok {
					continue
				}
				now, present := requirement(root, spec, r.Name)
				fine := present && now == was
				for _, g := range gives(root, behind(c, ch), ch.Touches[file], r.Name) {
					fine = fine || g.Section == "REMOVED" && !present || g.Section != "REMOVED" && present && g.Text == now
				}
				if !fine {
					ch.Diagnostics = append(ch.Diagnostics, doc.Diagnostic{Path: file, Line: r.Line,
						Message: fmt.Sprintf("requirement %q was %s since base %s; revise the change with propose", r.Name, map[bool]string{true: "rewritten", false: "removed"}[present], ch.Meta.Base[:min(7, len(ch.Meta.Base))])})
				}
			}
		}
	}
	return nil
}

// shapes reports each requirement of a maintained spec or an open change
// spec in c that lacks "shall" or a scenario, and each change heading left
// in a maintained spec.
func shapes(root fs.FS, c *doc.Collection) {
	report := func(d *doc.Document, file string) {
		b, _ := fs.ReadFile(root, file)
		for _, r := range doc.Requirements(string(b)) {
			for _, m := range r.Problems() {
				d.Diagnostics = append(d.Diagnostics, doc.Diagnostic{Path: file, Line: r.Line, Message: m})
			}
		}
	}
	for _, d := range c.Specs {
		report(d, d.Path)
		b, _ := fs.ReadFile(root, d.Path)
		for _, l := range doc.ChangeHeadings(string(b)) {
			d.Diagnostics = append(d.Diagnostics, doc.Diagnostic{Path: d.Path, Line: l, Message: "this heading belongs in a change spec; a maintained spec lists its requirements under ## Requirements"})
		}
	}
	for _, ch := range c.Changes {
		for _, file := range slices.Sorted(maps.Keys(ch.Touches)) {
			report(&ch.Document, file)
		}
	}
}

// stalePlans reports each member of a cycle in c, marked unattended and
// not done, that changed on the target since the cycle's plan. Only then is
// the target read, so a project with no such cycle checks without it.
func stalePlans(p project, c *doc.Collection) error {
	var target instruct.Project
	for _, cy := range c.Cycles {
		if !cy.Meta.Unattended || cy.Done || cy.Meta.Planned == "" {
			continue
		}
		if target.Agreed == nil {
			var err error
			if target, err = p.agreed(); err != nil {
				return err
			}
		}
		changed, err := changedSincePlan(p, target.Commit, target.Agreed, cy)
		if err != nil {
			cy.Diagnostics = append(cy.Diagnostics, doc.Diagnostic{Path: cy.Path, Message: "planned: " + cy.Meta.Planned + " is not a commit in this repository"})
		}
		for _, id := range changed {
			cy.Diagnostics = append(cy.Diagnostics, doc.Diagnostic{Path: cy.Path,
				Message: fmt.Sprintf("planned: %s changed on %s since this plan, so grove run skips this cycle; re-plan it with /grove-plan %s", id, p.Target, cy.ID)})
		}
	}
	return nil
}

// movedSinceBase reports each task of a change in c, in no cycle that is
// not done, whose named file differs on the target from the change's base.
// A member of such a cycle is its plan's to report instead (stalePlans).
func movedSinceBase(p project, c *doc.Collection) error {
	commit := ""
	for _, ch := range c.Changes {
		if inOpenCycle(c, ch.ID) || ch.Meta.Base == "" || !slices.ContainsFunc(ch.Tasks, func(t doc.Task) bool { return len(t.Files) > 0 }) {
			continue
		}
		if commit == "" {
			var err error
			if commit, err = gitx.ResolveBranch(p.Top, p.Target); err != nil {
				return err
			}
		}
		tasks, err := moved(p, commit, ch)
		if err != nil { // stale reports a base that is no commit
			continue
		}
		for _, t := range tasks {
			for _, f := range t.Files {
				ch.Diagnostics = append(ch.Diagnostics, doc.Diagnostic{Path: ch.Dir + "/tasks.md",
					Message: fmt.Sprintf("%s's %s changed on %s since base %s; revise the change with propose", t.ID, f, p.Target, ch.Meta.Base[:min(7, len(ch.Meta.Base))])})
			}
		}
	}
	return nil
}

// inOpenCycle reports whether a cycle in c that is not done lists change id.
func inOpenCycle(c *doc.Collection, id string) bool {
	return slices.ContainsFunc(c.Cycles, func(cy *doc.Cycle) bool { return !cy.Done && slices.Contains(cy.Meta.Changes, id) })
}

// moved returns ch's tasks that name a file differing at commit from it at
// ch's base, each with only those files.
func moved(p project, commit string, ch *doc.Change) ([]doc.Task, error) {
	diff, err := gitx.Run(p.Top, append([]string{"diff", "--name-only", "--end-of-options", ch.Meta.Base, commit, "--"}, named(ch)...)...)
	if err != nil {
		return nil, err
	}
	changed := strings.Split(diff, "\n")
	var out []doc.Task
	for _, t := range ch.Tasks {
		t.Files = slices.DeleteFunc(slices.Clone(t.Files), func(f string) bool { return !slices.Contains(changed, f) })
		if len(t.Files) > 0 {
			out = append(out, t)
		}
	}
	return out, nil
}

// named returns every file ch's tasks name.
func named(ch *doc.Change) []string {
	var out []string
	for _, t := range ch.Tasks {
		out = append(out, t.Files...)
	}
	return out
}

// changedSincePlan returns each member of cy that at, the documents at
// commit, holds open, whose folder, or a file its tasks name, a commit on
// commit's first-parent history since cy's planned commit changed without
// also editing the cycle file. A commit that edits both is the plan
// checking that member, so a squash or merge commit of a revision and its
// re-plan leaves the plan fresh.
func changedSincePlan(p project, commit string, at *doc.Collection, cy *doc.Cycle) ([]string, error) {
	planned, err := gitx.Run(p.Top, "rev-parse", "--verify", "--end-of-options", cy.Meta.Planned+"^{commit}")
	if err != nil {
		return nil, err
	}
	log, err := gitx.Run(p.Top, "log", "--first-parent", "--diff-merges=first-parent", "--name-only", "--format=%x00", strings.TrimSpace(planned)+".."+commit)
	if err != nil {
		return nil, err
	}
	var commits [][]string // the files each commit changed
	for _, c := range strings.Split(log, "\x00")[1:] {
		commits = append(commits, strings.Split(strings.TrimSpace(c), "\n"))
	}
	var out []string
	for _, id := range cy.Meta.Changes {
		ch := at.Change(id)
		if ch == nil {
			continue
		}
		touches := func(f string) bool {
			return strings.HasPrefix(f, p.Root+"/"+ch.Dir+"/") || slices.Contains(named(ch), f)
		}
		if slices.ContainsFunc(commits, func(files []string) bool {
			return slices.ContainsFunc(files, touches) && !slices.Contains(files, p.Root+"/"+cy.Path)
		}) {
			out = append(out, id)
		}
	}
	return out, nil
}

// behind returns ch and every change it depends on, directly or through
// others, open or archived in c.
func behind(c *doc.Collection, ch *doc.Change) []*doc.Change {
	all := map[string]*doc.Change{}
	for _, x := range slices.Concat(c.Changes, c.Archived) {
		all[x.ID] = x
	}
	var out []*doc.Change
	var walk func(x *doc.Change)
	walk = func(x *doc.Change) {
		if x == nil || slices.Contains(out, x) {
			return
		}
		out = append(out, x)
		for _, id := range x.Meta.DependsOn {
			walk(all[id])
		}
	}
	walk(ch)
	return out
}

// gives returns what the change specs of changes about capability say
// requirement name becomes: each ADDED, MODIFIED or REMOVED block of that
// name.
func gives(root fs.FS, changes []*doc.Change, capability, name string) []doc.Requirement {
	var out []doc.Requirement
	for _, ch := range changes {
		for _, file := range slices.Sorted(maps.Keys(ch.Touches)) {
			if ch.Touches[file] != capability {
				continue
			}
			b, _ := fs.ReadFile(root, file)
			for _, r := range doc.Requirements(string(b)) {
				if r.Name == name && r.Section != "" {
					out = append(out, r)
				}
			}
		}
	}
	return out
}

// requirement returns the text of requirement name in the spec file of
// fsys, and whether it has one.
func requirement(fsys fs.FS, file, name string) (string, bool) {
	b, _ := fs.ReadFile(fsys, file)
	for _, r := range doc.Requirements(string(b)) {
		if r.Name == name {
			return r.Text, true
		}
	}
	return "", false
}

// configFields are grove.yaml's fields, as project reads them.
var configFields = []string{"root", "target", "agents", "schema", "context", "rules"}

// unknownConfig returns a problem for each field and rules key of grove.yaml,
// rules being keyed by documents and skills, then by kind and by skill,
// that Grove does not know, and a schema that is not one. Other commands
// ignore the unknown fields and keys (D-58).
func unknownConfig(top string, ip instruct.Project) ([]string, error) {
	b, err := os.ReadFile(filepath.Join(top, "grove.yaml"))
	if err != nil {
		return nil, err
	}
	var root yaml.Node
	if err := yaml.Unmarshal(b, &root); err != nil || len(root.Content) == 0 {
		return nil, err
	}
	var skills []string
	for _, s := range instruct.Skills() {
		skills = append(skills, s.Name)
	}
	kinds, out := ip.Kinds(), []string(nil)
	fields := root.Content[0].Content
	for i := 0; i+1 < len(fields); i += 2 {
		key, value := fields[i], fields[i+1]
		if _, _, problem := ip.Layout(value.Value); key.Value == "schema" && problem != "" {
			out = append(out, fmt.Sprintf("grove.yaml:%d: %s", value.Line, problem))
		}
		if !slices.Contains(configFields, key.Value) {
			out = append(out, fmt.Sprintf("grove.yaml:%d: unknown field %q; grove.yaml takes %s", key.Line, key.Value, and(configFields)))
			continue
		}
		for j := 0; key.Value == "agents" && j+1 < len(value.Content); j += 2 {
			if _, problem := harnessFile(value.Content[j].Value, value.Content[j+1].Value); problem != "" {
				out = append(out, fmt.Sprintf("grove.yaml:%d: %s", value.Content[j].Line, problem))
			}
		}
		for j := 0; key.Value == "rules" && j+1 < len(value.Content); j += 2 {
			section, names := value.Content[j], value.Content[j+1]
			valid, noun := skills, "skills"
			switch section.Value {
			case "documents":
				valid, noun = kinds, "kinds"
			case "skills":
			default:
				out = append(out, fmt.Sprintf("grove.yaml:%d: unknown rules key %q; rules takes documents and skills", section.Line, section.Value))
				continue
			}
			for n := 0; n+1 < len(names.Content); n += 2 {
				if k := names.Content[n]; section.Value == "skills" && k.Value == "apply" {
					out = append(out, fmt.Sprintf("grove.yaml:%d: unknown rules.skills key \"apply\"; its lines go under rules.skills.implement", k.Line))
				} else if !slices.Contains(valid, k.Value) {
					out = append(out, fmt.Sprintf("grove.yaml:%d: unknown rules.%s key %q; %s are %s", k.Line, section.Value, k.Value, noun, strings.Join(valid, ", ")))
				}
			}
		}
	}
	return out, nil
}

func problems(stdout io.Writer, ds []doc.Diagnostic) {
	for _, d := range ds {
		fmt.Fprintf(stdout, "  problem: %s\n", d)
	}
}
