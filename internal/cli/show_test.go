package cli

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestShowJSONMatchesExactBytes(t *testing.T) {
	t.Parallel()
	root := projectFixture(t)
	crlf := "\ufeff" + strings.ReplaceAll(strings.Replace(question, "Which version?", "\"Wh\\u00efch — 版本?\"", 1), "\n", "\r\n")
	noFinalNewline := strings.TrimSuffix(crlf, "\r\n")
	write(t, root, "docs/records/questions/question.md", noFinalNewline)
	before := hashes(t, root)
	for _, id := range []string{"G-001", "G-002"} {
		var plain, jsonOut, errOut bytes.Buffer
		if code := Run([]string{"show", id}, root, &plain, &errOut); code != 0 {
			t.Fatal(errOut.String())
		}
		if code := Run([]string{"show", "--json", id}, root, &jsonOut, &errOut); code != 0 {
			t.Fatal(errOut.String())
		}
		var got map[string]any
		if err := json.Unmarshal(jsonOut.Bytes(), &got); err != nil || !bytes.HasSuffix(jsonOut.Bytes(), []byte("\n")) {
			t.Fatalf("invalid JSON %q: %v", jsonOut.String(), err)
		}
		sum := sha256.Sum256(plain.Bytes())
		want := map[string]any{"id": id, "path": "docs/records/work/renamed.md", "revision": "sha256:" + hex.EncodeToString(sum[:]), "source": plain.String()}
		if id == "G-002" {
			want["path"] = "docs/records/questions/question.md"
			if plain.String() != noFinalNewline {
				t.Fatal("plain show altered BOM, CRLF, or the missing final newline")
			}
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("got %v\nwant %v", got, want)
		}
		if !strings.Contains(errOut.String(), "File: ") {
			t.Fatal("file context must stay on stderr")
		}
	}
	if !reflect.DeepEqual(before, hashes(t, root)) {
		t.Fatal("show changed project files")
	}
	if _, err := os.Stat(filepath.Join(root, ".git")); !errors.Is(err, fs.ErrNotExist) {
		t.Fatal("fixture unexpectedly has Git state")
	}
}

func TestShowJSONUsage(t *testing.T) {
	t.Parallel()
	for _, args := range [][]string{{"list", "--json"}, {"check", "--json"}, {"show", "--json", "--json", "G-001"}, {"new", "work", "T", "--json"}} {
		var out, errOut bytes.Buffer
		if code := Run(args, t.TempDir(), &out, &errOut); code != 2 || out.Len() != 0 {
			t.Fatalf("%v: code=%d stderr=%s", args, code, errOut.String())
		}
	}
	root := projectFixture(t)
	var out, errOut bytes.Buffer
	if code := Run([]string{"--json", "show", "G-404"}, root, &out, &errOut); code != 1 || out.Len() != 0 {
		t.Fatalf("missing record: code=%d stdout=%q", code, out.String())
	}
}
