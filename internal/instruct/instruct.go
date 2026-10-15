// Package instruct holds Grove's embedded templates and guides and builds
// the instructions a session follows.
package instruct

import "embed"

//go:embed templates
var files embed.FS

// Template returns the embedded template called name, such as "brief.md".
func Template(name string) string {
	b, err := files.ReadFile("templates/" + name)
	if err != nil {
		panic(err) // embedded at build time: a missing name is a programming error
	}
	return string(b)
}
