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
		line("Git: %s at %s", inert(ref, false), head)
		line("  checkout: %s", inert(b.Git.Checkout, false))
		line("  common directory: %s", inert(b.Git.CommonDir, false))
	}
	line("Selected: %s", strings.Join(b.Selected, " "))
	line("Order: %s", strings.Join(b.Order, " "))
	line("Source bytes: %d of %d", b.SourceBytes, b.MaxBytes)
	line("\nRecords:")
	for _, r := range b.Records {
		role := "context"
		if r.Selected {
			role = "selected"
		}
		line("  %s  %s  %s  %s  %s", r.ID, r.Type, r.Status, role, inert(r.Path, false))
	}
	if len(b.Requirements) != 0 {
		line("\nRequirements (status as recorded here; done is not integration):")
		for _, r := range b.Requirements {
			selected := "not selected"
			if r.Selected {
				selected = "selected"
			}
			line("  %s depends on %s: %s, %s", r.Work, r.Prerequisite, r.Status, selected)
		}
	}
	if len(b.Questions) != 0 {
		line("\nBlocking questions:")
		for _, q := range b.Questions {
			line("  %s %s, blocks %s", q.ID, q.Status, strings.Join(q.Blocks, " "))
		}
	}
	if len(b.References) != 0 {
		line("\nLinks not included:")
		for _, r := range b.References {
			line("  %s -> %s: %s", inert(r.From, false), inert(r.Target, false), r.Reason)
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
