package update

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"unicode/utf8"

	"github.com/mascah/grove/internal/create"
	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/repo"
)

// ConvertRequest names one Markdown document outside the record root to bring
// in as a new record, with the Type and Title it becomes.
type ConvertRequest struct{ Source, Type, Title, Slug string }

// Conversion is the old-to-new mapping of one converted source. FromPath
// equals Source: the original document stays in place.
type Conversion struct{ From, FromPath, ID, Path string }

// Convert is the deliberate, one-source identity change that ordinary update
// refuses. A document becomes a new record with its bytes as the body and
// formerly: PATH; the original is left in place. Neither bodies nor Markdown
// links are ever rewritten. formerly makes a rerun a refusal, decided before
// any ID is reserved.
func Convert(root string, req ConvertRequest, report io.Writer) (Conversion, error) {
	p, ds := project.Load(root, root)
	if len(ds) != 0 {
		return Conversion{}, fmt.Errorf("the project is not valid; fix it before converting:\n%s", diagnostics(ds))
	}
	if _, err := source(p, req); err != nil {
		return Conversion{}, err
	}
	n, err := create.Allocate(p.Root, p.RecordDir, report)
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
	document, err := source(p, req)
	if err != nil {
		return reserved(err)
	}
	result := Conversion{From: req.Source, FromPath: req.Source, ID: id}
	slug, records := req.Slug, slices.Clone(p.Records)
	if slug == "" {
		slug = create.Slug(req.Title)
	}
	status := ""
	if t := project.Type(req.Type); len(t.Statuses) != 0 {
		status = "status: " + t.Statuses[0] + "\n"
	}
	// No dates: the document had none, and the record model does not invent them.
	content := fmt.Appendf(nil, "---\nid: %q\ntype: %s\ntitle: %q\n%sformerly: %q\n---\n\n%s", id, req.Type, strings.TrimSpace(req.Title), status, req.Source, document)
	result.Path = path.Join(filepath.ToSlash(p.RecordDir), id+"-"+slug+".md")
	converted, ds := project.ParseRecord(result.Path, content)
	if len(ds) != 0 {
		return reserved(fmt.Errorf("the converted record would be invalid:\n%s", diagnostics(ds)))
	}
	records = append(records, converted)
	if ds := project.Validate(records); len(ds) != 0 {
		return reserved(fmt.Errorf("the conversion would leave the project invalid:\n%s", diagnostics(ds)))
	}
	// Everything is validated; only now is anything written.
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
	if _, ds := project.Load(root, root); len(ds) != 0 {
		return result, fmt.Errorf("the project no longer validates:\n%s (%s was written as %s, so the conversion is incomplete; inspect with Git)", diagnostics(ds), result.Path, id)
	}
	return result, nil
}

// source resolves the request against p: the document's bytes. Every refusal
// that makes a rerun harmless is here, so it is decided before an ID is
// reserved and again under the lock.
func source(p *project.Project, req ConvertRequest) ([]byte, error) {
	if req.Slug != "" && !create.ValidSlug(req.Slug) {
		return nil, fmt.Errorf("slug must contain only lowercase ASCII letters, digits, and hyphens")
	}
	for _, r := range p.Records {
		// Without case, as the loader compares paths: a case-insensitive
		// filesystem opens docs/Plan.md as docs/plan.md.
		if strings.EqualFold(r.Formerly, req.Source) {
			return nil, fmt.Errorf("%s was already converted to %s in %s", req.Source, r.ID, r.Path)
		}
	}
	clean := path.Clean(filepath.ToSlash(req.Source))
	root := path.Clean(filepath.ToSlash(p.RecordDir)) + "/"
	switch t := project.Type(req.Type); {
	case clean != req.Source || !filepath.IsLocal(req.Source) || path.Ext(clean) != ".md":
		return nil, fmt.Errorf("the source must be a project-relative .md file as a clean path")
	case strings.HasPrefix(strings.ToLower(clean), strings.ToLower(root)) || strings.EqualFold(clean, p.Brief):
		return nil, fmt.Errorf("%s is not a legacy document: it is inside the record root or is the brief", req.Source)
	case t == nil || strings.TrimSpace(req.Title) == "":
		return nil, fmt.Errorf("converting a document requires --type (work, question, decision, term, plan, review, or page) and a nonempty --title")
	}
	document, err := project.ReadConfined(p.Root, clean)
	if err != nil {
		return nil, err
	}
	if !utf8.Valid(document) {
		return nil, fmt.Errorf("%s is not valid UTF-8", req.Source)
	}
	return bytes.TrimPrefix(document, []byte("\ufeff")), nil // a BOM belongs at a file's start, not after frontmatter
}
