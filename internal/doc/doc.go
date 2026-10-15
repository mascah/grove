// Package doc reads Grove's documents: Markdown files with a small YAML
// frontmatter block, through an fs.FS rooted at the document root.
//
// Under the root are maintained specs (specs/<id>/spec.md) and changes
// (changes/<id>-<slug>/, identified by their proposal.md). Problems inside a
// document are Diagnostics, never errors, so one broken document does not
// hide the others.
package doc

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io/fs"
	"path"
	"regexp"
	"strings"
	"time"
)

// Diagnostic is one problem found in a document.
type Diagnostic struct {
	Path    string // relative to the document root
	Line    int    // 1-based; 0 when the problem has no line
	Message string
}

func (d Diagnostic) String() string {
	if d.Line > 0 {
		return fmt.Sprintf("%s:%d: %s", d.Path, d.Line, d.Message)
	}
	return d.Path + ": " + d.Message
}

// Meta is a document's frontmatter.
type Meta struct {
	ID          string `yaml:"id"`
	Title       string `yaml:"title"`
	Created     string `yaml:"created,omitempty"`
	Description string `yaml:"description,omitempty"`
}

// Document is one loaded Markdown file.
type Document struct {
	Path        string // relative to the document root
	Source      string
	Meta        Meta
	Diagnostics []Diagnostic
}

// Change is one change folder, changes/<id>-<slug>/.
type Change struct {
	Dir      string // relative to the document root
	Proposal *Document
	Files    []string // every Markdown file in Dir, relative to the root, sorted
}

// Collection is everything under one document root.
type Collection struct {
	Specs   []*Document
	Changes []*Change
}

// Required fields of each kind, in written order; description is optional
// everywhere.
var (
	specFields   = []string{"id", "title"}
	changeFields = []string{"id", "title", "created"}
)

// Load reads the collection under fsys.
func Load(fsys fs.FS) (*Collection, error) {
	c := &Collection{}
	var err error
	if c.Specs, err = load(fsys, "specs/*/spec.md", specFields); err != nil {
		return nil, err
	}
	proposals, err := load(fsys, "changes/*/proposal.md", changeFields)
	if err != nil {
		return nil, err
	}
	for _, p := range proposals {
		ch := &Change{Dir: path.Dir(p.Path), Proposal: p}
		err := fs.WalkDir(fsys, ch.Dir, func(name string, d fs.DirEntry, err error) error {
			if err == nil && !d.IsDir() && strings.HasSuffix(name, ".md") {
				ch.Files = append(ch.Files, name)
			}
			return err
		})
		if err != nil {
			return nil, err
		}
		c.Changes = append(c.Changes, ch)
	}
	return c, nil
}

func load(fsys fs.FS, pattern string, required []string) ([]*Document, error) {
	names, err := fs.Glob(fsys, pattern) // sorted
	if err != nil {
		return nil, err
	}
	var docs []*Document
	for _, name := range names {
		b, err := fs.ReadFile(fsys, name)
		if err != nil {
			return nil, err
		}
		docs = append(docs, Parse(name, b, required))
	}
	return docs, nil
}

// Change returns the change whose ID is id, or nil.
func (c *Collection) Change(id string) *Change {
	for _, ch := range c.Changes {
		if ch.Proposal.Meta.ID == id {
			return ch
		}
	}
	return nil
}

// idAlphabet is Crockford base32 in lowercase: digits and letters without
// the easily confused i, l, o and u.
const idAlphabet = "0123456789abcdefghjkmnpqrstvwxyz"

// NewID returns a fresh change ID: the UTC date of now as
// YYMMDD, a hyphen, and five random base32 characters, such as 261014-k3m9q.
func NewID(now time.Time) string {
	var raw [4]byte
	rand.Read(raw[:])
	bits := binary.BigEndian.Uint32(raw[:]) >> 7 // the top 25 of 32
	var suffix [5]byte
	for i := range suffix {
		suffix[i] = idAlphabet[bits>>(5*(4-i))&31]
	}
	return now.UTC().Format("060102") + "-" + string(suffix[:])
}

var idPattern = regexp.MustCompile(`^[0-9]{6}-[0-9a-hjkmnp-tv-z]{5}$`)

// ValidID reports whether s is a change ID.
func ValidID(s string) bool { return idPattern.MatchString(s) }
