// Package update changes one record's frontmatter from the CLI while
// preserving the rest of the file and refusing stale or unsafe writes.
package update

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/repo"
)

// Request is one update to one record. Set entries keep request order.
type Request struct {
	ID, Expect string
	Set        []Field
	Unset      []string
}

type Field struct{ Name, Value string }

// Result describes the record after the request. Revision identifies the
// resulting bytes; for a no-op it equals the expected revision.
type Result struct {
	ID, Path, Revision string
	Changed            bool
}

// Failure retains publication state when an error happens after the rename
// succeeded, so a caller cannot mistake an applied update for a rejected one.
type Failure struct {
	Path, Revision string
	Err            error
}

func (f *Failure) Error() string {
	return fmt.Sprintf("%s (the update was applied to %s; revision %s)", f.Err, f.Path, f.Revision)
}

func (f *Failure) Unwrap() error { return f.Err }

// Fault lets tests inject failures or external writes at named steps: write,
// sync, close, compare, rename, dirsync, validate. A nil Fault never fails.
type Fault func(step string) error

var digits = regexp.MustCompile(`^[0-9]+$`)
var word = regexp.MustCompile(`^[a-z]+$`)

// Apply performs one update under the repository-wide write lock. Nothing is
// written unless the whole project, with the candidate substituted, validates.
func Apply(root string, req Request, now time.Time, fault Fault) (Result, error) {
	if fault == nil {
		fault = func(string) error { return nil }
	}
	common, _, err := repo.CommonDir(root)
	if err != nil {
		return Result{}, err
	}
	unlock, err := repo.WriteLock(common)
	if err != nil {
		return Result{}, err
	}
	defer unlock()
	p, ds := project.Load(root, root)
	if len(ds) != 0 {
		return Result{}, fmt.Errorf("the project is not valid; fix it before updating:\n%s", diagnostics(ds))
	}
	i := slices.IndexFunc(p.Records, func(r *project.Record) bool { return r.ID == req.ID })
	if i < 0 {
		return Result{}, fmt.Errorf("record %s not found in this project", req.ID)
	}
	r := p.Records[i]
	current := project.Revision(r.Source)
	if current != req.Expect {
		return Result{}, fmt.Errorf("%s changed since the expected revision; its current revision is %s", r.ID, current)
	}
	changes, err := plan(r, req)
	if err != nil {
		return Result{}, err
	}
	result := Result{ID: r.ID, Path: r.Path, Revision: current}
	if len(changes) == 0 {
		return result, nil
	}
	stamp := now.UTC().Truncate(time.Second)
	for _, existing := range []struct {
		name string
		date *time.Time
	}{{"created", r.Created}, {"updated", r.Updated}} {
		if existing.date != nil && stamp.Before(*existing.date) {
			return Result{}, fmt.Errorf("clock/date inconsistency: the current time %s precedes %s's %s %s; refusing to write an earlier updated time",
				stamp.Format(time.RFC3339), r.ID, existing.name, existing.date.Format(time.RFC3339))
		}
	}
	changes = append(changes, set("updated", strconv.Quote(stamp.Format("2006-01-02T15:04:05Z"))))
	candidate, err := Edit(r.Source, changes)
	if err != nil {
		return Result{}, fmt.Errorf("%s: %w", r.Path, err)
	}
	next, ds := project.ParseRecord(r.Path, r.Type, candidate)
	if len(ds) != 0 {
		return Result{}, fmt.Errorf("the update would leave %s invalid:\n%s", r.ID, diagnostics(ds))
	}
	if err := unchanged(r, next, changes); err != nil {
		return Result{}, err
	}
	records := slices.Clone(p.Records)
	records[i] = next
	if ds := project.Validate(records); len(ds) != 0 {
		return Result{}, fmt.Errorf("the update would leave the project invalid:\n%s", diagnostics(ds))
	}
	result.Revision = project.Revision(candidate)
	result.Changed = true
	return result, publish(root, p, i, candidate, fault)
}

// plan validates the request against the record's type and drops fields whose
// parsed meaning already matches, so a no-op never rewrites the file.
func plan(r *project.Record, req Request) ([]change, error) {
	allowed := map[string][]string{
		"work":     {"title", "status", "relates_to", "kind", "priority", "size", "members", "depends_on"},
		"question": {"title", "status", "relates_to", "blocks"},
		"decision": {"title", "status", "relates_to"},
	}[r.Type]
	lists := map[string][]string{"relates_to": r.RelatesTo, "members": r.Members, "depends_on": r.DependsOn, "blocks": r.Blocks}
	strs := map[string]string{"title": r.Title, "status": r.Status, "kind": r.Kind, "size": r.Size}
	check := func(name string) error {
		if slices.Contains([]string{"id", "type", "created", "updated"}, name) {
			return fmt.Errorf("%s cannot be changed by update", name)
		}
		if !slices.Contains(allowed, name) {
			return fmt.Errorf("%s is not a field that update accepts on %s records", name, r.Type)
		}
		return nil
	}
	var changes []change
	for _, f := range req.Set {
		if err := check(f.Name); err != nil {
			return nil, err
		}
		if !utf8.ValidString(f.Value) {
			return nil, fmt.Errorf("%s: value is not valid UTF-8", f.Name)
		}
		switch f.Name {
		case "priority":
			if !digits.MatchString(f.Value) {
				return nil, fmt.Errorf("priority must be decimal digits for an integer 1 through 5")
			}
			n, err := strconv.Atoi(f.Value)
			if err != nil {
				return nil, fmt.Errorf("priority: %w", err)
			}
			if r.Priority != nil && *r.Priority == n {
				continue
			}
			changes = append(changes, set("priority", strconv.Itoa(n)))
		case "relates_to", "members", "depends_on", "blocks":
			var ids []string
			if !strings.HasPrefix(strings.TrimSpace(f.Value), "[") || json.Unmarshal([]byte(f.Value), &ids) != nil {
				return nil, fmt.Errorf("%s must be a JSON array of record ID strings, such as [\"W-001\"]", f.Name)
			}
			if ids == nil {
				ids = []string{}
			}
			if lists[f.Name] != nil && slices.Equal(lists[f.Name], ids) {
				continue
			}
			quoted := make([]string, len(ids))
			for i, id := range ids {
				quoted[i] = strconv.Quote(id)
			}
			changes = append(changes, set(f.Name, "["+strings.Join(quoted, ", ")+"]"))
		default:
			if strings.TrimSpace(f.Value) == "" {
				return nil, fmt.Errorf("%s requires a nonempty string", f.Name)
			}
			if strs[f.Name] == f.Value {
				continue
			}
			value := strconv.Quote(f.Value)
			if f.Name != "title" && word.MatchString(f.Value) {
				value = f.Value // enumerated words stay plain; anything else is quoted for the schema check to reject
			}
			changes = append(changes, set(f.Name, value))
		}
	}
	for _, name := range req.Unset {
		if err := check(name); err != nil {
			return nil, err
		}
		if name == "title" || name == "status" {
			return nil, fmt.Errorf("%s is required and cannot be unset", name)
		}
		present := false
		switch name {
		case "priority":
			present = r.Priority != nil
		case "relates_to", "members", "depends_on", "blocks":
			present = lists[name] != nil
		default:
			present = strs[name] != ""
		}
		if present {
			changes = append(changes, unset(name))
		}
	}
	return changes, nil
}

// unchanged refuses a candidate whose untouched fields differ from the
// original, which would mean the editor altered an unrelated entry.
func unchanged(before, after *project.Record, changes []change) error {
	touched := map[string]bool{"updated": true}
	for _, c := range changes {
		touched[c.key] = true
	}
	b, a := fields(before), fields(after)
	for name := range b {
		if !touched[name] && a[name] != b[name] {
			return fmt.Errorf("%s: editing changed %s unexpectedly; refusing to write", before.Path, name)
		}
	}
	return nil
}

// fields renders every schema field so absent and empty values differ.
func fields(r *project.Record) map[string]string {
	list := func(ids []string) string { return fmt.Sprint(ids == nil, ids) }
	priority := "absent"
	if r.Priority != nil {
		priority = strconv.Itoa(*r.Priority)
	}
	created := "absent"
	if r.Created != nil {
		created = r.Created.Format(time.RFC3339)
	}
	return map[string]string{
		"id": r.ID, "type": r.Type, "title": r.Title, "status": r.Status, "kind": r.Kind, "size": r.Size,
		"priority": priority, "created": created,
		"relates_to": list(r.RelatesTo), "members": list(r.Members), "depends_on": list(r.DependsOn), "blocks": list(r.Blocks),
	}
}

// publish writes candidate beside the target, re-checks that nothing observed
// changed since the validated snapshot, and renames it into place. Errors
// before the rename leave the record untouched; errors after it are reported
// as *Failure with the applied revision. Old bytes are never restored: a later
// direct edit could be lost.
func publish(root string, snapshot *project.Project, idx int, candidate []byte, fault Fault) error {
	target := snapshot.Records[idx]
	full := filepath.Join(root, filepath.FromSlash(target.Path))
	before, err := os.Lstat(full)
	if err != nil {
		return err
	}
	if !before.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", target.Path)
	}
	dir := filepath.Dir(full)
	tmp, err := os.CreateTemp(dir, "."+filepath.Base(full)+".*.tmp")
	if err != nil {
		return fmt.Errorf("prepare temporary file: %w", err)
	}
	name := tmp.Name()
	discard := func(step string, err error) error {
		tmp.Close()
		os.Remove(name)
		return fmt.Errorf("%s: %w; %s is unchanged", step, err, target.Path)
	}
	for _, step := range []struct {
		name, fault string
		run         func() error
	}{
		{"write temporary file", "write", func() error { _, err := tmp.Write(candidate); return err }},
		{"sync temporary file", "sync", tmp.Sync},
		{"set permissions", "", func() error { return tmp.Chmod(before.Mode().Perm()) }},
		{"close temporary file", "close", tmp.Close},
		{"compare with the validated snapshot", "compare", func() error { return same(root, snapshot) }},
		{"verify the target file", "", func() error {
			after, err := os.Lstat(full)
			if err != nil {
				return err
			}
			if !after.Mode().IsRegular() || !os.SameFile(before, after) || after.Mode().Perm() != before.Mode().Perm() {
				return errors.New("the file was replaced or its permissions changed")
			}
			return nil
		}},
		{"rename into place", "rename", func() error { return os.Rename(name, full) }},
	} {
		err := fault(step.fault)
		if err == nil {
			err = step.run()
		}
		if err != nil {
			return discard(step.name, err)
		}
	}
	applied := &Failure{Path: target.Path, Revision: project.Revision(candidate)}
	err = fault("dirsync")
	if err == nil {
		err = syncDir(dir)
	}
	if err != nil {
		applied.Err = fmt.Errorf("directory sync failed, so durability is uncertain: %w", err)
		return applied
	}
	if err := fault("validate"); err != nil {
		applied.Err = err
		return applied
	}
	if _, ds := project.Load(root, root); len(ds) != 0 {
		applied.Err = fmt.Errorf("the project no longer validates:\n%s", diagnostics(ds))
		return applied
	}
	return nil
}

// same reports whether configuration, record inventory, and record bytes still
// match the snapshot the candidate was validated against.
func same(root string, snapshot *project.Project) error {
	p, ds := project.Load(root, root)
	if len(ds) != 0 {
		return fmt.Errorf("the project changed and no longer validates:\n%s", diagnostics(ds))
	}
	changed := p.RecordDir != snapshot.RecordDir || len(p.Records) != len(snapshot.Records)
	for i := 0; !changed && i < len(p.Records); i++ {
		changed = p.Records[i].Path != snapshot.Records[i].Path || !bytes.Equal(p.Records[i].Source, snapshot.Records[i].Source)
	}
	if changed {
		return errors.New("the project changed while the update was being prepared; inspect it and retry with a fresh revision")
	}
	return nil
}

func syncDir(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	return d.Sync()
}

func diagnostics(ds []project.Diagnostic) string {
	lines := make([]string, len(ds))
	for i, d := range ds {
		lines[i] = d.String()
	}
	return strings.Join(lines, "\n")
}
