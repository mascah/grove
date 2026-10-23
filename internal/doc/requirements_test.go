package doc

import (
	"reflect"
	"strings"
	"testing"
)

func TestRequirements(t *testing.T) {
	src := `# Changes to record-editing

## MODIFIED Requirements

### Requirement: Keep a category

When a user edits a record, the system shall keep its category.

'''sh
# edit a note
ledger edit 3 --note x
'''

#### Scenario: A note

- GIVEN a record

## REMOVED Requirements

### Requirement: Old rule

No longer wanted.

## Observed current behavior

Text.
`
	want := []Requirement{
		{Section: "MODIFIED", Name: "Keep a category", Line: 5, Text: "When a user edits a record, the system shall keep its category.\n\n```sh\n# edit a note\nledger edit 3 --note x\n```\n\n#### Scenario: A note\n\n- GIVEN a record"},
		{Section: "REMOVED", Name: "Old rule", Line: 20, Text: "No longer wanted."},
	}
	if got := Requirements(strings.ReplaceAll(src, "'''", "```")); !reflect.DeepEqual(got, want) {
		t.Errorf("Requirements:\n got %#v\nwant %#v", got, want)
	}
}

func TestShape(t *testing.T) {
	src := "## ADDED Requirements\n\n### Requirement: A\n\nIt SHALL hold.\n\n~~~\n## REMOVED Requirements\n~~~\n\n#### Scenario: S\n\n## REMOVED Requirements\n\n### Requirement: B\n\nGone.\n"
	if got := ChangeHeadings(src); !reflect.DeepEqual(got, []int{1, 13}) {
		t.Errorf("ChangeHeadings: %v", got)
	}
	for _, r := range Requirements(src) {
		if p := r.Problems(); p != nil {
			t.Errorf("%s: %v", r.Name, p)
		}
	}
	if got := (Requirement{Name: "C", Text: "Users can."}).Problems(); !reflect.DeepEqual(got, []string{`requirement "C" has no scenario`, `requirement "C" has no "shall"`}) {
		t.Errorf("Problems: %q", got)
	}
}
