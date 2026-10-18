// Package doc reads Grove's documents through an fs.FS rooted at the
// document root: Markdown files, and the .grove.yaml beside them.
//
// Under the root are maintained specs (specs/<id>/, identified by their
// spec.md) and changes (changes/<id>-<slug>/, identified by their
// proposal.md), each folder with a .grove.yaml. Problems inside a document
// are Diagnostics, never errors, so one broken document does not hide the
// others.
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

// Meta is a document's metadata.
type Meta struct{ Title, Created, Description string }

// Document is a maintained spec or a change: a folder, with its metadata.
type Document struct {
	ID          string
	Dir         string // relative to the document root
	Meta        Meta
	Diagnostics []Diagnostic
}

// Change is one change folder, changes/<id>-<slug>/.
type Change struct {
	Document
	Files []string // every Markdown file in Dir and its .grove.yaml, relative to the root, sorted
}

// Collection is everything under one document root.
type Collection struct {
	Specs   []*Document
	Changes []*Change
}

// Required fields of each kind, in written order; description is optional
// everywhere.
var (
	specFields   = []string{"title"}
	changeFields = []string{"title", "created"}
)

// Load reads the collection under fsys.
func Load(fsys fs.FS) (*Collection, error) {
	c := &Collection{}
	specs, err := fs.Glob(fsys, "specs/*/spec.md") // sorted
	if err != nil {
		return nil, err
	}
	for _, name := range specs {
		d := &Document{ID: path.Base(path.Dir(name)), Dir: path.Dir(name)}
		if err := d.readMeta(fsys, specFields, "spec's title"); err != nil {
			return nil, err
		}
		if err := d.heading(fsys, name); err != nil {
			return nil, err
		}
		c.Specs = append(c.Specs, d)
	}
	proposals, err := fs.Glob(fsys, "changes/*/proposal.md")
	if err != nil {
		return nil, err
	}
	for _, name := range proposals {
		ch := &Change{Document: Document{Dir: path.Dir(name)}}
		if folder := path.Base(ch.Dir); len(folder) >= 12 && ValidID(folder[:12]) && (len(folder) == 12 || folder[12] == '-') {
			ch.ID = folder[:12]
		} else {
			ch.report(ch.Dir+"/", 0, "the folder name must start with a change ID, such as 261015-k3m9q")
		}
		if err := ch.readMeta(fsys, changeFields, "change's title and created"); err != nil {
			return nil, err
		}
		if err := ch.heading(fsys, name); err != nil {
			return nil, err
		}
		err := fs.WalkDir(fsys, ch.Dir, func(name string, d fs.DirEntry, err error) error {
			if err == nil && !d.IsDir() && (strings.HasSuffix(name, ".md") || name == ch.Dir+"/.grove.yaml") {
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

// Change returns the change whose ID is id, or nil.
func (c *Collection) Change(id string) *Change {
	for _, ch := range c.Changes {
		if ch.ID == id {
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
