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
	Top    string            `yaml:"-"` // the top of the working tree
	Root   string            `yaml:"root"`
	Target string            `yaml:"target"`
	Agents map[string]string `yaml:"agents"` // harness to form, skills or commands
	Schema string            `yaml:"schema"`
	Rules  instruct.Rules    `yaml:"rules"`
}

func openProject() (project, error) {
	p := project{Root: "grove", Target: "main", Schema: "default"} // the documented defaults
	top, err := gitx.TopLevel(".")
	if err != nil {
		return p, err
	}
	p.Top = top
	b, err := os.ReadFile(filepath.Join(top, "grove.yaml"))
	if err != nil {
		return p, err
	}
	err = yaml.Unmarshal(b, &p)
	if len(p.Agents) == 0 { // set after reading, which would add to a default map rather than replace it
		p.Agents = defaultAgents
	}
	return p, err
}

// instruct returns what instructions say about the project, without the
// target's documents.
func (p project) instruct() instruct.Project {
	return instruct.Project{Root: p.Root, Target: p.Target, Rules: p.Rules,
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
	ip.Agreed, err = doc.Load(fsys, ip.Layout)
	return ip, err
}

func plural(n int, noun string) string {
	if n == 1 {
		return "1 " + noun
	}
	return fmt.Sprintf("%d %ss", n, noun)
}
