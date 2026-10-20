package cli

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/mascah/grove/internal/gitx"
	"github.com/mascah/grove/internal/instruct"
)

const config = `# Grove: where its documents live, which branch holds what is agreed, and
# the harnesses init writes skills for.
root: grove
target: main
agents: [claude, codex]
# The change schema a new change has unless the owner names another while
# proposing: Grove's default, or one of the project's in <root>/schemas/.
schema: default
# Facts about this project given to every skill: the language, how tests
# run, anything an agent cannot read from the code. Keep it short.
context: ""
# Extra lines for Grove's instructions: under documents, for one kind of
# document, given when it is written; under skills, for one skill, given
# with its guide. Such as
#   documents:
#     proposal:
#       - List who is affected.
#   skills:
#     apply:
#       - Commit after each task, naming its task ID.
# Kinds: proposal, change-spec, tasks, implementation, spec, and those of
# the project's own schemas.
# Skills: explore, propose, apply, archive.
rules: {}
`

// harnesses is where each harness reads a repository's skills, and what a
// skill writes for the words the user typed after its name (rule 10).
var harnesses = map[string]struct{ dir, args, note string }{
	"claude": {".claude/skills", "$ARGUMENTS", ""},
	"codex":  {".agents/skills", "<args>", " `<args>` is what the user typed after the skill's name, or nothing."},
}

// initProject writes grove.yaml at the top of the worktree unless it is
// there, then a SKILL.md for each skill and each harness it lists.
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
	for _, agent := range p.Agents {
		h, ok := harnesses[agent]
		if !ok {
			return fmt.Errorf("grove.yaml lists agent %q; Grove knows claude and codex", agent)
		}
		for _, s := range instruct.Skills() {
			args := strings.ReplaceAll(s.Args, "<args>", h.args)
			note := h.note
			if args == "" {
				note = ""
			}
			skill := fmt.Sprintf("---\nname: grove-%s\ndescription: %s\n---\n\nRun `grove instructions %s` and follow what it prints, in this conversation.%s\n",
				s.Name, s.Description, strings.TrimSpace(s.Name+" "+args), note)
			path := filepath.Join(top, h.dir, "grove-"+s.Name, "SKILL.md")
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
