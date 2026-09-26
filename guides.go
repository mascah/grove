// Package grove embeds the shared workflow guides, the review guide and the
// record model they cite, so a built binary carries the workflow of its own
// commit: the files stay the one editable owner, and grove guide prints them
// wherever it runs. It also owns the entrypoint templates init writes and the
// build identity that names all of them.
package grove

import "embed"

// Guides holds docs/work-execution.md, docs/work-shaping.md,
// docs/work-review.md and the record model they cite, docs/record-model.md,
// which names no record so it reads the same in any project.
//
//go:embed docs/work-execution.md docs/work-shaping.md docs/work-review.md docs/record-model.md
var Guides embed.FS

// GuideFiles maps each guide name grove guide takes to its file in Guides.
var GuideFiles = map[string]string{"work": "docs/work-execution.md", "shape": "docs/work-shaping.md", "review": "docs/work-review.md", "model": "docs/record-model.md"}
