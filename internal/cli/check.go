package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/mascah/grove-next2/internal/doc"
	"github.com/mascah/grove-next2/internal/instruct"
)

// check prints every problem in grove.yaml and in the Grove documents of
// this working tree.
func check(stdout io.Writer) error {
	p, err := openProject()
	if err != nil {
		return err
	}
	unknown, err := unknownConfig(p.Top)
	if err != nil {
		return err
	}
	for _, u := range unknown {
		fmt.Fprintln(stdout, u)
	}
	c, err := doc.Load(os.DirFS(filepath.Join(p.Top, p.Root)), instruct.MainDocument())
	if err != nil {
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

// configFields are grove.yaml's fields, as project reads them.
var configFields = []string{"root", "target", "agents", "context", "rules"}

// unknownConfig returns a problem for each field and rules key of grove.yaml
// that Grove does not know. Other commands ignore them (D-58).
func unknownConfig(top string) ([]string, error) {
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
	kinds, out := instruct.Kinds(), []string(nil)
	fields := root.Content[0].Content
	for i := 0; i+1 < len(fields); i += 2 {
		key, value := fields[i], fields[i+1]
		if !slices.Contains(configFields, key.Value) {
			out = append(out, fmt.Sprintf("grove.yaml:%d: unknown field %q; grove.yaml takes %s and rules", key.Line, key.Value, strings.Join(configFields[:4], ", ")))
			continue
		}
		for j := 0; key.Value == "rules" && j+1 < len(value.Content); j += 2 {
			if k := value.Content[j]; !slices.Contains(kinds, k.Value) && !slices.Contains(skills, k.Value) {
				out = append(out, fmt.Sprintf("grove.yaml:%d: unknown rules key %q; kinds are %s, and skills are %s",
					k.Line, k.Value, strings.Join(kinds, ", "), strings.Join(skills, ", ")))
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
