// Package grove embeds the shared workflow guides, so a built binary carries
// the workflow of its own commit: the files under docs/ stay the one editable
// owner, and grove guide prints them wherever the binary runs.
package grove

import "embed"

// Guides holds docs/work-execution.md and docs/work-shaping.md.
//
//go:embed docs/work-execution.md docs/work-shaping.md
var Guides embed.FS
