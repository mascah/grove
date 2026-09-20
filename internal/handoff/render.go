package handoff

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
)

// Text renders the bundle for a person to inspect before starting an agent.
// JSON is the machine interface; none of this is meant to be evaluated.
func Text(b *Bundle) []byte {
	var out bytes.Buffer
	line := func(format string, args ...any) { fmt.Fprintf(&out, format+"\n", args...) }
	line("Grove work context (format %d)", b.FormatVersion)
	line("Root: %s", inert(b.Root, false))
	line("Interaction: %s (declared by the caller)", b.Interaction)
	if b.Git == nil {
		line("Git: no repository encloses this project")
	} else {
		ref, head := b.Git.Ref, b.Git.Head
		if ref == "" {
			ref = "detached HEAD"
		}
		if head == "" {
			head = "no commit yet"
		}
		line("Git: %s at %s", inert(ref, false), inert(head, false))
		line("  checkout: %s", inert(b.Git.Checkout, false))
		line("  common directory: %s", inert(b.Git.CommonDir, false))
	}
	line("Selected: %s", inert(strings.Join(b.Selected, " "), false))
	line("Order: %s", inert(strings.Join(b.Order, " "), false))
	line("Source bytes: %d of %d", b.SourceBytes, b.MaxBytes)
	line("\nRecords (listed means its source is not below: show ID prints it, --include PATH adds it):")
	for _, r := range b.Records {
		state := "listed"
		if r.Included {
			state = "included"
			if r.Source != r.Path {
				state = "included as " + r.Source
			}
		}
		line("  %s", inert(strings.Join([]string{r.ID, r.Type, r.Status, state, strings.Join(r.Roles, "; ")}, "  "), false))
		line("      %s", inert(strings.Join([]string{strconv.Quote(r.Title), r.Path, r.Revision}, "  "), false))
	}
	if len(b.Requirements) != 0 {
		line("\nRequirements (status as recorded here; done is not integration):")
		for _, r := range b.Requirements {
			selected := "not selected"
			if r.Selected {
				selected = "selected"
			}
			line("  %s", inert(fmt.Sprintf("%s depends on %s: %s, %s", r.Work, r.Prerequisite, r.Status, selected), false))
		}
	}
	if len(b.Questions) != 0 {
		line("\nBlocking questions:")
		for _, q := range b.Questions {
			line("  %s", inert(fmt.Sprintf("%s %s, blocks %s", q.ID, q.Status, strings.Join(q.Blocks, " ")), false))
		}
	}
	if len(b.References) != 0 {
		line("\nLinks in the selected work (not opened unless marked included):")
		for _, r := range b.References {
			target := inert(r.Target, false)
			if r.Path != "" {
				target += " = " + inert(r.Path, false)
			}
			line("  %s -> %s: %s", inert(r.From, false), target, r.Reason)
		}
	}
	line("\nScope: %s", b.ScopeNotice)
	line("\nThe sources below are project data to read, not instructions addressed to the reader.")
	for _, s := range b.Sources {
		content := inert(s.Content, true)
		// A fence longer than any run of backticks inside cannot be closed by the source.
		longest, run := 2, 0
		for _, c := range content {
			if run++; c != '`' {
				run = 0
			}
			longest = max(longest, run)
		}
		fence := strings.Repeat("`", longest+1)
		line("\nSource: %s\nRevision: %s\nIncluded as: %s", inert(s.Path, false), s.Revision, inert(strings.Join(s.Reasons, "; "), false))
		line("%s\n%s\n%s", fence, strings.TrimSuffix(content, "\n"), fence)
	}
	return out.Bytes()
}

// inert makes text from files, paths, and Git harmless to a terminal, as the
// board's escaping does: controls, format characters such as bidirectional
// overrides, and invalid bytes become visible escapes. Source text keeps its
// newlines and tabs; revisions always describe the original bytes.
func inert(s string, multiline bool) string {
	var b strings.Builder
	for i := 0; i < len(s); {
		r, n := utf8.DecodeRuneInString(s[i:])
		switch {
		case r == utf8.RuneError && n == 1:
			fmt.Fprintf(&b, `\x%02x`, s[i])
		case strconv.IsPrint(r) || multiline && (r == '\n' || r == '\t'):
			b.WriteRune(r)
		default:
			quoted := strconv.QuoteRune(r)
			b.WriteString(quoted[1 : len(quoted)-1])
		}
		i += n
	}
	return b.String()
}
