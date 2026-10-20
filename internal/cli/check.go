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
	docs := append([]*doc.Document{}, c.Specs...)
	for _, ch := range c.Changes {
		docs = append(docs, &ch.Document)
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
// target rewrote or removed since the change's base. It compares with the
// target, not this working tree, so archiving a change that rewrites a spec
// does not report the change being archived.
func stale(p project, root fs.FS, c *doc.Collection) error {
	var target fs.FS
	for _, ch := range c.Changes {
		if ch.Meta.Base == "" {
			continue
		}
		if target == nil {
			commit, err := gitx.ResolveBranch(p.Top, p.Target)
			if err != nil {
				return err
			}
			if target, err = gitx.TreeFS(p.Top, commit, p.Root); err != nil {
				return err
			}
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
				what := ""
				switch now, ok := requirement(target, spec, r.Name); {
				case !ok:
					what = "removed"
				case now != was:
					what = "rewritten"
				}
				if what != "" {
					ch.Diagnostics = append(ch.Diagnostics, doc.Diagnostic{Path: file, Line: r.Line,
						Message: fmt.Sprintf("requirement %q was %s on %s since base %s; revise the change with propose", r.Name, what, p.Target, ch.Meta.Base[:min(7, len(ch.Meta.Base))])})
				}
			}
		}
	}
	return nil
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
			out = append(out, fmt.Sprintf("grove.yaml:%d: unknown field %q; grove.yaml takes %s and rules", key.Line, key.Value, strings.Join(configFields[:5], ", ")))
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
				if k := names.Content[n]; !slices.Contains(valid, k.Value) {
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
