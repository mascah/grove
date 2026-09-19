// Package create allocates shared sequential IDs and writes new records.
// The reader stays Git-free; Git access goes through internal/repo.
package create

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/repo"
)

type kind struct{ prefix, folder, status, body string }

var kinds = map[string]kind{
	"work":     {"W", "work", "proposed", "## Outcome\n\n## Constraints\n\n## Acceptance\n\n## Next\n"},
	"question": {"Q", "questions", "open", "## Question\n\n## Next\n"},
	"decision": {"D", "decisions", "proposed", "## Decision\n\n## Alternatives\n\n## Reconsideration\n"},
}

var (
	slugPattern = regexp.MustCompile(`^[a-z0-9-]+$`)
	idLine      = regexp.MustCompile(`^id:\s*["']?([WQD])-([0-9]+)["']?\s*$`)
)

// New allocates the next ID for kind, creates the record without overwriting
// anything, and reloads the project so an unreadable result fails loudly.
// It returns the created file's path relative to the project root.
func New(p *project.Project, kindName, title, slug string, now time.Time, report io.Writer) (string, error) {
	k, ok := kinds[kindName]
	if !ok {
		return "", fmt.Errorf("record type must be work, question, or decision")
	}
	title = strings.TrimSpace(title)
	if title == "" {
		return "", fmt.Errorf("a nonempty title is required")
	}
	if slug == "" {
		slug = Slug(title)
	} else if !slugPattern.MatchString(slug) {
		return "", fmt.Errorf("slug must contain only lowercase ASCII letters, digits, and hyphens")
	}
	common, showPrefix, err := repo.CommonDir(p.Root)
	if err != nil {
		return "", err
	}
	n, err := allocate(common, p.Root, p.RecordDir, showPrefix, k.prefix, report)
	if err != nil {
		return "", err
	}
	id := fmt.Sprintf("%s-%03d", k.prefix, n)
	relative := path.Join(filepath.ToSlash(p.RecordDir), k.folder, id+"-"+slug+".md")
	stamp := now.UTC().Format("2006-01-02T15:04:05Z")
	// %q emits Go escapes, a subset of YAML double-quoted escapes.
	content := fmt.Sprintf("---\nid: %q\ntype: %s\ntitle: %q\nstatus: %s\ncreated: %q\nupdated: %q\n---\n\n%s",
		id, kindName, title, k.status, stamp, stamp, k.body)
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

func diagnostics(ds []project.Diagnostic) string {
	lines := make([]string, len(ds))
	for i, d := range ds {
		lines[i] = d.String()
	}
	return strings.Join(lines, "\n")
}

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

// Allocate reserves the next number for prefix in the Git repository containing
// root. Every worktree and record type shares one lock and one counter file
// under the Git common directory; the issued number is never at or below an ID
// found in any local ref or worktree's live records.
func Allocate(root, recordDir, prefix string, report io.Writer) (int, error) {
	common, showPrefix, err := repo.CommonDir(root)
	if err != nil {
		return 0, err
	}
	return allocate(common, root, recordDir, showPrefix, prefix, report)
}

func allocate(common, root, recordDir, showPrefix, prefix string, report io.Writer) (int, error) {
	unlock, err := repo.AllocatorLock(common)
	if err != nil {
		return 0, err
	}
	defer unlock()
	counterPath := filepath.Join(common, "grove", "next-ids")
	counters, err := readCounters(counterPath)
	if err != nil {
		return 0, err
	}
	floor, err := highestUsed(root, recordDir, showPrefix, prefix)
	if err != nil {
		return 0, err
	}
	next := floor + 1
	current, known := counters[prefix]
	switch {
	case !known:
		fmt.Fprintf(report, "Initialized the %s counter at %d from records in local refs and worktrees; reservations for records never written or since deleted cannot be recovered.\n", prefix, next)
	case current < next:
		fmt.Fprintf(report, "The %s counter (%d) was below records in use; continuing at %d.\n", prefix, current, next)
	default:
		next = current
	}
	counters[prefix] = next + 1
	if err := writeCounters(counterPath, counters); err != nil {
		return 0, fmt.Errorf("no ID issued: %w", err)
	}
	return next, nil
}

func readCounters(path string) (map[string]int, error) {
	counters := map[string]int{}
	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return counters, nil
	}
	if err != nil {
		return nil, err
	}
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		fields := strings.Fields(line)
		n := 0
		if len(fields) == 2 {
			n, err = strconv.Atoi(fields[1])
		}
		if len(fields) != 2 || err != nil || n < 1 || len(fields[0]) != 1 || !strings.Contains("WQD", fields[0]) {
			return nil, fmt.Errorf("%s is corrupt (%q); fix or remove it to reinitialize from existing records", path, line)
		}
		counters[fields[0]] = n
	}
	return counters, nil
}

func writeCounters(path string, counters map[string]int) error {
	tmp, err := os.CreateTemp(filepath.Dir(path), "next-ids-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())
	for _, prefix := range slices.Sorted(maps.Keys(counters)) {
		if _, err := fmt.Fprintf(tmp, "%s %d\n", prefix, counters[prefix]); err != nil {
			tmp.Close()
			return err
		}
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
func highestUsed(root, recordDir, showPrefix, prefix string) (int, error) {
	highest := 0
	note := func(line string) {
		m := idLine.FindStringSubmatch(line)
		if m == nil || m[1] != prefix {
			return
		}
		if n, err := strconv.Atoi(m[2]); err == nil && n > highest {
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
	worktrees, err := repo.Git(root, "worktree", "list", "--porcelain")
	if err != nil {
		return 0, err
	}
	for _, line := range strings.Split(worktrees, "\n") {
		wt, ok := strings.CutPrefix(line, "worktree ")
		if !ok {
			continue
		}
		tree := filepath.Join(wt, filepath.FromSlash(showPrefix), recordDir)
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
