package project

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"regexp"
	"slices"
	"strings"
	"time"
	"unicode/utf8"

	"go.yaml.in/yaml/v3"
)

// Record retains the original Markdown so inspection never reserializes it.
type Record struct {
	ID, Type, Title, Status string
	Kind, Size              string
	Priority                *int
	Created, Updated        *time.Time
	DependsOn, Members      []string
	Blocks, RelatesTo       []string
	Path                    string
	Source                  []byte
}

type Diagnostic struct {
	Path, Field, Message string
	Line                 int
}

func (d Diagnostic) String() string {
	path := d.Path
	if d.Line > 0 {
		path += fmt.Sprintf(":%d", d.Line)
	}
	if d.Field != "" {
		return path + ": " + d.Field + ": " + d.Message
	}
	return path + ": " + d.Message
}

type metadata struct {
	path   string
	offset int
	fields map[string]*yaml.Node
	errors []Diagnostic
}

func (m *metadata) problem(field, message string) {
	line := 0
	if n := m.fields[field]; n != nil {
		line = n.Line + m.offset
	}
	m.errors = append(m.errors, Diagnostic{Path: m.path, Field: field, Message: message, Line: line})
}

func parseMapping(path string, source []byte, offset int) *metadata {
	m := &metadata{path: path, offset: offset, fields: make(map[string]*yaml.Node)}
	if !utf8.Valid(source) {
		m.problem("", "invalid UTF-8")
		return m
	}
	decoder := yaml.NewDecoder(bytes.NewReader(source))
	var doc yaml.Node
	if err := decoder.Decode(&doc); err != nil {
		m.problem("", "expected a YAML mapping: "+err.Error())
		return m
	}
	var extra yaml.Node
	if err := decoder.Decode(&extra); err != io.EOF {
		if err != nil {
			m.problem("", "invalid YAML document: "+err.Error())
		} else {
			m.problem("", "expected exactly one YAML document")
		}
	}
	if len(doc.Content) != 1 || doc.Content[0].Kind != yaml.MappingNode || doc.Content[0].Tag != "!!map" {
		m.problem("", "expected a YAML mapping")
		return m
	}
	nodes := doc.Content[0].Content
	for i := 0; i < len(nodes); i += 2 {
		key, value := nodes[i], nodes[i+1]
		if key.Kind != yaml.ScalarNode || key.Tag != "!!str" {
			m.problem("", "mapping keys must be strings; merge keys are not supported")
			continue
		}
		if _, exists := m.fields[key.Value]; exists {
			m.problem(key.Value, "duplicate YAML key")
			continue
		}
		m.fields[key.Value] = value
	}
	return m
}

func (m *metadata) stringField(key string, required bool) string {
	n, ok := m.fields[key]
	if !ok {
		if required {
			m.problem(key, "required field is missing")
		}
		return ""
	}
	if n.Kind == yaml.AliasNode {
		m.problem(key, "YAML aliases are not supported")
		return ""
	}
	if n.Kind != yaml.ScalarNode || n.Tag != "!!str" || strings.TrimSpace(n.Value) == "" {
		m.problem(key, "expected a nonempty string")
		return ""
	}
	return n.Value
}

func (m *metadata) integerField(key string, required bool) (int, bool) {
	n, ok := m.fields[key]
	if !ok {
		if required {
			m.problem(key, "required field is missing")
		}
		return 0, false
	}
	var value int
	if n.Kind != yaml.ScalarNode || n.Tag != "!!int" || n.Decode(&value) != nil {
		m.problem(key, "expected an integer")
		return 0, false
	}
	return value, true
}

func (m *metadata) listField(key string) []string {
	n, ok := m.fields[key]
	if !ok {
		return nil
	}
	if n.Kind != yaml.SequenceNode || n.Tag != "!!seq" {
		m.problem(key, "expected a list of record IDs")
		return nil
	}
	result := make([]string, 0, len(n.Content))
	seen := map[string]bool{}
	for _, item := range n.Content {
		if item.Kind != yaml.ScalarNode || item.Tag != "!!str" || strings.TrimSpace(item.Value) == "" {
			m.problem(key, "each target must be a nonempty string; aliases are not supported")
			continue
		}
		if seen[item.Value] {
			m.problem(key, "duplicate target "+item.Value)
		}
		seen[item.Value] = true
		result = append(result, item.Value)
	}
	return result
}

var datePattern = regexp.MustCompile("^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}Z$")
var idPattern = regexp.MustCompile("^[WQD]-[0-9]{3,}$")

func validID(id, kind string) bool {
	prefix := map[string]string{"work": "W-", "question": "Q-", "decision": "D-"}[kind]
	if prefix == "" || !strings.HasPrefix(id, prefix) || !idPattern.MatchString(id) {
		return false
	}
	number := id[2:]
	return strings.TrimLeft(number, "0") != "" && (len(number) == 3 || number[0] != '0')
}

func (m *metadata) dateField(key string) *time.Time {
	if _, ok := m.fields[key]; !ok {
		return nil
	}
	value := m.stringField(key, false)
	if value == "" {
		return nil
	}
	date, err := time.Parse(time.RFC3339, value)
	if !datePattern.MatchString(value) || err != nil {
		m.problem(key, "expected a quoted UTC timestamp YYYY-MM-DDTHH:MM:SSZ")
		return nil
	}
	return &date
}

// Revision identifies exact file content: "sha256:" plus the lowercase hex
// SHA-256 of every byte, including any BOM and line endings. Timestamps are an
// authoring convention; callers detecting stale input must compare content.
func Revision(source []byte) string {
	sum := sha256.Sum256(source)
	return "sha256:" + hex.EncodeToString(sum[:])
}

// ParseRecord validates one record's source against the schema for the type
// folder it sits in. Writers use it so the candidate they produce is judged by
// the same rules the reader applies.
func ParseRecord(path, folderType string, source []byte) (*Record, []Diagnostic) {
	r := &Record{Path: path, Source: source}
	fail := func(message string) (*Record, []Diagnostic) {
		return r, []Diagnostic{{Path: path, Field: "frontmatter", Message: message}}
	}
	if !utf8.Valid(source) {
		return fail("invalid UTF-8")
	}
	text := strings.TrimPrefix(string(source), "\ufeff")
	lines := strings.Split(strings.ReplaceAll(text, "\r\n", "\n"), "\n")
	if len(lines) == 0 || strings.TrimRight(lines[0], " \t") != "---" {
		return fail("expected an opening --- line")
	}
	end := 1
	for end < len(lines) && strings.TrimRight(lines[end], " \t") != "---" {
		end++
	}
	if end == len(lines) {
		return fail("missing closing --- line")
	}
	m := parseMapping(path, []byte(strings.Join(lines[1:end], "\n")), 1)
	r.ID = m.stringField("id", true)
	r.Type = m.stringField("type", true)
	r.Title = m.stringField("title", true)
	r.Status = m.stringField("status", true)
	if r.Type != folderType {
		m.problem("type", "must match type folder "+folderType)
	}
	if !validID(r.ID, r.Type) {
		m.problem("id", "expected a canonical positive ID with matching type prefix, e.g. W-001, Q-001, D-001")
	}
	lifecycles := map[string][]string{
		"work":     {"proposed", "active", "done", "abandoned"},
		"question": {"open", "resolved"}, "decision": {"proposed", "accepted", "rejected"},
	}
	if !slices.Contains(lifecycles[r.Type], r.Status) {
		m.problem("status", "unsupported lifecycle value for "+r.Type)
	}
	allowed := []string{"id", "type", "title", "status", "relates_to", "created", "updated"}
	if r.Type == "work" {
		allowed = append(allowed, "kind", "size", "priority", "members", "depends_on")
		r.Kind, r.Size = m.stringField("kind", false), m.stringField("size", false)
		if r.Kind != "" && !slices.Contains([]string{"feature", "fix", "refactor", "investigation", "tooling", "release"}, r.Kind) {
			m.problem("kind", "unsupported work kind")
		}
		if r.Size != "" && !slices.Contains([]string{"small", "medium", "large"}, r.Size) {
			m.problem("size", "expected small, medium, or large")
		}
		if priority, ok := m.integerField("priority", false); ok {
			r.Priority = &priority
			if priority < 1 || priority > 5 {
				m.problem("priority", "expected 1 (highest) through 5 (lowest)")
			}
		}
		r.Members, r.DependsOn = m.listField("members"), m.listField("depends_on")
	}
	if r.Type == "question" {
		allowed = append(allowed, "blocks")
		r.Blocks = m.listField("blocks")
	}
	for key := range m.fields {
		if !slices.Contains(allowed, key) {
			m.problem(key, "unknown field or not allowed on this record type")
		}
	}
	r.RelatesTo = m.listField("relates_to")
	r.Created, r.Updated = m.dateField("created"), m.dateField("updated")
	if r.Created != nil && r.Updated != nil && r.Updated.Before(*r.Created) {
		m.problem("updated", "must not precede created")
	}
	return r, m.errors
}
