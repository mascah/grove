package doc

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// readMeta reads d's .grove.yaml into d.Meta, reporting the file missing
// when it is, as one that holds holds.
func (d *Document) readMeta(fsys fs.FS, required []string, holds string) error {
	file := d.Dir + "/.grove.yaml"
	switch b, err := fs.ReadFile(fsys, file); {
	case errors.Is(err, fs.ErrNotExist):
		d.report(file, 0, "missing; it holds the "+holds)
	case err != nil:
		return err
	default:
		d.parse(file, b, required)
	}
	return nil
}

// parse reads the YAML mapping src of file into d.Meta, reporting every
// field of required that is missing or empty, and every field that is
// neither required nor description.
func (d *Document) parse(file string, src []byte, required []string) {
	var root yaml.Node
	if err := yaml.NewDecoder(bytes.NewReader(src)).Decode(&root); err != nil && !errors.Is(err, io.EOF) {
		line, msg := 0, strings.TrimPrefix(err.Error(), "yaml: ")
		if m := yamlErrorLine.FindStringSubmatch(err.Error()); m != nil {
			line, _ = strconv.Atoi(m[1])
			msg = m[2]
		}
		d.report(file, line, "invalid YAML: "+msg)
		return
	}
	values, lines := map[string]string{}, map[string]int{}
	if len(root.Content) > 0 {
		n := root.Content[0]
		if n.Kind != yaml.MappingNode {
			d.report(file, n.Line, "expected a YAML mapping of fields")
			return
		}
		for i := 0; i+1 < len(n.Content); i += 2 {
			key, value := n.Content[i], n.Content[i+1]
			switch first, seen := lines[key.Value]; {
			case seen:
				d.report(file, key.Line, fmt.Sprintf("duplicate key %q; first set on line %d", key.Value, first))
			case !slices.Contains(required, key.Value) && key.Value != "description":
				d.report(file, key.Line, fmt.Sprintf("unknown field %q; this file takes %s and description", key.Value, strings.Join(required, ", ")))
			case value.Kind != yaml.ScalarNode || strings.TrimSpace(value.Value) == "":
				d.report(file, key.Line, key.Value+": expected a non-empty value")
			default:
				values[key.Value], lines[key.Value] = value.Value, key.Line
			}
		}
	}
	for _, key := range required {
		if _, ok := values[key]; !ok && lines[key] == 0 {
			d.report(file, 0, key+": required field is missing")
		}
	}
	d.Meta = Meta{Title: values["title"], Created: values["created"], Description: values["description"]}
}

var yamlErrorLine = regexp.MustCompile(`^yaml: line ([0-9]+): (.*)$`)

func (d *Document) report(file string, line int, msg string) {
	d.Diagnostics = append(d.Diagnostics, Diagnostic{Path: file, Line: line, Message: msg})
}
