// Package doc reads Grove's documents through an fs.FS rooted at the
// document root: Markdown files, and the .grove.yaml beside them.
//
// Under the root are maintained specs (specs/<id>/) and changes
// (changes/<id>-<slug>/, archived in changes/archive/), each folder found by
// its .grove.yaml. Problems inside a document
// are Diagnostics, never errors, so one broken document does not hide the
// others.
package doc

import (
	"cmp"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io/fs"
	"maps"
	"path"
	"regexp"
	"slices"
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
type Meta struct {
	Title, Created, Description, Schema string
	DependsOn                           []string // the changes this one waits on
	Cycle                               string   // the cycle this change belongs to
	Base                                string   // the target commit this change was proposed from
	Goal, Start, End                    string   // a cycle's
}

// Document is a maintained spec, a change or a cycle, with its metadata.
type Document struct {
	ID          string
	Dir         string // relative to the document root; empty for a cycle
	Path        string // the main document, relative to the document root
	Meta        Meta
	Diagnostics []Diagnostic
	lines       map[string]int // the line of each field of its .grove.yaml
}

// Change is one change folder, changes/<id>-<slug>/ or, archived,
// changes/archive/<id>-<slug>/.
type Change struct {
	Document
	Files   []string          // every Markdown file in Dir and its .grove.yaml, relative to the root, sorted
	Touches map[string]string // file → the capability it is about, for each file at a schema path with <capability>
	Tasks   []Task            // from tasks.md, with their status from implementation.md
	Dropped bool              // implementation.md says Status: dropped
	WaitsOn []string          // what DependsOn names that is not archived as finished
}

// Task is one T-NN heading of a change's tasks.md, with the status
// implementation.md gives it, if any (rule 3).
type Task struct{ ID, Title, Status string }

// Cycle is one cycle file, cycles/<id>-<slug>.yaml. The changes that name
// it in their cycle field are in it.
type Cycle struct {
	Document
	Done bool // every change in it is archived
}

// Collection is everything under one document root.
type Collection struct {
	Specs    []*Document
	Changes  []*Change
	Archived []*Change
	Cycles   []*Cycle
}

// Layout returns where change schema schema puts a change's main document
// and the files about one capability each, as paths in the change folder
// such as specs/<capability>/spec.md, or a problem when the name is not a
// schema's.
type Layout func(schema string) (main string, capabilities []string, problem string)

// Load reads the collection under fsys.
func Load(fsys fs.FS, layout Layout) (*Collection, error) {
	c := &Collection{}
	specs, err := fs.Glob(fsys, "specs/*/.grove.yaml") // sorted
	if err != nil {
		return nil, err
	}
	for _, name := range specs {
		d := &Document{ID: path.Base(path.Dir(name)), Dir: path.Dir(name), Path: path.Dir(name) + "/spec.md"}
		if _, err := d.read(fsys, []string{"title"}, "description"); err != nil {
			return nil, err
		}
		if err := d.heading(fsys, d.Dir+"/spec.md", "missing; it is the spec's main document"); err != nil {
			return nil, err
		}
		c.Specs = append(c.Specs, d)
	}
	for _, glob := range []string{"changes/*/.grove.yaml", "changes/archive/*/.grove.yaml"} {
		metas, err := fs.Glob(fsys, glob)
		if err != nil {
			return nil, err
		}
		for _, name := range metas {
			ch, err := loadChange(fsys, name, layout)
			if err != nil {
				return nil, err
			}
			if strings.HasPrefix(name, "changes/archive/") {
				c.Archived = append(c.Archived, ch)
			} else {
				c.Changes = append(c.Changes, ch)
			}
		}
	}
	cycles, err := fs.Glob(fsys, "cycles/*.yaml")
	if err != nil {
		return nil, err
	}
	for _, name := range cycles {
		cy := &Cycle{Document: Document{ID: idOf(path.Base(name)), Path: name}}
		if cy.ID == "" {
			cy.report(name, 0, "the file name must start with a cycle ID, such as 261015-k3m9q")
		}
		b, err := fs.ReadFile(fsys, name)
		if err != nil {
			return nil, err
		}
		cy.lines = cy.parse(name, b, []string{"title"}, []string{"goal", "start", "end"})
		cy.Done = !slices.ContainsFunc(c.Changes, func(ch *Change) bool { return ch.Meta.Cycle != "" && ch.Meta.Cycle == cy.ID })
		c.Cycles = append(c.Cycles, cy)
	}
	finished := map[string]bool{}
	for _, ch := range c.Archived {
		finished[ch.ID] = !ch.Dropped
	}
	for _, ch := range c.Changes {
		for _, id := range ch.Meta.DependsOn {
			if !finished[id] {
				ch.WaitsOn = append(ch.WaitsOn, id)
			}
		}
		if ch.Meta.Cycle != "" && c.Cycle(ch.Meta.Cycle) == nil {
			ch.report(ch.Dir+"/.grove.yaml", ch.lines["cycle"], "cycle: "+ch.Meta.Cycle+" names no file in cycles/")
		}
	}
	c.dependencies()
	return c, nil
}

// dependencies reports each open change's depends_on that names no change,
// a dropped change, or one that waits on it in turn.
func (c *Collection) dependencies() {
	open, all := map[string]*Change{}, map[string]*Change{}
	for _, ch := range c.Changes {
		open[ch.ID], all[ch.ID] = ch, ch
	}
	for _, ch := range c.Archived {
		all[ch.ID] = ch
	}
	for _, ch := range c.Changes {
		// reaches reports whether from waits on ch, directly or through other open changes.
		seen := map[*Change]bool{}
		var reaches func(from *Change) bool
		reaches = func(from *Change) bool {
			if from == ch {
				return true
			}
			if seen[from] {
				return false
			}
			seen[from] = true
			for _, id := range from.Meta.DependsOn {
				if d := open[id]; d != nil && reaches(d) {
					return true
				}
			}
			return false
		}
		for _, id := range ch.Meta.DependsOn {
			switch d := all[id]; {
			case d == nil:
				ch.report(ch.Dir+"/.grove.yaml", ch.lines["depends_on"], "depends_on: "+id+" is no change, open or archived")
			case d.Dropped:
				ch.report(ch.Dir+"/.grove.yaml", ch.lines["depends_on"], "depends_on: "+id+" was dropped, so this change could never be applied")
			case open[id] != nil && reaches(d):
				ch.report(ch.Dir+"/.grove.yaml", ch.lines["depends_on"], "depends_on: "+id+" waits on this change in turn, so neither can be applied")
			}
		}
	}
}

func loadChange(fsys fs.FS, name string, layout Layout) (*Change, error) {
	ch := &Change{Document: Document{ID: idOf(path.Base(path.Dir(name))), Dir: path.Dir(name)}, Touches: map[string]string{}}
	if ch.ID == "" {
		ch.report(ch.Dir+"/", 0, "the folder name must start with a change ID, such as 261015-k3m9q")
	}
	lines, err := ch.read(fsys, []string{"title", "created"}, "schema", "description", "depends_on", "cycle", "base")
	if err != nil {
		return nil, err
	}
	ch.lines = lines
	ch.Meta.Schema = cmp.Or(ch.Meta.Schema, "default") // changes made before schemas
	main, capabilities, problem := layout(ch.Meta.Schema)
	switch {
	case problem != "":
		ch.report(name, lines["schema"], problem)
	case main != "":
		ch.Path = ch.Dir + "/" + main
		if err := ch.heading(fsys, ch.Dir+"/"+main, "missing; it is the change's main document in schema "+ch.Meta.Schema); err != nil {
			return nil, err
		}
	}
	err = fs.WalkDir(fsys, ch.Dir, func(name string, d fs.DirEntry, err error) error {
		if err == nil && !d.IsDir() && (strings.HasSuffix(name, ".md") || name == ch.Dir+"/.grove.yaml") {
			ch.Files = append(ch.Files, name)
			for _, p := range capabilities {
				if c := capability(p, strings.TrimPrefix(name, ch.Dir+"/")); c != "" {
					ch.Touches[name] = c
				}
			}
		}
		return err
	})
	ch.tasks(fsys)
	return ch, err
}

// capability returns the capability that file, in a change folder, is
// about when it is at pattern, such as specs/<capability>/spec.md, or "".
func capability(pattern, file string) string {
	p, f := strings.Split(pattern, "/"), strings.Split(file, "/")
	if len(p) != len(f) {
		return ""
	}
	c := ""
	for i := range p {
		if p[i] == "<capability>" {
			c = f[i]
		} else if p[i] != f[i] {
			return ""
		}
	}
	return c
}

// Capabilities returns the capabilities ch has a file about, sorted.
func (ch *Change) Capabilities() []string {
	return slices.Compact(slices.Sorted(maps.Values(ch.Touches)))
}

var taskHeading = regexp.MustCompile(`^#+ (T-\d+)\b[\s—–:-]*(.*)$`)

// tasks reads ch's tasks from tasks.md and their status from
// implementation.md, files Grove's contract names in every schema (rule 3).
func (ch *Change) tasks(fsys fs.FS) {
	tasks, _ := fs.ReadFile(fsys, ch.Dir+"/tasks.md")
	impl, _ := fs.ReadFile(fsys, ch.Dir+"/implementation.md")
	status, task := map[string]string{}, ""
	for _, l := range strings.Split(string(impl), "\n") {
		if m := taskHeading.FindStringSubmatch(l); m != nil {
			task = m[1]
		}
		if s, ok := strings.CutPrefix(l, "Status: "); ok {
			status[task] = strings.TrimSpace(s)
			ch.Dropped = ch.Dropped || status[task] == "dropped"
		}
	}
	for _, l := range strings.Split(string(tasks), "\n") {
		if m := taskHeading.FindStringSubmatch(l); m != nil {
			ch.Tasks = append(ch.Tasks, Task{ID: m[1], Title: m[2], Status: status[m[1]]})
		}
	}
}

// idOf returns the ID a change folder or cycle file name starts with, or "".
func idOf(name string) string {
	if len(name) >= 12 && ValidID(name[:12]) && (len(name) == 12 || name[12] == '-' || name[12] == '.') {
		return name[:12]
	}
	return ""
}

// Cycle returns the cycle whose ID is id, or nil.
func (c *Collection) Cycle(id string) *Cycle {
	for _, cy := range c.Cycles {
		if cy.ID == id {
			return cy
		}
	}
	return nil
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
