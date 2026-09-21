// Package project reads and validates one checkout without modifying its files.
package project

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
)

type Project struct {
	Root      string
	RecordDir string // configured record folder, relative to Root
	Schema    int    // schema_version; it decides which record types exist
	Config    []byte // exact grove.yaml bytes
	Brief     string // configured brief, relative to Root with forward slashes; "" when none
	Records   []*Record
}

// Load reads a whole project. Any diagnostics make the project unsuitable for
// presentation as valid; the returned records are for inspection by validators.
func Load(cwd, explicit string) (*Project, []Diagnostic) {
	root, d := Discover(cwd, explicit)
	if d != nil {
		return nil, []Diagnostic{*d}
	}
	p, ds := LoadFS(os.DirFS(root))
	p.Root = root
	// LoadFS judges only the brief's path, because a committed source holds just
	// grove.yaml and the record folder. A live checkout must have the file.
	if p.Brief != "" {
		if _, err := readConfined(root, p.Brief); err != nil {
			ds = sortedDiagnostics(append(ds, Diagnostic{Path: "grove.yaml", Field: "brief", Message: err.Error()}))
		}
	}
	return p, ds
}

// LoadFS validates the project at the root of fsys, which must implement
// fs.ReadLinkFS so symlinks are seen rather than followed. A live checkout and
// a committed Git tree go through this one path, so both obey the same rules.
func LoadFS(fsys fs.FS) (*Project, []Diagnostic) {
	p := &Project{}
	source, err := readRegular(fsys, "grove.yaml")
	if err != nil {
		return p, []Diagnostic{{Path: "grove.yaml", Message: err.Error()}}
	}
	p.Config = source
	config := parseMapping("grove.yaml", source, 0)
	version, ok := config.integerField("schema_version", true)
	if ok && version != 1 && version != 2 {
		config.problem("schema_version", fmt.Sprintf("unsupported version %d; expected 1 or 2", version))
	}
	recordDir := config.stringField("records", true)
	for key := range config.fields {
		if key != "schema_version" && key != "records" && !(key == "brief" && version >= 2) {
			config.problem(key, "unknown configuration key")
		}
	}
	brief := ""
	if version >= 2 {
		brief = config.stringField("brief", false)
	}
	if recordDir != "" {
		if !dedicated(recordDir) {
			config.problem("records", "must name a dedicated relative subdirectory without .. components")
		} else if err := checkRecordRoot(fsys, recordDir); err != nil {
			config.problem("records", err.Error())
		}
	}
	if len(config.errors) != 0 {
		return p, sortedDiagnostics(config.errors)
	}
	p.RecordDir, p.Schema = recordDir, version
	recordRoot := path.Clean(filepath.ToSlash(recordDir))
	folders := map[string]string{}
	var names []string
	for _, t := range Types {
		if t.Schema <= version {
			folders[t.Folder] = t.Name
			names = append(names, t.Folder)
		}
	}
	if brief != "" {
		clean := path.Clean(filepath.ToSlash(brief))
		inside := strings.Split(strings.TrimPrefix(clean, recordRoot+"/"), "/")
		switch {
		case !dedicated(brief) || clean != filepath.ToSlash(brief) || path.Ext(clean) != ".md":
			config.problem("brief", "must name a project-relative .md file as a clean path without .. components")
		case strings.HasPrefix(clean, recordRoot+"/") && len(inside) > 1 && folders[inside[0]] != "":
			config.problem("brief", "must not be inside a record type folder")
		}
		if len(config.errors) != 0 {
			return p, sortedDiagnostics(config.errors)
		}
		p.Brief = clean
	}
	var ds []Diagnostic
	err = fs.WalkDir(fsys, recordRoot, func(relative string, entry fs.DirEntry, walkErr error) error {
		problem := func(message string) {
			ds = append(ds, Diagnostic{Path: relative, Message: message})
		}
		if walkErr != nil {
			problem(walkErr.Error())
			return nil
		}
		if entry.Type()&fs.ModeSymlink != 0 {
			problem("symlink entries are not supported")
			return nil
		}
		if entry.IsDir() {
			return nil
		}
		if !entry.Type().IsRegular() {
			problem("expected a regular file or directory")
			return nil
		}
		inside := strings.TrimPrefix(relative, recordRoot+"/")
		parts := strings.Split(inside, "/")
		if len(parts) == 1 && folders[parts[0]] != "" {
			problem("type folder must be a directory")
			return nil
		}
		if path.Ext(relative) != ".md" {
			return nil
		}
		if relative == p.Brief {
			return nil
		}
		kind := folders[parts[0]]
		if len(parts) < 2 || kind == "" {
			problem("Markdown record must be inside a type folder: " + strings.Join(names, ", "))
			return nil
		}
		source, err := readRegular(fsys, relative)
		if err != nil {
			problem(err.Error())
			return nil
		}
		record, problems := ParseRecord(relative, kind, source)
		p.Records = append(p.Records, record)
		ds = append(ds, problems...)
		return nil
	})
	if err != nil {
		ds = append(ds, Diagnostic{Path: recordRoot, Message: err.Error()})
	}
	slices.SortFunc(p.Records, compareRecords)
	ds = append(ds, Validate(p.Records)...)
	return p, sortedDiagnostics(ds)
}

func dedicated(recordDir string) bool {
	return !filepath.IsAbs(recordDir) && path.Clean(recordDir) != "." && fs.ValidPath(path.Clean(recordDir)) && !slices.Contains(strings.Split(filepath.ToSlash(recordDir), "/"), "..")
}

// RecordRoot returns the record folder that LoadFS would walk for this
// grove.yaml, or "" when it would refuse the value without reading a folder.
// A caller that gives LoadFS part of a tree uses it to know which part: LoadFS
// reads grove.yaml, each component of this path, and everything below it.
func RecordRoot(config []byte) string {
	recordDir := parseMapping("grove.yaml", config, 0).stringField("records", false)
	if recordDir == "" || !dedicated(recordDir) {
		return ""
	}
	return path.Clean(filepath.ToSlash(recordDir))
}

// Discover finds the project root as Load does, without reading the project.
func Discover(cwd, explicit string) (string, *Diagnostic) {
	root, err := discover(cwd, explicit)
	if err != nil {
		return "", &Diagnostic{Path: cwd, Field: "grove.yaml", Message: err.Error()}
	}
	return root, nil
}

func discover(cwd, explicit string) (string, error) {
	start, err := filepath.Abs(cwd)
	if err != nil {
		return "", err
	}
	if explicit != "" {
		if !filepath.IsAbs(explicit) {
			explicit = filepath.Join(start, explicit)
		}
		start = filepath.Clean(explicit)
	} else {
		// Canonicalize the invocation directory (e.g. macOS /var -> /private/var)
		// before walking parents. Symlinks inside a project's tree are checked later.
		start, err = filepath.EvalSymlinks(start)
		if err != nil {
			return "", err
		}
	}
	for dir := start; ; dir = filepath.Dir(dir) {
		info, err := os.Stat(dir)
		if err != nil {
			return "", err
		}
		if !info.IsDir() {
			return "", fmt.Errorf("%s is not a directory", dir)
		}
		_, err = os.Lstat(filepath.Join(dir, "grove.yaml"))
		if err == nil {
			return dir, nil
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return "", err
		}
		if explicit != "" {
			return "", fmt.Errorf("no grove.yaml in explicitly selected project %s", dir)
		}
		_, gitErr := os.Lstat(filepath.Join(dir, ".git"))
		if gitErr == nil || dir == filepath.Dir(dir) {
			return "", fmt.Errorf("no grove.yaml found before project search boundary %s", dir)
		}
		if !errors.Is(gitErr, fs.ErrNotExist) {
			return "", gitErr
		}
	}
}

func checkRecordRoot(fsys fs.FS, relative string) error {
	current := ""
	for _, part := range strings.Split(path.Clean(filepath.ToSlash(relative)), "/") {
		current = path.Join(current, part)
		info, err := fs.Lstat(fsys, current)
		if err != nil {
			return err
		}
		if info.Mode()&fs.ModeSymlink != 0 {
			return fmt.Errorf("%s is a symlink", relative)
		}
		if !info.IsDir() {
			return fmt.Errorf("%s must be a directory", relative)
		}
	}
	return nil
}

// ReadBrief returns the configured brief's bytes from the live checkout.
func (p *Project) ReadBrief() ([]byte, error) {
	if p.Brief == "" {
		return nil, errors.New("grove.yaml names no brief; add a brief: PATH key (schema_version 2)")
	}
	return readConfined(p.Root, p.Brief)
}

// readConfined reads a regular file below root, refusing a symlink in any
// component so the path cannot leave the project.
func readConfined(root, name string) ([]byte, error) {
	fsys := os.DirFS(root)
	for i, c := range name {
		if c == '/' {
			if info, err := fs.Lstat(fsys, name[:i]); err != nil {
				return nil, err
			} else if info.Mode()&fs.ModeSymlink != 0 {
				return nil, fmt.Errorf("%s is a symlink", name[:i])
			}
		}
	}
	return readRegular(fsys, name)
}

func readRegular(fsys fs.FS, name string) ([]byte, error) {
	info, err := fs.Lstat(fsys, name)
	if err != nil {
		return nil, err
	}
	if info.Mode()&fs.ModeSymlink != 0 {
		return nil, fmt.Errorf("symlink files are not supported")
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("expected a regular file")
	}
	return fs.ReadFile(fsys, name)
}

func compareRecords(a, b *Record) int {
	if (a.Created == nil) != (b.Created == nil) {
		if a.Created == nil {
			return 1
		}
		return -1
	}
	if a.Created != nil && b.Created != nil {
		if c := a.Created.Compare(*b.Created); c != 0 {
			return c
		}
	}
	// Canonical IDs use the same two-character prefix and zero padding.
	// Length comparison avoids both lexical W-1000 < W-999 and integer overflow.
	if len(a.ID) >= 2 && len(b.ID) >= 2 && a.ID[:2] == b.ID[:2] && len(a.ID) != len(b.ID) {
		if len(a.ID) < len(b.ID) {
			return -1
		}
		return 1
	}
	if c := strings.Compare(a.ID, b.ID); c != 0 {
		return c
	}
	return strings.Compare(a.Path, b.Path)
}

func sortedDiagnostics(ds []Diagnostic) []Diagnostic {
	slices.SortFunc(ds, func(a, b Diagnostic) int {
		if c := strings.Compare(a.Path, b.Path); c != 0 {
			return c
		}
		if a.Line != b.Line {
			return a.Line - b.Line
		}
		if c := strings.Compare(a.Field, b.Field); c != 0 {
			return c
		}
		return strings.Compare(a.Message, b.Message)
	})
	return ds
}
