package doc

import (
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
	heading = regexp.MustCompile(`^#{1,3} `)
	section = regexp.MustCompile(`^## (ADDED|MODIFIED|REMOVED) Requirements\s*$`)
)

// Requirements returns the requirement blocks of src, each running to the
// next heading of level 3 or above.
func Requirements(src string) []Requirement {
	var out []Requirement
	sec, cur, text := "", -1, []string(nil)
	flush := func() {
		if cur >= 0 {
			out[cur].Text = strings.TrimSpace(strings.Join(text, "\n"))
		}
		cur, text = -1, nil
	}
	for i, l := range strings.Split(src, "\n") {
		if !heading.MatchString(l) {
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
