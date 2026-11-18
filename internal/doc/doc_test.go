package doc

import (
	"reflect"
	"testing"
	"testing/fstest"
)

func TestOpenQuestions(t *testing.T) {
	comment := "<!-- Say which must be answered before review and which can wait.\n## Not a heading -->\n"
	for src, want := range map[string]string{
		"# T\n\n## Open questions\n\n" + comment + "\nNone.\n\n## Related documents\n\n- x\n": "",
		"# T\n\n## Open questions\n\n" + comment + "\n## Related documents\n":                 "",
		"# T\n\n## Open questions\n\nnone\n":                                                  "",
		"# T\n\n## Problem\n\nNone of it.\n":                                                  "",
		"# T\n\n## Open questions\n\n" + comment + "Should IDs be six digits? Blocks T-02.\n\n### Asked\n\nby whom\n\n## Related documents\n": "Should IDs be six digits? Blocks T-02.\n\n### Asked\n\nby whom",
		"# T\n\n## Open questions\n\nNone, but see T-02.\n": "None, but see T-02.",
	} {
		if got := openQuestions(src); got != want {
			t.Errorf("openQuestions(%q) = %q, want %q", src, got, want)
		}
	}
}

func TestFront(t *testing.T) {
	for _, c := range []struct{ src, fields, notes string }{
		{"---\ntitle: K\n---\n\n## Found\n\n- B moved\n", "---\ntitle: K\n", "## Found\n\n- B moved"},
		{"---\ntitle: K\n---\n", "---\ntitle: K\n", ""},
		{"# K\n\ntitle: K\n", "", ""},
		{"---\ntitle: K\n", "", ""},
	} {
		fields, notes, ok := Front(c.src)
		if fields != c.fields || notes != c.notes || ok != (c.fields != "") {
			t.Errorf("Front(%q) = %q, %q, %v", c.src, fields, notes, ok)
		}
	}
}

func TestShares(t *testing.T) {
	fsys := fstest.MapFS{}
	change := func(id, depends, capability, section, name string) {
		dir := "changes/" + id + "-x/"
		fsys[dir+".grove.yaml"] = &fstest.MapFile{Data: []byte("title: X\ncreated: 2026-10-14T00:00:00Z\ndepends_on: [" + depends + "]\n")}
		fsys[dir+"specs/"+capability+"/spec.md"] = &fstest.MapFile{Data: []byte("## " + section + " Requirements\n\n### Requirement: " + name + "\n")}
	}
	change("261001-aaaaa", "", "record-editing", "MODIFIED", "Keep a category")
	change("261001-bbbbb", "", "record-editing", "REMOVED", "Keep a category")
	change("261001-ccccc", "261001-ddddd", "record-editing", "MODIFIED", "Keep a category") // through d, after a
	change("261001-ddddd", "261001-aaaaa", "record-editing", "ADDED", "Other")
	change("261001-eeeee", "", "reports", "ADDED", "Keep a category")
	c, err := Load(fsys, func(string) (string, []string, string) {
		return "", []string{"specs/<capability>/spec.md"}, ""
	})
	if err != nil {
		t.Fatal(err)
	}
	keep := func(ids ...string) (out []Share) {
		for _, id := range ids {
			out = append(out, Share{id, "record-editing", "Keep a category"})
		}
		return out
	}
	for id, want := range map[string][]Share{
		"261001-aaaaa": keep("261001-bbbbb"),
		"261001-bbbbb": keep("261001-aaaaa", "261001-ccccc"),
		"261001-ccccc": keep("261001-bbbbb"),
		"261001-ddddd": nil,
		"261001-eeeee": nil,
	} {
		if got := c.Change(id).Shares; !reflect.DeepEqual(got, want) {
			t.Errorf("%s shares %v, want %v", id, got, want)
		}
	}
}
