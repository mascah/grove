package cli

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"github.com/mascah/grove/internal/gitx"
	"github.com/mascah/grove/internal/instruct"
)

const config = `# Grove: where its documents live, which branch holds what is agreed, and
# the harnesses init writes Grove's skills for.
root: grove
target: main
# Each harness init writes Grove's skills for (claude, codex, cursor,
# opencode, pi, hermes), as skills or, where the harness has them, as
# commands. Run grove init again after changing it.
agents:
%s# The change schema a new change has unless the owner names another while
# proposing: Grove's default, or one of the project's in <root>/schemas/.
schema: default
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
# Skills: explore, prototype, propose, implement, apply, archive.
rules: {}
`

const sharedSkills = ".agents/skills/grove-<skill>/SKILL.md"

// harnesses are the agents init writes Grove's skills for (rule 10): what
// at the top of a repository shows one is used, and the file it reads a
// skill from, <skill> being its name, as a skill and as a command, which
// not every harness has.
var harnesses = map[string]struct {
	marks            []string
	skills, commands string
}{
	"claude":   {[]string{".claude", "CLAUDE.md"}, ".claude/skills/grove-<skill>/SKILL.md", ".claude/commands/grove/<skill>.md"},
	"codex":    {[]string{".codex", "AGENTS.md"}, sharedSkills, ""},
	"cursor":   {[]string{".cursor"}, sharedSkills, ""},
	"hermes":   {[]string{".hermes"}, sharedSkills, ""},
	"opencode": {[]string{".opencode", "opencode.json"}, sharedSkills, ".opencode/commands/grove-<skill>.md"},
	"pi":       {[]string{".pi"}, sharedSkills, ".pi/prompts/grove-<skill>.md"},
}

// defaultAgents are the harnesses when the repository shows none, and
// when grove.yaml leaves agents out.
var defaultAgents = map[string]string{"claude": "skills", "codex": "skills"}

// initProject writes grove.yaml at the top of the worktree unless it is
// there, with the harnesses the repository shows, then each listed
// harness's skill files, and removes Grove's skill files it no longer lists.
func initProject(stdout io.Writer) error {
	top, err := gitx.TopLevel(".")
	if err != nil {
		return err
	}
	detected := map[string]string{}
	for name, h := range harnesses {
		for _, m := range h.marks {
			if _, err := os.Stat(filepath.Join(top, m)); err == nil {
				detected[name] = "skills"
			}
		}
	}
	if len(detected) == 0 {
		detected = defaultAgents
	}
	agents := ""
	for _, name := range slices.Sorted(maps.Keys(detected)) {
		agents += "  " + name + ": " + detected[name] + "\n"
	}
	wrote := "Wrote grove.yaml and"
	if err := create(filepath.Join(top, "grove.yaml"), fmt.Sprintf(config, agents)); errors.Is(err, fs.ErrExist) {
		wrote = "Kept grove.yaml and wrote"
	} else if err != nil {
		return err
	}
	p, err := openProject()
	if err != nil {
		return err
	}
	files, forms := map[string]string{}, map[string][]string{}
	for _, name := range slices.Sorted(maps.Keys(p.Agents)) {
		form := p.Agents[name]
		file, problem := harnessFile(name, form)
		if problem != "" {
			return errors.New("grove.yaml: " + problem)
		}
		forms[form] = append(forms[form], name)
		for _, s := range instruct.Skills() {
			files[strings.ReplaceAll(file, "<skill>", s.Name)] = skillFile(s, form)
		}
	}
	keep, places := map[string]bool{}, map[string]bool{}
	for file := range files {
		keep[unit(file)] = true
	}
	for _, h := range harnesses {
		for _, file := range []string{h.skills, h.commands} {
			if file != "" {
				places[unit(strings.ReplaceAll(file, "<skill>", "*"))] = true
			}
		}
	}
	// Removing before writing keeps a skill folder that is a link to another,
	// such as .claude/skills to .agents/skills, from losing what was just written.
	removed := map[string][]string{}
	for _, place := range slices.Sorted(maps.Keys(places)) {
		matches, _ := fs.Glob(os.DirFS(top), place) // places are valid patterns
		for _, rel := range matches {
			if keep[rel] {
				continue
			}
			if err := os.RemoveAll(filepath.Join(top, rel)); err != nil {
				return err
			}
			removed[path.Dir(rel)] = append(removed[path.Dir(rel)], path.Base(rel))
		}
	}
	for file, content := range files {
		if err := os.MkdirAll(filepath.Join(top, path.Dir(file)), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(top, file), []byte(content), 0o644); err != nil {
			return err
		}
	}
	var what, gone []string
	for _, form := range []string{"skills", "commands"} {
		if len(forms[form]) > 0 {
			what = append(what, form+" for "+and(forms[form]))
		}
	}
	for _, dir := range slices.Sorted(maps.Keys(removed)) {
		gone = append(gone, and(removed[dir])+" from "+dir)
	}
	msg := wrote + " the Grove " + and(what)
	if len(gone) > 0 {
		msg += ", and removed " + and(gone)
	}
	fmt.Fprintln(stdout, msg+"; nothing is committed.")
	return nil
}

// harnessFile returns the file harness name reads a skill from in form, or
// what is wrong with them.
func harnessFile(name, form string) (file, problem string) {
	h, ok := harnesses[name]
	if !ok {
		return "", fmt.Sprintf("unknown harness %q; agents takes %s", name, and(slices.Sorted(maps.Keys(harnesses))))
	}
	switch {
	case form == "skills":
		return h.skills, ""
	case form == "commands" && h.commands != "":
		return h.commands, ""
	}
	has := "skills"
	if h.commands != "" {
		has = "skills and commands"
	}
	return "", fmt.Sprintf("harness %q has no form %q; it has %s", name, form, has)
}

// skillFile is skill s as a harness reads it in form. Every skill file has
// one body, so harnesses sharing .agents/skills/ agree on it; Claude Code
// appends the user's words to it. A command has $ARGUMENTS replaced.
func skillFile(s instruct.Skill, form string) string {
	head, args, note := "name: grove-"+s.Name+"\n", s.Args, " `<args>` is what the user typed after the skill's name, or nothing."
	if form == "commands" {
		head, args, note = "", strings.ReplaceAll(s.Args, "<args>", "$ARGUMENTS"), ""
	}
	if args == "" {
		note = ""
	}
	return fmt.Sprintf("---\n%sdescription: %s\n---\n\nRun `grove instructions %s` and follow what it prints, in this conversation.%s\n",
		head, s.Description, strings.TrimSpace(s.Name+" "+args), note)
}

// unit is what init writes and removes for a skill file: a skill's folder,
// or a command's file.
func unit(file string) string {
	if path.Base(file) == "SKILL.md" {
		return path.Dir(file)
	}
	return file
}

// and lists words as a sentence does: a, b and c.
func and(words []string) string {
	if len(words) < 2 {
		return strings.Join(words, "")
	}
	return strings.Join(words[:len(words)-1], ", ") + " and " + words[len(words)-1]
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
