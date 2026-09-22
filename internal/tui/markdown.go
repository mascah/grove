package tui

import (
	"regexp"
	"strconv"
	"strings"

	"charm.land/glamour/v2"
	gansi "charm.land/glamour/v2/ansi"
	"charm.land/glamour/v2/styles"
	"github.com/charmbracelet/x/ansi"
)

// body returns a record's Markdown after its frontmatter: what the detail
// renders. Bytes that are not a record (no frontmatter) are returned whole.
func body(source []byte) string {
	s := strings.TrimPrefix(string(source), "\xef\xbb\xbf")
	lines := strings.Split(s, "\n")
	if strings.TrimRight(lines[0], "\r") != "---" {
		return s
	}
	for i := 1; i < len(lines); i++ {
		if strings.TrimRight(lines[i], "\r") == "---" {
			return strings.Join(lines[i+1:], "\n")
		}
	}
	return s
}

// style is glamour's plain style with bold headings and ANSI 16 accents,
// which follow the terminal's own theme. It never queries the terminal, so a
// test renders the same bytes a screen gets. Meaning stays in the text: a
// heading keeps its # prefix, code its backticks.
var style = func() gansi.StyleConfig {
	c := styles.ASCIIStyleConfig
	yes, one := true, uint(1)
	c.Document.Margin = &one
	c.Heading.Bold = &yes
	c.Heading.Color = ptr("4")
	c.Strong.Bold = &yes
	c.Emph.Italic = &yes
	c.Code.Color = ptr("3")
	c.CodeBlock.Color = ptr("3")
	c.Link.Color = ptr("4")
	c.Link.Underline = &yes
	c.LinkText.Bold = &yes
	c.HorizontalRule.Faint = &yes
	c.BlockQuote.Faint = &yes
	return c
}()

func ptr(s string) *string { return &s }

// osc8 matches the terminal hyperlinks glamour puts around links.
var osc8 = regexp.MustCompile(`\x1b\]8;[^\x07\x1b]*(?:\x07|\x1b\\)`)

// sgr matches one of glamour's own styles: the only sequences a rendered
// row may hold.
var sgr = regexp.MustCompile(`\x1b\[[0-9;]*m`)

// reference matches an HTML character reference, which Markdown decodes
// after the text was escaped: `&#x1b;` would come out of glamour as a real
// escape byte. escapeLines makes the ampersand literal instead.
var reference = regexp.MustCompile(`&(#[0-9]+|#[xX][0-9a-fA-F]+|[a-zA-Z][a-zA-Z0-9]*);`)

// render turns Markdown from a record into rows of exactly w cells. The text
// is escaped before glamour sees it, so the only sequences in the rows are
// glamour's own styles. Its terminal hyperlinks are removed: G-017 lets no
// file-provided sequence reach the terminal, and the owner chose plain
// links on 2026-09-22. Markdown glamour cannot render is shown wrapped.
func render(markdown string, w int) []string {
	w = max(w, 1)
	r, err := glamour.NewTermRenderer(glamour.WithStyles(style), glamour.WithWordWrap(w))
	var out string
	if err == nil {
		out, err = r.Render(escapeLines(markdown))
	}
	if err != nil {
		return wrapAll(markdown, w)
	}
	rows := strings.Split(strings.Trim(osc8.ReplaceAllString(out, ""), "\n"), "\n")
	for i := range rows {
		rows[i] = clip(styledOnly(rows[i]), w)
	}
	return rows
}

// escapeLines escapes each line of a text as safe does, keeping the line
// breaks that give Markdown its structure, and makes character references
// literal. Tabs become spaces, as wrap does.
func escapeLines(text string) string {
	lines := strings.Split(strings.ReplaceAll(text, "\t", "    "), "\n")
	for i, l := range lines {
		lines[i] = reference.ReplaceAllString(safe(strings.TrimRight(l, "\r")), "&amp;$1;")
	}
	return strings.Join(lines, "\n")
}

// styledOnly is the second layer behind escapeLines: whatever the renderer
// emits, only its styles pass, and every other byte is escaped as safe does.
func styledOnly(s string) string {
	var b strings.Builder
	last := 0
	for _, loc := range sgr.FindAllStringIndex(s, -1) {
		b.WriteString(safe(s[last:loc[0]]))
		b.WriteString(s[loc[0]:loc[1]])
		last = loc[1]
	}
	b.WriteString(safe(s[last:]))
	return b.String()
}

// clip fits an already styled row to w cells; line does the same for text
// that still needs escaping.
func clip(s string, w int) string {
	s = ansi.Truncate(strings.TrimRight(s, " "), w, "…")
	return s + strings.Repeat(" ", max(w-ansi.StringWidth(s), 0))
}

// rendered returns the rows for key's Markdown at width w, rendering once
// per key and width for the current result: a scroll never re-renders.
func (m *Model) rendered(key, markdown string, w int) []string {
	k := key + "\x00" + strconv.Itoa(w)
	rows, ok := m.md[k]
	if !ok {
		if m.md == nil {
			m.md = map[string][]string{}
		}
		rows = render(markdown, w)
		m.md[k] = rows
	}
	return rows
}
