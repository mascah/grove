package versions

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/repo"
)

// Commit is one commit that touched a record's file, with the status the
// record held there. Subject is Git's text: escape it before display.
type Commit struct {
	ID, Subject string
	Status      string // "" when the record's status could not be read there
	When        time.Time
}

// HistoryContext lists the commits reachable from commit that touched the
// record at path (relative to root, as Version.Path is), newest first,
// following renames. It says what happened on that line of history and nothing
// about any other branch. Once ctx is done, running Git processes are killed
// and the error is ctx.Err().
func HistoryContext(ctx context.Context, root, commit, path string) ([]Commit, error) {
	// --raw names the record's blob at each commit, so a rename needs no path
	// read back from Git. A commit line starts with NUL, which no subject holds.
	out, err := repo.GitContext(ctx, root, "log", "--follow", "--raw", "--no-abbrev", "--diff-merges=first-parent",
		"--no-show-signature", "--no-color", "--format=%x00%H %at %s", commit, "--", ":(literal)"+path)
	if err != nil {
		return nil, err
	}
	var commits []Commit
	var blobs []string
	for line := range strings.SplitSeq(out, "\n") {
		switch {
		case strings.HasPrefix(line, "\x00"):
			fields := strings.SplitN(line[1:], " ", 3)
			if len(fields) < 2 {
				return nil, fmt.Errorf("git log: unexpected line %q", line)
			}
			at, err := strconv.ParseInt(fields[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("git log: unexpected line %q", line)
			}
			c := Commit{ID: fields[0], When: time.Unix(at, 0)}
			if len(fields) == 3 {
				c.Subject = fields[2]
			}
			commits, blobs = append(commits, c), append(blobs, "")
		case strings.HasPrefix(line, ":") && len(commits) != 0:
			// ":<old mode> <new mode> <old ID> <new ID> <change>\t<path>"
			if fields := strings.SplitN(line, " ", 5); len(fields) == 5 {
				blobs[len(blobs)-1] = fields[3]
			}
		}
	}
	o := newObjects(ctx, root, "")
	defer o.close()
	for i, id := range blobs {
		if id == "" || strings.Trim(id, "0") == "" { // no diff shown, or deleted there
			continue
		}
		data, err := o.blob(id)
		if o.failed != nil {
			return nil, o.failed
		}
		if err == nil {
			// Today's validation may reject an old record; its status field
			// is reported regardless.
			r, _ := project.ParseRecord(path, "", data)
			commits[i].Status = r.Status
		}
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return commits, nil
}
