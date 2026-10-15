package cli

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/mascah/grove-next2/internal/gitx"
	"github.com/mascah/grove-next2/internal/instruct"
)

const config = `# Grove: where its documents live, which branch holds what is agreed, and
# the harnesses init writes skills for.
root: grove
target: main
agents: [claude, codex]
# Facts about this project given to every skill: the language, how tests
# run, anything an agent cannot read from the code. Keep it short.
context: ""
# Extra lines for one kind of document, given when it is written, such as
#   proposal:
#     - List who is affected.
# Kinds: proposal, change-spec, tasks, implementation, spec.
rules: {}
`

// harnesses is where each harness reads a repository's skills, and what a
// skill writes for the words the user typed after its name (rule 10).
var harnesses = map[string]struct{ dir, args, note string }{
	"claude": {".claude/skills", "$ARGUMENTS", ""},
	"codex":  {".agents/skills", "<args>", " `<args>` is what the user typed after the skill's name, or nothing."},
}

// initProject writes grove.yaml at the top of the worktree unless it is
// there, then the skills of every workflow for each harness it lists.
func initProject(stdout io.Writer) error {
	top, err := gitx.TopLevel(".")
	if err != nil {
		return err
	}
	wrote := "Wrote grove.yaml and"
	if err := create(filepath.Join(top, "grove.yaml"), config); errors.Is(err, fs.ErrExist) {
		wrote = "Kept grove.yaml and wrote"
	} else if err != nil {
		return err
	}
	p, err := openProject()
	if err != nil {
		return err
	}
	exe, err := os.Executable()
	if err != nil {
		return err
	}
	for _, agent := range p.Agents {
		h, ok := harnesses[agent]
		if !ok {
			return fmt.Errorf("grove.yaml lists agent %q; Grove knows claude and codex", agent)
		}
		for _, w := range instruct.Workflows() {
			args := strings.ReplaceAll(w.Args, "<args>", h.args)
			note := h.note
			if args == "" {
				note = ""
			}
			skill := fmt.Sprintf("---\nname: grove-%s\ndescription: %s\n---\n\nRun `%s instructions %s` and follow what it prints, in this conversation.%s\n",
				w.Name, w.Skill, exe, strings.TrimSpace(w.Name+" "+args), note)
			path := filepath.Join(top, h.dir, "grove-"+w.Name, "SKILL.md")
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				return err
			}
			if err := os.WriteFile(path, []byte(skill), 0o644); err != nil {
				return err
			}
		}
	}
	fmt.Fprintf(stdout, "%s the Grove skills for %s; nothing is committed.\n", wrote, strings.Join(p.Agents, " and "))
	return nil
}

func create(path, content string) error {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return err
	}
	if _, err := f.WriteString(content); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}
