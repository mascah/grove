package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/mascah/grove/internal/doc"
	"github.com/mascah/grove/internal/instruct"
)

// check prints every problem in grove.yaml and in the Grove documents of
// this working tree.
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
	c, err := doc.Load(os.DirFS(filepath.Join(p.Top, p.Root)), ip.Layout)
	if err != nil {
		return err
	}
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

func problems(stdout io.Writer, indent string, ds []doc.Diagnostic) {
	for _, d := range ds {
		fmt.Fprintf(stdout, "%sproblem: %s\n", indent, d)
	}
}
