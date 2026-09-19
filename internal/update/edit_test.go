package update

import (
	"strings"
	"testing"
)

const block = `---
id: "W-001"   # identity
type: work
"title": 'It''s an ''Ünïcode'' title'
status:	proposed # tab after colon
kind: feature
members:
- W-002
- W-003 # trailing item comment
# standalone comment between entries
depends_on:
  - "W-004"
  -   W-005
relates_to: [ "D-001",
  Q-001 ] # after list
priority: 2
size: |
  large
  # looks like a comment but is content
created: "2026-09-19T12:00:00Z"
# final comment
---

Body with --- inside.
`

func TestEditBlockForms(t *testing.T) {
	for _, tc := range []struct {
		name    string
		changes []change
		want    string
	}{
		{"same line plain keeps inline comment", []change{set("status", "active")},
			strings.Replace(block, "status:\tproposed # tab after colon", "status:\tactive # tab after colon", 1)},
		{"quoted value with quoted key", []change{set("title", `"New \"title\" ✓"`)},
			strings.Replace(block, "\"title\": 'It''s an ''Ünïcode'' title'", "\"title\": \"New \\\"title\\\" ✓\"", 1)},
		{"double quoted with inline comment", []change{set("id", `"W-001"`)}, block},
		{"block sequence replaced from the colon keeps the last item's comment", []change{set("members", `["W-009"]`)},
			strings.Replace(block, "members:\n- W-002\n- W-003 # trailing item comment\n", "members: [\"W-009\"] # trailing item comment\n", 1)},
		{"indented block sequence", []change{set("depends_on", "[]")},
			strings.Replace(block, "depends_on:\n  - \"W-004\"\n  -   W-005\n", "depends_on: []\n", 1)},
		{"multiline flow list keeps trailing comment", []change{set("relates_to", `["D-002"]`)},
			strings.Replace(block, "relates_to: [ \"D-001\",\n  Q-001 ] # after list", "relates_to: [\"D-002\"] # after list", 1)},
		{"literal block scalar", []change{set("size", "small")},
			strings.Replace(block, "size: |\n  large\n  # looks like a comment but is content\n", "size: small\n", 1)},
		{"priority", []change{set("priority", "5")}, strings.Replace(block, "priority: 2", "priority: 5", 1)},
		{"unset removes lines and inline comment, keeps standalone comments", []change{unset("members"), unset("id"), unset("size")},
			strings.NewReplacer("id: \"W-001\"   # identity\n", "", "members:\n- W-002\n- W-003 # trailing item comment\n", "",
				"size: |\n  large\n  # looks like a comment but is content\n", "").Replace(block)},
		{"unset absent is ignored and append goes before the closing delimiter", []change{unset("blocks"), set("updated", `"2026-09-19T13:00:00Z"`)},
			strings.Replace(block, "# final comment\n---", "# final comment\nupdated: \"2026-09-19T13:00:00Z\"\n---", 1)},
		{"multiple changes in one pass", []change{set("status", "done"), unset("priority"), set("members", "[]"), set("updated", `"2026-09-19T13:00:00Z"`)},
			strings.NewReplacer("status:\tproposed", "status:\tdone", "priority: 2\n", "",
				"members:\n- W-002\n- W-003 # trailing item comment\n", "members: [] # trailing item comment\n",
				"# final comment\n---", "# final comment\nupdated: \"2026-09-19T13:00:00Z\"\n---").Replace(block)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Edit([]byte(block), tc.changes)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != tc.want {
				t.Fatalf("got:\n%s\nwant:\n%s", got, tc.want)
			}
		})
	}
}

func TestEditPreservesBOMCRLFIndentationAndBody(t *testing.T) {
	src := "\ufeff---\r\n  id: W-001\r\n  type: work\r\n  title: first\r\n    second\r\n\r\n    third\r\n  status: proposed\r\n---\r\nbody\r\nno final newline"
	got, err := Edit([]byte(src), []change{set("title", `"one"`), set("status", "done"), set("updated", `"2026-09-19T13:00:00Z"`)})
	if err != nil {
		t.Fatal(err)
	}
	want := "\ufeff---\r\n  id: W-001\r\n  type: work\r\n  title: \"one\"\r\n  status: done\r\n  updated: \"2026-09-19T13:00:00Z\"\r\n---\r\nbody\r\nno final newline"
	if string(got) != want {
		t.Fatalf("got %q\nwant %q", got, want)
	}
	// A folded scalar and a plain multi-line scalar end before the next key.
	src = "---\ntitle: >-\n  folded\n  text\n\nstatus: open # c\nnote: a\n b\n---\n"
	got, err = Edit([]byte(src), []change{set("title", "t"), set("note", "n")})
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "---\ntitle: t\n\nstatus: open # c\nnote: n\n---\n" {
		t.Fatalf("got %q", got)
	}
}

func TestEditFlowMapping(t *testing.T) {
	src := "---\n{id: W-001, \"tïtle\": 'x', status: open, blocks: [a, \"b]\"],\n priority: 2 }\n---\nbody\n"
	for _, tc := range []struct {
		name    string
		changes []change
		want    string
	}{
		{"set unicode key value", []change{set("tïtle", `"y"`)}, strings.Replace(src, "'x'", `"y"`, 1)},
		{"set list with bracket in string", []change{set("blocks", "[]")}, strings.Replace(src, "[a, \"b]\"]", "[]", 1)},
		{"unset first", []change{unset("id")}, strings.Replace(src, "id: W-001, ", "", 1)},
		{"unset last uses preceding separator", []change{unset("priority")}, strings.Replace(src, ",\n priority: 2 }", " }", 1)},
		{"append after last value", []change{set("updated", `"2026-09-19T13:00:00Z"`), set("size", "small")},
			strings.Replace(src, "priority: 2 }", "priority: 2, updated: \"2026-09-19T13:00:00Z\", size: small }", 1)},
		{"unset last and append", []change{unset("priority"), set("size", "small")}, strings.Replace(src, ",\n priority: 2 }", ", size: small }", 1)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Edit([]byte(src), tc.changes)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != tc.want {
				t.Fatalf("got:\n%s\nwant:\n%s", got, tc.want)
			}
		})
	}
}

func TestEditRefusesUnsupportedSpans(t *testing.T) {
	for _, tc := range []struct{ name, src, key string }{
		{"tagged value", "---\nid: !!str W-001\nstatus: open\n---\n", "id"},
		{"anchored value", "---\nid: &a W-001\nstatus: open\n---\n", "id"},
		{"tagged key", "---\n!!str id: W-001\nstatus: open\n---\n", "id"},
		{"mapping value", "---\nid: {a: 1}\nstatus: open\n---\n", "id"},
		{"no closing delimiter", "---\nid: W-001\n", "id"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := Edit([]byte(tc.src), []change{set(tc.key, "x")}); err == nil {
				t.Fatal("expected a refusal")
			}
		})
	}
	// Unrelated tagged entries only serve as bounds and do not block edits.
	got, err := Edit([]byte("---\nid: !!str W-001\nstatus: open\n---\n"), []change{set("status", "resolved")})
	if err != nil || string(got) != "---\nid: !!str W-001\nstatus: resolved\n---\n" {
		t.Fatalf("got %q, %v", got, err)
	}
}
