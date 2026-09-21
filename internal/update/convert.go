package update

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/mascah/grove/internal/create"
	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/repo"
)

// ConvertRequest names one source to bring under a neutral ID: a record's
// typed ID, or the project-relative path of a Markdown document outside the
// record root, which also needs the Type and Title of the record it becomes.
type ConvertRequest struct{ Source, Type, Title, Slug string }

// Conversion is the old-to-new mapping of one converted source. FromPath is
// the file it came from; a converted record's is gone, a document's remains.
type Conversion struct{ From, FromPath, ID, Path string }

// Convert is the deliberate, one-source identity change that ordinary update
// refuses. A record keeps every byte but its id and gains formerly; its file
// moves flat into the record root, and relationship lists naming the old ID
// are rewritten. A document becomes a new record with its bytes as the body.
// Neither path touches updated, bodies, or Markdown links. formerly makes a
// rerun a refusal, decided before any ID is reserved.
func Convert(root string, req ConvertRequest, report io.Writer) (Conversion, error) {
	p, ds := project.Load(root, root)
	if len(ds) != 0 {
		return Conversion{}, fmt.Errorf("the project is not valid; fix it before converting:\n%s", diagnostics(ds))
	}
	if _, _, err := source(p, req); err != nil {
		return Conversion{}, err
	}
	n, err := create.Allocate(p.Root, p.RecordDir, project.NeutralPrefix, report)
	if err != nil {
		return Conversion{}, err
	}
	id := fmt.Sprintf("%s-%03d", project.NeutralPrefix, n)
	reserved := func(err error) (Conversion, error) {
		return Conversion{}, fmt.Errorf("%s reserved but nothing converted: %w", id, err)
	}
	common, _, err := repo.CommonDir(p.Root)
	if err != nil {
		return reserved(err)
	}
	unlock, err := repo.WriteLock(common)
	if err != nil {
		return reserved(err)
	}
	defer unlock()
	if p, ds = project.Load(root, root); len(ds) != 0 {
		return reserved(fmt.Errorf("the project no longer validates:\n%s", diagnostics(ds)))
	}
	old, document, err := source(p, req)
	if err != nil {
		return reserved(err)
	}
	result := Conversion{From: req.Source, FromPath: req.Source, ID: id}
	slug, records := req.Slug, slices.Clone(p.Records)
	rewritten := map[int][]byte{} // referrers, by index in records
	var content []byte
	if old == nil {
		if slug == "" {
			slug = create.Slug(req.Title)
		}
		status := ""
		if t := project.Type(req.Type); len(t.Statuses) != 0 {
			status = "status: " + t.Statuses[0] + "\n"
		}
		// No dates: the document had none, and the record model does not invent them.
		content = fmt.Appendf(nil, "---\nid: %q\ntype: %s\ntitle: %q\n%sformerly: %q\n---\n\n%s", id, req.Type, strings.TrimSpace(req.Title), status, req.Source, document)
	} else {
		result.FromPath = old.Path
		if slug == "" {
			name := strings.TrimSuffix(path.Base(old.Path), ".md")
			slug = create.Slug(strings.TrimPrefix(name, old.ID))
		}
		changes := []change{set("id", strconv.Quote(id)), set("formerly", strconv.Quote(old.ID))}
		if content, err = Edit(old.Source, changes); err != nil {
			return reserved(fmt.Errorf("%s: %w", old.Path, err))
		}
		for i, r := range p.Records {
			var changes []change
			for _, rel := range []struct {
				field string
				ids   []string
			}{{"depends_on", r.DependsOn}, {"members", r.Members}, {"blocks", r.Blocks}, {"relates_to", r.RelatesTo}, {"work", r.Work}} {
				if !slices.Contains(rel.ids, old.ID) {
					continue
				}
				quoted := make([]string, len(rel.ids))
				for j, target := range rel.ids {
					if target == old.ID {
						target = id
					}
					quoted[j] = strconv.Quote(target)
				}
				changes = append(changes, set(rel.field, "["+strings.Join(quoted, ", ")+"]"))
			}
			if len(changes) == 0 {
				continue
			}
			edited, err := Edit(r.Source, changes)
			if err != nil {
				return reserved(fmt.Errorf("%s: %w", r.Path, err))
			}
			next, ds := project.ParseRecord(r.Path, "", p.Schema, edited)
			if len(ds) != 0 {
				return reserved(fmt.Errorf("rewriting %s would leave it invalid:\n%s", r.Path, diagnostics(ds)))
			}
			if err := unchanged(r, next, changes); err != nil {
				return reserved(err)
			}
			records[i], rewritten[i] = next, edited
		}
	}
	result.Path = path.Join(filepath.ToSlash(p.RecordDir), id+"-"+slug+".md")
	converted, ds := project.ParseRecord(result.Path, "", p.Schema, content)
	if len(ds) != 0 {
		return reserved(fmt.Errorf("the converted record would be invalid:\n%s", diagnostics(ds)))
	}
	if old == nil {
		records = append(records, converted)
	} else {
		if err := unchanged(old, converted, []change{set("id", ""), set("formerly", "")}); err != nil {
			return reserved(err)
		}
		records[slices.Index(p.Records, old)] = converted
	}
	if ds := project.Validate(records); len(ds) != 0 {
		return reserved(fmt.Errorf("the conversion would leave the project invalid:\n%s", diagnostics(ds)))
	}
	// Everything is validated; only now is anything written. The new file comes
	// first because it is the step most likely to be refused, which then leaves
	// the project untouched.
	// ponytail: not atomic across files. An interruption from here on leaves a
	// project that fails check (formerly names a record that still exists, or a
	// list names a missing ID); recover with Git. Journal it if that ever bites.
	full := filepath.Join(p.Root, filepath.FromSlash(result.Path))
	f, err := os.OpenFile(full, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return reserved(err)
	}
	if _, err = f.Write(content); err == nil {
		err = f.Sync()
	}
	if closeErr := f.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		os.Remove(full) // O_EXCL made it ours, and nothing else has been touched yet
		return reserved(fmt.Errorf("%s: %w", result.Path, err))
	}
	// The mapping is returned with the error: it is already true of the files.
	partial := func(err error) (Conversion, error) {
		return result, fmt.Errorf("%w (%s was written as %s, so the conversion is incomplete; inspect with Git)", err, result.Path, id)
	}
	if old != nil {
		if err := os.Remove(filepath.Join(p.Root, filepath.FromSlash(old.Path))); err != nil {
			return partial(err)
		}
	}
	for i, edited := range rewritten {
		if err := replace(filepath.Join(p.Root, filepath.FromSlash(p.Records[i].Path)), edited); err != nil {
			return partial(err)
		}
	}
	if _, ds := project.Load(root, root); len(ds) != 0 {
		return partial(fmt.Errorf("the project no longer validates:\n%s", diagnostics(ds)))
	}
	return result, nil
}

// source resolves the request against p: the record to convert, or, for a
// document, nil and its bytes. Every refusal that makes a rerun harmless is
// here, so it is decided before an ID is reserved and again under the lock.
func source(p *project.Project, req ConvertRequest) (*project.Record, []byte, error) {
	if p.Schema < 3 {
		return nil, nil, fmt.Errorf("convert needs schema_version 3 in grove.yaml; this project is schema %d", p.Schema)
	}
	if req.Slug != "" && !create.ValidSlug(req.Slug) {
		return nil, nil, fmt.Errorf("slug must contain only lowercase ASCII letters, digits, and hyphens")
	}
	for _, r := range p.Records {
		// Without case, as the loader compares paths: a case-insensitive
		// filesystem opens docs/Plan.md as docs/plan.md.
		if strings.EqualFold(r.Formerly, req.Source) {
			return nil, nil, fmt.Errorf("%s was already converted to %s in %s", req.Source, r.ID, r.Path)
		}
	}
	if project.IDPattern.MatchString(req.Source) {
		if req.Type != "" || req.Title != "" {
			return nil, nil, fmt.Errorf("--type and --title apply only to converting a document; a record keeps its own")
		}
		if strings.HasPrefix(req.Source, project.NeutralPrefix+"-") {
			return nil, nil, fmt.Errorf("%s already has a neutral ID", req.Source)
		}
		i := slices.IndexFunc(p.Records, func(r *project.Record) bool { return r.ID == req.Source })
		if i < 0 {
			return nil, nil, fmt.Errorf("record %s not found in this project", req.Source)
		}
		return p.Records[i], nil, nil
	}
	clean := path.Clean(filepath.ToSlash(req.Source))
	root := path.Clean(filepath.ToSlash(p.RecordDir)) + "/"
	switch t := project.Type(req.Type); {
	case clean != req.Source || !filepath.IsLocal(req.Source) || path.Ext(clean) != ".md":
		return nil, nil, fmt.Errorf("the source must be a record ID or a project-relative .md file as a clean path")
	case strings.HasPrefix(strings.ToLower(clean), strings.ToLower(root)) || strings.EqualFold(clean, p.Brief):
		return nil, nil, fmt.Errorf("%s is not a legacy document: it is inside the record root or is the brief", req.Source)
	case t == nil || strings.TrimSpace(req.Title) == "":
		return nil, nil, fmt.Errorf("converting a document requires --type (work, question, decision, term, plan, review, or page) and a nonempty --title")
	}
	document, err := project.ReadConfined(p.Root, clean)
	if err != nil {
		return nil, nil, err
	}
	if !utf8.Valid(document) {
		return nil, nil, fmt.Errorf("%s is not valid UTF-8", req.Source)
	}
	return nil, bytes.TrimPrefix(document, []byte("\ufeff")), nil // a BOM belongs at a file's start, not after frontmatter
}

// replace swaps a file's bytes through a synced temporary file beside it.
func replace(full string, data []byte) error {
	info, err := os.Lstat(full)
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(full), "."+filepath.Base(full)+".*.tmp")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err = tmp.Write(data); err == nil {
		err = tmp.Sync()
	}
	if err == nil {
		err = tmp.Chmod(info.Mode().Perm())
	}
	if closeErr := tmp.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}
	return os.Rename(tmp.Name(), full)
}
