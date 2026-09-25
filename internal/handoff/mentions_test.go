package handoff

import (
	"reflect"
	"testing"

	"github.com/mascah/grove/internal/project"
)

// Mentions finds real links to project paths and code spans, each with its
// line, and nothing that only looks like one.
func TestMentions(t *testing.T) {
	t.Parallel()
	r := &project.Record{Path: "grove/G-001-x.md", Source: []byte("---\nid: G-001\n---\n\n" +
		"See [search](../internal/tui/search.go) and `ws.rs` here.\n" +
		"[`run.py`](../evals/run.py#top), [web](https://example.com), [up](../../out.md), ![img](../a.png)\n\n" +
		"```\n`fenced` [not](../fenced.go)\n```\n\n" +
		"Two ``a ` b`` spans and `  ` blank.\n")}
	want := []Mention{
		{Path: "internal/tui/search.go", Line: "See [search](../internal/tui/search.go) and `ws.rs` here."},
		{Span: "ws.rs", Line: "See [search](../internal/tui/search.go) and `ws.rs` here."},
		{Path: "evals/run.py", Line: "[`run.py`](../evals/run.py#top), [web](https://example.com), [up](../../out.md), ![img](../a.png)"},
		{Span: "run.py", Line: "[`run.py`](../evals/run.py#top), [web](https://example.com), [up](../../out.md), ![img](../a.png)"},
		{Span: "a ` b", Line: "Two ``a ` b`` spans and `  ` blank."},
	}
	if got := Mentions(r); !reflect.DeepEqual(got, want) {
		t.Fatalf("Mentions:\n got %q\nwant %q", got, want)
	}
}
