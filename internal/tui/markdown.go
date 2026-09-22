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
		rows[i] = clip(rows[i], w)
	}
	return rows
}

// escapeLines escapes each line of a text as safe does, keeping the line
// breaks that give Markdown its structure. Tabs become spaces, as wrap does.
func escapeLines(text string) string {
	lines := strings.Split(strings.ReplaceAll(text, "\t", "    "), "\n")
	for i, l := range lines {
		lines[i] = safe(strings.TrimRight(l, "\r"))
	}
	return strings.Join(lines, "\n")
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
