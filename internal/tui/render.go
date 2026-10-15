// Rendering, copied from grove-next's internal/render. Document content is
// untrusted: Markdown returns printable characters, "\n", and only the SGR
// sequences its own style uses (bold, faint, italic, underline, yellow, blue,
// reset), so nothing a file holds can move the cursor, retitle the window or
// write the clipboard. styledOnly enforces this on the output by itself;
// escapeLines makes sure the renderer has nothing to decode into a control
// character.
package tui

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"charm.land/glamour/v2"
	gansi "charm.land/glamour/v2/ansi"
	"charm.land/glamour/v2/styles"
)

// minWidth is the narrowest width Markdown wraps to: below it, margins and
// list indentation leave no room for words.
const minWidth = 20

// Markdown renders source as styled terminal text wrapped to width cells.
// Lines carry no trailing spaces or final newline.
func Markdown(source string, width int) string {
	return strings.Join(rows(source, max(width, minWidth)), "\n")
}

// Sanitize makes one line of untrusted text, such as a title, an ID or a path,
// inert for a terminal. Control characters (ESC, C0 and C1, including newline
// and tab), format characters such as bidirectional overrides, and invalid
// UTF-8 become visible escapes like \x1b and \u202e. Letters, marks and wide
// characters are kept. Sanitize is idempotent.
//
// ponytail: this also escapes the joiners inside emoji sequences and
// non-ASCII spaces; allow those individually if titles need them.
func Sanitize(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); {
		r, n := utf8.DecodeRuneInString(s[i:])
		switch {
		case r == utf8.RuneError && n == 1:
			fmt.Fprintf(&b, `\x%02x`, s[i])
		case strconv.IsPrint(r):
			b.WriteRune(r)
		default:
			quoted := strconv.QuoteRune(r)
			b.WriteString(quoted[1 : len(quoted)-1])
		}
		i += n
	}
	return b.String()
}

// style is glamour's plain style with bold headings, backticks around inline
// code and ANSI 16 accents, which follow the terminal's own theme. It never
// queries the terminal, so a test renders the same bytes a screen gets.
var style = func() gansi.StyleConfig {
	c := styles.ASCIIStyleConfig
	yes, one, hanging := true, uint(1), uint(2)
	blue, yellow := "4", "3"
	c.Document.Margin = &one
	c.Heading.Bold = &yes
	c.Heading.Color = &blue
	c.Strong.Bold = &yes
	c.Strong.BlockPrefix, c.Strong.BlockSuffix = strongOpen+strongOpen, strongClose+strongClose
	c.Emph.Italic = &yes
	c.Emph.BlockPrefix, c.Emph.BlockSuffix = emphOpen, emphClose
	c.List.Indent = &hanging
	c.Item.BlockPrefix = bullet + " "
	c.Enumeration.BlockPrefix = enumeration + " "
	c.Task.Ticked, c.Task.Unticked = task+"x] ", task+" ] "
	c.Code.Prefix, c.Code.Suffix = "`", "`"
	c.Code.Color = &yellow
	c.CodeBlock.Color = &yellow
	c.Link.Color = &blue
	c.Link.Underline = &yes
	c.LinkText.Bold = &yes
	c.HorizontalRule.Faint = &yes
	c.BlockQuote.Faint = &yes
	return c
}()

// amp stands in for "&" while the renderer runs. It is a private-use
// character, which Sanitize escapes, so no escaped line contains one.
const amp = "\ue000"

// The style writes these private-use characters, which no escaped line
// contains either, where it means the emphasis markers and list item
// prefixes, so that rows can tell them from text. Each is one cell wide, like
// the character it stands for.
//
// glamour puts an emphasis's markers around each piece of text in it, and
// none around a code span, so "**a `x` b**" comes out as "**a **`x`** b**";
// rows joins the pieces. It wraps every line of a list item to the list's
// indent; style indents lists by the width of a bullet and rows moves each
// item's first line back by as much, so that the lines after it hang under
// its text.
//
// ponytail: they hang under the third cell of "1. " and "[ ] " as well; a
// per-kind indent needs glamour to indent lists by kind.
const (
	strongOpen, strongClose = "\ue001", "\ue002" // "*", twice each
	emphOpen, emphClose     = "\ue003", "\ue004" // "*"
	bullet                  = "\ue005"           // "•"
	enumeration             = "\ue006"           // "." after the number
	task                    = "\ue007"           // "[" of "[ ]" and "[x]"
)

// restore puts back what amp and the style's own characters stand for,
// including an "&" glamour percent-encoded in the path or fragment of a
// relative link.
//
// ponytail: a document that spells out "%EE%80%80" itself is shown an "&"
// there too.
var restore = strings.NewReplacer(amp, "&", "%EE%80%80", "&",
	strongOpen, "*", strongClose, "*", emphOpen, "*", emphClose, "*", bullet, "•", enumeration, ".", task, "[")

var (
	// joined matches where one emphasis's markers close and open again with
	// only styles and code spans between: between two pieces of the same
	// emphasis. The pieces are joined by dropping both markers.
	//
	// ponytail: "**a**`x`**b**", two emphases around a code span with no
	// space, reads as one, "**a`x`b**".
	joined = regexp.MustCompile(strongClose + strongClose + `((?:` + anyStyle + `|` + codeSpan + `)*)` + strongOpen + strongOpen +
		`|` + emphClose + `((?:` + anyStyle + `|` + codeSpan + `)*)` + emphOpen)

	// item matches the start of a row that starts a list item: what comes
	// before the list's indent, the indent, and the start of the item's
	// prefix.
	item = regexp.MustCompile(`^((?:` + anyStyle + `|[ |])*?)  ([0-9]*[` + bullet + enumeration + task + `])`)
)

// The sequences a row may hold around text, for joined.
const (
	anyStyle = `\x1b\[[0-9;]*m`
	codeSpan = `\x1b\[33m[^\x1b\n]*\x1b\[m`
)

var (
	// osc8 matches the terminal hyperlinks glamour puts around links. They
	// are removed: no sequence that carries file-provided text may reach the
	// terminal.
	osc8 = regexp.MustCompile(`\x1b\]8;[^\x07\x1b]*(?:\x07|\x1b\\)`)

	// sgr matches the sequences style can produce and no others: a reset, or
	// any of bold (1), faint (2), italic (3), underline (4), yellow (33) and
	// blue (34). Conceal, other colors and every other sequence do not match.
	sgr = regexp.MustCompile(`\x1b\[(?:(?:[1-4]|3[34])(?:;(?:[1-4]|3[34]))*)?m`)
)

// rows renders source at width w and returns its lines holding nothing but
// text and glamour's styles, without trailing spaces or blank lines at either
// end. A line may be wider than w when it cannot wrap.
func rows(source string, w int) []string {
	r, _ := glamour.NewTermRenderer(glamour.WithStyles(style), glamour.WithWordWrap(w)) // fails only for a bad style file
	out, _ := r.Render(escapeLines(source))                                             // renders into a buffer, which does not fail
	out = joined.ReplaceAllString(osc8.ReplaceAllString(out, ""), "$1$2")
	rows := strings.Split(out, "\n")
	for i, row := range rows {
		row = item.ReplaceAllString(row, "$1$2")
		rows[i] = strings.TrimRight(styledOnly(restore.Replace(row)), " ")
	}
	// Blank rows at either end are margins, or what is left of a comment.
	for len(rows) > 1 && rows[0] == "" {
		rows = rows[1:]
	}
	for len(rows) > 1 && rows[len(rows)-1] == "" {
		rows = rows[:len(rows)-1]
	}
	return rows
}

// escapeLines escapes each line as Sanitize does, keeping the line breaks
// that give Markdown its structure, and replaces every "&" with amp; rows puts
// them back after rendering. Tabs become spaces.
//
// Without an "&" the renderer sees no character reference in any form:
// glamour decodes them with html.UnescapeString, which accepts "&#27" with
// no semicolon, and joins text from several nodes before decoding, so no
// pattern over the source can find them all. References are therefore shown
// as written, "&amp;" included.
//
// ponytail: an "&" also no longer counts as punctuation next to emphasis
// markers, after a backslash, or in a bare URL being linked, which ends
// there (its text is kept). glamour does not expose its parser; if a later
// version does, give text nodes back their "&" in an AST transformer.
func escapeLines(text string) string {
	lines := strings.Split(strings.ReplaceAll(text, "\t", "    "), "\n")
	for i, l := range lines {
		lines[i] = strings.ReplaceAll(Sanitize(strings.TrimRight(l, "\r")), "&", amp)
	}
	return strings.Join(lines, "\n")
}

// styledOnly enforces the package invariant on one rendered row: the
// sequences sgr matches pass, and every other byte is escaped as Sanitize
// does, an escape byte that starts anything else included. It does not rely
// on escapeLines or on the renderer.
func styledOnly(s string) string {
	var b strings.Builder
	last := 0
	for _, loc := range sgr.FindAllStringIndex(s, -1) {
		b.WriteString(Sanitize(s[last:loc[0]]))
		b.WriteString(s[loc[0]:loc[1]])
		last = loc[1]
	}
	b.WriteString(Sanitize(s[last:]))
	return b.String()
}
