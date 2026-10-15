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
	Context string              `yaml:"context"`
	Rules   map[string][]string `yaml:"rules"`
}

func openProject() (project, error) {
	p := project{Root: "grove", Target: "main", Agents: []string{"claude", "codex"}} // the documented defaults
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

// agreed returns the target's commit and the documents it holds.
func (p project) agreed() (instruct.Project, error) {
	ip := instruct.Project{Root: p.Root, Target: p.Target, Context: p.Context, Rules: p.Rules}
	var err error
	if ip.Commit, err = gitx.ResolveBranch(p.Top, p.Target); err != nil {
		return ip, err
	}
	fsys, err := gitx.TreeFS(p.Top, ip.Commit, p.Root)
	if err != nil {
		return ip, err
	}
	ip.Agreed, err = doc.Load(fsys)
	return ip, err
}

func plural(n int, noun string) string {
	if n == 1 {
		return "1 " + noun
	}
	return fmt.Sprintf("%d %ss", n, noun)
}
