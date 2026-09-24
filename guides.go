// Package grove embeds the shared workflow guides and the reviewer definition,
// so a built binary carries the workflow of its own commit: the files stay the
// one editable owner, and grove guide prints the guides wherever it runs.
package grove

import "embed"

// Guides holds docs/work-execution.md and docs/work-shaping.md.
//
//go:embed docs/work-execution.md docs/work-shaping.md
var Guides embed.FS

// Reviewer is the grove-reviewer agent definition the work guide's step 6
// names, which init writes verbatim: this repository's copy is the one
// owner, and init here reports it unchanged.
//
//go:embed .claude/agents/grove-reviewer.md
var Reviewer string
