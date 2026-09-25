// Package grove embeds the shared workflow guides, the record model they cite
// and the reviewer definition, so a built binary carries the workflow of its
// own commit: the files stay the one editable owner, and grove guide prints the
// guides and the model wherever it runs.
package grove

import "embed"

// Guides holds docs/work-execution.md, docs/work-shaping.md and the record
// model they cite, docs/record-model.md, which links to no record so it reads
// the same in any project.
//
//go:embed docs/work-execution.md docs/work-shaping.md docs/record-model.md
var Guides embed.FS

// Reviewer is the grove-reviewer agent definition the work guide's step 6
// names, which init writes verbatim: this repository's copy is the one
// owner, and init here reports it unchanged.
//
//go:embed .claude/agents/grove-reviewer.md
var Reviewer string
