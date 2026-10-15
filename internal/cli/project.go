package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/mascah/grove-next2/internal/doc"
	"github.com/mascah/grove-next2/internal/gitx"
	"github.com/mascah/grove-next2/internal/instruct"
)

// project is the repository Grove was run in.
type project struct {
	Main   string `yaml:"-"` // the main checkout, which holds .grove/worktrees
	Root   string `yaml:"root"`
	Target string `yaml:"target"`
}

func openProject() (project, error) {
	var p project
	top, err := gitx.TopLevel(".")
	if err != nil {
		return p, err
	}
	worktrees, err := gitx.Worktrees(top)
	if err != nil {
		return p, err
	}
	p.Main = worktrees[0].Path
	b, err := os.ReadFile(filepath.Join(p.Main, "grove.yaml"))
	if err != nil {
		return p, err
	}
	return p, yaml.Unmarshal(b, &p)
}

// agreed returns the target's commit and the documents it holds.
func (p project) agreed() (instruct.Project, error) {
	exe, _ := os.Executable()
	ip := instruct.Project{Root: p.Root, Target: p.Target, Exe: exe}
	var err error
	if ip.Commit, err = gitx.ResolveBranch(p.Main, p.Target); err != nil {
		return ip, err
	}
	fsys, err := gitx.TreeFS(p.Main, ip.Commit, p.Root)
	if err != nil {
		return ip, err
	}
	ip.Agreed, err = doc.Load(fsys)
	return ip, err
}

func (p project) worktrees() string { return filepath.Join(p.Main, ".grove", "worktrees") }

var notSlug = regexp.MustCompile(`[^a-z0-9]+`)

// slug turns a description into a short file-name part.
func slug(s string) string {
	s = strings.Trim(notSlug.ReplaceAllString(strings.ToLower(s), "-"), "-")
	for len(s) > 40 {
		i := strings.LastIndexByte(s[:40], '-')
		if i < 0 {
			i = 40
		}
		s = s[:i]
	}
	return s
}

// shown is path as the owner would type it from the working directory.
func shown(path string) string {
	wd, _ := os.Getwd()
	wd, _ = filepath.EvalSymlinks(wd) // Git reports /private/tmp where the shell says /tmp
	if rel, err := filepath.Rel(wd, path); err == nil && !strings.HasPrefix(rel, "..") {
		return rel
	}
	return path
}

func files(n int) string { return plural(n, "file") }

func plural(n int, noun string) string {
	if n == 1 {
		return "1 " + noun
	}
	return fmt.Sprintf("%d %ss", n, noun)
}
