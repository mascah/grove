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

// parse reads the YAML mapping src of file into d.Meta and returns its
// values, reporting every field of required that is missing or empty, and
// every field that is neither required nor description. offset is the number
// of lines in file before src; a problem with no line of its own is reported
// on line offset.
func (d *Document) parse(file string, src []byte, offset int, required []string) map[string]string {
	values, lines := map[string]string{}, map[string]int{}
	var root yaml.Node
	if err := yaml.NewDecoder(bytes.NewReader(src)).Decode(&root); err != nil && !errors.Is(err, io.EOF) {
		line, msg := offset, strings.TrimPrefix(err.Error(), "yaml: ")
		if m := yamlErrorLine.FindStringSubmatch(err.Error()); m != nil {
			n, _ := strconv.Atoi(m[1])
			line, msg = n+offset, m[2]
		}
		d.report(file, line, "invalid YAML: "+msg)
		return values
	}
	if len(root.Content) > 0 {
		n := root.Content[0]
		if n.Kind != yaml.MappingNode {
			d.report(file, n.Line+offset, "expected a YAML mapping of fields")
			return values
		}
		for i := 0; i+1 < len(n.Content); i += 2 {
			key, value := n.Content[i], n.Content[i+1]
			line := key.Line + offset
			switch first, seen := lines[key.Value]; {
			case seen:
				d.report(file, line, fmt.Sprintf("duplicate key %q; first set on line %d", key.Value, first))
			case !slices.Contains(required, key.Value) && key.Value != "description":
				d.report(file, line, fmt.Sprintf("unknown field %q; this file takes %s and description", key.Value, strings.Join(required, ", ")))
			case value.Kind != yaml.ScalarNode || strings.TrimSpace(value.Value) == "":
				d.report(file, line, key.Value+": expected a non-empty value")
			default:
				values[key.Value], lines[key.Value] = value.Value, line
			}
		}
	}
	for _, key := range required {
		if _, ok := values[key]; !ok && lines[key] == 0 {
			d.report(file, offset, key+": required field is missing")
		}
	}
	d.Meta = Meta{Title: values["title"], Created: values["created"], Description: values["description"]}
	return values
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

func (d *Document) report(file string, line int, msg string) {
	d.Diagnostics = append(d.Diagnostics, Diagnostic{Path: file, Line: line, Message: msg})
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
