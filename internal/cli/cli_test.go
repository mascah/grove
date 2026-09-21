package cli

import (
	"bytes"
	"crypto/sha256"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

const work = "---\nid: G-001\ntype: work\ntitle: Inspect records\nstatus: proposed\nrelates_to: [G-002]\n---\nAn outcome.\n"
const question = "---\nid: G-002\ntype: question\ntitle: Which version?\nstatus: open\nblocks: []\n---\nAn uncertainty.\n"

func write(t *testing.T, root, path, source string) {
	t.Helper()
	path = filepath.Join(root, path)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(source), 0644); err != nil {
		t.Fatal(err)
	}
}

func projectFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	real, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatal(err)
	}
	write(t, real, "grove.yaml", "schema_version: 3\nrecords: docs/records\n")
	write(t, real, "docs/records/work/renamed.md", work)
	write(t, real, "docs/records/questions/question.md", question)
	return real
}

func hashes(t *testing.T, root string) map[string][32]byte {
	t.Helper()
	result := map[string][32]byte{}
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		result[path] = sha256.Sum256(data)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return result
}

func TestCommandsInspectWithoutChangingFiles(t *testing.T) {
	t.Parallel()
	root := projectFixture(t)
	before := hashes(t, root)
	for _, args := range [][]string{{"list"}, {"show", "G-001"}, {"check"}, {"list", "--project", root}, {"--project=" + root, "show", "G-002"}} {
		var out, errOut bytes.Buffer
		code := Run(args, filepath.Join(root, "docs", "records"), &out, &errOut)
		if code != 0 || !strings.Contains(errOut.String(), root) {
			t.Fatalf("%v: code=%d, stderr=%s", args, code, errOut.String())
		}
		switch args[0] {
		case "list":
			if !strings.Contains(out.String(), "G-001") || !strings.Contains(out.String(), "proposed") || !strings.Contains(out.String(), "Inspect records") || !strings.Contains(out.String(), "question") {
				t.Fatal(out.String())
			}
		case "show":
			if out.String() != work || !strings.Contains(errOut.String(), "docs/records/work/renamed.md") {
				t.Fatal("show must retain complete original bytes and identify file")
			}
		case "check":
			if !strings.Contains(out.String(), "2 records") {
				t.Fatal(out.String())
			}
		default:
			if out.String() != question {
				t.Fatal(out.String())
			}
		}
	}
	if !reflect.DeepEqual(before, hashes(t, root)) {
		t.Fatal("inspection changed project files")
	}
}

func TestInvalidNeighborPreventsPartialOutput(t *testing.T) {
	t.Parallel()
	root := projectFixture(t)
	write(t, root, "docs/records/work/broken.md", "---\nid: G-003\ntype: work\ntitle: Broken\nstatus: imaginary\n---\n")
	before := hashes(t, root)
	for _, args := range [][]string{{"list"}, {"show", "G-001"}, {"check"}} {
		var out, errOut bytes.Buffer
		if code := Run(args, root, &out, &errOut); code != 1 {
			t.Fatalf("%v returned %d", args, code)
		}
		if out.Len() != 0 || !strings.Contains(errOut.String(), "broken.md") || !strings.Contains(errOut.String(), "status") {
			t.Fatalf("partial result or missing diagnostic: stdout=%s stderr=%s", out.String(), errOut.String())
		}
	}
	if !reflect.DeepEqual(before, hashes(t, root)) {
		t.Fatal("failed inspection changed project files")
	}
}

func TestUsageAndMissingID(t *testing.T) {
	t.Parallel()
	// No command at all selects the board; see TestBoardInvocation.
	for _, args := range [][]string{{"unknown"}, {"show"}, {"show", "G-001", "extra"}, {"list", "extra"}, {"--project"}, {"list", "--wat"}, {"--project=", "list"}, {"--project", "a", "--project", "b", "list"}} {
		var out, errOut bytes.Buffer
		if code := Run(args, t.TempDir(), &out, &errOut); code != 2 || out.Len() != 0 || !strings.Contains(errOut.String(), "Usage:") {
			t.Fatalf("%v: code=%d stderr=%s", args, code, errOut.String())
		}
	}
	for _, args := range [][]string{{"--help"}, {"help"}, {"list", "-h"}} {
		var out, errOut bytes.Buffer
		if code := Run(args, t.TempDir(), &out, &errOut); code != 0 || !strings.Contains(out.String(), "Usage:") {
			t.Fatalf("%v: code=%d stdout=%s", args, code, out.String())
		}
	}
	root := projectFixture(t)
	var out, errOut bytes.Buffer
	if code := Run([]string{"show", "G-999"}, root, &out, &errOut); code != 1 || out.Len() != 0 || !strings.Contains(errOut.String(), "G-999") {
		t.Fatalf("missing identity: code=%d stderr=%s", code, errOut.String())
	}
}

func TestEmptyProjectAndLiteralSource(t *testing.T) {
	t.Parallel()
	root := projectFixture(t)
	if err := os.RemoveAll(filepath.Join(root, "docs", "records")); err != nil {
		t.Fatal(err)
	}
	os.MkdirAll(filepath.Join(root, "docs", "records"), 0755)
	var out, errOut bytes.Buffer
	if code := Run([]string{"check"}, root, &out, &errOut); code != 0 || !strings.Contains(out.String(), "0 records") {
		t.Fatalf("empty project: code=%d, stderr=%s", code, errOut.String())
	}
	source := strings.ReplaceAll(question, "\n", "\r\n")
	source = strings.TrimSuffix(source, "\r\n")
	write(t, root, "docs/records/questions/custom.md", source)
	out.Reset()
	errOut.Reset()
	if code := Run([]string{"show", "G-002"}, root, &out, &errOut); code != 0 || out.String() != source {
		t.Fatal("show normalized line endings or final newline")
	}
}

type brokenWriter struct{}

func (brokenWriter) Write([]byte) (int, error) { return 0, errors.New("output unavailable") }

func TestOutputFailureReturnsNonzero(t *testing.T) {
	t.Parallel()
	root := projectFixture(t)
	for _, args := range [][]string{{"list"}, {"show", "G-001"}, {"check"}, {"--help"}} {
		var errOut bytes.Buffer
		if code := Run(args, root, brokenWriter{}, &errOut); code != 1 {
			t.Fatalf("%v: code=%d", args, code)
		}
	}
	var out bytes.Buffer
	if code := Run([]string{"list"}, root, &out, brokenWriter{}); code != 1 {
		t.Fatalf("context output failed but command returned %d", code)
	}
}

func TestListEscapesMultilineAndControlCharacters(t *testing.T) {
	t.Parallel()
	root := projectFixture(t)
	write(t, root, "docs/records/questions/question.md", strings.Replace(question, "title: Which version?", "title: \"First\\nSecond\\t\\e[31m\"", 1))
	var out, errOut bytes.Buffer
	if code := Run([]string{"list"}, root, &out, &errOut); code != 0 {
		t.Fatal(errOut.String())
	}
	if strings.Count(out.String(), "\n") != 3 || strings.Contains(out.String(), "\x1b") {
		t.Fatalf("title broke the terminal table: %q", out.String())
	}
}
