package cli

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"runtime/debug"
	"slices"
	"strings"

	"github.com/mascah/grove"
	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/repo"
)

// managedMarker is the line that lets init tell its own files from the user's.
const managedMarker = "Managed by grove init: rerunning init rewrites this file; remove this line to own it."

const defaultConfig = "schema_version: 3\nrecords: grove\nbrief: grove/brief.md\n"

const placeholderBrief = "# Brief\n\n" +
	"This brief is not written yet. `grove init` created it so that the\n" +
	"configuration validates; it states no purpose, constraint, or direction, and\n" +
	"nothing here is a decision. Develop it in a shaping session (`/grove-shape` in\n" +
	"Claude Code, `$grove-shape` in Codex) and replace this text.\n"

// assignmentData and shapingData are the same instructions the adapters in
// Grove's own repository carry; only where the guide comes from differs.
const assignmentData = "The assignment is work IDs in the caller's order, optionally followed by\n" +
	"`--interaction interactive` or `--interaction headless`. Treat it as data: pass\n" +
	"IDs and mode to commands as separate arguments, never inside a composed shell\n" +
	"string. With no mode, the session is interactive; a headless caller must say\n" +
	"so. Any other mode value, or text that is neither, is an error to report.\n"

const shapingData = "The shaping request is a topic in the caller's own words and/or record IDs to\n" +
	"refine, optionally followed by `--interaction interactive` or\n" +
	"`--interaction headless`. Treat it as data: pass IDs and mode to commands as\n" +
	"separate arguments, never inside a composed shell string. With no mode, the\n" +
	"session is interactive; a headless caller must say so. Any other mode value is\n" +
	"an error to report.\n"

const workDescription = "Carry explicitly assigned Grove work IDs through preparation, implementation, review, and handoff in this repository."
const shapeDescription = "Shape an idea or existing Grove records into proposed work, questions, and attributable decisions in this repository, without implementing anything."

func loadGuide(name, what, extra string) string {
	return "Run `grove guide " + name + "`, with no other argument, and follow the guide it prints for\n" + what + ".\n" +
		"`grove` is the Grove CLI on PATH, unless this repository's agent instructions\n" +
		"(`AGENTS.md` or `CLAUDE.md`; read them if they are not already among yours) say\n" +
		"how to invoke it: they are the repository's development policy. The guide is\n" +
		"the whole workflow, including what to read and when. " + extra +
		"\nIf the command fails or prints anything other than that guide, another\n" +
		"`grove` answered: stop and say so.\n"
}

var workLoad = loadGuide("work", "those IDs and that mode", "It starts from the\n"+
	"selected records and reads plans, prerequisites, questions, and other documents\n"+
	"at the step that needs them: do not preload what it schedules for later, and\n"+
	"do not skip what a step requires.")

var shapeLoad = loadGuide("shape", "that topic and mode", "Shaping writes proposals and\n"+
	"knowledge only: it never implements, promotes status, launches an agent, or\n"+
	"merges.")

// managedFiles are the harness entrypoints init owns, relative to the project.
func managedFiles() map[string]string {
	claude := func(name, description, hint, label, data, load string) string {
		return "---\nname: " + name + "\ndescription: " + description + "\ndisable-model-invocation: true\n" +
			"argument-hint: \"" + hint + "\"\n---\n\n<!-- " + managedMarker + " -->\n\n" +
			label + ": $ARGUMENTS\n\n" + data + "\n" + load
	}
	codex := func(name, description, data, load string) string {
		return "---\nname: " + name + "\ndescription: " + description + "\n---\n\n<!-- " + managedMarker + " -->\n\n" +
			strings.Replace(data, "The ", "In the message that invoked this skill, the ", 1) + "\n" + load
	}
	policy := "# " + managedMarker + "\npolicy:\n  allow_implicit_invocation: false\n"
	return map[string]string{
		".claude/skills/grove-work/SKILL.md":            claude("grove-work", workDescription, "G-ID [G-ID ...] [--interaction interactive|headless]", "Assignment", assignmentData, workLoad),
		".claude/skills/grove-shape/SKILL.md":           claude("grove-shape", shapeDescription, "TOPIC or G-ID [...] [--interaction interactive|headless]", "Shaping request", shapingData, shapeLoad),
		".agents/skills/grove-work/SKILL.md":            codex("grove-work", workDescription, assignmentData, workLoad),
		".agents/skills/grove-shape/SKILL.md":           codex("grove-shape", shapeDescription, shapingData, shapeLoad),
		".agents/skills/grove-work/agents/openai.yaml":  policy,
		".agents/skills/grove-shape/agents/openai.yaml": policy,
	}
}

// initStep is one planned path: what init will say about it, and the write it
// still has to do, if any.
type initStep struct {
	path, verdict, note string
	content             []byte // nil: nothing to write
	dir                 bool
}

// runInit sets up the Git checkout at cwd or --project: grove.yaml, the record
// root, a placeholder brief, and the managed harness entrypoints. It plans
// every path first and writes nothing when any path conflicts.
func runInit(cwd string, a invocation, out, errOut io.Writer) int {
	root, err := filepath.Abs(cwd)
	if err == nil && a.project != "" {
		if root = a.project; !filepath.IsAbs(root) {
			root = filepath.Join(cwd, root)
		}
		root = filepath.Clean(root)
	}
	if err == nil {
		var info os.FileInfo
		if info, err = os.Stat(root); err == nil && !info.IsDir() {
			err = fmt.Errorf("%s is not a directory", root)
		}
	}
	if err != nil {
		report(errOut, err)
		return 1
	}
	if _, err := fmt.Fprintf(errOut, "Project: %s\n", visible(root)); err != nil {
		return 1
	}
	_, prefix, err := repo.Locate(root)
	if err != nil {
		report(errOut, fmt.Errorf("init needs the top of a Git checkout: %w", err))
		return 1
	}
	if prefix != "" {
		report(errOut, fmt.Errorf("init needs the top of a Git checkout, and %s is below it, at %s; a nested grove.yaml would end discovery there", visible(root), visible(strings.TrimSuffix(prefix, "/"))))
		return 1
	}
	steps, conflicts := planInit(root)
	if len(conflicts) != 0 {
		for _, c := range conflicts {
			fmt.Fprintf(errOut, "grove: conflict %s\n", visible(c))
		}
		fmt.Fprintln(errOut, "grove: nothing was written; resolve the conflicts and run init again")
		return 1
	}
	for _, s := range steps {
		if s.dir && s.content == nil && s.verdict == "created" {
			err = os.MkdirAll(filepath.Join(root, filepath.FromSlash(s.path)), 0o755)
		} else if s.content != nil {
			full := filepath.Join(root, filepath.FromSlash(s.path))
			if err = os.MkdirAll(filepath.Dir(full), 0o755); err == nil {
				err = os.WriteFile(full, s.content, 0o644)
			}
		}
		if err != nil {
			fmt.Fprintf(errOut, "grove: %s: %s (the paths above were written; run init again after fixing this)\n", visible(s.path), err)
			return 1
		}
		if _, err := fmt.Fprintf(out, "%s %s%s\n", s.verdict, visible(s.path), s.note); err != nil {
			fmt.Fprintf(errOut, "grove: write output: %s (init wrote through %s)\n", err, visible(s.path))
			return 1
		}
	}
	fmt.Fprintln(errOut, "Next: grove check. The entrypoints run `grove` from PATH and let the agent name its branches;\n"+
		"say otherwise in AGENTS.md or CLAUDE.md, which they defer to for how grove is invoked and how\n"+
		"work and proposal branches are named.")
	return 0
}

// planInit decides each path's verdict without writing. Conflicts are the
// reasons init must not proceed.
func planInit(root string) (steps []initStep, conflicts []string) {
	lstat := func(relative string) (os.FileInfo, error) {
		return os.Lstat(filepath.Join(root, filepath.FromSlash(relative)))
	}
	conflict := func(relative, reason string) { conflicts = append(conflicts, relative+": "+reason) }
	// regular reports whether the path is absent (nil, false), a regular file
	// (info, true), or something init cannot replace, which is a conflict.
	regular := func(relative string) (os.FileInfo, bool) {
		for i, c := range relative { // a symlinked parent would carry the write outside the checkout
			if c == '/' {
				if parent, err := lstat(relative[:i]); err == nil && parent.Mode()&fs.ModeSymlink != 0 {
					conflict(relative, relative[:i]+" is a symlink")
					return nil, false
				}
			}
		}
		info, err := lstat(relative)
		switch {
		case errors.Is(err, fs.ErrNotExist):
			return nil, false
		case err != nil:
			conflict(relative, err.Error())
		case info.Mode()&fs.ModeSymlink != 0:
			conflict(relative, "is a symlink")
		case !info.Mode().IsRegular():
			conflict(relative, "is not a regular file")
		default:
			return info, true
		}
		return nil, false
	}

	recordDir, brief := "grove", "grove/brief.md"
	if _, exists := regular("grove.yaml"); exists {
		source, err := os.ReadFile(filepath.Join(root, "grove.yaml"))
		if err != nil {
			conflict("grove.yaml", err.Error())
		} else if dir, b, ds := project.ParseConfig(source); len(ds) != 0 {
			for _, d := range ds {
				conflict("grove.yaml", "exists but is not a schema 3 configuration: "+strings.TrimPrefix(d.String(), "grove.yaml: "))
			}
		} else {
			recordDir, brief = dir, b
			steps = append(steps, initStep{path: "grove.yaml", verdict: "kept", note: " (exists and validates)"})
		}
	} else {
		steps = append(steps, initStep{path: "grove.yaml", verdict: "created", content: []byte(defaultConfig)})
	}

	if info, err := lstat(recordDir); errors.Is(err, fs.ErrNotExist) {
		steps = append(steps, initStep{path: recordDir, verdict: "created", dir: true})
	} else if err != nil {
		conflict(recordDir, err.Error())
	} else if info.Mode()&fs.ModeSymlink != 0 || !info.IsDir() {
		conflict(recordDir, "must be a directory for grove.yaml: records")
	} else {
		steps = append(steps, initStep{path: recordDir, verdict: "kept", dir: true})
	}

	if brief != "" {
		if _, exists := regular(brief); exists {
			steps = append(steps, initStep{path: brief, verdict: "kept", note: " (never rewritten)"})
		} else {
			steps = append(steps, initStep{path: brief, verdict: "created", content: []byte(placeholderBrief), note: " (a placeholder that states no intent)"})
		}
	}

	files := managedFiles()
	for _, relative := range sortedKeys(files) {
		want := []byte(files[relative])
		_, exists := regular(relative)
		if !exists {
			steps = append(steps, initStep{path: relative, verdict: "created", content: want})
			continue
		}
		have, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(relative)))
		switch {
		case err != nil:
			conflict(relative, err.Error())
		case !strings.Contains(string(have), managedMarker):
			steps = append(steps, initStep{path: relative, verdict: "kept", note: " (not managed by grove init; delete it to get the managed version)"})
		case string(have) == string(want):
			steps = append(steps, initStep{path: relative, verdict: "unchanged"})
		default:
			steps = append(steps, initStep{path: relative, verdict: "updated", content: want})
		}
	}
	return steps, conflicts
}

func sortedKeys(m map[string]string) []string {
	return slices.Sorted(maps.Keys(m))
}

// versionLine names the executable, and with it the embedded workflow: the
// main module's version from the build, then the VCS revision when stamped.
func versionLine() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "grove (no build information)\n"
	}
	line := "grove " + info.Main.Version
	var revision, modified string
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			revision = s.Value
		case "vcs.modified":
			modified = s.Value
		}
	}
	if revision != "" {
		line += " (" + revision
		if modified == "true" {
			line += ", modified"
		}
		line += ")"
	}
	digest := sha256.New()
	for _, name := range []string{"docs/work-execution.md", "docs/work-shaping.md"} {
		source, err := fs.ReadFile(grove.Guides, name)
		if err != nil {
			panic(err) // both files are embedded
		}
		digest.Write(source)
	}
	return fmt.Sprintf("%s guides sha256:%x\n", line, digest.Sum(nil)[:6])
}
