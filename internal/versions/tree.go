package versions

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"testing/fstest"

	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/repo"
)

// tree is one commit's project at the selected prefix, loaded once per commit.
type tree struct {
	present bool // grove.yaml exists at the prefix
	project *project.Project
	ds      []project.Diagnostic
	err     error // Git could not supply the tree
}

func (t *tree) valid() bool { return t.err == nil && t.present && len(t.ds) == 0 }

func loadTree(ctx context.Context, root, commit string, cache map[string]*tree) *tree {
	if t, ok := cache[commit]; ok {
		return t
	}
	t := &tree{}
	cache[commit] = t
	fsys, present, err := treeFS(ctx, root, commit)
	if err != nil {
		t.err = err
		return t
	}
	t.present = present
	if present {
		t.project, t.ds = project.LoadFS(fsys)
	}
	return t
}

// treeFS reads the committed tree below root's repository prefix into an
// in-memory fs.ReadLinkFS. ls-tree run from root limits output to that prefix
// and names entries relative to it. Symlinks, submodules, and directories keep
// their type bits so the shared loader judges them as it judges a checkout.
// ponytail: reads every .md blob under the prefix, not only the record root;
// restrict to the configured folder if projects with large doc trees make
// versions slow.
func treeFS(ctx context.Context, root, commit string) (fstest.MapFS, bool, error) {
	out, err := repo.GitContext(ctx, root, "ls-tree", "-r", "-t", "-z", commit)
	if err != nil {
		return nil, false, err
	}
	fsys := fstest.MapFS{}
	var shas, names []string
	for _, entry := range strings.Split(out, "\x00") {
		if entry == "" {
			continue
		}
		meta, name, ok := strings.Cut(entry, "\t")
		// With -t, the prefix directory and each of its parents are listed
		// too ("./", "../", "../../") when root is below the repository root.
		// Only those end in a slash; a real entry named "..." does not.
		if strings.HasSuffix(name, "/") && strings.Trim(name, "./") == "" {
			continue
		}
		fields := strings.Fields(meta)
		if !ok || len(fields) != 3 || !fs.ValidPath(name) {
			return nil, false, fmt.Errorf("git ls-tree: unexpected entry %q", entry)
		}
		file := &fstest.MapFile{}
		switch fields[0] {
		case "040000":
			file.Mode = fs.ModeDir
		case "120000":
			file.Mode = fs.ModeSymlink
		case "100644", "100755":
			if name == "grove.yaml" || path.Ext(name) == ".md" {
				shas, names = append(shas, fields[2]), append(names, name)
			}
		default:
			file.Mode = fs.ModeIrregular
		}
		fsys[name] = file
	}
	if _, ok := fsys["grove.yaml"]; !ok {
		return fsys, false, nil
	}
	contents, err := catFile(ctx, root, shas)
	if err != nil {
		return nil, false, err
	}
	for i, name := range names {
		fsys[name].Data = contents[i]
	}
	return fsys, true, nil
}

// catFile fetches blobs in one git process and returns them in request order.
func catFile(ctx context.Context, root string, shas []string) ([][]byte, error) {
	if len(shas) == 0 {
		return nil, nil
	}
	cmd := exec.CommandContext(ctx, "git", "-C", root, "cat-file", "--batch")
	cmd.WaitDelay = repo.GitWaitDelay
	cmd.Stdin = strings.NewReader(strings.Join(shas, "\n") + "\n")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	if err != nil {
		return nil, fmt.Errorf("git cat-file: %s", strings.TrimSpace(stderr.String()))
	}
	contents := make([][]byte, len(shas))
	for i, sha := range shas {
		header, rest, ok := bytes.Cut(out, []byte("\n"))
		fields := strings.Fields(string(header))
		size := -1
		if len(fields) == 3 {
			size, _ = strconv.Atoi(fields[2])
		}
		if !ok || len(fields) != 3 || fields[0] != sha || fields[1] != "blob" || size < 0 || len(rest) < size+1 {
			return nil, fmt.Errorf("git cat-file: unexpected reply %q for %s", header, sha)
		}
		contents[i], out = rest[:size:size], rest[size+1:]
	}
	return contents, nil
}
