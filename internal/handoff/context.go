// Package handoff assembles read-only context for explicitly selected work
// from one checkout. The selected records, the configuration, and the files
// the caller names are read in full with their revisions; prerequisites,
// blocking questions, related records, and links are listed compactly so the
// caller can retrieve each when its activity needs it. It reports facts. It
// does not authorize work, decide readiness, or start anything.
package handoff

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"
	"slices"
	"strings"

	"github.com/mascah/grove/internal/project"
)

const (
	DefaultMaxBytes = 262144
	LimitMaxBytes   = 8388608

	scopeNotice = "Read in full, with revisions: grove.yaml, the selected work, and caller includes. " +
		"Listed only: transitive depends_on prerequisites, questions blocking the selected work or a prerequisite, " +
		"records the selected work names in relates_to or members or links to, and every link in the selected work's bodies. " +
		"A listed record's title, status, and revision are what this checkout held; its constraints are in its body, which show ID prints. " +
		"A listed link was not opened or checked, not even for existence; --include PATH adds a file with its revision and fails if it is missing. " +
		"Nothing here was summarized or truncated, and a listing is not a reading: read what the current activity depends on before acting on it. " +
		"A status is what a record says in this checkout: done is not integration, and assembled context is not readiness, acceptance, or authorization."
)

type Options struct {
	Interaction string   // "interactive" (default) or "headless": the caller's declared intent, not a detected capability
	MaxBytes    int      // budget for unique included source bytes; 0 selects DefaultMaxBytes
	Include     []string // clean project-relative files the caller requires
}

type Source struct {
	Path     string   `json:"path"`
	Revision string   `json:"revision"`
	Reasons  []string `json:"reasons"`
	Content  string   `json:"content"`
}

// Record is an observation of a record the loader read and validated. Its
// full source is among the sources only when Included says so.
type Record struct {
	ID       string   `json:"id"`
	Path     string   `json:"path"`
	Type     string   `json:"type"`
	Title    string   `json:"title"`
	Status   string   `json:"status"`
	Revision string   `json:"revision"` // of the record as loaded, included or not
	Roles    []string `json:"roles"`    // why it is listed
	Selected bool     `json:"selected"`
	Included bool     `json:"included"`
}

// Requirement is one depends_on edge of selected or prerequisite work, with
// the status the prerequisite records in this checkout.
type Requirement struct {
	Work         string `json:"work"`
	Prerequisite string `json:"prerequisite"`
	Status       string `json:"status"`
	Selected     bool   `json:"selected"`
}

type Question struct {
	ID     string   `json:"id"`
	Status string   `json:"status"`
	Blocks []string `json:"blocks"` // the selected or prerequisite work it blocks
}

// Reference is a link in a selected record's body. Path is the project file it
// resolves to, or "" when it is not an in-project document. Unless Reason says
// the file is included, nothing was opened: the file may not exist.
type Reference struct {
	From   string `json:"from"`
	Target string `json:"target"` // as the record wrote it
	Path   string `json:"path"`
	Reason string `json:"reason"`
}

type Git struct {
	Checkout  string `json:"checkout"`
	CommonDir string `json:"common_dir"`
	Ref       string `json:"ref"`  // full ref name; "" when HEAD is detached
	Head      string `json:"head"` // "" when the branch has no commit yet
}

type Bundle struct {
	FormatVersion int           `json:"format_version"` // of this output, not the record schema
	Root          string        `json:"root"`
	Interaction   string        `json:"interaction"`
	Selected      []string      `json:"selected"` // as requested
	Order         []string      `json:"order"`    // the same IDs, prerequisites first
	Git           *Git          `json:"git"`      // null outside a Git repository
	Records       []Record      `json:"records"`
	Requirements  []Requirement `json:"requirements"`
	Questions     []Question    `json:"questions"`
	Sources       []Source      `json:"sources"`
	References    []Reference   `json:"references"`
	ScopeNotice   string        `json:"scope_notice"`
	SourceBytes   int           `json:"source_bytes"`
	MaxBytes      int           `json:"max_bytes"`
}

// betweenReads lets tests change the checkout between Build's two readings.
var betweenReads = func() {}

// Build reads everything twice and refuses a checkout that changed in between.
// That is an optimistic check, not a snapshot or a lease: whoever acts on the
// bundle revalidates the revisions it names. Nothing is written or started.
func Build(ctx context.Context, root string, ids []string, opts Options) (*Bundle, error) {
	if opts.Interaction == "" {
		opts.Interaction = "interactive"
	}
	if opts.Interaction != "interactive" && opts.Interaction != "headless" {
		return nil, fmt.Errorf("interaction must be interactive or headless, not %q", opts.Interaction)
	}
	if opts.MaxBytes == 0 {
		opts.MaxBytes = DefaultMaxBytes
	}
	if opts.MaxBytes < 1 || opts.MaxBytes > LimitMaxBytes {
		return nil, fmt.Errorf("the source budget must be 1 through %d bytes", LimitMaxBytes)
	}
	// The held directory confines every document read, and proves afterwards
	// that root still names the directory that was read.
	dir, err := os.OpenRoot(root)
	if err != nil {
		return nil, err
	}
	defer dir.Close()
	first, err := assemble(ctx, dir, root, ids, opts)
	if err != nil {
		return nil, err
	}
	betweenReads()
	second, err := assemble(ctx, dir, root, ids, opts)
	if err != nil {
		return nil, fmt.Errorf("the context could not be read a second time to confirm it (%w); rerun to read it again", err)
	}
	held, err := dir.Stat(".")
	if err != nil {
		return nil, err
	}
	if now, err := os.Stat(root); err != nil || !os.SameFile(held, now) {
		return nil, fmt.Errorf("%s was replaced while its context was being read; rerun to read it again", root)
	}
	if !reflect.DeepEqual(first, second) {
		return nil, fmt.Errorf("the checkout changed while its context was being read (%s); rerun to read it again", difference(first, second))
	}
	return second, nil
}

func difference(a, b *Bundle) string {
	revisions := map[string]string{}
	for _, s := range a.Sources {
		revisions[s.Path] = s.Revision
	}
	for _, s := range b.Sources {
		if revisions[s.Path] != s.Revision {
			return s.Path
		}
		delete(revisions, s.Path)
	}
	for path := range revisions {
		return path
	}
	return "Git state or record relationships"
}

func assemble(ctx context.Context, dir *os.Root, root string, ids []string, opts Options) (*Bundle, error) {
	p, ds := project.Load(root, root)
	if len(ds) != 0 {
		lines := make([]string, len(ds))
		for i, d := range ds {
			lines[i] = d.String()
		}
		return nil, errors.New(strings.Join(lines, "\n"))
	}
	byID := map[string]*project.Record{}
	for _, r := range p.Records {
		byID[r.ID] = r
	}
	order, reach, err := selection(byID, ids)
	if err != nil {
		return nil, err
	}
	git, err := gitIdentity(ctx, root)
	if err != nil {
		return nil, err
	}
	b := &Bundle{
		FormatVersion: 2, Root: root, Interaction: opts.Interaction, Selected: ids, Order: order, Git: git,
		Requirements: []Requirement{}, Questions: []Question{}, References: []Reference{},
		ScopeNotice: scopeNotice, MaxBytes: opts.MaxBytes,
	}
	s := &sources{dir: dir, remaining: opts.MaxBytes, byPath: map[string]*Source{}}
	if err := s.addLoaded("grove.yaml", p.Config, "project configuration"); err != nil {
		return nil, err
	}

	// roles records why each record is listed. Only the selected work is read in
	// full by default; everything else is an observation until the caller asks.
	roles := map[string][]string{}
	role := func(id, why string) {
		if !slices.Contains(roles[id], why) {
			roles[id] = append(roles[id], why)
		}
	}
	for _, id := range ids {
		role(id, "selected work")
		if err := s.addLoaded(byID[id].Path, byID[id].Source, "selected work"); err != nil {
			return nil, err
		}
	}
	// scope is the selected work plus everything it transitively depends on.
	var scope []*project.Record
	for _, r := range p.Records {
		var needs []string
		for _, id := range ids {
			if reach[id][r.ID] {
				needs = append(needs, id)
			}
		}
		if needs != nil {
			role(r.ID, "prerequisite of "+strings.Join(needs, ", "))
		}
		if needs != nil || slices.Contains(ids, r.ID) {
			scope = append(scope, r)
		}
	}
	for _, r := range scope {
		for _, id := range r.DependsOn {
			b.Requirements = append(b.Requirements, Requirement{
				Work: r.ID, Prerequisite: id, Status: byID[id].Status, Selected: slices.Contains(ids, id),
			})
		}
	}
	for _, r := range p.Records {
		var blocks []string
		for _, w := range scope {
			if slices.Contains(r.Blocks, w.ID) {
				blocks = append(blocks, w.ID)
			}
		}
		if blocks != nil {
			b.Questions = append(b.Questions, Question{ID: r.ID, Status: r.Status, Blocks: blocks})
			role(r.ID, "question blocking "+strings.Join(blocks, ", "))
		}
	}
	// Relationships of the selected work are context, never more selection.
	for _, id := range ids {
		for _, related := range byID[id].RelatesTo {
			role(related, "related to "+id)
		}
		for _, member := range byID[id].Members {
			role(member, "member of "+id)
		}
	}
	for _, name := range opts.Include {
		if err := s.addInclude(name); err != nil {
			return nil, err
		}
	}
	atPath := map[string]*project.Record{}
	for _, r := range p.Records {
		atPath[r.Path] = r
	}
	for _, id := range ids {
		refs, err := s.references(byID[id])
		if err != nil {
			return nil, err
		}
		for _, ref := range refs {
			if linked := atPath[ref.Path]; linked != nil {
				role(linked.ID, "linked from "+id)
			}
		}
		b.References = append(b.References, refs...)
	}

	for _, r := range p.Records {
		if roles[r.ID] != nil {
			b.Records = append(b.Records, Record{
				ID: r.ID, Path: r.Path, Type: r.Type, Title: r.Title, Status: r.Status, Revision: project.Revision(r.Source),
				Roles: roles[r.ID], Selected: slices.Contains(ids, r.ID), Included: s.holds(r.Path),
			})
		}
	}
	for _, source := range s.byPath {
		b.Sources = append(b.Sources, *source)
	}
	slices.SortFunc(b.Sources, func(x, y Source) int { return strings.Compare(x.Path, y.Path) })
	b.SourceBytes = opts.MaxBytes - s.remaining
	return b, nil
}

// selection validates the requested IDs and orders them so that each follows
// the selected work it depends on, directly or through unselected work; ties
// keep the requested order. reach maps each ID to its transitive prerequisites.
// Load has already refused dependency cycles and non-work targets.
func selection(byID map[string]*project.Record, ids []string) (order []string, reach map[string]map[string]bool, err error) {
	if len(ids) == 0 {
		return nil, nil, errors.New("context requires at least one work ID")
	}
	reach = map[string]map[string]bool{}
	for _, id := range ids {
		r := byID[id]
		switch {
		case reach[id] != nil:
			return nil, nil, fmt.Errorf("%s is selected more than once", id)
		case r == nil:
			return nil, nil, fmt.Errorf("work %s is not in this checkout; context takes explicit work IDs and reads only the selected checkout", id)
		case r.Type != "work":
			return nil, nil, fmt.Errorf("%s is a %s; only work can be selected (related records are included as context)", id, r.Type)
		}
		reach[id] = map[string]bool{}
		var visit func(*project.Record)
		visit = func(r *project.Record) {
			for _, next := range r.DependsOn {
				if !reach[id][next] {
					reach[id][next] = true
					visit(byID[next])
				}
			}
		}
		visit(r)
	}
	placed := map[string]bool{}
	for len(order) < len(ids) {
		before := len(order)
		for _, id := range ids {
			ready := !placed[id]
			for _, other := range ids {
				ready = ready && (!reach[id][other] || placed[other])
			}
			if ready {
				placed[id] = true
				order = append(order, id)
				break
			}
		}
		if len(order) == before { // Load refuses cycles; never spin if one gets here
			return nil, nil, errors.New("the selected work depends on itself in a cycle")
		}
	}
	return order, reach, nil
}
