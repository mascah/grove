package doc

import "testing"

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
