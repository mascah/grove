package versions

import (
	"bufio"
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"testing/fstest"

	"github.com/mascah/grove/internal/project"
	"github.com/mascah/grove/internal/repo"
)

// tree is one commit's project at the selected prefix.
type tree struct {
	present bool // grove.yaml exists at the prefix
	project *project.Project
	ds      []project.Diagnostic
	err     error // Git could not supply the tree
}

func (t *tree) valid() bool { return t.err == nil && t.present && len(t.ds) == 0 }

// entry is one line of a Git tree object.
type entry struct{ mode, name, id string }

func (e entry) isDir() bool { return e.mode == "40000" }

func (e entry) file() *fstest.MapFile {
	switch e.mode {
	case "40000":
		return &fstest.MapFile{Mode: fs.ModeDir}
	case "120000":
		return &fstest.MapFile{Mode: fs.ModeSymlink}
	case "100644", "100755":
		return &fstest.MapFile{}
	}
	return &fstest.MapFile{Mode: fs.ModeIrregular} // a submodule
}

// objects reads one inspection's committed projects through a single git
// cat-file process, asked only for object IDs, so a branch costs a few
// lookups rather than processes and no path ever needs quoting. It reads only
// what the project loader reads (grove.yaml, the path to the record folder,
// and that folder), so the size of the rest of the repository does not matter,
// and it loads each distinct project once however many branches hold it.
type objects struct {
	ctx          context.Context
	root, prefix string
	cmd          *exec.Cmd
	in           io.WriteCloser
	out          *bufio.Reader
	stderr       bytes.Buffer
	failed       error // the process is unusable; every later read fails alike
	idLen        int   // bytes in a raw object ID: 20, or 32 under SHA-256

	commits  map[string]*tree // by commit
	projects map[string]*tree // by what the loader would read
	trees    map[string][]entry
	blobs    map[string][]byte
}

func newObjects(ctx context.Context, root, prefix string) *objects {
	return &objects{ctx: ctx, root: root, prefix: prefix, commits: map[string]*tree{}, projects: map[string]*tree{}, trees: map[string][]entry{}, blobs: map[string][]byte{}}
}

// close ends the process, if one was started. Closing its input is how
// cat-file is asked to finish; a cancelled context has killed it already.
func (o *objects) close() {
	if o.cmd != nil {
		o.in.Close()
		o.cmd.Wait()
	}
}

// read returns the object named by spec, which is an object ID or an ID with
// a peeling suffix, never a path.
func (o *objects) read(spec, kind string) ([]byte, error) {
	if o.failed == nil && o.cmd == nil {
		cmd := repo.Command(o.ctx, o.root, "cat-file", "--batch")
		cmd.WaitDelay = repo.WaitDelay(o.ctx)
		cmd.Stderr = &o.stderr
		in, err := cmd.StdinPipe()
		var out io.ReadCloser
		if err == nil {
			out, err = cmd.StdoutPipe()
		}
		if err == nil {
			err = cmd.Start()
		}
		if err != nil {
			o.failed = fmt.Errorf("git cat-file: %w", err)
		} else {
			o.cmd, o.in, o.out = cmd, in, bufio.NewReader(out)
		}
	}
	if o.failed != nil {
		return nil, o.failed
	}
	data, err := o.exchange(spec, kind)
	if o.failed != nil && o.ctx.Err() != nil {
		o.failed = o.ctx.Err()
	}
	if o.failed != nil {
		return nil, o.failed
	}
	return data, err
}

// exchange asks for one object. A reply that is not the wanted object fails
// only this read; losing the conversation sets failed.
func (o *objects) exchange(spec, kind string) ([]byte, error) {
	broken := func(err error) {
		// Git's words are complete, and safe to read, only once it has ended.
		o.close()
		o.cmd = nil
		if said := strings.TrimSpace(o.stderr.String()); said != "" {
			err = fmt.Errorf("%s", said)
		}
		o.failed = fmt.Errorf("git cat-file: %w", err)
	}
	if _, err := io.WriteString(o.in, spec+"\n"); err != nil {
		broken(err)
		return nil, nil
	}
	header, err := o.out.ReadString('\n')
	if err != nil {
		broken(err)
		return nil, nil
	}
	fields := strings.Fields(header)
	if len(fields) != 3 { // "<spec> missing", or ambiguous
		return nil, fmt.Errorf("git cat-file: %s", strings.TrimSpace(header))
	}
	size, err := strconv.Atoi(fields[2])
	if err != nil || size < 0 {
		broken(fmt.Errorf("unexpected reply %q", strings.TrimSpace(header)))
		return nil, nil
	}
	data := make([]byte, size+1)
	if _, err := io.ReadFull(o.out, data); err != nil {
		broken(err)
		return nil, nil
	}
	if fields[1] != kind {
		return nil, fmt.Errorf("git cat-file: %s is a %s, not a %s", spec, fields[1], kind)
	}
	return data[:size:size], nil
}

// entries lists the tree named by spec, cached by its ID when spec is one.
func (o *objects) entries(spec string) ([]entry, error) {
	if es, ok := o.trees[spec]; ok {
		return es, nil
	}
	data, err := o.read(spec, "tree")
	if err != nil {
		return nil, err
	}
	// Each entry is "<mode> <name>\0<raw ID>".
	raw := o.idLen
	var es []entry
	for len(data) != 0 {
		mode, rest, ok := bytes.Cut(data, []byte(" "))
		name, rest, ok2 := bytes.Cut(rest, []byte{0})
		if !ok || !ok2 || len(rest) < raw {
			return nil, fmt.Errorf("git cat-file: malformed tree %s", spec)
		}
		es = append(es, entry{string(mode), string(name), hex.EncodeToString(rest[:raw])})
		data = rest[raw:]
	}
	o.trees[spec] = es
	return es, nil
}

func (o *objects) blob(id string) ([]byte, error) {
	if data, ok := o.blobs[id]; ok {
		return data, nil
	}
	data, err := o.read(id, "blob")
	if err == nil {
		o.blobs[id] = data
	}
	return data, err
}

func entryNamed(es []entry, name string) (entry, bool) {
	for _, e := range es {
		if e.name == name {
			return e, true
		}
	}
	return entry{}, false
}

// loadTree returns commit's project at the prefix. Branches that agree on
// everything the loader reads share one loaded project.
func (o *objects) loadTree(commit string) *tree {
	if t, ok := o.commits[commit]; ok {
		return t
	}
	t := &tree{}
	o.commits[commit], o.idLen = t, len(commit)/2
	fsys, key, err := o.treeFS(commit)
	switch {
	case err != nil:
		t.err = err
	case fsys == nil: // no grove.yaml at the prefix
	case o.projects[key] != nil:
		*t = *o.projects[key]
	default:
		t.present = true
		t.project, t.ds = project.LoadFS(fsys)
		o.projects[key] = t
	}
	return t
}

// treeFS builds the part of commit's tree that the project loader reads, as an
// in-memory fs.ReadLinkFS whose symlinks, submodules, and directories keep
// their type bits so the loader judges them as it judges a checkout. key
// names that part: equal keys mean equal content, types, and names. A nil
// file system means the prefix holds no grove.yaml.
func (o *objects) treeFS(commit string) (fsys fstest.MapFS, key string, err error) {
	dir, err := o.entries(commit + "^{tree}")
	if err != nil {
		return nil, "", err
	}
	for _, part := range strings.Split(strings.TrimSuffix(o.prefix, "/"), "/") {
		if part == "" {
			continue
		}
		e, ok := entryNamed(dir, part)
		if !ok || !e.isDir() {
			return nil, "", nil
		}
		if dir, err = o.entries(e.id); err != nil {
			return nil, "", err
		}
	}
	config, ok := entryNamed(dir, "grove.yaml")
	if !ok {
		return nil, "", nil
	}
	fsys = fstest.MapFS{"grove.yaml": config.file()}
	key = config.mode + " " + config.id
	if !fsys["grove.yaml"].Mode.IsRegular() {
		return fsys, key, nil
	}
	if fsys["grove.yaml"].Data, err = o.blob(config.id); err != nil {
		return nil, "", err
	}
	recordRoot := project.RecordRoot(fsys["grove.yaml"].Data)
	if recordRoot == "" {
		return fsys, key, nil
	}
	// The path to the record folder matters by type only; the folder itself
	// by its ID, which covers everything below it.
	at := ""
	for _, part := range strings.Split(recordRoot, "/") {
		e, ok := entryNamed(dir, part)
		if !ok {
			return fsys, key + "\x00missing", nil
		}
		at = path.Join(at, part)
		fsys[at] = e.file()
		key += "\x00" + e.mode
		if !e.isDir() {
			return fsys, key, nil
		}
		if at == recordRoot {
			key += " " + e.id
		}
		if dir, err = o.entries(e.id); err != nil {
			return nil, "", err
		}
	}
	return fsys, key, o.fill(fsys, recordRoot, dir)
}

// fill adds everything below a directory, with the content of Markdown files.
func (o *objects) fill(fsys fstest.MapFS, at string, dir []entry) error {
	for _, e := range dir {
		name := path.Join(at, e.name)
		if !fs.ValidPath(name) || path.Base(name) != e.name {
			return fmt.Errorf("git cat-file: unexpected entry %q in %s", e.name, at)
		}
		file := e.file()
		fsys[name] = file
		switch {
		case e.isDir():
			sub, err := o.entries(e.id)
			if err == nil {
				err = o.fill(fsys, name, sub)
			}
			if err != nil {
				return err
			}
		case file.Mode.IsRegular() && path.Ext(name) == ".md":
			data, err := o.blob(e.id)
			if err != nil {
				return err
			}
			file.Data = data
		}
	}
	return nil
}
