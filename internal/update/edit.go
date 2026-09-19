package update

import (
	"bytes"
	"errors"
	"fmt"
	"slices"
	"strings"
	"unicode/utf8"

	"go.yaml.in/yaml/v3"
)

// change is one frontmatter edit: set key to the canonical YAML text value,
// or remove the key. Keys absent from the file are appended in change order.
type change struct {
	key, value string
	remove     bool
}

func set(key, value string) change { return change{key: key, value: value} }
func unset(key string) change      { return change{key: key, remove: true} }

// span is a byte range of the frontmatter to replace.
type span struct {
	start, end int
	text       string
}

// Edit applies changes to one record's source by replacing only the byte
// ranges of the edited entries. Everything else, including BOM, line endings,
// comments outside the edited values, key order, and the body, is retained.
// Positions come from yaml.v3's node lines and columns; value ends are found
// by a scanner per style, bounded by the next key. Any byte that does not
// match the expected syntax at a computed position refuses the edit.
func Edit(source []byte, changes []change) ([]byte, error) {
	fmStart, fmEnd, newline, err := frontmatter(source)
	if err != nil {
		return nil, err
	}
	e := &editor{fm: source[fmStart:fmEnd]}
	var doc yaml.Node
	if err := yaml.Unmarshal(e.fm, &doc); err != nil {
		return nil, fmt.Errorf("frontmatter: %w", err)
	}
	if len(doc.Content) != 1 || doc.Content[0].Kind != yaml.MappingNode {
		return nil, errors.New("frontmatter: expected a YAML mapping")
	}
	mapping := doc.Content[0]
	e.flow = mapping.Style&yaml.FlowStyle != 0
	e.indent = mapping.Column - 1
	e.end = len(e.fm)
	if e.flow {
		off, err := e.offset(mapping.Line, mapping.Column)
		if err != nil {
			return nil, err
		}
		if e.end, err = e.scanFlow(off); err != nil {
			return nil, fmt.Errorf("frontmatter: %w", err)
		}
		e.end-- // the closing brace
	}
	for i := 0; i < len(mapping.Content); i += 2 {
		k, v := mapping.Content[i], mapping.Content[i+1]
		keyOff, err := e.offset(k.Line, k.Column)
		if err != nil {
			return nil, err
		}
		valOff, err := e.offset(v.Line, v.Column)
		if err != nil {
			return nil, err
		}
		if !e.keyStartsAt(k, keyOff) {
			return nil, fmt.Errorf("frontmatter: cannot locate key %s", k.Value)
		}
		e.entries = append(e.entries, entry{key: k, value: v, start: e.entryStart(keyOff), keyOff: keyOff, valOff: valOff})
	}
	var spans []span
	var appends []change
	var removed []int // flow entries, whose separators are planned together
	for _, c := range changes {
		idx := -1
		for i, en := range e.entries {
			if en.key.Value == c.key {
				idx = i
			}
		}
		switch {
		case idx < 0 && !c.remove:
			appends = append(appends, c)
		case idx < 0:
		case c.remove && e.flow:
			removed = append(removed, idx)
		default:
			s, err := e.edit(idx, c)
			if err != nil {
				return nil, err
			}
			spans = append(spans, s)
		}
	}
	if e.flow {
		planned, err := e.planFlow(removed, appends)
		if err != nil {
			return nil, err
		}
		spans = append(spans, planned...)
	} else if len(appends) != 0 {
		s, err := e.append(appends)
		if err != nil {
			return nil, err
		}
		spans = append(spans, s)
	}
	for i := range spans {
		for j := range spans {
			if i != j && spans[i].start < spans[j].end && spans[j].start < spans[i].end {
				return nil, errors.New("frontmatter: overlapping edits")
			}
		}
	}
	// Apply from the end so earlier offsets stay valid.
	slices.SortFunc(spans, func(a, b span) int { return b.start - a.start })
	out := append([]byte{}, e.fm...)
	for _, s := range spans {
		text := strings.ReplaceAll(s.text, "\n", newline)
		out = append(out[:s.start], append([]byte(text), out[s.end:]...)...)
	}
	result := make([]byte, 0, len(source)+len(out)-len(e.fm))
	result = append(result, source[:fmStart]...)
	result = append(result, out...)
	result = append(result, source[fmEnd:]...)
	return result, nil
}

type entry struct {
	key, value            *yaml.Node
	start, keyOff, valOff int // start is the explicit-key "?" when there is one
}

type editor struct {
	fm      []byte
	flow    bool
	indent  int // the mapping's indentation
	end     int // where the mapping's entries end
	entries []entry
}

// frontmatter returns the byte range between the opening and closing --- lines
// and the newline style of the opening line.
func frontmatter(source []byte) (start, end int, newline string, err error) {
	pos := 0
	if bytes.HasPrefix(source, []byte("\ufeff")) {
		pos = 3
	}
	nl := bytes.IndexByte(source[pos:], '\n')
	if nl < 0 || strings.TrimRight(string(source[pos:pos+nl]), " \t\r") != "---" {
		return 0, 0, "", errors.New("frontmatter: expected an opening --- line")
	}
	newline = "\n"
	if bytes.HasSuffix(source[pos:pos+nl], []byte("\r")) {
		newline = "\r\n"
	}
	start = pos + nl + 1
	for line := start; line <= len(source); {
		nl := bytes.IndexByte(source[line:], '\n')
		lineEnd := len(source)
		if nl >= 0 {
			lineEnd = line + nl
		}
		if strings.TrimRight(string(source[line:lineEnd]), " \t\r") == "---" {
			return start, line, newline, nil
		}
		if nl < 0 {
			break
		}
		line = lineEnd + 1
	}
	return 0, 0, "", errors.New("frontmatter: missing closing --- line")
}

// offset converts yaml.v3's 1-based line and rune column to a byte offset.
func (e *editor) offset(line, column int) (int, error) {
	pos := 0
	for l := 1; l < line; l++ {
		nl := bytes.IndexByte(e.fm[pos:], '\n')
		if nl < 0 {
			return 0, errors.New("frontmatter: position beyond end")
		}
		pos += nl + 1
	}
	for c := 1; c < column; c++ {
		if pos >= len(e.fm) || e.fm[pos] == '\n' {
			return 0, errors.New("frontmatter: position beyond line end")
		}
		_, w := utf8.DecodeRune(e.fm[pos:])
		pos += w
	}
	return pos, nil
}

func (e *editor) keyStartsAt(k *yaml.Node, off int) bool {
	rest := e.fm[off:]
	switch k.Style {
	case yaml.DoubleQuotedStyle:
		return len(rest) > 0 && rest[0] == '"'
	case yaml.SingleQuotedStyle:
		return len(rest) > 0 && rest[0] == '\''
	default:
		return k.Style&yaml.TaggedStyle == 0 && k.Anchor == "" && bytes.HasPrefix(rest, []byte(k.Value)) ||
			bytes.HasPrefix(rest, []byte(k.Value[:0]+string(e.fm[off:e.skipProperties(off)])))
	}
}

// skipProperties steps over explicit tag and anchor tokens, which yaml.v3
// positions a node at, to the value or key text itself.
func (e *editor) skipProperties(off int) int {
	for off < len(e.fm) && (e.fm[off] == '!' || e.fm[off] == '&') {
		for off < len(e.fm) && !isSpace(e.fm[off]) {
			off++
		}
		for off < len(e.fm) && (e.fm[off] == ' ' || e.fm[off] == '\t') {
			off++
		}
	}
	return off
}

// entryStart returns the offset of the "?" indicator before an explicit key
// at keyOff, or keyOff itself for an ordinary key.
func (e *editor) entryStart(keyOff int) int {
	j := keyOff
	for j > 0 && (e.fm[j-1] == ' ' || e.fm[j-1] == '\t') {
		j--
	}
	if j > 0 && j < keyOff && e.fm[j-1] == '?' && (j == 1 || isSpace(e.fm[j-2]) || e.fm[j-2] == '{' || e.fm[j-2] == ',') {
		return j - 1
	}
	return keyOff
}

// bound is the first offset an entry's value may not reach: the next entry, or
// the mapping's end.
func (e *editor) bound(idx int) int {
	if idx+1 < len(e.entries) {
		return e.entries[idx+1].start
	}
	return e.end
}

// entryEnd returns the offset just past the value of entry idx.
func (e *editor) entryEnd(idx int) (int, error) {
	en := e.entries[idx]
	end, err := e.valueEnd(en.value, en.valOff, e.flow, e.indent)
	if err != nil {
		return 0, fmt.Errorf("frontmatter: %s: %w", en.key.Value, err)
	}
	if end <= en.valOff || end > e.bound(idx) {
		return 0, fmt.Errorf("frontmatter: %s: cannot safely identify the value span", en.key.Value)
	}
	return end, nil
}

// edit plans a replaced value, or a removed entry of a block mapping.
func (e *editor) edit(idx int, c change) (span, error) {
	en := e.entries[idx]
	end, err := e.entryEnd(idx)
	if err != nil {
		return span{}, err
	}
	if c.remove {
		start := bytes.LastIndexByte(e.fm[:en.start], '\n') + 1
		lineEnd := bytes.IndexByte(e.fm[end:], '\n')
		if lineEnd < 0 {
			return span{start, len(e.fm), ""}, nil
		}
		return span{start, end + lineEnd + 1, ""}, nil
	}
	if en.value.Line == en.key.Line {
		return span{en.valOff, end, c.value}, nil
	}
	after, err := e.afterColon(en)
	if err != nil {
		return span{}, err
	}
	if after > en.valOff {
		return span{}, fmt.Errorf("frontmatter: %s: cannot locate the key's colon", c.key)
	}
	if !bytes.Contains(e.fm[after:en.valOff], []byte("#")) {
		return span{after, end, " " + c.value}, nil
	}
	// Comments sit between the colon and the value: keep them and put the new
	// value where the old one began. Only a block sequence may start at its
	// key's indentation, so any other replacement is indented further there.
	pad := ""
	if column := en.valOff - (bytes.LastIndexByte(e.fm[:en.valOff], '\n') + 1); !e.flow && column <= e.indent {
		pad = strings.Repeat(" ", e.indent+2-column)
	}
	return span{en.valOff, end, pad + c.value}, nil
}

// planFlow plans the removed entries and appended fields of a flow mapping
// together, so each separator is removed at most once and an append knows
// whether a separator survives before it. An entry takes the comma that
// follows it; the last entry instead takes the comma before the removed run it
// ends, when only whitespace separates them. A comma hidden behind a comment
// stays as a trailing comma, which YAML allows. A removed entry's inline
// comment goes with it; standalone comment lines stay.
func (e *editor) planFlow(removed []int, appends []change) ([]span, error) {
	parts := make([]string, len(appends))
	for i, c := range appends {
		parts[i] = c.key + ": " + c.value
	}
	added := strings.Join(parts, ", ")
	if len(e.entries) == 0 {
		if added == "" {
			return nil, nil
		}
		return []span{{e.end, e.end, added}}, nil
	}
	last := len(e.entries) - 1
	slices.Sort(removed)
	removed = slices.Compact(removed)
	var spans []span
	for n, idx := range removed {
		end, err := e.entryEnd(idx)
		if err != nil {
			return nil, err
		}
		s := span{start: e.entries[idx].start, end: end}
		after := e.skipBlanks(end)
		following := after < len(e.fm) && e.fm[after] == ','
		if following {
			after = e.skipBlanks(after + 1)
			s.end = after
		}
		if after < len(e.fm) && e.fm[after] == '#' {
			for after < len(e.fm) && e.fm[after] != '\n' && e.fm[after] != '\r' {
				after++
			}
			s.end = after
		}
		consumed := false
		if !following {
			// ponytail: a separator on a later line than its entry's value is
			// refused; scan comments and line breaks for it if such files appear.
			if idx != last {
				return nil, fmt.Errorf("frontmatter: %s: cannot locate the separator after the entry", e.entries[idx].key.Value)
			}
			first := n // of the removed run that ends the mapping
			for first > 0 && removed[first-1] == removed[first]-1 {
				first--
			}
			j := e.entries[removed[first]].start
			for j > 0 && isSpace(e.fm[j-1]) {
				j--
			}
			if j > 0 && e.fm[j-1] == ',' {
				consumed = true
				if first == n {
					s.start = j - 1
				} else {
					spans[first].start = j - 1
				}
			}
		}
		if idx == last && added != "" {
			switch {
			case following:
				s.text = added + ", "
			case consumed:
				s.text = ", " + added
			default:
				s.text = added
			}
			added = ""
		}
		if s.text == "" {
			s = e.wholeLine(s)
		}
		spans = append(spans, s)
	}
	if added != "" {
		end, err := e.entryEnd(last)
		if err != nil {
			return nil, err
		}
		spans = append(spans, span{end, end, ", " + added})
	}
	return spans, nil
}

func (e *editor) skipBlanks(i int) int {
	for i < len(e.fm) && (e.fm[i] == ' ' || e.fm[i] == '\t') {
		i++
	}
	return i
}

// wholeLine widens a removal that leaves only whitespace on its line to the
// line itself, so no blank line remains.
func (e *editor) wholeLine(s span) span {
	lineStart := bytes.LastIndexByte(e.fm[:s.start], '\n') + 1
	end := s.end
	if end < len(e.fm) && e.fm[end] == '\r' {
		end++
	}
	if len(bytes.TrimLeft(e.fm[lineStart:s.start], " \t")) == 0 && end < len(e.fm) && e.fm[end] == '\n' {
		return span{lineStart, end + 1, ""}
	}
	return s
}

// afterColon returns the offset just past the key's colon, for values that
// start on a later line than their key. An explicit key's colon is on a later
// line than the key, after optional comments.
func (e *editor) afterColon(en entry) (int, error) {
	var err error
	end := e.skipProperties(en.keyOff)
	switch en.key.Style &^ yaml.TaggedStyle {
	case yaml.DoubleQuotedStyle:
		end, err = e.scanDoubleQuoted(end)
	case yaml.SingleQuotedStyle:
		end, err = e.scanSingleQuoted(end)
	default:
		end += len(en.key.Value)
	}
	if err != nil {
		return 0, err
	}
	for end < len(e.fm) && (isSpace(e.fm[end]) || e.fm[end] == '#') {
		if e.fm[end] == '#' {
			for end < len(e.fm) && e.fm[end] != '\n' {
				end++
			}
			continue
		}
		end++
	}
	if end >= len(e.fm) || e.fm[end] != ':' {
		return 0, fmt.Errorf("frontmatter: %s: cannot locate the key's colon", en.key.Value)
	}
	return end + 1, nil
}

// append adds fields after the last line of a block mapping.
func (e *editor) append(appends []change) (span, error) {
	var b strings.Builder
	if len(e.fm) != 0 && e.fm[len(e.fm)-1] != '\n' {
		return span{}, errors.New("frontmatter: mapping does not end with a newline")
	}
	for _, c := range appends {
		fmt.Fprintf(&b, "%s%s: %s\n", strings.Repeat(" ", e.indent), c.key, c.value)
	}
	return span{len(e.fm), len(e.fm), b.String()}, nil
}

// valueEnd returns the offset just past a value's syntax, excluding trailing
// whitespace and any comment. indent is the enclosing block's indentation.
func (e *editor) valueEnd(v *yaml.Node, start int, flow bool, indent int) (int, error) {
	// A replaced value drops its explicit tag or anchor along with the old
	// text; the reader accepts no aliases, so nothing can refer to the anchor.
	start = e.skipProperties(start)
	if start >= len(e.fm) {
		return 0, errors.New("value beyond end")
	}
	switch v.Kind {
	case yaml.ScalarNode:
		switch v.Style &^ yaml.TaggedStyle {
		case yaml.DoubleQuotedStyle:
			return e.scanDoubleQuoted(start)
		case yaml.SingleQuotedStyle:
			return e.scanSingleQuoted(start)
		case yaml.LiteralStyle, yaml.FoldedStyle:
			if e.fm[start] != '|' && e.fm[start] != '>' {
				return 0, errors.New("expected a block scalar indicator")
			}
			return e.scanBlockScalar(start, indent), nil
		case 0:
			end := e.scanPlain(start, flow, indent)
			raw := strings.Join(strings.Fields(string(e.fm[start:end])), " ")
			if raw != strings.Join(strings.Fields(v.Value), " ") {
				return 0, errors.New("plain scalar span does not match its parsed value")
			}
			return end, nil
		}
	case yaml.SequenceNode:
		if v.Style&yaml.FlowStyle != 0 {
			if e.fm[start] != '[' {
				return 0, errors.New("expected [")
			}
			return e.scanFlow(start)
		}
		if e.fm[start] != '-' || len(v.Content) == 0 {
			return 0, errors.New("expected a block sequence")
		}
		last := v.Content[len(v.Content)-1]
		off, err := e.offset(last.Line, last.Column)
		if err != nil {
			return 0, err
		}
		return e.valueEnd(last, off, false, v.Column-1)
	}
	return 0, errors.New("unsupported value form")
}

func (e *editor) scanDoubleQuoted(start int) (int, error) {
	if e.fm[start] != '"' {
		return 0, errors.New("expected \"")
	}
	for i := start + 1; i < len(e.fm); i++ {
		switch e.fm[i] {
		case '\\':
			i++
		case '"':
			return i + 1, nil
		}
	}
	return 0, errors.New("unterminated double-quoted scalar")
}

func (e *editor) scanSingleQuoted(start int) (int, error) {
	if e.fm[start] != '\'' {
		return 0, errors.New("expected '")
	}
	for i := start + 1; i < len(e.fm); i++ {
		if e.fm[i] == '\'' {
			if i+1 < len(e.fm) && e.fm[i+1] == '\'' {
				i++
				continue
			}
			return i + 1, nil
		}
	}
	return 0, errors.New("unterminated single-quoted scalar")
}

func (e *editor) scanFlow(start int) (int, error) {
	depth := 0
	for i := start; i < len(e.fm); i++ {
		var err error
		switch c := e.fm[i]; {
		case c == '"':
			if i, err = e.scanDoubleQuoted(i); err != nil {
				return 0, err
			}
			i--
		case c == '\'':
			if i, err = e.scanSingleQuoted(i); err != nil {
				return 0, err
			}
			i--
		case c == '[' || c == '{':
			depth++
		case c == ']' || c == '}':
			depth--
			if depth == 0 {
				return i + 1, nil
			}
		case c == '#' && (i == start || isSpace(e.fm[i-1])):
			for i < len(e.fm) && e.fm[i] != '\n' {
				i++
			}
		}
	}
	return 0, errors.New("unterminated flow collection")
}

func isSpace(c byte) bool { return c == ' ' || c == '\t' || c == '\r' || c == '\n' }

// scanPlain ends a plain scalar at a comment, a mapping indicator, a flow
// terminator, or a line that is not an indented continuation.
func (e *editor) scanPlain(start int, flow bool, indent int) int {
	end := start
	for i := start; i < len(e.fm); i++ {
		c := e.fm[i]
		if c == '\n' || c == '\r' {
			next, ok := e.continuation(i, flow, indent)
			if !ok {
				return end
			}
			i = next - 1
			continue
		}
		if c == '#' && (i == start || e.fm[i-1] == ' ' || e.fm[i-1] == '\t') {
			return end
		}
		if c == ':' && (i+1 >= len(e.fm) || isSpace(e.fm[i+1]) || (flow && strings.IndexByte(",]}", e.fm[i+1]) >= 0)) {
			return end
		}
		if flow && (c == ',' || c == ']' || c == '}') {
			return end
		}
		if !isSpace(c) {
			end = i + 1
		}
	}
	return end
}

// continuation reports whether a plain scalar continues after the line break
// at i, returning the offset of the continuation's first character.
func (e *editor) continuation(i int, flow bool, indent int) (int, bool) {
	j := i
	for j < len(e.fm) && isSpace(e.fm[j]) {
		j++
	}
	if j >= len(e.fm) {
		return 0, false
	}
	if flow {
		return j, strings.IndexByte(",]}#", e.fm[j]) < 0
	}
	lineStart := bytes.LastIndexByte(e.fm[:j], '\n') + 1
	return j, j-lineStart > indent && e.fm[j] != '#'
}

// scanBlockScalar returns the end of the last content line of a literal or
// folded scalar whose header starts at start.
func (e *editor) scanBlockScalar(start, indent int) int {
	nl := bytes.IndexByte(e.fm[start:], '\n')
	if nl < 0 {
		return len(e.fm)
	}
	end := start + nl
	if end > start && e.fm[end-1] == '\r' {
		end--
	}
	for pos := start + nl + 1; pos < len(e.fm); {
		lineEnd := len(e.fm)
		if nl := bytes.IndexByte(e.fm[pos:], '\n'); nl >= 0 {
			lineEnd = pos + nl
		}
		content := lineEnd
		if content > pos && e.fm[content-1] == '\r' {
			content--
		}
		line := e.fm[pos:content]
		if strings.TrimLeft(string(line), " \t") != "" {
			if len(line)-len(bytes.TrimLeft(line, " ")) <= indent {
				break
			}
			end = content
		}
		pos = lineEnd + 1
	}
	return end
}
