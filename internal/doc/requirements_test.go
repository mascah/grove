package doc

import (
	"reflect"
	"testing"
)

func TestRequirements(t *testing.T) {
	src := `# Changes to record-editing

## MODIFIED Requirements

### Requirement: Keep a category

When a user edits a record, the system shall keep its category.

#### Scenario: A note

- GIVEN a record

## REMOVED Requirements

### Requirement: Old rule

No longer wanted.

## Observed current behavior

Text.
`
	want := []Requirement{
		{Section: "MODIFIED", Name: "Keep a category", Line: 5, Text: "When a user edits a record, the system shall keep its category.\n\n#### Scenario: A note\n\n- GIVEN a record"},
		{Section: "REMOVED", Name: "Old rule", Line: 15, Text: "No longer wanted."},
	}
	if got := Requirements(src); !reflect.DeepEqual(got, want) {
		t.Errorf("Requirements:\n got %#v\nwant %#v", got, want)
	}
}
