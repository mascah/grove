// Package doc reads Grove's documents through an fs.FS rooted at the
// document root: Markdown files, and the .grove.yaml beside them.
//
// Under the root are maintained specs (specs/<id>/, found by its spec.md),
// changes (changes/<id>-<slug>/, archived in changes/archive/, found by its
// .grove.yaml), and cycles (cycles/<id>-<slug>.yaml). Problems inside a document
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

// Meta is a document's metadata. The JSON tags are for grove changes --json.
type Meta struct {
	Title       string   `json:"title"`
	Created     string   `json:"created,omitempty"`
	Description string   `json:"description,omitempty"`
	Schema      string   `json:"schema,omitempty"`
	DependsOn   []string `json:"depends_on,omitempty"` // the changes this one waits on
	Base        string   `json:"base,omitempty"`       // the target commit this change was proposed from

	// A cycle's, which grove changes --json does not print.
	Goal       string            `json:"-"`
	Changes    []string          `json:"-"` // its members
	Planned    string            `json:"-"` // the target commit its plan checked
	NeedsHuman map[string]string `json:"-"` // member to the owner's reason it needs a person
	Unattended bool              `json:"-"` // handed to grove run
}

// Document is a maintained spec, a change or a cycle, with its metadata.
type Document struct {
	ID   string `json:"id"`
	Dir  string `json:"-"`    // relative to the document root; empty for a cycle
	Path string `json:"path"` // the main document, relative to the document root
	Meta
	Diagnostics []Diagnostic   `json:"-"`
	lines       map[string]int // the line of each field of its .grove.yaml
}

// Change is one change folder, changes/<id>-<slug>/ or, archived,
// changes/archive/<id>-<slug>/.
type Change struct {
	Document
	State   string            `json:"state"`              // approved, or archived or dropped once archived
	WaitsOn []string          `json:"waits_on,omitempty"` // what DependsOn names that is not archived as finished
	Tasks   []Task            `json:"tasks,omitempty"`    // from tasks.md, with their status from implementation.md
	Files   []string          `json:"-"`                  // every Markdown file in Dir and its .grove.yaml, relative to the root, sorted
	Touches map[string]string `json:"-"`                  // file → the capability it is about, for each file at a schema path with <capability>
	Dropped bool              `json:"-"`                  // implementation.md says Status: dropped
}

// Task is one T-NN heading of a change's tasks.md, with the status
// implementation.md gives it, if any (rule 3).
type Task struct {
	ID     string   `json:"id"`
	Title  string   `json:"title"`
	Status string   `json:"status,omitempty"`
	Files  []string `json:"-"` // the repo-relative paths its Files: line names, a folder ending in /
}

// Cycle is one cycle file: approved changes the owner means to build
// together, listed in its changes field.
type Cycle struct {
	Document
	Done    bool     // every member is archived, finished or dropped
	WaitsOn []string // what DependsOn names that is not a cycle done
}

// Names reports whether file, a path a Files: line names, is path or a
// folder above it.
func Names(file, path string) bool {
	return file == path || strings.HasSuffix(file, "/") && strings.HasPrefix(path, file)
}

// Collection is everything under one document root.
type Collection struct {
	Specs    []*Document
	Changes  []*Change
	Archived []*Change
	Cycles   []*Cycle // oldest first, by created

	Diagnostics []Diagnostic // problems in no document: a spec folder's leftover .grove.yaml
}

// Layout returns where change schema schema puts a change's main document
// and the files about one capability each, as paths in the change folder
// such as specs/<capability>/spec.md, or a problem when the name is not a
// schema's.
type Layout func(schema string) (main string, capabilities []string, problem string)

// Load reads the collection under fsys.
func Load(fsys fs.FS, layout Layout) (*Collection, error) {
	c := &Collection{}
	specs, err := fs.Glob(fsys, "specs/*/spec.md") // sorted
	if err != nil {
		return nil, err
	}
	for _, name := range specs {
		b, err := fs.ReadFile(fsys, name)
		if err != nil {
			return nil, err
		}
		d := &Document{ID: path.Base(path.Dir(name)), Dir: path.Dir(name), Path: name}
		if d.Meta.Title, d.Meta.Description = specMeta(string(b)); d.Meta.Title == "" {
			d.report(name, 0, `no "# " heading to give the spec its title`)
		}
		c.Specs = append(c.Specs, d)
	}
	leftovers, err := fs.Glob(fsys, "specs/*/.grove.yaml")
	if err != nil {
		return nil, err
	}
	for _, name := range leftovers {
		c.Diagnostics = append(c.Diagnostics, Diagnostic{Path: name, Message: "no longer read, since a spec's title and description come from its spec.md; delete it"})
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
			if ch.State = "approved"; strings.HasPrefix(name, "changes/archive/") {
				ch.State = map[bool]string{false: "archived", true: "dropped"}[ch.Dropped]
				c.Archived = append(c.Archived, ch)
			} else {
				c.Changes = append(c.Changes, ch)
			}
		}
	}
	for _, ch := range c.Changes {
		for _, id := range ch.Meta.DependsOn {
			if !c.Finished(id) {
				ch.WaitsOn = append(ch.WaitsOn, id)
			}
		}
	}
	c.dependencies()
	cycles, err := fs.Glob(fsys, "cycles/*.yaml")
	if err != nil {
		return nil, err
	}
	for _, name := range cycles {
		cy := &Cycle{Document: Document{ID: idOf(strings.TrimSuffix(path.Base(name), ".yaml")), Path: name}}
		b, err := fs.ReadFile(fsys, name)
		if err != nil {
			return nil, err
		}
		cy.lines = cy.parse(name, b, []string{"title", "created", "changes", "planned"}, []string{"goal", "depends_on", "needs_human", "unattended"})
		cy.Done = !slices.ContainsFunc(cy.Meta.Changes, func(id string) bool { return !c.InArchive(id) })
		c.Cycles = append(c.Cycles, cy)
	}
	slices.SortStableFunc(c.Cycles, func(a, b *Cycle) int { return cmp.Compare(a.Meta.Created, b.Meta.Created) })
	for _, cy := range c.Cycles {
		for _, id := range cy.Meta.DependsOn {
			if d := c.Cycle(id); d == nil || !d.Done {
				cy.WaitsOn = append(cy.WaitsOn, id)
			}
		}
	}
	c.members()
	c.follows()
	return c, nil
}

// follows reports each depends_on of a cycle not done that names no cycle,
// or a cycle that waits on it in turn, directly or through other cycles.
func (c *Collection) follows() {
	for _, cy := range c.Cycles {
		if cy.Done {
			continue
		}
		seen := map[*Cycle]bool{}
		var reaches func(from *Cycle) bool
		reaches = func(from *Cycle) bool {
			if from == cy {
				return true
			}
			if seen[from] {
				return false
			}
			seen[from] = true
			for _, id := range from.WaitsOn {
				if d := c.Cycle(id); d != nil && reaches(d) {
					return true
				}
			}
			return false
		}
		for _, id := range cy.Meta.DependsOn {
			switch d := c.Cycle(id); {
			case d == nil:
				cy.report(cy.Path, cy.lines["depends_on"], "depends_on: "+id+" is no cycle")
			case !d.Done && reaches(d):
				cy.report(cy.Path, cy.lines["depends_on"], "depends_on: "+id+" waits on this cycle in turn, so neither can run")
			}
		}
	}
}

// members reports each cycle's member that is no change, or that a cycle
// not done before it lists too, and each member needs_human names while the
// cycle is marked unattended.
func (c *Collection) members() {
	all, in := map[string]bool{}, map[string]string{} // in: change → the first cycle not done to list it
	for _, ch := range slices.Concat(c.Changes, c.Archived) {
		all[ch.ID] = true
	}
	for _, cy := range c.Cycles {
		for _, id := range cy.Meta.Changes {
			switch {
			case !all[id]:
				cy.report(cy.Path, cy.lines["changes"], "changes: "+id+" is no change, open or archived")
			case cy.Done:
			case in[id] != "":
				cy.report(cy.Path, cy.lines["changes"], "changes: "+id+" is also in cycle "+in[id]+", which is not done")
			default:
				in[id] = cy.ID
			}
			if _, ok := cy.Meta.NeedsHuman[id]; ok && cy.Meta.Unattended {
				cy.report(cy.Path, cy.lines["unattended"], "unattended: "+id+" needs a person, as needs_human says, so grove run skips this cycle; remove unattended, or build the cycle with /grove-implement "+cy.ID)
			}
		}
	}
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
	lines, err := ch.read(fsys, []string{"title", "created"}, "schema", "description", "depends_on", "base")
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

var (
	taskHeading = regexp.MustCompile(`^#+ (T-\d+)\b[\s—–:-]*(.*)$`)
	filesLine   = regexp.MustCompile("^`[^`]+`(\\s*,\\s*`[^`]+`)*$")
	fileBullet  = regexp.MustCompile("^\\s*[-*] `[^`]+`\\s*$")
	quoted      = regexp.MustCompile("`([^`]+)`")
	field       = regexp.MustCompile(`^[A-Z]\w*:`)
)

// tasks reads ch's tasks from tasks.md, with the files each names, and
// their status from implementation.md, files Grove's contract names in
// every schema (rule 3).
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
	lines := strings.Split(string(tasks), "\n")
	for i, l := range lines {
		if m := taskHeading.FindStringSubmatch(l); m != nil {
			ch.Tasks = append(ch.Tasks, Task{ID: m[1], Title: m[2], Status: status[m[1]]})
		}
		files, ok := strings.CutPrefix(l, "Files:")
		if !ok || len(ch.Tasks) == 0 {
			continue
		}
		bullets := strings.TrimSpace(files) == ""
		for j := i + 1; j < len(lines) && strings.TrimSpace(lines[j]) != "" && !field.MatchString(lines[j]); j++ { // it wraps until a blank line or the next field
			files += " " + lines[j]
			bullets = bullets && fileBullet.MatchString(lines[j])
		}
		t := &ch.Tasks[len(ch.Tasks)-1]
		files = strings.TrimSpace(files)
		for _, m := range quoted.FindAllStringSubmatch(files, -1) {
			t.Files = append(t.Files, m[1])
		}
		if !(bullets && len(t.Files) > 0 || filesLine.MatchString(files)) || slices.ContainsFunc(t.Files, func(f string) bool {
			f = strings.TrimSuffix(f, "/")
			return !fs.ValidPath(f) || f == "."
		}) {
			ch.report(ch.Dir+"/tasks.md", i+1, t.ID+"'s Files: line should list repo-relative paths in backticks, one per bullet or separated by commas, and a folder ends in /")
			t.Files = nil
		}
	}
}

// idOf returns the ID a change folder name starts with, or "".
func idOf(name string) string {
	if len(name) >= 12 && ValidID(name[:12]) && (len(name) == 12 || name[12] == '-') {
		return name[:12]
	}
	return ""
}

// Finished reports whether c holds change id archived as finished.
func (c *Collection) Finished(id string) bool {
	for _, ch := range c.Archived {
		if ch.ID == id {
			return !ch.Dropped
		}
	}
	return false
}

// Order returns changes in build order, each after the changes it depends
// on and otherwise oldest first, a dependency being met when finished
// reports it or an earlier change is it; then the changes that cannot be
// ordered, each waiting on a change that is neither.
func Order(changes []*Change, finished func(id string) bool) (ordered, waiting []*Change) {
	done := map[string]bool{}
	waiting = slices.SortedStableFunc(slices.Values(changes), func(a, b *Change) int { return cmp.Compare(a.Meta.Created, b.Meta.Created) })
	for len(waiting) > 0 {
		i := slices.IndexFunc(waiting, func(ch *Change) bool {
			return !slices.ContainsFunc(ch.Meta.DependsOn, func(id string) bool { return !done[id] && !finished(id) })
		})
		if i < 0 {
			break
		}
		ordered, done[waiting[i].ID] = append(ordered, waiting[i]), true
		waiting = slices.Delete(waiting, i, i+1)
	}
	return ordered, waiting
}

// InArchive reports whether c holds change id archived, finished or dropped.
func (c *Collection) InArchive(id string) bool {
	return slices.ContainsFunc(c.Archived, func(ch *Change) bool { return ch.ID == id })
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
