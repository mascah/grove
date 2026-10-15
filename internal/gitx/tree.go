package gitx

import (
	"bytes"
	"io/fs"
	"strconv"
	"strings"
	"testing/fstest"
)

// TreeFS returns a read-only file system over the regular files committed at
// commit under subdir, a slash-separated path from the repository root.
// Symbolic links and submodules are left out. A missing subdir gives an
// empty file system.
func TreeFS(dir, commit, subdir string) (fs.FS, error) {
	out, err := run(command(dir, "ls-tree", "-r", "-z", "--full-tree", "--end-of-options", commit, subdir))
	if err != nil {
		return nil, err
	}
	// ponytail: the whole subtree is loaded into memory up front and served
	// by fstest.MapFS. Fine for a folder of small Markdown documents; read
	// blobs on Open through one long-lived `git cat-file --batch` if a root
	// ever holds large files.
	files := fstest.MapFS{}
	var names []string
	var ids strings.Builder
	// Each record is "<mode> <type> <id>\t<path>"; -r lists leaves only.
	for record := range strings.SplitSeq(string(out), "\x00") {
		meta, name, ok := strings.Cut(record, "\t")
		fields := strings.Fields(meta)
		if !ok || (fields[0] != "100644" && fields[0] != "100755") {
			continue // the empty piece after the last terminator, a link or a submodule
		}
		name = strings.TrimPrefix(name, subdir+"/")
		files[name] = &fstest.MapFile{Mode: 0o644}
		names = append(names, name)
		ids.WriteString(fields[2] + "\n")
	}
	if len(names) == 0 {
		return files, nil
	}
	batch := command(dir, "cat-file", "--batch")
	batch.Stdin = strings.NewReader(ids.String())
	if out, err = run(batch); err != nil {
		return nil, err
	}
	// Each reply is "<id> blob <size>\n<contents>\n", in the order asked.
	for _, name := range names {
		header, rest, _ := bytes.Cut(out, []byte("\n"))
		size, _ := strconv.Atoi(string(header[bytes.LastIndexByte(header, ' ')+1:]))
		files[name].Data = rest[:size:size]
		out = rest[size+1:]
	}
	return files, nil
}
