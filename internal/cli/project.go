package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/mascah/grove-next2/internal/doc"
	"github.com/mascah/grove-next2/internal/gitx"
	"github.com/mascah/grove-next2/internal/instruct"
)

// project is the repository Grove was run in.
type project struct {
	Top     string              `yaml:"-"` // the top of the working tree
	Root    string              `yaml:"root"`
	Target  string              `yaml:"target"`
	Agents  []string            `yaml:"agents"`
	Schema  string              `yaml:"schema"`
	Context string              `yaml:"context"`
	Rules   map[string][]string `yaml:"rules"`
}

func openProject() (project, error) {
	p := project{Root: "grove", Target: "main", Agents: []string{"claude", "codex"}, Schema: "default"} // the documented defaults
	top, err := gitx.TopLevel(".")
	if err != nil {
		return p, err
	}
	p.Top = top
	b, err := os.ReadFile(filepath.Join(top, "grove.yaml"))
	if err != nil {
		return p, err
	}
	return p, yaml.Unmarshal(b, &p)
}

// instruct returns what instructions say about the project, without the
// target's documents.
func (p project) instruct() instruct.Project {
	return instruct.Project{Root: p.Root, Target: p.Target, Context: p.Context, Rules: p.Rules,
		Schemas: os.DirFS(filepath.Join(p.Top, p.Root, "schemas")), Schema: p.Schema}
}

// agreed returns the target's commit and the documents it holds.
func (p project) agreed() (instruct.Project, error) {
	ip := p.instruct()
	var err error
	if ip.Commit, err = gitx.ResolveBranch(p.Top, p.Target); err != nil {
		return ip, err
	}
	fsys, err := gitx.TreeFS(p.Top, ip.Commit, p.Root)
	if err != nil {
		return ip, err
	}
	ip.Agreed, err = doc.Load(fsys, ip.MainDocument)
	return ip, err
}

func plural(n int, noun string) string {
	if n == 1 {
		return "1 " + noun
	}
	return fmt.Sprintf("%d %ss", n, noun)
}
