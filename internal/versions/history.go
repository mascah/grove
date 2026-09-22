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
	Status      string // "-" where the commit deleted the file, "" when unreadable
	When        time.Time
	Source      []byte // the record's bytes at that commit; nil where deleted or unreadable
}

// HistoryContext lists the commits reachable from commit that changed the
// record at path (relative to root, as Version.Path is), newest first,
// following renames: what plain git log --follow lists. It says what happened
// on that line of history and nothing about any other branch. Merges are never
// listed, so content that a merge itself gave the record (a conflict
// resolution) has no row, and the first row need not be the record as it is
// at commit; a caller that knows the record there can tell. Once ctx is done, running Git processes are killed and
// the error is ctx.Err().
func HistoryContext(ctx context.Context, root, commit, path string) ([]Commit, error) {
	// --raw names the record's blob at each commit, so a rename needs no path
	// read back from Git. A commit line starts with NUL, which no subject holds.
	// Asking Git to diff merges as well makes --follow take a rename seen from a
	// merge's other parent for the record's own, and lose commits. In date order
	// no commit comes before one made from it, whatever their clocks said.
	out, err := repo.GitContext(ctx, root, "log", "--follow", "--raw", "--no-abbrev", "--no-merges", "--date-order",
		"--no-show-signature", "--no-color", "--format=%x00%H%x00%at%x00%s", commit, "--", ":(literal)"+path)
	if err != nil {
		return nil, err
	}
	var commits []Commit
	var blobs []string
	for line := range strings.SplitSeq(out, "\n") {
		switch {
		case strings.HasPrefix(line, "\x00"):
			fields := strings.SplitN(line[1:], "\x00", 3)
			if len(fields) != 3 {
				return nil, fmt.Errorf("git log: unexpected line %q", line)
			}
			at, err := strconv.ParseInt(fields[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("git log: unexpected line %q", line)
			}
			commits, blobs = append(commits, Commit{ID: fields[0], When: time.Unix(at, 0), Subject: fields[2]}), append(blobs, "")
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
		if id == "" {
			continue
		}
		if strings.Trim(id, "0") == "" {
			commits[i].Status = "-" // the commit deleted the file
			continue
		}
		data, err := o.blob(id)
		if o.failed != nil {
			return nil, o.failed
		}
		if err == nil {
			// Today's validation may reject an old record; its status field
			// is reported regardless.
			r, _ := project.ParseRecord(path, data)
			commits[i].Status, commits[i].Source = r.Status, data
		}
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return commits, nil
}
