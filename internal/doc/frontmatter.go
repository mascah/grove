package doc

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Parse reads one document and reports every field of required that is
// missing or empty, and every field that is neither required nor
// description. The result is never nil.
func Parse(path string, src []byte, required []string) *Document {
	d := &Document{Path: path, Source: string(src)}
	meta, problem := splitFrontmatter(string(src))
	if problem != "" {
		d.report(1, problem)
		return d
	}
	var root yaml.Node
	if err := yaml.NewDecoder(bytes.NewReader(meta)).Decode(&root); err != nil && !errors.Is(err, io.EOF) {
		line, msg := 1, strings.TrimPrefix(err.Error(), "yaml: ")
		if m := yamlErrorLine.FindStringSubmatch(err.Error()); m != nil {
			n, _ := strconv.Atoi(m[1])
			line, msg = n+1, m[2]
		}
		d.report(line, "invalid YAML frontmatter: "+msg)
		return d
	}
	values, lines := map[string]string{}, map[string]int{}
	if len(root.Content) > 0 {
		n := root.Content[0]
		if n.Kind != yaml.MappingNode {
			d.report(n.Line+1, "frontmatter must be a YAML mapping")
			return d
		}
		for i := 0; i+1 < len(n.Content); i += 2 {
			key, value := n.Content[i], n.Content[i+1]
			line := key.Line + 1 // the opening fence is line 1
			switch first, seen := lines[key.Value]; {
			case seen:
				d.report(line, fmt.Sprintf("duplicate key %q; first set on line %d", key.Value, first))
			case !slices.Contains(required, key.Value) && key.Value != "description":
				d.report(line, fmt.Sprintf("unknown field %q; this document takes %s and description", key.Value, strings.Join(required, ", ")))
			case value.Kind != yaml.ScalarNode || strings.TrimSpace(value.Value) == "":
				d.report(line, key.Value+": expected a non-empty value")
			default:
				values[key.Value], lines[key.Value] = value.Value, line
			}
		}
	}
	for _, key := range required {
		if _, ok := values[key]; !ok && lines[key] == 0 {
			d.report(1, key+": required field is missing")
		}
	}
	d.Meta = Meta{ID: values["id"], Title: values["title"], Created: values["created"],
		Exploration: values["exploration"], Description: values["description"]}
	return d
}

var yamlErrorLine = regexp.MustCompile(`^yaml: line ([0-9]+): (.*)$`)

// splitFrontmatter returns the lines between a "---" first line and the next
// "---" line. A non-empty problem means there is no frontmatter block.
func splitFrontmatter(text string) (meta []byte, problem string) {
	isFence := func(line string) bool { return strings.TrimRight(line, " \t\r") == "---" }
	first, rest, _ := strings.Cut(strings.TrimPrefix(text, "\ufeff"), "\n")
	if !isFence(first) {
		return nil, "expected frontmatter: the first line must be ---"
	}
	for off := 0; off < len(rest); {
		l, _, more := strings.Cut(rest[off:], "\n")
		if isFence(l) {
			return []byte(rest[:off]), ""
		}
		off += len(l)
		if more {
			off++
		}
	}
	return nil, "frontmatter is not closed: expected a second --- line"
}

func (d *Document) report(line int, msg string) {
	d.Diagnostics = append(d.Diagnostics, Diagnostic{Path: d.Path, Line: line, Message: msg})
}

// Format returns a document with meta as its frontmatter, then body.
func Format(meta Meta, body string) string {
	out, _ := yaml.Marshal(meta) // a struct of strings always marshals
	return "---\n" + string(out) + "---\n" + body
}

// Body returns src without its frontmatter block.
func Body(src string) string {
	meta, problem := splitFrontmatter(src)
	if problem != "" {
		return src
	}
	_, rest, _ := strings.Cut(src, "\n")
	return strings.TrimPrefix(rest[len(meta):], "---")
}
