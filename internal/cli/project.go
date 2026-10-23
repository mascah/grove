package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/mascah/grove/internal/doc"
	"github.com/mascah/grove/internal/gitx"
	"github.com/mascah/grove/internal/instruct"
)

// project is the repository Grove was run in.
type project struct {
	Top     string            `yaml:"-"` // the top of the working tree
	Root    string            `yaml:"root"`
	Target  string            `yaml:"target"`
	Agents  map[string]string `yaml:"agents"` // harness to form, skills or commands
	Schema  string            `yaml:"schema"`
	Context string            `yaml:"context"`
	Rules   instruct.Rules    `yaml:"rules"`
}

func openProject() (project, error) {
	top, err := gitx.TopLevel(".")
	if err != nil {
		return project{}, err
	}
	b, err := os.ReadFile(filepath.Join(top, "grove.yaml"))
	if err != nil {
		return project{}, err
	}
	return parseProject(top, b)
}

// parseProject reads grove.yaml's contents b for the working tree at top.
func parseProject(top string, b []byte) (project, error) {
	p := project{Top: top, Root: "grove", Target: "main", Schema: "default"} // the documented defaults
	err := yaml.Unmarshal(b, &p)
	if len(p.Agents) == 0 { // set after reading, which would add to a default map rather than replace it
		p.Agents = defaultAgents
	}
	return p, err
}

// instruct returns what instructions say about the project, without the
// target's documents.
func (p project) instruct() instruct.Project {
	return instruct.Project{Root: p.Root, Target: p.Target, Context: p.Context, Rules: p.Rules,
		Schemas: os.DirFS(filepath.Join(p.Top, p.Root, "schemas")), Schema: p.Schema}
}

// agreed returns the target's commit and the documents it holds.
func (p project) agreed() (instruct.Project, error) {
	commit, err := gitx.ResolveBranch(p.Top, p.Target)
	if err != nil {
		return p.instruct(), err
	}
	return p.at(commit)
}

// at returns commit and the documents it holds.
func (p project) at(commit string) (instruct.Project, error) {
	ip := p.instruct()
	ip.Commit = commit
	fsys, err := gitx.TreeFS(p.Top, ip.Commit, p.Root)
	if err != nil {
		return ip, err
	}
	ip.Agreed, err = doc.Load(fsys, ip.Layout)
	return ip, err
}

func plural(n int, noun string) string {
	if n == 1 {
		return "1 " + noun
	}
	return fmt.Sprintf("%d %ss", n, noun)
}
