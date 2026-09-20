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
	"github.com/yuin/goldmark/util"

	"github.com/mascah/grove/internal/project"
)

// sources collects each included file once, within the byte budget. Nothing is
// truncated or dropped to fit: a source that does not fit is an error.
type sources struct {
	dir       *os.Root
	remaining int
	byPath    map[string]*Source
	files     []sourceFile // every source, by identity: a case-insensitive filesystem or a hard link reaches one file by several names
}

type sourceFile struct {
	info   fs.FileInfo
	source *Source
}

// find returns the source that already holds this file, under this name or,
// given its identity, any other.
func (s *sources) find(name string, info fs.FileInfo) *Source {
	if existing := s.byPath[name]; existing != nil {
		return existing
	}
	for _, f := range s.files {
		if info != nil && os.SameFile(f.info, info) {
			return f.source
		}
	}
	return nil
}

// known adds reason to the source that already holds this file and reports
// whether there was one. Callers ask before reading or charging, so another
// name for an included file costs nothing.
func (s *sources) known(name string, info fs.FileInfo, reason string) bool {
	existing := s.find(name, info)
	if existing != nil && !slices.Contains(existing.Reasons, reason) {
		existing.Reasons = append(existing.Reasons, reason)
	}
	return existing != nil
}

func (s *sources) add(name string, info fs.FileInfo, content []byte, reason string) error {
	if len(content) > s.remaining {
		return oversize(name, s.remaining)
	}
	s.remaining -= len(content)
	source := &Source{Path: name, Revision: project.Revision(content), Reasons: []string{reason}, Content: string(content)}
	s.byPath[name] = source
	s.files = append(s.files, sourceFile{info, source})
	return nil
}

// addLoaded includes bytes the project loader already read and checked.
func (s *sources) addLoaded(name string, content []byte, reason string) error {
	if s.known(name, nil, reason) {
		return nil
	}
	info, err := s.dir.Lstat(name)
	if err != nil {
		return err
	}
	if s.known(name, info, reason) {
		return nil
	}
	return s.add(name, info, content, reason)
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
	if s.known(name, nil, reason) {
		return nil
	}
	info, err := statConfined(s.dir, name)
	if errors.Is(err, fs.ErrNotExist) {
		return fmt.Errorf("%s names %s, which does not exist in this checkout", referrer, name)
	}
	if err != nil {
		return err
	}
	if s.known(name, info, reason) {
		return nil
	}
	content, err := readConfined(s.dir, name, info, s.remaining)
	if err != nil {
		return err
	}
	return s.add(name, info, content, reason)
}

// holds reports whether the file at name is a source, under any spelling.
func (s *sources) holds(name string) bool {
	info, _ := s.dir.Lstat(name)
	return s.find(name, info) != nil
}

// references lists every link in r's body with what it resolves to. It opens
// nothing: a linked document is a source only if it was included another way.
func (s *sources) references(r *project.Record) ([]Reference, error) {
	var refs []Reference
	for _, destination := range links(body(r.Source)) {
		target, reason, err := resolve(r.Path, destination)
		if err != nil {
			return nil, err
		}
		if reason == "" && s.byPath[target] != nil {
			reason = "included in full; a fragment or query was not applied or checked"
		} else if reason == "" {
			reason = "in-project path; not opened or checked, --include adds it"
		}
		if ref := (Reference{From: r.Path, Target: destination, Path: target, Reason: reason}); !slices.Contains(refs, ref) {
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

// resolve turns a link written in the record at from into the project path it
// names, or the reason it names none. Nothing is read or checked here.
func resolve(from, destination string) (target, reason string, err error) {
	// Goldmark keeps the destination as written. Markdown's backslash escapes
	// and entities come off first, as its own renderer does; url.Parse then
	// decodes percent escapes once.
	decoded := util.ResolveEntityNames(util.ResolveNumericReferences(util.UnescapePunctuations([]byte(destination))))
	u, parseErr := url.Parse(string(decoded))
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
	switch {
	case target == ".." || strings.HasPrefix(target, "../"):
		return "", "outside the selected project; not followed", nil
	case gitMetadata(target):
		return "", "Git metadata; not followed", nil
	}
	return target, "", nil
}

// gitMetadata covers a .git directory and a linked worktree's .git file, in
// any case a case-insensitive filesystem would accept.
func gitMetadata(name string) bool {
	return slices.ContainsFunc(strings.Split(name, "/"), func(part string) bool { return strings.EqualFold(part, ".git") })
}

// statConfined returns the identity of the regular file name below dir,
// refusing a symlink in any component, as the record reader does. os.Root
// keeps a concurrently swapped parent from leading outside dir.
func statConfined(dir *os.Root, name string) (fs.FileInfo, error) {
	var seen fs.FileInfo
	for i, c := range name + "/" {
		if c != '/' {
			continue
		}
		info, err := dir.Lstat(name[:i])
		if err != nil {
			return nil, err
		}
		if info.Mode()&fs.ModeSymlink != 0 {
			return nil, fmt.Errorf("%s: %s is a symlink; symlinked sources are not supported", name, name[:i])
		}
		seen = info
	}
	if !seen.Mode().IsRegular() {
		return nil, fmt.Errorf("%s: expected a regular file", name)
	}
	return seen, nil
}

// readConfined reads the UTF-8 file that statConfined saw, of at most limit
// bytes. Opening without blocking means a file swapped for a FIFO is refused
// by the descriptor check instead of hanging.
func readConfined(dir *os.Root, name string, seen fs.FileInfo, limit int) ([]byte, error) {
	f, err := dir.OpenFile(name, os.O_RDONLY|syscall.O_NOFOLLOW|syscall.O_NONBLOCK, 0)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	opened, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if !opened.Mode().IsRegular() || !os.SameFile(seen, opened) {
		return nil, fmt.Errorf("%s: changed while it was being read", name)
	}
	if opened.Size() > int64(limit) {
		return nil, oversize(name, limit)
	}
	content, err := io.ReadAll(io.LimitReader(f, int64(limit)+1))
	if err != nil {
		return nil, err
	}
	if len(content) > limit {
		return nil, oversize(name, limit)
	}
	if !utf8.Valid(content) {
		return nil, fmt.Errorf("%s: invalid UTF-8", name)
	}
	return content, nil
}
