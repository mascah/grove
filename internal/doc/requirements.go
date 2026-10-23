package doc

import (
	"fmt"
	"regexp"
	"strings"
)

// Requirement is one "### Requirement: <name>" block of a maintained spec
// or a change spec.
type Requirement struct {
	Section string // ADDED, MODIFIED or REMOVED in a change spec; empty in a maintained spec
	Name    string
	Text    string // the block below its heading, trimmed
	Line    int    // of its heading
}

var (
	heading  = regexp.MustCompile(`^#{1,3} `)
	section  = regexp.MustCompile(`^## (ADDED|MODIFIED|REMOVED) Requirements\s*$`)
	shall    = regexp.MustCompile(`(?i)\bshall\b`)
	scenario = regexp.MustCompile(`(?m)^#### Scenario:`)
)

// Problems returns what r lacks of what every requirement not under REMOVED
// has: "shall" and a scenario.
func (r Requirement) Problems() []string {
	if r.Section == "REMOVED" {
		return nil
	}
	var out []string
	if !scenario.MatchString(r.Text) {
		out = append(out, fmt.Sprintf("requirement %q has no scenario", r.Name))
	}
	if !shall.MatchString(r.Text) {
		out = append(out, fmt.Sprintf("requirement %q has no \"shall\"", r.Name))
	}
	return out
}

// ChangeHeadings returns the line of each ADDED, MODIFIED or REMOVED
// Requirements heading in src outside fenced code.
func ChangeHeadings(src string) []int {
	var out []int
	fenced := false
	for i, l := range strings.Split(src, "\n") {
		if fenced = fenced != fence(l); !fenced && section.MatchString(l) {
			out = append(out, i+1)
		}
	}
	return out
}

func fence(l string) bool {
	f := strings.TrimLeft(l, " ")
	return strings.HasPrefix(f, "```") || strings.HasPrefix(f, "~~~")
}

// Requirements returns the requirement blocks of src, each running to the
// next heading of level 3 or above. A line in a fenced code block is never a
// heading, so a shell comment in an example does not end its requirement.
func Requirements(src string) []Requirement {
	var out []Requirement
	sec, cur, text, fenced := "", -1, []string(nil), false
	flush := func() {
		if cur >= 0 {
			out[cur].Text = strings.TrimSpace(strings.Join(text, "\n"))
		}
		cur, text = -1, nil
	}
	for i, l := range strings.Split(src, "\n") {
		fenced = fenced != fence(l)
		if fenced || !heading.MatchString(l) {
			text = append(text, l)
			continue
		}
		flush()
		if strings.HasPrefix(l, "## ") {
			sec = ""
			if m := section.FindStringSubmatch(l); m != nil {
				sec = m[1]
			}
		}
		if name, ok := strings.CutPrefix(l, "### Requirement: "); ok {
			out = append(out, Requirement{Section: sec, Name: strings.TrimSpace(name), Line: i + 1})
			cur = len(out) - 1
		}
	}
	flush()
	return out
}
