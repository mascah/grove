// Package project reads and validates one checkout without modifying its files.
package project

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type Project struct {
	Root      string
	RecordDir string // configured record folder, relative to Root
	Records   []*Record
}

// Load reads a whole project. Any diagnostics make the project unsuitable for
// presentation as valid; the returned records are for inspection by validators.
func Load(cwd, explicit string) (*Project, []Diagnostic) {
	root, err := discover(cwd, explicit)
	if err != nil {
		return nil, []Diagnostic{{Path: cwd, Field: "grove.yaml", Message: err.Error()}}
	}
	p := &Project{Root: root}
	configPath := filepath.Join(root, "grove.yaml")
	source, err := readRegular(configPath)
	if err != nil {
		return p, []Diagnostic{{Path: "grove.yaml", Message: err.Error()}}
	}
	config := parseMapping("grove.yaml", source, 0)
	version, ok := config.integerField("schema_version", true)
	if ok && version != 1 {
		config.problem("schema_version", fmt.Sprintf("unsupported version %d; expected 1", version))
	}
	recordDir := config.stringField("records", true)
	for key := range config.fields {
		if key != "schema_version" && key != "records" {
			config.problem(key, "unknown configuration key")
		}
	}
	if recordDir != "" {
		if filepath.IsAbs(recordDir) || filepath.Clean(recordDir) == "." || slices.Contains(strings.Split(filepath.ToSlash(recordDir), "/"), "..") {
			config.problem("records", "must name a dedicated relative subdirectory without .. components")
		} else if err := checkRecordRoot(root, recordDir); err != nil {
			config.problem("records", err.Error())
		}
	}
	if len(config.errors) != 0 {
		return p, sortedDiagnostics(config.errors)
	}
	p.RecordDir = recordDir
	recordRoot := filepath.Join(root, recordDir)
	var ds []Diagnostic
	folders := map[string]string{"work": "work", "questions": "question", "decisions": "decision"}
	err = filepath.WalkDir(recordRoot, func(path string, entry fs.DirEntry, walkErr error) error {
		relative, _ := filepath.Rel(root, path)
		relative = filepath.ToSlash(relative)
		problem := func(message string) {
			ds = append(ds, Diagnostic{Path: relative, Message: message})
		}
		if walkErr != nil {
			problem(walkErr.Error())
			return nil
		}
		if entry.Type()&os.ModeSymlink != 0 {
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
		inside, _ := filepath.Rel(recordRoot, path)
		parts := strings.Split(filepath.ToSlash(inside), "/")
		if len(parts) == 1 && folders[parts[0]] != "" {
			problem("type folder must be a directory")
			return nil
		}
		if filepath.Ext(path) != ".md" {
			return nil
		}
		kind := folders[parts[0]]
		if len(parts) < 2 || kind == "" {
			problem("Markdown record must be inside a work, questions, or decisions type folder")
			return nil
		}
		source, err := readRegular(path)
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
		ds = append(ds, Diagnostic{Path: filepath.ToSlash(recordDir), Message: err.Error()})
	}
	slices.SortFunc(p.Records, compareRecords)
	ds = append(ds, Validate(p.Records)...)
	return p, sortedDiagnostics(ds)
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

func checkRecordRoot(root, relative string) error {
	current := root
	for _, part := range strings.Split(filepath.Clean(relative), string(filepath.Separator)) {
		current = filepath.Join(current, part)
		info, err := os.Lstat(current)
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("%s is a symlink", relative)
		}
		if !info.IsDir() {
			return fmt.Errorf("%s must be a directory", relative)
		}
	}
	return nil
}

func readRegular(path string) ([]byte, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return nil, err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return nil, fmt.Errorf("symlink files are not supported")
	}
	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("expected a regular file")
	}
	return os.ReadFile(path)
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
