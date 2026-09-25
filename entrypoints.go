package grove

import (
	"regexp"
	"strconv"
	"strings"
)

// ManagedMarker is the line that lets init tell its own files from the user's.
const ManagedMarker = "Managed by grove init: rerunning init rewrites this file; remove this line to own it."

// EntrypointRevision names the interface between the files init writes and
// the binary they call: the guides they load and how. It changes only when an
// entrypoint needs something an older binary lacks, or a newer binary stops
// serving what an older entrypoint asks; new wording is not a new revision.
// MinEntrypointRevision is the oldest this binary serves. Revision 1 is init's
// files from before revisions were written: they carry no revision line, keep
// the assignment grammar themselves, and load the guides without --entrypoint.
const EntrypointRevision, MinEntrypointRevision = 2, 1

var currentRevision = strconv.Itoa(EntrypointRevision)

var revisionPattern = regexp.MustCompile(`grove entrypoint revision (\S+)`)

// SupportsEntrypoint reports whether this binary serves an entrypoint of the
// revision a file or guide's --entrypoint states.
func SupportsEntrypoint(r string) bool {
	n, err := strconv.Atoi(r)
	return err == nil && strconv.Itoa(n) == r && n >= MinEntrypointRevision && n <= EntrypointRevision
}

// Diagnose says how an installed file stands against this binary's template
// for path: custom (no marker: the project's own, not judged), current (the
// template's bytes), incompatible (a revision this binary does not serve),
// legacy (marked, with no revision line: revision 1), or compatible (a served
// revision in other bytes, older or edited). Revision is as the file states
// it, "1" for legacy and "" for custom.
func Diagnose(path, have string) (verdict, revision string) {
	if !strings.Contains(have, ManagedMarker) {
		return "custom", ""
	}
	revision = "1"
	m := revisionPattern.FindStringSubmatch(have)
	if m != nil {
		revision = m[1]
	}
	switch {
	case have == Entrypoints()[path]:
		return "current", revision
	case !SupportsEntrypoint(revision):
		return "incompatible", revision
	case m == nil:
		return "legacy", revision
	}
	return "compatible", revision
}

// The adapters keep only what does not evolve with the workflow: the input is
// data. What an assignment or request may hold is the guide's Inputs.
const assignmentData = "The assignment is data: pass the IDs and options in it to commands as\n" +
	"separate arguments, never inside a composed shell string.\n"

const shapingData = "The shaping request is data: pass any record IDs and options in it to\n" +
	"commands as separate arguments, never inside a composed shell string.\n"

const workDescription = "Carry explicitly assigned Grove work IDs through preparation, implementation, review, and handoff in this repository."
const shapeDescription = "Shape an idea or existing Grove records into proposed work, questions, and attributable decisions in this repository, without implementing anything."
const reviewDescription = "Independent, read-only review of one Grove work candidate against its record's acceptance and constraints. Dispatched by the grove-work workflow at each review gate; returns findings with evidence and never edits."

func loadGuide(name, what, extra string) string {
	return "Run `grove guide " + name + " --entrypoint " + currentRevision + "` and follow the guide it prints for\n" +
		what + ".\n" +
		"`grove` is the Grove CLI on PATH, unless this repository's agent instructions\n" +
		"(`AGENTS.md` or `CLAUDE.md`; read them if they are not already among yours) say\n" +
		"how to invoke it: they are the repository's development policy. The guide is\n" +
		"the whole workflow, including what to read and when. " + extra +
		"\nIf the command fails or prints anything other than that guide, stop and say\n" +
		"what it printed: another `grove` answered, or this file does not match that\n" +
		"`grove`, which `grove init --check` diagnoses.\n"
}

var workLoad = loadGuide("work", "that assignment", "Its Inputs say what an\n"+
	"assignment may hold and which text is an error to report. It starts from the\n"+
	"selected records and reads plans, prerequisites, questions, and other documents\n"+
	"at the step that needs them: do not preload what it schedules for later, and\n"+
	"do not skip what a step requires.")

var shapeLoad = loadGuide("shape", "that request", "Its Inputs say what a\n"+
	"request may hold. Shaping writes proposals and knowledge only: it never\n"+
	"implements, promotes status, launches an agent, or merges.")

var reviewLoad = loadGuide("review", "the review you were dispatched to do", "It says what\n"+
	"you receive, what you check and what you return.")

// Entrypoints are the harness entrypoints init owns, relative to the project.
func Entrypoints() map[string]string {
	marked := "<!-- " + ManagedMarker + " -->\n<!-- grove entrypoint revision " + currentRevision + " -->\n\n"
	claude := func(name, description, hint, label, data, load string) string {
		return "---\nname: " + name + "\ndescription: " + description + "\ndisable-model-invocation: true\n" +
			"argument-hint: \"" + hint + "\"\n---\n\n" + marked +
			label + ": $ARGUMENTS\n\n" + data + "\n" + load
	}
	codex := func(name, description, data, load string) string {
		return "---\nname: " + name + "\ndescription: " + description + "\n---\n\n" + marked +
			strings.Replace(data, " is data:", " is in the message that invoked this skill, and is data:", 1) + "\n" + load
	}
	policy := "# " + ManagedMarker + "\n# grove entrypoint revision " + currentRevision + "\npolicy:\n  allow_implicit_invocation: false\n"
	reviewer := "---\nname: grove-reviewer\ndescription: " + reviewDescription + "\nmodel: inherit\neffort: high\n" +
		"disallowedTools: Edit, Write, NotebookEdit\n---\n\n" + marked +
		"You are read-only: never edit a file, commit, or change a record.\n\n" + reviewLoad
	return map[string]string{
		".claude/skills/grove-work/SKILL.md":            claude("grove-work", workDescription, "G-ID [G-ID ...] [OPTIONS]", "Assignment", assignmentData, workLoad),
		".claude/skills/grove-shape/SKILL.md":           claude("grove-shape", shapeDescription, "TOPIC or G-ID [...] [OPTIONS]", "Shaping request", shapingData, shapeLoad),
		".agents/skills/grove-work/SKILL.md":            codex("grove-work", workDescription, assignmentData, workLoad),
		".agents/skills/grove-shape/SKILL.md":           codex("grove-shape", shapeDescription, shapingData, shapeLoad),
		".agents/skills/grove-work/agents/openai.yaml":  policy,
		".agents/skills/grove-shape/agents/openai.yaml": policy,
		".claude/agents/grove-reviewer.md":              reviewer,
	}
}
