// Package create allocates shared sequential IDs and writes new records.
// The reader stays Git-free; Git access goes through internal/repo.
package create

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/repo"
)

// bodies are the skeletons new writes; the rest of a type is project.Types.
var bodies = map[string]string{
	"work":     "## Outcome\n\n## Constraints\n\n## Acceptance\n\n## Next\n",
	"question": "## Question\n\n## Next\n",
	"decision": "## Decision\n\n## Alternatives\n\n## Reconsideration\n",
	"term":     "## Meaning\n\n## Relationships\n",
	"plan":     "## Design\n\n## Steps\n",
	"review":   "## Examined\n\n## Findings\n\n## Disposition\n",
	"page":     "",
}

var (
	slugPattern = regexp.MustCompile(`^[a-z0-9-]+$`)
	idLine      = regexp.MustCompile(`^id:\s*["']?` + project.NeutralPrefix + `-([0-9]+)["']?\s*$`)
)

// New allocates the next ID for kind, creates the record without overwriting
// anything, and reloads the project so an unreadable result fails loudly.
// It returns the created file's path relative to the project root.
func New(p *project.Project, kindName, title, slug string, now time.Time, report io.Writer) (string, error) {
	k := project.Type(kindName)
	if k == nil {
		return "", fmt.Errorf("record type must be work, question, decision, term, plan, review, or page")
	}
	title = strings.TrimSpace(title)
	if title == "" {
		return "", fmt.Errorf("a nonempty title is required")
	}
	// A term's title is its identity, so new can collide where other types
	// cannot. Refuse before reserving; a term that appears after this load is
	// refused again under the write lock, before anything is written.
	if err := definedTerm(p, kindName, title); err != nil {
		return "", err
	}
	if slug == "" {
		slug = Slug(title)
	} else if !ValidSlug(slug) {
		return "", fmt.Errorf("slug must contain only lowercase ASCII letters, digits, and hyphens")
	}
	common, showPrefix, err := repo.CommonDir(p.Root)
	if err != nil {
		return "", err
	}
	// Every type creates flat in the record root under the one neutral ID.
	n, err := allocate(common, p.Root, p.RecordDir, showPrefix, report)
	if err != nil {
		return "", err
	}
	id := fmt.Sprintf("%s-%03d", project.NeutralPrefix, n)
	relative := path.Join(filepath.ToSlash(p.RecordDir), id+"-"+slug+".md")
	stamp := now.UTC().Format("2006-01-02T15:04:05Z")
	status := ""
	if len(k.Statuses) != 0 {
		status = "status: " + k.Statuses[0] + "\n"
	}
	// %q emits Go escapes, a subset of YAML double-quoted escapes.
	content := fmt.Sprintf("---\nid: %q\ntype: %s\ntitle: %q\n%screated: %q\nupdated: %q\n---\n\n%s",
		id, kindName, title, status, stamp, stamp, bodies[kindName])
	full := filepath.Join(p.Root, filepath.FromSlash(relative))
	// The reservation is already durable, so a failure from here on consumes
	// it. Publication is serialized with update through the shared write lock,
	// which is only taken after the allocator lock was released.
	unlock, err := repo.WriteLock(common)
	if err != nil {
		return "", fmt.Errorf("%s reserved but not created: %w", id, err)
	}
	defer unlock()
	current, ds := project.Load(p.Root, p.Root)
	if len(ds) != 0 {
		return "", fmt.Errorf("%s reserved but not created: the project no longer validates:\n%s", id, diagnostics(ds))
	}
	if current.RecordDir != p.RecordDir {
		return "", fmt.Errorf("%s reserved but not created: the record root changed from %s to %s during allocation", id, p.RecordDir, current.RecordDir)
	}
	if err := definedTerm(current, kindName, title); err != nil {
		return "", fmt.Errorf("%s reserved but not created: %w", id, err)
	}
	if !bytes.Equal(current.Config, p.Config) {
		return "", fmt.Errorf("%s reserved but not created: grove.yaml changed during allocation; inspect it and retry", id)
	}
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return "", fmt.Errorf("%s reserved but not created: %w", id, err)
	}
	f, err := os.OpenFile(full, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return "", fmt.Errorf("%s reserved but not created: %w", id, err)
	}
	if _, err := f.WriteString(content); err != nil {
		f.Close()
		return "", fmt.Errorf("%s: %w", relative, err)
	}
	if err := f.Close(); err != nil {
		return "", fmt.Errorf("%s: %w", relative, err)
	}
	if _, ds := project.Load(p.Root, p.Root); len(ds) != 0 {
		return "", fmt.Errorf("created %s but the project no longer validates:\n%s", relative, diagnostics(ds))
	}
	return relative, nil
}

func definedTerm(p *project.Project, kindName, title string) error {
	for _, r := range p.Records {
		if kindName == "term" && r.Type == "term" && project.TermKey(r.Title) == project.TermKey(title) {
			return fmt.Errorf("the term %s is already defined by %s in %s", title, r.ID, r.Path)
		}
	}
	return nil
}

func diagnostics(ds []project.Diagnostic) string {
	lines := make([]string, len(ds))
	for i, d := range ds {
		lines[i] = d.String()
	}
	return strings.Join(lines, "\n")
}

// ValidSlug reports whether a caller-supplied slug is acceptable in a filename.
func ValidSlug(slug string) bool { return slugPattern.MatchString(slug) }

// Slug derives a short filename fragment from a title per the record model.
func Slug(title string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(title) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		case b.Len() > 0 && !strings.HasSuffix(b.String(), "-"):
			b.WriteByte('-')
		}
	}
	s := b.String()
	if len(s) > 32 {
		s = s[:32]
	}
	s = strings.TrimRight(s, "-")
	if s == "" {
		return "record"
	}
	return s
}

// Allocate reserves the next number in the Git repository containing root.
// Every worktree shares one lock under the Git common directory; the issued
// number is never at or below an ID found in any local ref or worktree's live
// records. An existing next-ids file from before the one neutral counter is
// simply ignored.
func Allocate(root, recordDir string, report io.Writer) (int, error) {
	common, showPrefix, err := repo.CommonDir(root)
	if err != nil {
		return 0, err
	}
	return allocate(common, root, recordDir, showPrefix, report)
}

func allocate(common, root, recordDir, showPrefix string, report io.Writer) (int, error) {
	unlock, err := repo.AllocatorLock(common)
	if err != nil {
		return 0, err
	}
	defer unlock()
	counterPath := filepath.Join(common, "grove", "neutral-ids")
	current, known, err := readCounter(counterPath)
	if err != nil {
		return 0, err
	}
	floor, err := highestUsed(root, recordDir, showPrefix)
	if err != nil {
		return 0, err
	}
	next := floor + 1
	switch {
	case !known:
		fmt.Fprintf(report, "Initialized the %s counter at %d from records in local refs and worktrees; reservations for records never written or since deleted cannot be recovered.\n", project.NeutralPrefix, next)
	case current < next:
		fmt.Fprintf(report, "The %s counter (%d) was below records in use; continuing at %d.\n", project.NeutralPrefix, current, next)
	default:
		next = current
	}
	if err := writeCounter(counterPath, next+1); err != nil {
		return 0, fmt.Errorf("no ID issued: %w", err)
	}
	return next, nil
}

func readCounter(path string) (int, bool, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return 0, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	fields := strings.Fields(strings.TrimSpace(string(data)))
	var n int
	if len(fields) == 2 {
		n, err = strconv.Atoi(fields[1])
	}
	if len(fields) != 2 || err != nil || n < 1 || fields[0] != project.NeutralPrefix {
		return 0, false, fmt.Errorf("%s is corrupt (%q); fix or remove it to reinitialize from existing records", path, strings.TrimSpace(string(data)))
	}
	return n, true, nil
}

func writeCounter(path string, n int) error {
	tmp, err := os.CreateTemp(filepath.Dir(path), "neutral-ids-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	if _, err := fmt.Fprintf(tmp, "%s %d\n", project.NeutralPrefix, n); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}

// highestUsed scans committed records on every local ref and live records in
// every worktree. Lines that merely look like IDs only raise the floor.
// ponytail: one git grep over all refs per allocation; scan only on
// initialization or mismatch if repositories with many refs make this slow.
func highestUsed(root, recordDir, showPrefix string) (int, error) {
	highest := 0
	note := func(line string) {
		m := idLine.FindStringSubmatch(line)
		if m == nil {
			return
		}
		if n, err := strconv.Atoi(m[1]); err == nil && n > highest {
			highest = n
		}
	}
	refs, err := repo.Git(root, "for-each-ref", "--format=%(objectname)", "refs/heads", "refs/remotes", "refs/tags")
	if err != nil {
		return 0, err
	}
	if trees := strings.Fields(refs); len(trees) != 0 {
		// git grep only pre-filters; note validates every line.
		args := append([]string{"-C", root, "grep", "-h", "-I", "-e", "^id:"}, trees...)
		cmd := exec.Command("git", append(args, "--", recordDir)...)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		out, err := cmd.Output()
		var exit *exec.ExitError
		if err != nil && !(errors.As(err, &exit) && exit.ExitCode() == 1) { // 1 means no matches
			return 0, fmt.Errorf("git grep: %s", strings.TrimSpace(stderr.String()))
		}
		for _, line := range strings.Split(string(out), "\n") {
			note(line)
		}
	}
	worktrees, err := repo.Worktrees(root)
	if err != nil {
		return 0, err
	}
	for _, w := range worktrees {
		tree := filepath.Join(w.Path, filepath.FromSlash(showPrefix), recordDir)
		err := filepath.WalkDir(tree, func(p string, entry fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() || filepath.Ext(p) != ".md" {
				return nil
			}
			data, err := os.ReadFile(p)
			if err != nil {
				return err
			}
			for _, line := range strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n") {
				note(line)
			}
			return nil
		})
		// A worktree without the record folder is normal; anything else would
		// silently lower the floor and risk reissuing an ID.
		if err != nil && !(errors.Is(err, fs.ErrNotExist) && strings.HasPrefix(err.Error(), "lstat "+tree)) {
			return 0, fmt.Errorf("cannot scan worktree records: %w", err)
		}
	}
	return highest, nil
}
