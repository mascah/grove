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
// about any other branch: a merge is listed only when it gave the record
// content that no listed commit did, as a conflict resolution does. Once ctx is
// done, running Git processes are killed and the error is ctx.Err().
func HistoryContext(ctx context.Context, root, commit, path string) ([]Commit, error) {
	// --raw names the record's blob at each commit, so a rename needs no path
	// read back from Git. A commit line starts with NUL, which no subject holds.
	// Diffing a merge against its first parent gives it a blob too, and keeps
	// merges that path limiting would drop; those are removed below.
	out, err := repo.GitContext(ctx, root, "log", "--follow", "--raw", "--no-abbrev", "--diff-merges=first-parent",
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
			commits = append(commits, Commit{ID: fields[0], When: time.Unix(at, 0), Subject: fields[3]})
			blobs, merge = append(blobs, ""), append(merge, strings.Contains(fields[2], " "))
		case strings.HasPrefix(line, ":") && len(commits) != 0:
			// ":<old mode> <new mode> <old ID> <new ID> <change>\t<path>"
			if fields := strings.SplitN(line, " ", 5); len(fields) == 5 {
				blobs[len(blobs)-1] = fields[3]
			}
		}
	}
	// A merge that only brought in a listed commit's content repeats that
	// commit, and would name another branch where nothing else does.
	authored := map[string]bool{}
	for i, id := range blobs {
		authored[id] = authored[id] || !merge[i]
	}
	kept := 0
	for i := range commits {
		if !merge[i] || !authored[blobs[i]] {
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
