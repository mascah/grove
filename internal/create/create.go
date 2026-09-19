// Package create allocates shared sequential IDs and writes new records.
// It is the only package that runs Git; the reader stays Git-free.
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
	gitIDLine   = `^id:[[:space:]]*["']?%s-[0-9]+["']?[[:space:]]*$`
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
	n, err := Allocate(p.Root, p.RecordDir, k.prefix, report)
	if err != nil {
		return "", err
	}
	id := fmt.Sprintf("%s-%03d", k.prefix, n)
	relative := path.Join(filepath.ToSlash(p.RecordDir), k.folder, id+"-"+slug+".md")
	stamp := now.UTC().Format("2006-01-02T15:04:05Z")
	content := fmt.Sprintf("---\nid: %q\ntype: %s\ntitle: %s\nstatus: %s\ncreated: %q\nupdated: %q\n---\n\n%s",
		id, kindName, strconv.Quote(title), k.status, stamp, stamp, k.body)
	full := filepath.Join(p.Root, filepath.FromSlash(relative))
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
		lines := make([]string, len(ds))
		for i, d := range ds {
			lines[i] = d.String()
		}
		return "", fmt.Errorf("created %s but the project no longer validates:\n%s", relative, strings.Join(lines, "\n"))
	}
	return relative, nil
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
	if _, err := exec.LookPath("git"); err != nil {
		return 0, errors.New("creating records requires Git on PATH; allocation state lives in the repository's common directory")
	}
	out, err := gitOutput(root, "rev-parse", "--path-format=absolute", "--git-common-dir", "--show-prefix")
	if err != nil {
		return 0, fmt.Errorf("creating records requires a Git repository; allocation state lives in its common directory (%w)", err)
	}
	lines := strings.SplitN(strings.TrimRight(out, "\n"), "\n", 2)
	common, showPrefix := lines[0], ""
	if len(lines) == 2 {
		showPrefix = lines[1]
	}
	dir := filepath.Join(common, "grove")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return 0, err
	}
	unlock, err := lock(filepath.Join(dir, "lock"))
	if err != nil {
		return 0, err
	}
	defer unlock()
	counterPath := filepath.Join(dir, "next-ids")
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

func gitOutput(dir string, args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git %s: %s", args[0], strings.TrimSpace(stderr.String()))
	}
	return string(out), nil
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
		if len(fields) != 2 || err != nil || n < 1 || !slices.Contains([]string{"W", "Q", "D"}, fields[0]) {
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
	refs, err := gitOutput(root, "for-each-ref", "--format=%(objectname)", "refs/heads", "refs/remotes")
	if err != nil {
		return 0, err
	}
	if trees := strings.Fields(refs); len(trees) != 0 {
		args := append([]string{"-C", root, "grep", "-h", "-I", "-E", fmt.Sprintf(gitIDLine, prefix)}, trees...)
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
	worktrees, err := gitOutput(root, "worktree", "list", "--porcelain")
	if err != nil {
		return 0, err
	}
	for _, line := range strings.Split(worktrees, "\n") {
		wt, ok := strings.CutPrefix(line, "worktree ")
		if !ok {
			continue
		}
		tree := filepath.Join(wt, filepath.FromSlash(showPrefix), recordDir)
		_ = filepath.WalkDir(tree, func(p string, entry fs.DirEntry, err error) error {
			if err != nil || entry.IsDir() || filepath.Ext(p) != ".md" {
				return nil
			}
			data, err := os.ReadFile(p)
			if err != nil {
				return nil
			}
			for _, line := range strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n") {
				note(line)
			}
			return nil
		})
	}
	return highest, nil
}
