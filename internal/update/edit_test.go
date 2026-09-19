package update

import (
	"reflect"
	"strings"
	"testing"

	"go.yaml.in/yaml/v3"
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

func TestEditTaggedAndAnchoredEntries(t *testing.T) {
	// The reader accepts explicit standard tags and unreferenced anchors; a
	// replaced value drops them, and an unrelated one stays untouched.
	for _, tc := range []struct{ name, src, key, value, want string }{
		{"tagged value", "---\nid: !!str W-001\ntitle: !!str hello # c\n---\n", "title", `"t"`, "---\nid: !!str W-001\ntitle: \"t\" # c\n---\n"},
		{"anchored value", "---\ntitle: &t hello\nstatus: open\n---\n", "title", "t", "---\ntitle: t\nstatus: open\n---\n"},
		{"tagged quoted value", "---\ntitle: !!str \"he\\\"llo\"\nstatus: open\n---\n", "status", "resolved", "---\ntitle: !!str \"he\\\"llo\"\nstatus: resolved\n---\n"},
		{"tagged key with later-line value", "---\n!!str members:\n- a\nstatus: open\n---\n", "members", "[]", "---\n!!str members: []\nstatus: open\n---\n"},
		{"anchored flow list", "---\nmembers: &m [a, b] # c\nstatus: open\n---\n", "members", `["c"]`, "---\nmembers: [\"c\"] # c\nstatus: open\n---\n"},
		{"unset tagged", "---\nid: W-001\nsize: !!str small\nstatus: open\n---\n", "size", "", "---\nid: W-001\nstatus: open\n---\n"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := set(tc.key, tc.value)
			if tc.value == "" {
				c = unset(tc.key)
			}
			got, err := Edit([]byte(tc.src), []change{c})
			if err != nil || string(got) != tc.want {
				t.Fatalf("got %q, %v\nwant %q", got, err, tc.want)
			}
		})
	}
}

func TestEditRefusesUnsupportedSpans(t *testing.T) {
	for _, tc := range []struct{ name, src, key string }{
		{"mapping value", "---\nid: {a: 1}\nstatus: open\n---\n", "id"},
		{"no closing delimiter", "---\nid: W-001\n", "id"},
		{"no opening delimiter", "id: W-001\n---\n", "id"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := Edit([]byte(tc.src), []change{set(tc.key, "x")}); err == nil {
				t.Fatal("expected a refusal")
			}
		})
	}
	// A comment line inside a removed list goes with the list; one after it stays.
	src := "---\nmembers:\n# why these members\n- W-002\n# after the list\nstatus: open\n---\n"
	got, err := Edit([]byte(src), []change{unset("members")})
	if err != nil || string(got) != "---\n# after the list\nstatus: open\n---\n" {
		t.Fatalf("got %q, %v", got, err)
	}
	// A brace in a trailing comment does not extend a flow mapping's bound.
	src = "---\n{id: W-001, status: open}\n# a } brace\n---\n"
	got, err = Edit([]byte(src), []change{set("status", "resolved"), set("size", "small")})
	if err != nil || string(got) != "---\n{id: W-001, status: resolved, size: small}\n# a } brace\n---\n" {
		t.Fatalf("got %q, %v", got, err)
	}
}

// editTable runs exact-byte cases and checks each result still parses.
func editTable(t *testing.T, cases []struct {
	name, src string
	changes   []change
	want      string
}) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := Edit([]byte(tc.src), tc.changes)
			if err != nil {
				t.Fatal(err)
			}
			if string(got) != tc.want {
				t.Fatalf("got:\n%q\nwant:\n%q", got, tc.want)
			}
			start, end, _, err := frontmatter(got)
			if err != nil {
				t.Fatal(err)
			}
			var parsed map[string]any
			if err := yaml.Unmarshal(got[start:end], &parsed); err != nil {
				t.Fatalf("the result must stay valid YAML: %v\n%s", err, got)
			}
		})
	}
}

// A comment between a key's colon and a value that starts on a later line is
// outside the value, so replacing the value keeps it.
func TestEditKeepsCommentsBeforeLaterLineValues(t *testing.T) {
	editTable(t, []struct {
		name, src string
		changes   []change
		want      string
	}{
		{"key-line comment", "---\ntitle: # retain\n  First\nstatus: open\n---\n", []change{set("title", `"New"`)},
			"---\ntitle: # retain\n  \"New\"\nstatus: open\n---\n"},
		{"standalone comment and a multi-line value", "---\ntitle:\n  # standalone\n  First\n  second line\nstatus: open # s\n---\nbody\n", []change{set("title", `"New"`)},
			"---\ntitle:\n  # standalone\n  \"New\"\nstatus: open # s\n---\nbody\n"},
		{"BOM and CRLF with a quoted value", "\ufeff---\r\ntitle: # retain\r\n  'First'\r\nstatus: open\r\n---\r\nbody", []change{set("title", `"New"`)},
			"\ufeff---\r\ntitle: # retain\r\n  \"New\"\r\nstatus: open\r\n---\r\nbody"},
		{"quoted key with a list at the key's indentation", "---\n\"members\": # why\n- W-002\n- W-003 # last\nstatus: open\n---\n", []change{set("members", `["W-009"]`)},
			"---\n\"members\": # why\n  [\"W-009\"] # last\nstatus: open\n---\n"},
		{"indented list", "---\nmembers: # why\n  - W-002\n  - W-003\nstatus: open\n---\n", []change{set("members", `["W-009"]`)},
			"---\nmembers: # why\n  [\"W-009\"]\nstatus: open\n---\n"},
		{"indented mapping", "---\n  id: W-001\n  members: # why\n  - W-002\n---\n", []change{set("members", "[]")},
			"---\n  id: W-001\n  members: # why\n    []\n---\n"},
		{"flow mapping", "---\n{id: W-001, title: # c\n  First, status: open}\n---\n", []change{set("title", `"New"`)},
			"---\n{id: W-001, title: # c\n  \"New\", status: open}\n---\n"},
		{"no comment keeps the accepted collapse", "---\ntitle:\n  First\nstatus: open\n---\n", []change{set("title", `"New"`)},
			"---\ntitle: \"New\"\nstatus: open\n---\n"},
		{"unset removes the entry with its key-line comment", "---\nid: W-001\n# above\nsize: # c\n  small\n# below\nstatus: open\n---\n", []change{unset("size")},
			"---\nid: W-001\n# above\n# below\nstatus: open\n---\n"},
	})
}

// The reader accepts explicit keys, so the editor must edit them.
func TestEditExplicitKeys(t *testing.T) {
	editTable(t, []struct {
		name, src string
		changes   []change
		want      string
	}{
		{"set", "---\n? status\n: proposed\nid: W-001\n---\n", []change{set("status", "active")},
			"---\n? status\n: active\nid: W-001\n---\n"},
		{"set keeps key and value comments", "---\nid: W-001\n? \"title\" # key comment\n: # value comment\n  hello\n---\n", []change{set("title", `"t"`)},
			"---\nid: W-001\n? \"title\" # key comment\n: # value comment\n  \"t\"\n---\n"},
		{"CRLF", "---\r\n? status\r\n: proposed # c\r\nid: W-001\r\n---\r\n", []change{set("status", "active")},
			"---\r\n? status\r\n: active # c\r\nid: W-001\r\n---\r\n"},
		{"unset tagged", "---\nid: W-001\n? !!str size\n: small # c\nstatus: open\n---\n", []change{unset("size")},
			"---\nid: W-001\nstatus: open\n---\n"},
		{"unset multi-line value", "---\nid: W-001\n# keep\n? members\n: - W-002\n  - W-003\nstatus: open\n---\n", []change{unset("members")},
			"---\nid: W-001\n# keep\nstatus: open\n---\n"},
		{"set then append", "---\n? status\n: proposed\n---\n", []change{set("status", "active"), set("updated", `"2026-09-19T13:00:00Z"`)},
			"---\n? status\n: active\nupdated: \"2026-09-19T13:00:00Z\"\n---\n"},
		{"flow set", "---\n{id: W-001, ? size : small, status: open}\n---\n", []change{set("size", "large")},
			"---\n{id: W-001, ? size : large, status: open}\n---\n"},
		{"flow unset", "---\n{id: W-001, ? size : small, status: open}\n---\n", []change{unset("size")},
			"---\n{id: W-001, status: open}\n---\n"},
		{"flow unset last", "---\n{id: W-001, ? size : small}\n---\n", []change{unset("size")},
			"---\n{id: W-001}\n---\n"},
	})
}

// Separators are planned for the whole request: each comma is removed once.
func TestEditFlowRemovals(t *testing.T) {
	const stamp = `"2026-09-19T13:00:00Z"`
	line := "---\n{id: W-001, type: work, title: T, status: proposed, kind: fix, size: small}\n---\nBody\n"
	trimmed := "---\n{id: W-001, type: work, title: T, status: proposed}\n---\nBody\n"
	multi := "---\n{id: W-001, # identity\n # standalone stays\n kind: fix, # why fix\n status: open,\n size: small # why small\n}\n---\n"
	editTable(t, []struct {
		name, src string
		changes   []change
		want      string
	}{
		{"final two", line, []change{unset("kind"), unset("size")}, trimmed},
		{"final two reversed", line, []change{unset("size"), unset("kind")}, trimmed},
		{"final two and an absent updated", line, []change{unset("kind"), unset("size"), set("updated", stamp)},
			"---\n{id: W-001, type: work, title: T, status: proposed, updated: " + stamp + "}\n---\nBody\n"},
		{"final two reversed and an absent updated", line, []change{set("updated", stamp), unset("size"), unset("kind")},
			"---\n{id: W-001, type: work, title: T, status: proposed, updated: " + stamp + "}\n---\nBody\n"},
		{"final two and an existing updated", "---\n{id: W-001, updated: \"old\", kind: fix, size: small}\n---\n", []change{unset("kind"), unset("size"), set("updated", stamp)},
			"---\n{id: W-001, updated: " + stamp + "}\n---\n"},
		{"first, middle, and last", "---\n{kind: fix, id: W-001, size: small, status: open, priority: 2}\n---\n", []change{unset("priority"), unset("kind"), unset("size")},
			"---\n{id: W-001, status: open}\n---\n"},
		{"every entry", "---\n{kind: fix, size: small}\n---\n", []change{unset("kind"), unset("size")}, "---\n{}\n---\n"},
		{"every entry and an append", "---\n{kind: fix, size: small}\n---\n", []change{unset("kind"), unset("size"), set("updated", stamp)}, "---\n{updated: " + stamp + "}\n---\n"},
		{"trailing comma", "---\n{id: W-001, kind: fix, size: small, }\n---\n", []change{unset("kind"), unset("size")}, "---\n{id: W-001, }\n---\n"},
		{"trailing comma and an append", "---\n{id: W-001, kind: fix, size: small, }\n---\n", []change{unset("kind"), unset("size"), set("updated", stamp)},
			"---\n{id: W-001, updated: " + stamp + ", }\n---\n"},
		{"retained last entry with a trailing comma and an append", "---\n{id: W-001, kind: fix, }\n---\n", []change{set("updated", stamp)},
			"---\n{id: W-001, kind: fix, updated: " + stamp + ", }\n---\n"},
		{"set and unset", "---\n{id: W-001, status: open, size: small}\n---\n", []change{set("status", "done"), unset("size")},
			"---\n{id: W-001, status: done}\n---\n"},
		{"multi-line: inline comments go, standalone comments and neighbors stay", multi, []change{unset("kind"), unset("size"), set("updated", stamp)},
			"---\n{id: W-001, # identity\n # standalone stays\n status: open, updated: " + stamp + "\n}\n---\n"},
		{"multi-line without an append", multi, []change{unset("size"), unset("kind")},
			"---\n{id: W-001, # identity\n # standalone stays\n status: open\n}\n---\n"},
		{"multi-line CRLF", strings.ReplaceAll(multi, "\n", "\r\n"), []change{unset("kind")},
			strings.ReplaceAll(strings.Replace(multi, " kind: fix, # why fix\n", "", 1), "\n", "\r\n")},
		{"a comment hides the preceding comma, which stays as a valid trailing one", "---\n{id: W-001, # identity\n size: small}\n---\n", []change{unset("size")},
			"---\n{id: W-001, # identity\n }\n---\n"},
		{"a comment hides the preceding comma and an append reuses it", "---\n{id: W-001, # identity\n size: small}\n---\n", []change{unset("size"), set("updated", stamp)},
			"---\n{id: W-001, # identity\n updated: " + stamp + "}\n---\n"},
		{"inline comment of a removed middle entry", "---\n{id: W-001, kind: fix, # why\n status: open}\n---\n", []change{unset("kind")},
			"---\n{id: W-001, \n status: open}\n---\n"},
	})
}

// Every subset of removals, with and without an append, over several flow
// layouts must leave valid YAML holding exactly the expected keys and values.
func TestEditFlowRemovalsExhaustive(t *testing.T) {
	keys := []string{"a", "b", "c", "d"}
	for _, src := range []string{
		"---\n{a: 1, b: 2, c: 3, d: 4}\n---\n",
		"---\n{a: 1, b: 2, c: 3, d: 4, }\n---\n",
		"---\n{ ? a : 1 , b: 2 ,c: 3,d: 4 }\n---\n",
		"---\n{a: 1,\n b: 2,\n c: 3,\n d: 4}\n---\n",
		"---\n{a: 1, # one\n # standalone\n b: 2, # two\n c: 3, # three\n d: 4 # four\n}\n---\n",
		"---\n{\r\n  a: 1,\r\n  b: [x, \"y,}\"], # two\r\n  c: 'it''s',\r\n  d: 4,\r\n}\r\n---\r\n",
	} {
		var before map[string]any
		start, end, _, _ := frontmatter([]byte(src))
		if err := yaml.Unmarshal([]byte(src)[start:end], &before); err != nil {
			t.Fatal(err)
		}
		for mask := 0; mask < 1<<len(keys); mask++ {
			for _, appendToo := range []bool{false, true} {
				want := map[string]any{}
				var changes []change
				for i, k := range keys {
					if mask&(1<<i) != 0 {
						changes = append(changes, unset(k))
					} else {
						want[k] = before[k]
					}
				}
				if appendToo {
					changes = append(changes, set("e", "5"), set("f", `"six"`))
					want["e"], want["f"] = 5, "six"
				}
				got, err := Edit([]byte(src), changes)
				if err != nil {
					t.Fatalf("%q mask %04b append %v: %v", src, mask, appendToo, err)
				}
				start, end, _, err := frontmatter(got)
				if err != nil {
					t.Fatal(err)
				}
				parsed := map[string]any{}
				if err := yaml.Unmarshal(got[start:end], &parsed); err != nil || !reflect.DeepEqual(parsed, want) {
					t.Fatalf("%q mask %04b append %v:\n%s\nparsed %v (%v), want %v", src, mask, appendToo, got, parsed, err, want)
				}
				if strings.Contains(src, "# standalone") != strings.Contains(string(got), "# standalone") {
					t.Fatalf("standalone comment lost:\n%s", got)
				}
			}
		}
	}
}
