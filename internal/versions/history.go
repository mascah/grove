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
}

// HistoryContext lists the commits reachable from commit that touched the
// record at path (relative to root, as Version.Path is), newest first,
// following renames. It says what happened on that line of history and nothing
// about any other branch: a merge is listed only when the record's content
// differs from the row beneath it, as after a conflict resolution. Once ctx is
// done, running Git processes are killed and the error is ctx.Err().
func HistoryContext(ctx context.Context, root, commit, path string) ([]Commit, error) {
	// --raw names the record's blob at each commit, so a rename needs no path
	// read back from Git. A commit line starts with NUL, which no subject holds.
	// A merge is diffed against each parent, so that one differing from any of
	// them is listed (once per such parent) with its blob: a conflict resolved
	// by keeping the first parent's side differs only from the second. Merges
	// that change nothing a reader can see are removed below. In date order no
	// commit comes before one made from it, whatever their clocks said.
	out, err := repo.GitContext(ctx, root, "log", "--follow", "--raw", "--no-abbrev", "--diff-merges=separate", "--date-order",
		"--no-show-signature", "--no-color", "--format=%x00%H%x00%at%x00%P%x00%s", commit, "--", ":(literal)"+path)
	if err != nil {
		return nil, err
	}
	var commits []Commit
	var blobs []string
	var merge []bool
	for line := range strings.SplitSeq(out, "\n") {
		switch {
		case strings.HasPrefix(line, "\x00"):
			fields := strings.SplitN(line[1:], "\x00", 4)
			if len(fields) != 4 {
				return nil, fmt.Errorf("git log: unexpected line %q", line)
			}
			at, err := strconv.ParseInt(fields[1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("git log: unexpected line %q", line)
			}
			if len(commits) != 0 && commits[len(commits)-1].ID == fields[0] {
				continue // the same merge, against its next parent
			}
			commits = append(commits, Commit{ID: fields[0], When: time.Unix(at, 0), Subject: fields[3]})
			blobs, merge = append(blobs, ""), append(merge, strings.Contains(fields[2], " "))
		case strings.HasPrefix(line, ":") && len(commits) != 0:
			// ":<old mode> <new mode> <old ID> <new ID> <change>\t<path>"
			if fields := strings.SplitN(line, " ", 5); len(fields) == 5 {
				blobs[len(blobs)-1] = fields[3]
			}
		}
	}
	// A merge holding the same content as the row beneath it changed nothing
	// a reader can see: it repeats that row, and would name another branch.
	// Any other merge stays, so that each row still follows from the one below
	// and the first is the record as it is at commit.
	kept := 0
	for i := range commits {
		if !merge[i] || i+1 == len(blobs) || blobs[i] != blobs[i+1] {
			commits[kept], blobs[kept] = commits[i], blobs[i]
			kept++
		}
	}
	commits, blobs = commits[:kept], blobs[:kept]
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
			r, _ := project.ParseRecord(path, "", data)
			commits[i].Status = r.Status
		}
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return commits, nil
}
