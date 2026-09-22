package tui

import (
	"strconv"
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
)

func TestBodyDropsTheFrontmatter(t *testing.T) {
	t.Parallel()
	for _, c := range [][2]string{
		{"---\nid: G-001\n---\n\n# Hi\n", "\n# Hi\n"},
		{"\xef\xbb\xbf---\r\nid: G-001\r\n---\r\nbody\r\n", "body\r\n"},
		{"---\n---\nempty\n", "empty\n"},
		{"no frontmatter\n", "no frontmatter\n"},
		{"---\nnever closed\n", "---\nnever closed\n"},
	} {
		if got := body([]byte(c[0])); got != c[1] {
			t.Errorf("body(%q) = %q, want %q", c[0], got, c[1])
		}
	}
}

// Every byte from a record is escaped before glamour renders it, so the rows
// carry glamour's styles and nothing a file could have planted.
func TestRenderedRowsHoldOnlyGlamourStyles(t *testing.T) {
	t.Parallel()
	md := "## Outcome\n\nA **bold** claim with \x1b[31mred\x1b[m, a title \x1b]0;pwned\x07, C1 \u009b31m, " +
		"an override \u202e and a [link](G-093-current-view-plan.md) to http://example.com/x.\n\n" +
		"- one 日本語の長いタイトルがここにあります\n- two\n\n```sh\ngo run ./cmd/grove\n```\n\n| a | b |\n|---|---|\n| 1 | 2 |\n" +
		// Character references decode after the escaping, in text, code
		// spans, headings, HTML and cells: they must stay literal.
		"\n&#27;]52;c;cHduZWQ=&#7; &#x1b;[8mhidden &#x202e;bidi &rlm;named a &amp; b\n\n`code &#x1b;[31m`\n\n" +
		"### &#x1b;]2;title&#x7;\n\n<div>&#27;[31m</div>\n\n| &#x1b;x | &#X1B;y |\n|---|---|\n| &#7; | z |\n"
	for _, w := range []int{20, 40, 80} {
		rows := render(md, w)
		if len(rows) < 8 {
			t.Fatalf("width %d: %d rows", w, len(rows))
		}
		text := ansi.Strip(strings.Join(rows, "\n"))
		for _, row := range rows {
			if got := ansi.StringWidth(row); got != w {
				t.Errorf("width %d: row is %d cells: %q", w, got, row)
			}
			if rest := sgr.ReplaceAllString(row, ""); strings.ContainsRune(rest, 0x1b) {
				t.Errorf("width %d: a sequence besides a style: %q", w, row)
			}
			if strings.Contains(row, "\x1b[31m") || strings.Contains(row, "\x1b[8m") {
				t.Errorf("width %d: a planted style: %q", w, row)
			}
			for _, r := range ansi.Strip(row) {
				if !strconv.IsPrint(r) && r != ' ' {
					t.Errorf("width %d: control %q in %q", w, r, row)
				}
			}
		}
		for _, want := range []string{"## Outcome", `\x1b[31mred`, `\x1b]0;pwned\a`, `\u009b31m`, `\u202e`, "link", "G-093-", "plan.md", "go run", "日本語",
			"&#27;]52;c;", "[8mhidden", "&#x202e;bidi", "&rlm;named", "&amp;", "code &#x1b;", "&#x1b;]2;title", "&#27;[31m", "&#x1b;x", "&#X1B;y"} {
			if !strings.Contains(text, want) {
				t.Errorf("width %d: %q missing from\n%s", w, want, text)
			}
		}
	}
	if rows := render("plain", 0); len(rows) != 1 || ansi.StringWidth(rows[0]) != 1 {
		t.Errorf("width 0 renders one cell: %q", rows)
	}
}

func TestRenderedIsCachedPerKeyAndWidth(t *testing.T) {
	t.Parallel()
	m := &Model{}
	a := m.rendered("k", "# A", 30)
	b := m.rendered("k", "# B", 30) // same key: the cached rows
	c := m.rendered("k", "# A", 31)
	if ansi.Strip(a[0]) != ansi.Strip(b[0]) || len(m.md) != 2 || ansi.StringWidth(c[0]) != 31 {
		t.Errorf("cache: %q %q %q, %d entries", a, b, c, len(m.md))
	}
}
