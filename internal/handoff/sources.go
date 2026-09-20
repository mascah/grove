package handoff

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/url"
	"os"
	"path"
	"slices"
	"strings"
	"syscall"
	"unicode/utf8"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"

	"github.com/mascah/grove/internal/project"
)

// sources collects each included file once, within the byte budget. Nothing is
// truncated or dropped to fit: a source that does not fit is an error.
type sources struct {
	dir       *os.Root
	remaining int
	byPath    map[string]*Source
	files     []readFile // by identity: a case-insensitive filesystem reaches one file by several spellings
}

type readFile struct {
	info fs.FileInfo
	name string
}

func (s *sources) add(name string, content []byte, reason string) error {
	if existing := s.byPath[name]; existing != nil {
		if !slices.Contains(existing.Reasons, reason) {
			existing.Reasons = append(existing.Reasons, reason)
		}
		return nil
	}
	if len(content) > s.remaining {
		return oversize(name, s.remaining)
	}
	s.remaining -= len(content)
	s.byPath[name] = &Source{Path: name, Revision: project.Revision(content), Reasons: []string{reason}, Content: string(content)}
	return nil
}

func oversize(name string, remaining int) error {
	return fmt.Errorf("the selection does not fit the source budget: %d bytes were left when %s was reached; nothing is truncated, so raise --max-bytes (at most %d) or select less", remaining, name, LimitMaxBytes)
}

func (s *sources) addInclude(name string) error {
	if !fs.ValidPath(name) || name == "." || strings.ContainsRune(name, 0) {
		return fmt.Errorf("--include %s: expected a clean project-relative file path", name)
	}
	if gitMetadata(name) {
		return fmt.Errorf("--include %s: Git metadata is not a source", name)
	}
	return s.read(name, "included by the caller", "--include")
}

func (s *sources) read(name, reason, referrer string) error {
	if s.byPath[name] != nil {
		return s.add(name, nil, reason)
	}
	content, info, err := readConfined(s.dir, name, s.remaining)
	if errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("%s names %s, which does not exist in this checkout", referrer, name)
	}
	if err != nil {
		return err
	}
	for _, f := range s.files {
		if os.SameFile(f.info, info) {
			return s.add(f.name, nil, reason)
		}
	}
	s.files = append(s.files, readFile{info, name})
	return s.add(name, content, reason)
}

// addLinked includes the .md and .txt documents that r's body links to
// directly, and returns every other link with the reason it was left out.
func (s *sources) addLinked(r *project.Record) ([]Reference, error) {
	var refs []Reference
	for _, destination := range links(body(r.Source)) {
		target, reason, err := resolve(r.Path, destination)
		if err != nil {
			return nil, err
		}
		if reason == "" {
			if err := s.read(target, "linked from "+r.Path, r.Path); err != nil {
				return nil, err
			}
			if !strings.ContainsAny(destination, "#?") {
				continue
			}
			reason = "the whole document is included; the fragment or query was not applied or checked"
		}
		if ref := (Reference{From: r.Path, Target: destination, Reason: reason}); !slices.Contains(refs, ref) {
			refs = append(refs, ref)
		}
	}
	slices.SortFunc(refs, func(a, b Reference) int { return strings.Compare(a.Target, b.Target) })
	return refs, nil
}

// body is the record after its frontmatter, under the reader's delimiter rules.
func body(source []byte) []byte {
	_, rest, _ := bytes.Cut(bytes.TrimPrefix(source, []byte("\xef\xbb\xbf")), []byte("\n"))
	for len(rest) > 0 {
		var line []byte
		line, rest, _ = bytes.Cut(rest, []byte("\n"))
		if strings.TrimRight(string(line), " \t\r") == "---" {
			return rest
		}
	}
	return nil
}

// links returns the destinations of real Markdown links. Going through the
// parser keeps code spans, fences, images, and HTML from becoming file reads.
func links(markdown []byte) []string {
	var destinations []string
	doc := goldmark.DefaultParser().Parse(text.NewReader(markdown))
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) { // the walker returns no error
		if link, ok := n.(*ast.Link); ok && entering {
			destinations = append(destinations, string(link.Destination))
		}
		return ast.WalkContinue, nil
	})
	return destinations
}

// resolve turns a link written in the record at from into a project path to
// include, or the reason it is only a reference. Nothing is read here.
func resolve(from, destination string) (target, reason string, err error) {
	u, parseErr := url.Parse(destination)
	if parseErr != nil || strings.ContainsRune(u.Path, 0) {
		return "", "", fmt.Errorf("%s: malformed link destination %q", from, destination)
	}
	switch {
	case u.Scheme != "" || u.Host != "":
		return "", "external URL; not fetched", nil
	case u.Path == "":
		return "", "fragment only", nil
	case strings.HasPrefix(u.Path, "/"):
		return "", "absolute path; not followed", nil
	}
	target = path.Join(path.Dir(from), u.Path) // url.Parse decoded the path once
	switch ext := strings.ToLower(path.Ext(target)); {
	case target == ".." || strings.HasPrefix(target, "../"):
		return "", "outside the selected project; not followed", nil
	case gitMetadata(target):
		return "", "Git metadata; not followed", nil
	case ext != ".md" && ext != ".txt":
		return "", "not a .md or .txt document; not included", nil
	}
	return target, "", nil
}

// gitMetadata covers a .git directory and a linked worktree's .git file, in
// any case a case-insensitive filesystem would accept.
func gitMetadata(name string) bool {
	return slices.ContainsFunc(strings.Split(name, "/"), func(part string) bool { return strings.EqualFold(part, ".git") })
}

// readConfined reads a regular UTF-8 file of at most limit bytes below dir.
// os.Root keeps a concurrently swapped parent from leading outside dir; the
// component checks refuse symlinks that stay inside it too, as the record
// reader does. Opening without blocking means a file swapped for a FIFO is
// refused by the descriptor check instead of hanging.
func readConfined(dir *os.Root, name string, limit int) ([]byte, fs.FileInfo, error) {
	var seen fs.FileInfo
	for i, c := range name + "/" {
		if c != '/' {
			continue
		}
		info, err := dir.Lstat(name[:i])
		if err != nil {
			return nil, nil, err
		}
		if info.Mode()&fs.ModeSymlink != 0 {
			return nil, nil, fmt.Errorf("%s: %s is a symlink; symlinked sources are not supported", name, name[:i])
		}
		seen = info
	}
	if !seen.Mode().IsRegular() {
		return nil, nil, fmt.Errorf("%s: expected a regular file", name)
	}
	f, err := dir.OpenFile(name, os.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0)
	if err != nil {
		return nil, nil, err
	}
	defer f.Close()
	opened, err := f.Stat()
	if err != nil {
		return nil, nil, err
	}
	if !opened.Mode().IsRegular() || !os.SameFile(seen, opened) {
		return nil, nil, fmt.Errorf("%s: changed while it was being read", name)
	}
	if opened.Size() > int64(limit) {
		return nil, nil, oversize(name, limit)
	}
	content, err := io.ReadAll(io.LimitReader(f, int64(limit)+1))
	if err != nil {
		return nil, nil, err
	}
	if len(content) > limit {
		return nil, nil, oversize(name, limit)
	}
	if !utf8.Valid(content) {
		return nil, nil, fmt.Errorf("%s: invalid UTF-8", name)
	}
	return content, opened, nil
}
