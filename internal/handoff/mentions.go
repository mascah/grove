package handoff

import (
	"bytes"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"

	"github.com/mascah/grove/internal/project"
)

// Mention is something a record's body names that a search for a file can
// find (G-153): a real link that resolves to a project path, or a code span,
// with the body line it starts on. Exactly one of Path and Span is set.
type Mention struct {
	Path, Span, Line string
}

// Mentions lists the links and code spans in r's body, in order. It reads
// nothing: a link is resolved as context resolves it, and one that names no
// project path, or is malformed, is not a mention.
func Mentions(r *project.Record) []Mention {
	src := body(r.Source)
	var out []Mention
	doc := goldmark.DefaultParser().Parse(text.NewReader(src))
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) { // the walker returns no error
		if !entering {
			return ast.WalkContinue, nil
		}
		switch n := n.(type) {
		case *ast.Link:
			if target, _, err := resolve(r.Path, string(n.Destination)); err == nil && target != "" {
				out = append(out, Mention{Path: target, Line: lineAt(src, n, string(n.Destination))})
			}
		case *ast.CodeSpan:
			var span strings.Builder
			for c := n.FirstChild(); c != nil; c = c.NextSibling() {
				switch c := c.(type) {
				case *ast.Text:
					span.Write(c.Segment.Value(src))
				case *ast.String:
					span.Write(c.Value)
				}
			}
			if s := strings.TrimSpace(span.String()); s != "" {
				out = append(out, Mention{Span: s, Line: lineAt(src, n, s)})
			}
			return ast.WalkSkipChildren, nil
		}
		return ast.WalkContinue, nil
	})
	return out
}

// lineAt is the body line holding n's first text, or failing that the first
// line containing fallback, trimmed.
func lineAt(src []byte, n ast.Node, fallback string) string {
	at := -1
	_ = ast.Walk(n, func(c ast.Node, entering bool) (ast.WalkStatus, error) {
		if t, ok := c.(*ast.Text); ok && entering {
			at = t.Segment.Start
			return ast.WalkStop, nil
		}
		return ast.WalkContinue, nil
	})
	if at < 0 {
		at = bytes.Index(src, []byte(fallback))
	}
	if at < 0 {
		return ""
	}
	start := bytes.LastIndexByte(src[:at], '\n') + 1
	end := bytes.IndexByte(src[at:], '\n')
	if end < 0 {
		end = len(src) - at
	}
	return strings.TrimSpace(string(src[start : at+end]))
}
